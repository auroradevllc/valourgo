package store

import (
	"errors"

	valour "github.com/auroradevllc/valourgo"
)

var ErrNotFound = errors.New("item not found in store")

type Resettable interface {
	Reset() error
}

type MeStore interface {
	Resettable

	Me() (*valour.User, error)
	MyselfSet(u valour.User, update bool) error
}

type PlanetStore interface {
	Resettable

	Planet(id valour.PlanetID) (*valour.Planet, error)

	Planets() ([]valour.Planet, error)

	PlanetSet(c *valour.Planet, update bool) error
	PlanetRemove(id valour.PlanetID) error
}

type MemberStore interface {
	Resettable

	Member(valour.MemberID) (*valour.Member, error)
	MemberByUser(valour.PlanetID, valour.UserID) (*valour.Member, error)
	Members(valour.PlanetID) ([]valour.Member, error)

	MemberSet(m *valour.Member, update bool) error
	MemberRemove(valour.MemberID) error
}

type ChannelStore interface {
	Resettable

	Channel(id valour.ChannelID) (*valour.Channel, error)

	Channels(id valour.PlanetID) ([]valour.Channel, error)

	ChannelSet(c *valour.Channel, update bool) error
	ChannelRemove(c *valour.Channel) error
}

type RoleStore interface {
	Resettable

	Role(planetID valour.PlanetID, roleID valour.RoleID) (*valour.Role, error)
	Roles(planetID valour.PlanetID) ([]valour.Role, error)

	RoleSet(c *valour.Role, update bool) error
	RoleRemove(planetID valour.PlanetID, roleID valour.RoleID) error
}

type EmojiStore interface {
	Resettable

	Emoji(planetID valour.PlanetID, emojiID valour.EmojiID) (*valour.Emoji, error)
	Emojis(planetID valour.PlanetID) ([]valour.Emoji, error)

	EmojiSet(planetID valour.PlanetID, emojis []valour.Emoji, update bool) error
}
