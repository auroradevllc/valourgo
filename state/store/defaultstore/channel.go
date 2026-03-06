package defaultstore

import (
	"slices"

	valour "github.com/auroradevllc/valourgo"
	"github.com/auroradevllc/valourgo/state/store"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type Channel struct {
	channels       cmap.ConcurrentMap[valour.ChannelID, valour.Channel]
	planetChannels cmap.ConcurrentMap[valour.PlanetID, []valour.ChannelID]
}

var _ store.ChannelStore = (*Channel)(nil)

func NewChannel() *Channel {
	return &Channel{
		channels:       cmap.NewStringer[valour.ChannelID, valour.Channel](),
		planetChannels: cmap.NewStringer[valour.PlanetID, []valour.ChannelID](),
	}
}

func (s *Channel) Reset() error {
	s.channels.Clear()
	s.planetChannels.Clear()
	return nil
}

func (s *Channel) Channel(id valour.ChannelID) (*valour.Channel, error) {
	item, ok := s.channels.Get(id)

	if !ok {
		return nil, store.ErrNotFound
	}

	return &item, nil
}

func (s *Channel) Channels(planet valour.PlanetID) ([]valour.Channel, error) {
	ids, ok := s.planetChannels.Get(planet)

	if !ok {
		return nil, store.ErrNotFound
	}

	var channels = make([]valour.Channel, 0, len(ids))

	for _, id := range ids {
		ch, ok := s.channels.Get(id)

		if !ok {
			return nil, store.ErrNotFound
		}

		channels = append(channels, ch)
	}

	return channels, nil
}

func (s *Channel) ChannelSet(c *valour.Channel, update bool) error {
	cpy := *c

	s.channels.Set(c.ID, cpy)

	list, _ := s.planetChannels.Get(c.PlanetID)

	if !slices.Contains(list, c.ID) {
		list = append(list, c.ID)

		s.planetChannels.Set(c.PlanetID, list)
	}

	return nil
}

func (s *Channel) ChannelRemove(c *valour.Channel) error {
	s.channels.Remove(c.ID)

	list, _ := s.planetChannels.Get(c.PlanetID)

	if !slices.Contains(list, c.ID) {
		list = append(list, c.ID)

		idx := slices.Index(list, c.ID)

		if idx != -1 {
			s.planetChannels.Set(c.PlanetID, slices.Delete(list, idx, idx))
		}
	}

	return nil
}
