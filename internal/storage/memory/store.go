package memory

type Storage struct {
	Store map[string]string
}

func NewMemoryStorage() *Storage {
	return &Storage{Store: make(map[string]string)}
}

func (m *Storage) Put(key, value string) error {
	m.Store[key] = value
	return nil
}

func (m *Storage) Get(key string) (string, error) {
	value, ok := m.Store[key]
	if !ok {
		return "", ErrorNoSuchKey
	}
	return value, nil
}

func (m *Storage) Delete(key string) error {
	delete(m.Store, key)
	return nil
}
