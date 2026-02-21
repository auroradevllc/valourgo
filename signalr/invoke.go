package signalr

import (
	"encoding/json"
	"sync"
)

type InvocationManager struct {
	mu      sync.Mutex
	nextID  int
	waiters map[string]chan json.RawMessage
}

func NewInvocationManager() *InvocationManager {
	return &InvocationManager{
		waiters: make(map[string]chan json.RawMessage),
	}
}

func (m *InvocationManager) New() (string, chan json.RawMessage) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := string(rune(m.nextID))
	m.nextID++

	ch := make(chan json.RawMessage, 1)
	m.waiters[id] = ch
	return id, ch
}

func (m *InvocationManager) Resolve(id string, payload json.RawMessage) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ch, ok := m.waiters[id]; ok {
		ch <- payload
		close(ch)
		delete(m.waiters, id)
	}
}
