package valourgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/auroradevllc/handler"
	"github.com/auroradevllc/valourgo/signalr"
	log "github.com/sirupsen/logrus"
)

type RTCState int

const (
	RTCStateUnknown RTCState = iota
	RTCStateConnected
	RTCStateUnauthorized
	RTCStateAuthorized
)

var (
	ErrInvalidPing = errors.New("invalid ping response")
)

type RTC struct {
	client  *signalr.Client
	handler handler.HandlerInterface
	state   RTCState
}

type BaseRTCResponse struct {
	Success   bool    `json:"Success"`
	Message   *string `json:"Message"`
	ErrorCode *int    `json:"ErrorCode"`
}

type RTCAuthResponse struct {
	BaseRTCResponse
}

func ConnectRTC(ctx context.Context, name, address string, handler handler.HandlerInterface) (*RTC, error) {
	h := make(http.Header)
	h.Set("X-Server-Select", name)

	r := &RTC{
		handler: handler,
	}

	r.client = signalr.NewClient(address,
		signalr.WithHTTPHeaders(h),
		signalr.WithDefaultHandler(r.defaultHandler))

	if err := r.client.Connect(ctx); err != nil {
		_ = r.client.Close()
		return nil, err
	}

	r.state = RTCStateConnected

	return r, nil
}

func (r *RTC) defaultHandler(target string, args []json.RawMessage) {
	switch target {
	case "Channel-State":
		decodeAndCall[ChannelStateEvent](args[0], r.handler)
	case "User-Update":
		decodeAndCall[UserUpdateEvent](args[0], r.handler)
	case "Relay":
		decodeAndCall[MessageCreateEvent](args[0], r.handler)
	case "RelayEdit":
		decodeAndCall[MessageEditEvent](args[0], r.handler)
	case "DeleteMessage":
		decodeAndCall[MessageDeleteEvent](args[0], r.handler)
	case "Channel-Watching-Update":
		decodeAndCall[ChannelWatchingUpdate](args[0], r.handler)
	case "Channel-CurrentlyTyping-Update":
		decodeAndCall[ChannelCurrentlyTypingUpdate](args[0], r.handler)
	case "PlanetMember-Update":
		decodeAndCall[PlanetMemberUpdate](args[0], r.handler)
	case "MessageReactionAdd":
		decodeAndCall[MessageReactionAddedEvent](args[0], r.handler)
	case "MessageReactionRemove":
		decodeAndCall[MessageReactionRemovedEvent](args[0], r.handler)
	default:
		log.WithField("target", target).Debug("No handler registered for target")
	}
}

func decodeAndCall[V any](b json.RawMessage, h handler.HandlerInterface) {
	var e V

	if err := json.Unmarshal(b, &e); err != nil {
		fmt.Println("Failed to decode type " + reflect.TypeOf(e).Elem().Name() + ": " + err.Error())
		return
	}

	go h.Call(&e)
}

// Ping sends a ping request to the SignalR hub
func (r *RTC) Ping() error {
	res, err := r.client.Invoke("ping", true)

	if err != nil {
		return err
	}

	v := <-res

	var response string

	if err := json.Unmarshal(v, &response); err != nil {
		return err
	}

	if response != "pong" {
		return ErrInvalidPing
	}

	return nil
}

// Start simply starts the ping ticker, sent every 60 seconds
func (r *RTC) Start() {
	t := time.NewTicker(60 * time.Second)

	for {
		r.Ping()
		<-t.C
	}
}

// Authorize sends our token to the SignalR hub, used as authentication
func (r *RTC) Authorize(token string) error {
	res, err := r.client.Invoke("Authorize", token)

	if err != nil {
		return err
	}

	if err := r.checkInvokeError(res); err != nil {
		r.state = RTCStateUnauthorized
		return err
	}

	r.state = RTCStateAuthorized

	return nil
}

// State returns our current state
func (r *RTC) State() RTCState {
	return r.state
}

// JoinUser will join the user update channel
func (r *RTC) JoinUser() error {
	res, err := r.client.Invoke("JoinUser", true)

	if err != nil {
		return err
	}

	return r.checkInvokeError(res)
}

// JoinPlanet will subscribe to the planet channel to receive updates for a planet
func (r *RTC) JoinPlanet(planet PlanetID) error {
	res, err := r.client.Invoke("JoinPlanet", planet)

	if err != nil {
		return err
	}

	return r.checkInvokeError(res)
}

// LeavePlanet removes our subscription to the planet channel
func (r *RTC) LeavePlanet(planet PlanetID) error {
	res, err := r.client.Invoke("JoinPlanet", planet)

	if err != nil {
		return err
	}

	return r.checkInvokeError(res)
}

// JoinChannel subscribes to the channel's updates/messages
func (r *RTC) JoinChannel(channel ChannelID) error {
	res, err := r.client.Invoke("JoinChannel", channel)

	if err != nil {
		return err
	}

	return r.checkInvokeError(res)
}

// LeaveChannel unsubscribes from channel updates/messages
func (r *RTC) LeaveChannel(channel ChannelID) error {
	res, err := r.client.Invoke("LeaveChannel", channel)

	if err != nil {
		return err
	}

	return r.checkInvokeError(res)
}

// checkInvokeError will validate a message, decoding it as an RTC Response, and returning an error if one occurred
func (r *RTC) checkInvokeError(c <-chan json.RawMessage) error {
	b := <-c

	var result BaseRTCResponse

	if err := json.Unmarshal(b, &result); err != nil {
		return err
	}

	if result.Success {
		return nil
	}

	return RTCError{
		ErrorCode: result.ErrorCode,
		Message:   result.Message,
	}
}

// Close will close the signalr client
func (r *RTC) Close() error {
	return r.client.Close()
}

type RTCError struct {
	ErrorCode *int
	Message   *string
}

func (e RTCError) Error() string {
	if e.Message != nil {
		return *e.Message
	}

	return fmt.Sprintf("unknown error %d", *e.ErrorCode)
}
