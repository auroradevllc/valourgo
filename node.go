package valourgo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/auroradevllc/handler"
	"github.com/orcaman/concurrent-map/v2"
	log "github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc/pool"
	"golang.org/x/sync/errgroup"
)

var (
	ErrAlreadyOpen         = errors.New("connection already open")
	ErrInvalidResponseType = errors.New("invalid response type for request")
)

type Node struct {
	*handler.Handler
	client         *http.Client
	baseAddress    string
	token          string
	rtc            *RTC
	me             *User
	members        cmap.ConcurrentMap[PlanetID, Member]
	planetNodeList cmap.ConcurrentMap[PlanetID, string]
	childNodes     cmap.ConcurrentMap[string, *Node]

	Name    string
	Primary *Node
}

type NodeOption func(*Node)

func WithNodeHandler(h *handler.Handler) NodeOption {
	return func(n *Node) {
		n.Handler = h
	}
}

func NewNode(baseAddress, name, token string, opts ...NodeOption) (*Node, error) {
	client := &http.Client{
		Transport: &headerRoundTripper{
			headers: http.Header{
				"X-Server-Select": {name},
				"Authorization":   {token},
				"User-Agent":      {"Valour-Go (" + runtime.Version() + ")"},
			},
			rt: http.DefaultTransport,
		},
	}

	n := &Node{
		client:         client,
		token:          token,
		baseAddress:    baseAddress,
		members:        cmap.NewStringer[PlanetID, Member](),
		planetNodeList: cmap.NewStringer[PlanetID, string](),
		Name:           name,
	}

	for _, opt := range opts {
		opt(n)
	}

	if n.IsPrimary() {
		n.childNodes = cmap.New[*Node]()
	}

	if n.Handler == nil {
		n.Handler = handler.New()
	}

	// Validate node connection by using api/node/name
	name, err := n.NodeName()

	if err != nil {
		return nil, err
	}

	// Assign the name, just in case we haven't already assigned it
	n.Name = name

	return n, nil
}

// NodeName requests the current node name from the API, guaranteeing an accurate result
// This shouldn't be needed, Node.Name should be plenty for everyday use.
func (n *Node) NodeName() (string, error) {
	b, err := n.requestBytes(http.MethodGet, "api/node/name", nil)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Version requests the current server version
func (n *Node) Version() (string, error) {
	b, err := n.requestBytes(http.MethodGet, "api/version", nil)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Open will open a signalr websocket to the server for real-time events
func (n *Node) Open(ctx context.Context) error {
	if n.rtc != nil {
		return ErrAlreadyOpen
	}

	log.WithField("node", n.Name).Debug("Opening node connection")

	rtc, err := ConnectRTC(ctx, n.Name, n.baseAddress+"/hubs/core", n)

	if err != nil {
		return err
	}

	n.rtc = rtc

	if err = rtc.Authorize(n.token); err != nil {
		_ = rtc.Close()
		return err
	}

	if err = rtc.JoinUser(); err != nil {
		_ = rtc.Close()
		return err
	}

	go rtc.Start()

	return nil
}

// Connected checks whether the node is connected to SignalR
func (n *Node) Connected() bool {
	return n.rtc != nil
}

// JoinAllChannels will join all channels the account has access to.
// This will always be called on the primary node.
func (n *Node) JoinAllChannels(ctx context.Context) error {
	if !n.IsPrimary() {
		return n.Primary.JoinAllChannels(ctx)
	}

	planets, err := n.Planets()

	if err != nil {
		return err
	}

	wg := pool.New().
		WithErrors().
		WithContext(ctx).
		WithMaxGoroutines(4)

	for _, planet := range planets {
		// This may be slow, but we need to avoid a race condition with instances/Open
		// In the future, we could preload all node names ahead of time?
		node, err := n.NodeForPlanet(planet.ID)

		if err != nil {
			return err
		}

		// Make sure the node is connected
		if !node.Connected() {
			if err := node.Open(ctx); err != nil {
				return err
			}
		}

		wg.Go(func(ctx context.Context) error {
			log.WithField("planet", planet.Name).Debug("Getting nodes")

			channels, err := node.Channels(planet.ID)

			if err != nil {
				fmt.Println("Failed to get channels for planet", err)
				return err
			}

			var w errgroup.Group

			for _, channel := range channels {
				w.Go(func() error {
					log.WithFields(log.Fields{
						"planet":  planet.ID,
						"channel": channel.ID,
					}).Debug("Joining channel")

					return node.rtc.JoinChannel(channel.ID)
				})
			}

			return w.Wait()
		})
	}

	return wg.Wait()
}

// IsPrimary checks whether a node is the primary Valour node
func (n *Node) IsPrimary() bool {
	return n.Primary == nil
}

// Close will close any open rtc connections.
// If this is called on the primary node, all child nodes will also be closed.
func (n *Node) Close() error {
	var wg errgroup.Group

	if n.rtc != nil {
		wg.Go(n.rtc.Close)
	}

	if n.IsPrimary() {
		n.childNodes.IterCb(func(name string, v *Node) {
			if v.rtc == nil {
				return
			}

			wg.Go(v.Close)
		})
	}

	return wg.Wait()
}

func (n *Node) request(method, uri string, body any) (*http.Response, error) {
	var r io.Reader
	var bodyLen int

	if body != nil {
		b, err := json.Marshal(body)

		if err != nil {
			return nil, err
		}

		bodyLen = len(b)

		r = bytes.NewReader(b)
	}

	log.WithFields(log.Fields{
		"node":   n.Name,
		"method": method,
		"uri":    n.baseAddress + "/" + uri,
	}).Debug("Sending request to node")

	req, err := http.NewRequest(method, n.baseAddress+"/"+uri, r)

	if err != nil {
		return nil, err
	}

	if r != nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", strconv.Itoa(bodyLen))
	}

	res, err := n.client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Node) requestBytes(method, uri string, body any) ([]byte, error) {
	res, err := n.request(method, uri, body)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func (n *Node) requestJSON(method, uri string, body any, dest any) error {
	res, err := n.request(method, uri, body)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	contentType := res.Header.Get("Content-Type")

	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = contentType[:idx]
	}

	switch contentType {
	case "application/json":
		if err = json.NewDecoder(res.Body).Decode(dest); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%v: %s", ErrInvalidResponseType, contentType)
	}

	return nil
}
