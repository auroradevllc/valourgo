package defaultstore

import (
	valour "github.com/auroradevllc/valourgo"
	"github.com/auroradevllc/valourgo/state/store"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type emojis = cmap.ConcurrentMap[valour.EmojiID, valour.Emoji]

func NewEmoji() *Emoji {
	return &Emoji{
		planets: cmap.NewStringer[valour.PlanetID, emojis](),
	}
}

var _ store.EmojiStore = (*Emoji)(nil)

type Emoji struct {
	planets cmap.ConcurrentMap[valour.PlanetID, emojis]
}

func (s *Emoji) Reset() error {
	s.planets.Clear()
	return nil
}

func (s *Emoji) Emoji(planetID valour.PlanetID, emojiID valour.EmojiID) (*valour.Emoji, error) {
	planet, ok := s.planets.Get(planetID)

	if !ok {
		return nil, store.ErrNotFound
	}

	emoji, ok := planet.Get(emojiID)

	if !ok {
		return nil, store.ErrNotFound
	}

	return &emoji, nil
}

func (s *Emoji) Emojis(planetID valour.PlanetID) ([]valour.Emoji, error) {
	planet, ok := s.planets.Get(planetID)

	if !ok {
		return nil, store.ErrNotFound
	}

	emojis := make([]valour.Emoji, 0, planet.Count())

	for i := range planet.IterBuffered() {
		emojis = append(emojis, i.Val)
	}

	return emojis, nil
}

func (s *Emoji) EmojiSet(planetID valour.PlanetID, emojis []valour.Emoji, update bool) error {
	planet, ok := s.planets.Get(planetID)

	if ok && !update {
		return nil
	}

	planet.Clear()

	for _, emoji := range emojis {
		planet.Set(emoji.ID, emoji)
	}

	return nil
}
