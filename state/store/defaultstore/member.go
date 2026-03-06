package defaultstore

import (
	valour "github.com/auroradevllc/valourgo"
	"github.com/auroradevllc/valourgo/state/store"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type planetMembers struct {
	memberIDs cmap.ConcurrentMap[valour.UserID, valour.MemberID]
}

var _ store.MemberStore = (*Member)(nil)

func NewMember() *Member {
	return &Member{
		members: cmap.NewStringer[valour.MemberID, valour.Member](),
		planets: cmap.NewStringer[valour.PlanetID, *planetMembers](),
	}
}

type Member struct {
	members cmap.ConcurrentMap[valour.MemberID, valour.Member]
	planets cmap.ConcurrentMap[valour.PlanetID, *planetMembers]
}

func (s *Member) Reset() error {
	s.members.Clear()
	s.planets.Clear()
	return nil
}

func (s *Member) Member(id valour.MemberID) (*valour.Member, error) {
	m, ok := s.members.Get(id)

	if !ok {
		return nil, store.ErrNotFound
	}

	return &m, nil
}

func (s *Member) MemberByUser(id valour.PlanetID, userID valour.UserID) (*valour.Member, error) {
	planet, ok := s.planets.Get(id)

	if !ok {
		return nil, store.ErrNotFound
	}

	memberID, ok := planet.memberIDs.Get(userID)

	if !ok {
		return nil, store.ErrNotFound
	}

	return s.Member(memberID)
}

func (s *Member) Members(id valour.PlanetID) ([]valour.Member, error) {
	planet, ok := s.planets.Get(id)

	if !ok {
		return nil, store.ErrNotFound
	}

	var members = make([]valour.Member, 0, planet.memberIDs.Count())

	for t := range planet.memberIDs.IterBuffered() {
		member, err := s.Member(t.Val)

		if err != nil {
			return nil, err
		}

		members = append(members, *member)
	}

	return members, nil
}

func (s *Member) MemberSet(m *valour.Member, update bool) error {
	if !s.members.Has(m.ID) || update {
		s.members.Set(m.ID, *m)

		planet, ok := s.planets.Get(m.PlanetID)

		if !ok {
			s.planets.Set(m.PlanetID, &planetMembers{
				memberIDs: cmap.NewStringer[valour.UserID, valour.MemberID](),
			})
		} else {
			planet.memberIDs.Set(m.UserID, m.ID)
		}
	}

	return nil
}

func (s *Member) MemberRemove(memberID valour.MemberID) error {
	// Utilize members cache to remove from planet
	s.members.RemoveCb(memberID, func(_ valour.MemberID, member valour.Member, exists bool) bool {
		planets, ok := s.planets.Get(member.PlanetID)

		if ok {
			planets.memberIDs.Remove(member.UserID)
		}

		return true
	})

	return nil
}
