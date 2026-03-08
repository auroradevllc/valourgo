package valour

import (
	"context"

	"github.com/auroradevllc/handler"
)

type Client interface {
	handler.HandlerInterface

	NodeName() (string, error)
	Version() (string, error)
	Open(ctx context.Context) error
	Connected() bool
	JoinAllChannels(ctx context.Context) error
	Close() error

	NodeForPlanet(planetID PlanetID) (*Node, error)
	GetNodeNameForPlanet(planetID PlanetID) (string, error)
	Planets() ([]Planet, error)
	Planet(id PlanetID) (*Planet, error)
	CreatePlanet(planet CreatePlanetData) (*Planet, error)
	UpdatePlanet(id PlanetID, data EditPlanetData) (*Planet, error)
	DeletePlanet(id PlanetID) error
	PlanetInitialData(id PlanetID) (*PlanetInitialData, error)
	JoinPlanet(planet PlanetID, inviteCode string) error

	Channel(planetID PlanetID, channelID ChannelID) (*Channel, error)
	Channels(id PlanetID) ([]Channel, error)

	Me() (*User, error)
	MyMember(planetID PlanetID) (*Member, error)
	Member(id MemberID) (*Member, error)
	MemberByUser(planetID PlanetID, id UserID) (*Member, error)
}
