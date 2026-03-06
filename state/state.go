package state

import (
	"github.com/auroradevllc/handler"
	valour "github.com/auroradevllc/valourgo"
	"github.com/auroradevllc/valourgo/state/store"
	"github.com/auroradevllc/valourgo/state/store/defaultstore"
)

type State struct {
	valour.Client
	*store.Cabinet
	*handler.Handler
}

var _ valour.Client = (*State)(nil)

func New(token string, opts ...valour.Option) (valour.Client, error) {
	c, err := valour.NewClient(token, opts...)

	if err != nil {
		return nil, err
	}

	return NewWithClient(c), nil
}

func NewWithClient(c valour.Client) *State {
	s := &State{
		Client:  c,
		Cabinet: defaultstore.New(),
		Handler: handler.New(),
	}

	s.hookEvents()

	return s
}

func (s *State) Me() (*valour.User, error) {
	me, err := s.Cabinet.Me()

	if err == nil {
		return me, nil
	}

	me, err = s.Client.Me()

	if err == nil {
		s.Cabinet.MyselfSet(*me, false)
	}

	return me, err
}

func (s *State) Planet(id valour.PlanetID) (*valour.Planet, error) {
	p, err := s.Cabinet.Planet(id)

	if err == nil {
		return p, nil
	}

	p, err = s.Client.Planet(id)

	if err == nil {
		s.Cabinet.PlanetSet(p, false)
	}

	return p, err
}

func (s *State) Planets() ([]valour.Planet, error) {
	planets, err := s.Cabinet.Planets()

	if err == nil {
		return planets, nil
	}

	planets, err = s.Client.Planets()

	if err == nil {
		for i := range planets {
			s.Cabinet.PlanetSet(&planets[i], false)
		}
	}

	return planets, err
}

func (s *State) MyMember(planetID valour.PlanetID) (*valour.Member, error) {
	me, err := s.Me()

	if err != nil {
		return nil, err
	}

	member, err := s.MemberByUser(planetID, me.ID)

	if err != nil {
		return nil, err
	}

	return member, nil
}

func (s *State) Member(id valour.MemberID) (*valour.Member, error) {
	member, err := s.Cabinet.Member(id)

	if err == nil {
		return member, nil
	}

	member, err = s.Client.Member(id)

	if err == nil {
		s.Cabinet.MemberSet(member, false)
	}

	return member, err
}

func (s *State) MemberByUser(planetID valour.PlanetID, id valour.UserID) (*valour.Member, error) {
	member, err := s.Cabinet.MemberByUser(planetID, id)

	if err == nil {
		return member, nil
	}

	member, err = s.Client.MemberByUser(planetID, id)

	if err == nil {
		s.Cabinet.MemberSet(member, false)
	}

	return member, err
}
