package defaultstore

import (
	valour "github.com/auroradevllc/valourgo"
	"github.com/auroradevllc/valourgo/state/store"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type Planet struct {
	planets cmap.ConcurrentMap[valour.PlanetID, valour.Planet]
}

func NewPlanet() *Planet {
	return &Planet{
		planets: cmap.NewStringer[valour.PlanetID, valour.Planet](),
	}
}

var _ store.PlanetStore = (*Planet)(nil)

func (s *Planet) Reset() error {
	s.planets.Clear()
	return nil
}

func (s *Planet) Planet(id valour.PlanetID) (*valour.Planet, error) {
	p, ok := s.planets.Get(id)

	if !ok {
		return nil, store.ErrNotFound
	}

	return &p, nil
}

func (s *Planet) Planets() ([]valour.Planet, error) {
	var planets = make([]valour.Planet, 0, len(s.planets.Keys()))

	s.planets.IterCb(func(_ valour.PlanetID, p valour.Planet) {
		planets = append(planets, p)
	})

	return planets, nil
}

func (s *Planet) PlanetSet(c *valour.Planet, update bool) error {
	return nil
}

func (s *Planet) PlanetRemove(id valour.PlanetID) error {
	s.planets.Remove(id)
	return nil
}
