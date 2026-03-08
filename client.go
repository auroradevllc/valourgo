package valour

import (
	"context"

	"github.com/auroradevllc/handler"
)

type Client interface {
	handler.HandlerInterface
	Planets
	Messages
	Channels
	Nodes
	Roles

	JoinAllChannels(ctx context.Context) error

	Me() (*User, error)
	MyMember(planetID PlanetID) (*Member, error)
	Member(id MemberID) (*Member, error)
	MemberByUser(planetID PlanetID, id UserID) (*Member, error)
}

type Nodes interface {
	NodeName() (string, error)
	Version() (string, error)
	Open(ctx context.Context) error
	Connected() bool
	Close() error
}
