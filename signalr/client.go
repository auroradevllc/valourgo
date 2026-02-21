package signalr

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Option func(*Client)

type BackoffFunc func(attempt int) time.Duration

func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) {
		cl.httpClient = c
	}
}

func WithHTTPHeaders(h http.Header) Option {
	return func(cl *Client) {
		cl.headers = h
	}
}

func WithDefaultHandler(h HandlerFunc) Option {
	return func(cl *Client) {
		cl.defaultHandler = h
	}
}

func DefaultBackoff(attempt int) time.Duration {
	d := time.Second * time.Duration(1<<attempt)
	if d > 30*time.Second {
		return 30 * time.Second
	}
	return d
}

func WithConnectHandler(f func()) Option {
	return func(cl *Client) {
		cl.onConnect = f
	}
}

type HandlerFunc func(target string, args []json.RawMessage)

type Client struct {
	url        string
	httpClient *http.Client
	headers    http.Header
	serializer Serializer

	conn    *WSConn
	ctx     context.Context
	cancel  context.CancelFunc
	backoff BackoffFunc

	started bool

	mu             sync.RWMutex
	handlers       map[string]HandlerFunc
	onConnect      func()
	defaultHandler HandlerFunc
	invokes        *InvocationManager
	connected      chan struct{}
	disconnected   chan error
}

func NewClient(url string, opts ...Option) *Client {
	c := &Client{
		url:        url,
		httpClient: http.DefaultClient,
		serializer: jsonSerializer,
		backoff:    DefaultBackoff,
		handlers:   make(map[string]HandlerFunc),
		invokes:    NewInvocationManager(),
		connected:  make(chan struct{}),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) On(method string, fn HandlerFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[method] = fn
}

// Connect will start the SignalR connection
// Note that this behaves like a persistent connection, which will always try to reconnect
// The only exception is once the context is closed/done, it will exit.
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	if c.started {
		c.mu.Unlock()
		return nil
	}
	c.started = true
	c.mu.Unlock()

	// First connection attempt (blocking)
	if err := c.connectOnce(ctx); err != nil {
		c.started = false
		return err
	}

	// Background lifecycle
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go c.run()

	return nil
}

func (c *Client) run() {
	attempt := 0

	for {
		select {
		case err := <-c.disconnected:
			if err == nil {
				return
			}

			delay := c.backoff(attempt)
			attempt++

			time.Sleep(delay)

			for {
				select {
				case <-c.ctx.Done():
					return
				default:
				}

				if err := c.connectOnce(context.Background()); err != nil {
					time.Sleep(c.backoff(attempt))
					attempt++
					continue
				}

				attempt = 0
				break
			}

		case <-c.ctx.Done():
			return
		}
	}
}

// connectOnce attempts a single connection to the server
// The lifecycle is negotiate (which requests info about the SignalR hub), dial, handshake.
func (c *Client) connectOnce(ctx context.Context) error {
	proto, err := c.negotiate(ctx, c.url)

	if err != nil {
		return err
	}

	ws, err := DialWS(ctx, proto.WebSocketURL(), c.headers)

	if err != nil {
		return err
	}

	c.disconnected = make(chan error, 1)
	c.conn = ws

	if err := c.sendHandshake(); err != nil {
		_ = ws.Close()
		return err
	}

	go c.readLoop()
	return nil
}

type Handshake struct {
	Protocol string `json:"protocol"`
	Version  int    `json:"version"`
}

// sendHandshake will send a handshake message to the server
func (c *Client) sendHandshake() error {
	if err := c.Write(Handshake{
		Protocol: "json",
		Version:  1,
	}); err != nil {
		return err
	}

	msg, err := c.conn.Read()

	if err != nil {
		return err
	}

	c.handleMessage(msg)

	return nil
}

// readLoop reads from the connection until an error is handled, which is then sent back to the disconnected chan
func (c *Client) readLoop() {
	for {
		msg, err := c.conn.Read()

		if err != nil {
			select {
			case c.disconnected <- err:
			default:
			}
			return
		}

		c.handleMessage(msg)
	}
}

type Invocation struct {
	Type         int    `json:"type"`
	InvocationID string `json:"invocationId"`
	Target       string `json:"target"`
	Arguments    []any  `json:"arguments"`
}

// Invoke will invoke a method/target, then listen for a response
// In an actual SignalR instance, this wouldn't be json.RawMessage, but something like []byte or similar
func (c *Client) Invoke(method string, args ...any) (<-chan json.RawMessage, error) {
	id, ch := c.invokes.New()

	msg := Invocation{
		Type:         messageTypeInvocation,
		InvocationID: id,
		Target:       method,
		Arguments:    args,
	}

	if err := c.Write(msg); err != nil {
		return nil, err
	}

	return ch, nil
}

// Write will serialize and write an object to the hub
func (c *Client) Write(v any) error {
	var b []byte
	var ok bool

	if b, ok = v.([]byte); !ok {
		serialized, err := c.serializer.Serialize(v)

		if err != nil {
			return err
		}

		b = serialized
	}

	return c.conn.Send(b)
}

// handleMessage handles a message from the hub
func (c *Client) handleMessage(data []byte) {
	frames, err := parseFrames(data)

	if err != nil {
		return
	}

	for _, f := range frames {
		switch f.Type {
		case messageTypeInvocation: // invocation
			if c.defaultHandler != nil {
				c.defaultHandler(f.Target, f.Arguments)
			}

			c.mu.RLock()
			if h := c.handlers[f.Target]; h != nil {
				h(f.Target, f.Arguments)
			}
			c.mu.RUnlock()

		case messageTypeCompletion: // completion
			c.invokes.Resolve(f.InvocationID, f.Result)

		case messageTypePing: // ping
			// ignore
		}
	}
}

// Close will close our connection
func (c *Client) Close() error {
	if c.cancel != nil {
		c.cancel()
	}

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
