package defaultstore

import "github.com/auroradevllc/valourgo/state/store"

func New() *store.Cabinet {
	return &store.Cabinet{
		MeStore:      NewMe(),
		ChannelStore: NewChannel(),
		PlanetStore:  NewPlanet(),
		MemberStore:  NewMember(),
		RoleStore:    NewRole(),
		EmojiStore:   NewEmoji(),
	}
}
