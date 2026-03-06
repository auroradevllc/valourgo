package state

import valour "github.com/auroradevllc/valourgo"

func (s *State) hookEvents() {
	s.Client.AddSyncHandler(func(event interface{}) {
		// Handle events to populate the store before calling the other handler
		s.onEvent(event)

		s.Handler.Call(event)
	})
}

func (s *State) onEvent(e interface{}) {
	switch ev := e.(type) {
	case *valour.PlanetJoinEvent:
		s.retrieveInitialPlanet(ev.PlanetID)
	case *valour.PlanetUpdateEvent:
		if err := s.Cabinet.PlanetSet(&ev.Planet, true); err != nil {
			s.logError(err)
		}
	case *valour.PlanetDeleteEvent:
		if err := s.Cabinet.PlanetRemove(ev.PlanetID); err != nil {
			s.logError(err)
		}
	}
}

// retrieveInitialPlanet stores a planet we retrieved on RTC join
func (s *State) retrieveInitialPlanet(id valour.PlanetID) error {
	// Call Planet to ensure the initial planet exists
	_, err := s.Planet(id)

	// Retrieve initial data from the API (channels, roles, emojis, voice channels)
	data, err := s.Client.PlanetInitialData(id)

	if err != nil {
		return err
	}

	for _, ch := range data.Channels {
		if err := s.Cabinet.ChannelSet(&ch, false); err != nil {
			// Handle error
		}
	}

	// handle data.Roles

	// handle data.Emojis
	return nil
}

func (s *State) logError(err error) {
	// TODO: log error
}
