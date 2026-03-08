package valour

import (
	"net/http"
	"time"
)

type ChannelType int

const (
	PlanetChat ChannelType = iota
	PlanetCategory
	PlanetVoice

	DirectChat
	DirectVoice

	GroupChat
	GroupVoice
)

type Channel struct {
	ID             ChannelID       `json:"id"`
	PlanetID       PlanetID        `json:"planetId"`
	ParentID       ChannelID       `json:"parentId"`
	ChannelType    ChannelType     `json:"channelType"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	InheritsPerms  bool            `json:"inheritsPerms"`
	IsDefault      bool            `json:"isDefault"`
	NSFW           bool            `json:"nsfw"`
	RawPosition    int64           `json:"rawPosition"`
	Position       ChannelPosition `json:"position"`
	LastUpdateTime time.Time       `json:"lastUpdateTime"`
}

type ChannelPosition struct {
	RawPosition   int64 `json:"rawPosition"`
	Depth         int   `json:"depth"`
	LocalPosition int   `json:"localPosition"`
}

type Channels interface {
	Channel(planetID PlanetID, channelID ChannelID) (*Channel, error)
	Channels(id PlanetID) ([]Channel, error)
}

func (n *Node) Channel(planetID PlanetID, channelID ChannelID) (*Channel, error) {
	var channel Channel

	if err := n.requestJSON(http.MethodGet, planetID.Route("channels", channelID.String()), nil, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

// Channels gets a planet's channels
func (n *Node) Channels(id PlanetID) ([]Channel, error) {
	var channels []Channel

	node, err := n.NodeForPlanet(id)

	if err != nil {
		return nil, err
	}

	if err := node.requestJSON(http.MethodGet, id.Route("channels"), nil, &channels); err != nil {
		return nil, err
	}

	return channels, nil
}
