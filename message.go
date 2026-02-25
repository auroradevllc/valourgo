package valourgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// LatestMessageIndex is long.MaxValue in C#, this will return our latest messages only
const (
	LatestMessageIndex = MessageID(9223372036854775807)
	maxMessageLimit    = 64
)

var ErrInvalidCount = errors.New("invalid message count")

type Message struct {
	ID              MessageID           `json:"id"`
	PlanetID        PlanetID            `json:"planetId"`
	ChannelID       ChannelID           `json:"channelId"`
	ReplyToID       *UserID             `json:"replyToId"`
	ReplyTo         *Message            `json:"replyTo"`
	AuthorID        UserID              `json:"authorUserId"`
	MemberID        MemberID            `json:"authorMemberId"`
	Content         string              `json:"content"`
	TimeSent        time.Time           `json:"timeSent"`
	EditedTime      *time.Time          `json:"editedTime"`
	Fingerprint     string              `json:"fingerprint"`
	Reactions       []Reaction          `json:"reactions"`
	Attachments     []MessageAttachment `json:"attachments"`
	AttachmentsData string              `json:"attachmentsData"`
}

func (m *Message) decodeAttachments() error {
	if m.AttachmentsData != "" {
		return json.Unmarshal([]byte(m.AttachmentsData), &m.Attachments)
	}

	return nil
}

type SendMessageData struct {
	AuthorMemberID  MemberID             `json:"authorMemberId"`
	PlanetID        PlanetID             `json:"planetId"`
	ChannelID       ChannelID            `json:"channelId"`
	ReplyToID       *MessageID           `json:"replyToId"`
	Content         string               `json:"content"`
	Attachments     []*MessageAttachment `json:"attachments,omitempty"`
	AttachmentsData string               `json:"attachmentsData,omitempty"`
	EmbedData       string               `json:"embedData,omitempty"`
	Fingerprint     string               `json:"fingerprint"`
}

type EditMessageData struct {
	ID       MessageID `json:"id"`
	PlanetID PlanetID  `json:"planetId"`
	Content  *string   `json:"content"`
}

// Messages retrieves the latest x messages
func (n *Node) Messages(planetID PlanetID, channelID ChannelID, limit uint) ([]Message, error) {
	return n.MessagesBefore(planetID, channelID, LatestMessageIndex, limit)
}

// MessagesBefore retrieves messages before a specific message
func (n *Node) MessagesBefore(planetID PlanetID, channelID ChannelID, index MessageID, limit uint) ([]Message, error) {
	msgs := make([]Message, 0, limit)

	fetch := uint(maxMessageLimit)

	unlimited := limit == 0

	for limit > 0 || unlimited {
		if limit > 0 {
			fetch = uint(intMin(maxMessageLimit, int(limit)))
			limit -= maxMessageLimit
		}

		m, err := n.messagesBefore(planetID, channelID, index, fetch)

		if err != nil {
			return msgs, err
		}

		msgs = append(msgs, m...)

		if len(m) < maxMessageLimit {
			break
		}

		index = m[len(m)-1].ID
	}

	if len(msgs) == 0 {
		return nil, nil
	}

	return msgs, nil
}

// messagesBefore is called to retrieve messages, used with MessagesBefore to append to a slice
func (n *Node) messagesBefore(planetID PlanetID, channelID ChannelID, index MessageID, limit uint) ([]Message, error) {
	switch {
	case limit == 0:
		limit = 50
	case limit > 64:
		limit = 64
	}

	var messages []Message

	v := make(url.Values)
	v.Set("index", index.String())
	v.Set("count", strconv.FormatUint(uint64(limit), 10))

	if err := n.requestJSON(http.MethodGet, planetID.Route("channels", channelID.String(), "messages"), nil, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// Message retrieves a single message
func (n *Node) Message(id MessageID) (*Message, error) {
	var message Message

	if err := n.requestJSON(http.MethodGet, id.Route(), nil, &message); err != nil {
		return nil, err
	}

	if err := message.decodeAttachments(); err != nil {
		return nil, err
	}

	return &message, nil
}

// EditMessage updates a message
func (n *Node) EditMessage(id MessageID, m EditMessageData) (*Message, error) {
	m.ID = id

	var updatedMessage Message

	if err := n.requestJSON(http.MethodPut, id.Route(), m, &updatedMessage); err != nil {
		return nil, err
	}

	if err := updatedMessage.decodeAttachments(); err != nil {
		return nil, err
	}

	return &updatedMessage, nil
}

// DeleteMessage deletes a message
func (n *Node) DeleteMessage(id MessageID) error {
	res, err := n.request(http.MethodDelete, id.Route(), nil)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("unknown status %d", res.StatusCode)
}

// SendMessage sends a simple text message
func (n *Node) SendMessage(planetID PlanetID, channelID ChannelID, content string) (*Message, error) {
	return n.SendMessageComplex(planetID, channelID, SendMessageData{
		Content: content,
	})
}

// SendMessageComplex sends a message with optional text, attachments, and embeds
func (n *Node) SendMessageComplex(planetID PlanetID, channelID ChannelID, send SendMessageData) (*Message, error) {
	send.PlanetID = planetID
	send.ChannelID = channelID

	// Fingerprints are required
	u, err := uuid.NewV7()

	if err != nil {
		return nil, err
	}

	if len(send.Attachments) > 0 {
		b, err := json.Marshal(send.Attachments)

		if err != nil {
			return nil, err
		}

		send.AttachmentsData = string(b)
	}

	send.Fingerprint = u.String()

	// Until this PR is merged, this is required: https://github.com/Valour-Software/Valour/pull/1426
	if send.AuthorMemberID == 0 {
		myMember, err := n.MyMember(planetID)

		if err != nil {
			return nil, err
		}

		send.AuthorMemberID = myMember.ID
	}

	res, err := n.request(http.MethodPost, apiMessageBase, send)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("unknown status %d %s", res.StatusCode, string(b))
	}

	var m Message

	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, err
	}

	if err := m.decodeAttachments(); err != nil {
		return nil, err
	}

	return &m, nil
}

// intMin is a simple math.Min for integers
func intMin(a, b int) int {
	if a < b {
		return a
	}

	return b
}
