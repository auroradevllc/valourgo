package valourgo

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Planet is Valour's representation of a server/group
type Planet struct {
	ID                  PlanetID `json:"id"`
	OwnerID             UserID   `json:"ownerId"`
	Name                string   `json:"name"`
	NodeName            string   `json:"nodeName"`
	HasCustomIcon       bool     `json:"hasCustomIcon"`
	HasAnimatedIcon     bool     `json:"hasAnimatedIcon"`
	Description         string   `json:"description"`
	Public              bool     `json:"public"`
	Discoverable        bool     `json:"discoverable"`
	NSFW                bool     `json:"nsfw"`
	Version             int      `json:"version"`
	HasCustomBackground bool     `json:"hasCustomBackground"`
	Tags                []Tag    `json:"tags"`
}

// PlanetInitialData contains planet data
type PlanetInitialData struct {
	Channels []Channel `json:"channels"`
	Roles    []Role    `json:"roles"`
	Emojis   []Emoji   `json:"emojis"`
}

func (n *Node) NodeForPlanet(planetID PlanetID) (*Node, error) {
	// Always pass this call up to the primary node
	if !n.IsPrimary() {
		return n.Primary.NodeForPlanet(planetID)
	}

	name, exists := n.planetNodeList.Get(planetID)

	if !exists {
		primary := n

		if !n.IsPrimary() {
			primary = n.Primary
		}

		nodeName, err := primary.GetNodeNameForPlanet(planetID)

		if err != nil {
			return nil, err
		}

		name = nodeName

		n.planetNodeList.Set(planetID, nodeName)
	}

	// If the name matches the current node (primary) return the node
	if name == n.Name {
		return n, nil
	}

	// Otherwise, try to get the node from our list of already initialized child nodes
	if node, exists := n.childNodes.Get(name); exists {
		return node, nil
	}

	node, err := NewNode(n.baseAddress, name, n.token, WithNodeHandler(n.Handler))

	if err != nil {
		return nil, err
	}

	node.Primary = n

	n.childNodes.Set(name, node)

	return node, nil
}

// GetNodeNameForPlanet retrieves the node name for the specified planet
func (n *Node) GetNodeNameForPlanet(planetID PlanetID) (string, error) {
	res, err := n.requestBytes(http.MethodGet, "api/node/planet/"+planetID.String(), nil)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(res)), nil
}

// Planets returns the user's planets
// This request always goes to the primary node
func (n *Node) Planets() ([]Planet, error) {
	if n.Primary != nil {
		return n.Primary.Planets()
	}

	var planets []Planet

	if err := n.requestJSON(http.MethodGet, "api/users/me/planets", nil, &planets); err != nil {
		return nil, err
	}

	// Pre-populate node list for planets
	for _, planet := range planets {
		n.planetNodeList.Set(planet.ID, planet.NodeName)
	}

	return planets, nil
}

// Planet retrieves a planet by a specified ID
// This request always goes to the primary node
func (n *Node) Planet(id PlanetID) (*Planet, error) {
	if !n.IsPrimary() {
		return n.Primary.Planet(id)
	}

	var planet Planet

	if err := n.requestJSON(http.MethodGet, fmt.Sprintf("api/planets/%d", id), nil, &planet); err != nil {
		return nil, err
	}

	n.planetNodeList.Set(planet.ID, planet.NodeName)

	return &planet, nil
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

type CreatePlanetData struct {
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	Public       bool    `json:"public"`
	Discoverable bool    `json:"discoverable"`
	NSFW         bool    `json:"nsfw"`
}

// CreatePlanet will create a new planet
// This request always goes to the primary node
func (n *Node) CreatePlanet(planet CreatePlanetData) (*Planet, error) {
	if !n.IsPrimary() {
		return n.Primary.CreatePlanet(planet)
	}

	var newPlanet Planet

	if err := n.requestJSON(http.MethodPost, "api/planets", planet, &newPlanet); err != nil {
		return nil, err
	}

	return &newPlanet, nil
}

type EditPlanetData struct {
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	Public       bool    `json:"public"`
	Discoverable bool    `json:"discoverable"`
	NSFW         bool    `json:"nsfw"`
	Tags         []Tag   `json:"tags"`
}

// UpdatePlanet will update an existing planet with the specified data
// This request always goes to the primary node
func (n *Node) UpdatePlanet(id PlanetID, data EditPlanetData) (*Planet, error) {
	if !n.IsPrimary() {
		return n.Primary.UpdatePlanet(id, data)
	}

	planet, err := n.Planet(id)

	if err != nil {
		return nil, err
	}

	// API requires OwnerID to exist, but the planet owner cannot be changed. Don't give this option to users.
	fields := struct {
		EditPlanetData
		OwnerID UserID `json:"ownerId"`
	}{
		EditPlanetData: data,
		OwnerID:        planet.OwnerID,
	}

	var newPlanet Planet

	if err := n.requestJSON(http.MethodPut, id.Route(), fields, &newPlanet); err != nil {
		return nil, err
	}

	return &newPlanet, nil
}

// DeletePlanet will delete a planet
func (n *Node) DeletePlanet(id PlanetID) error {
	res, err := n.request(http.MethodDelete, id.Route(), nil)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

// PlanetInitialData retrieves initial data for a planet, such as channels, roles, and emojis
func (n *Node) PlanetInitialData(id PlanetID) (*PlanetInitialData, error) {
	node, err := n.NodeForPlanet(id)

	if err != nil {
		return nil, err
	}

	var data PlanetInitialData

	if err := node.requestJSON(http.MethodGet, id.Route(apiPlanetInitialData), nil, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// JoinPlanet allows you to join a planet, with optional invite code
func (n *Node) JoinPlanet(planet PlanetID, inviteCode string) error {
	if !n.IsPrimary() {
		return n.Primary.JoinPlanet(planet, inviteCode)
	}

	uri := planet.Route("join")

	if inviteCode != "" {
		q := make(url.Values)
		q.Set("inviteCode", inviteCode)
		uri += "?" + q.Encode()
	}

	res, err := n.request(http.MethodPost, uri, nil)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("unknown status %d", res.StatusCode)
}
