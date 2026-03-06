package defaultstore

import (
	valour "github.com/auroradevllc/valourgo"
	"github.com/auroradevllc/valourgo/state/store"
	cmap "github.com/orcaman/concurrent-map/v2"
)

func NewRole() *Role {
	return &Role{
		planets: cmap.NewStringer[valour.PlanetID, roleMap](),
	}
}

var _ store.RoleStore = (*Role)(nil)

type roleMap = cmap.ConcurrentMap[valour.RoleID, valour.Role]

type Role struct {
	planets cmap.ConcurrentMap[valour.PlanetID, roleMap]
}

func (s *Role) Reset() error {
	s.planets.Clear()
	return nil
}

func (s *Role) Role(planetID valour.PlanetID, roleID valour.RoleID) (*valour.Role, error) {
	planet, ok := s.planets.Get(planetID)

	if !ok {
		return nil, store.ErrNotFound
	}

	role, ok := planet.Get(roleID)

	if !ok {
		return nil, store.ErrNotFound
	}

	return &role, nil
}

func (s *Role) Roles(planetID valour.PlanetID) ([]valour.Role, error) {
	planet, ok := s.planets.Get(planetID)

	if !ok {
		return nil, store.ErrNotFound
	}

	roles := make([]valour.Role, 0, planet.Count())

	for v := range planet.IterBuffered() {
		roles = append(roles, v.Val)
	}

	return roles, nil
}

func (s *Role) RoleSet(c *valour.Role, update bool) error {
	planet, ok := s.planets.Get(c.PlanetID)

	if !ok {
		planet = cmap.NewStringer[valour.RoleID, valour.Role]()

		s.planets.Set(c.PlanetID, planet)
	}

	if !planet.Has(c.ID) || update {
		planet.Set(c.ID, *c)
	}

	return nil
}

func (s *Role) RoleRemove(planetID valour.PlanetID, roleID valour.RoleID) error {
	planet, ok := s.planets.Get(planetID)

	if !ok {
		return nil
	}

	planet.Remove(roleID)
	return nil
}
