package defaultstore

import (
	"sync"

	valour "github.com/auroradevllc/valourgo"
	"github.com/auroradevllc/valourgo/state/store"
)

func NewMe() *Me {
	return &Me{}
}

var _ store.MeStore = (*Me)(nil)

type Me struct {
	mut sync.RWMutex
	me  valour.User
}

func (s *Me) Reset() error {
	s.mut.Lock()
	s.me = valour.User{}
	s.mut.Unlock()
	return nil
}

func (s *Me) Me() (*valour.User, error) {
	s.mut.RLock()
	me := s.me
	s.mut.RUnlock()

	if !me.ID.IsValid() {
		return nil, store.ErrNotFound
	}

	return &s.me, nil
}

func (s *Me) MyselfSet(u valour.User, update bool) error {
	s.me = u
	return nil
}
