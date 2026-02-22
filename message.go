package valourgo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID          MessageID  `json:"id"`
	PlanetID    PlanetID   `json:"planetId"`
	ChannelID   ChannelID  `json:"channelId"`
	ReplyToID   *UserID    `json:"replyToId"`
	ReplyTo     *Message   `json:"replyTo"`
	AuthorID    UserID     `json:"authorUserId"`
	MemberID    MemberID   `json:"authorMemberId"`
	Content     string     `json:"content"`
	TimeSent    time.Time  `json:"timeSent"`
	EditedTime  *time.Time `json:"editedTime"`
	Fingerprint string     `json:"fingerprint"`
	Reactions   []Reaction `json:"reactions"`
}

type SendMessageData struct {
	AuthorMemberID  MemberID   `json:"authorMemberId"`
	PlanetID        PlanetID   `json:"planetId"`
	ChannelID       ChannelID  `json:"channelId"`
	ReplyToID       *MessageID `json:"replyToId"`
	Content         string     `json:"content"`
	Attachments     any        `json:"attachments"`
	AttachmentsData string     `json:"attachmentsData"`
	Fingerprint     string     `json:"fingerprint"`
}

type EditMessageData struct {
	ID       MessageID `json:"id"`
	PlanetID PlanetID  `json:"planetId"`
	Content  *string   `json:"content"`
}

func (n *Node) Message(id MessageID) (*Message, error) {
	var message Message

	if err := n.requestJSON(http.MethodGet, id.Route(), nil, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (n *Node) EditMessage(id MessageID, m EditMessageData) (*Message, error) {
	m.ID = id

	var updatedMessage Message

	if err := n.requestJSON(http.MethodPut, id.Route(), m, &updatedMessage); err != nil {
		return nil, err
	}

	return &updatedMessage, nil
}

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

func (n *Node) SendMessage(planetID PlanetID, channelID ChannelID, content string) (*Message, error) {
	return n.SendMessageComplex(planetID, channelID, SendMessageData{
		Content: content,
	})
}

func (n *Node) SendMessageComplex(planetID PlanetID, channelID ChannelID, send SendMessageData) (*Message, error) {
	send.PlanetID = planetID
	send.ChannelID = channelID

	u, err := uuid.NewV7()

	if err != nil {
		return nil, err
	}

	send.Fingerprint = u.String()

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

	return &m, nil
}
