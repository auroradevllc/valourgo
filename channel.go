package valour

import "net/http"

func (n *Node) Channel(planetID PlanetID, channelID ChannelID) (*Channel, error) {
	var channel Channel

	if err := n.requestJSON(http.MethodGet, planetID.Route("channels", channelID.String()), nil, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}
