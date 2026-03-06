package store

type Cabinet struct {
	MeStore
	ChannelStore
	PlanetStore
	MemberStore
	RoleStore
	EmojiStore
}

func (c *Cabinet) Reset() error {
	errs := []error{
		c.MeStore.Reset(),
		c.ChannelStore.Reset(),
		c.PlanetStore.Reset(),
		c.MemberStore.Reset(),
		c.RoleStore.Reset(),
		c.EmojiStore.Reset(),
	}

	return errs[0]
}
