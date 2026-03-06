package valour

import (
	"fmt"
	"net/http"
)

type Role struct {
	ID                  RoleID   `json:"id"`
	PlanetID            PlanetID `json:"planetId"`
	Name                string   `json:"name"`
	Position            int      `json:"position"`
	IsDefault           bool     `json:"isDefault"`
	Permissions         int64    `json:"permissions"`
	ChatPermissions     int      `json:"chatPermissions"`
	CategoryPermissions int      `json:"categoryPermissions"`
	VoicePermissions    int      `json:"voicePermissions"`
	Color               string   `json:"color"`
	Bold                bool     `json:"bold"`
	Italics             bool     `json:"italics"`
	FlagBitIndex        int      `json:"flagBitIndex"`
	AnyoneCanMention    bool     `json:"anyoneCanMention"`
	IsAdmin             bool     `json:"isAdmin"`
}

func (n *Node) Role(planetID PlanetID, roleID RoleID) (*Role, error) {
	var role Role

	if err := n.requestJSON(http.MethodGet, planetID.Route("roles", roleID.String()), nil, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

func (n *Node) UpdateRole(planetID PlanetID, role Role) (*Role, error) {
	var newRole Role

	if err := n.requestJSON(http.MethodPut, planetID.Route("roles", role.ID.String()), role, &newRole); err != nil {
		return nil, err
	}

	return &newRole, nil
}

func (n *Node) DeleteRole(planetID PlanetID, roleID RoleID) error {
	res, err := n.request(http.MethodDelete, planetID.Route("roles", roleID.String()), nil)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unknown status: %s", res.Status)
	}

	return nil
}
