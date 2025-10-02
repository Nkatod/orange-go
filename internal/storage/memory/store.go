package memory

import (
	"sync"

	"orange-go/internal/storage"
)

type Storage struct {
	Store map[string]string
	mu    sync.RWMutex
}

func NewMemoryStorage() storage.IMemoryStorage {
	return &Storage{Store: make(map[string]string)}
}

func (m *Storage) Put(key, value string) error {
	m.mu.Lock()
	m.Store[key] = value
	m.mu.Unlock()
	return nil
}

func (m *Storage) Get(key string) (string, error) {
	m.mu.RLock()
	value, ok := m.Store[key]
	m.mu.RUnlock()
	if !ok {
		return "", ErrorNoSuchKey
	}
	return value, nil
}

func (m *Storage) Delete(key string) error {
	m.mu.Lock()
	delete(m.Store, key)
	m.mu.Unlock()
	return nil
}
