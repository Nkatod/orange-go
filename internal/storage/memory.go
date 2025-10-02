package storage

type IMemoryStorage interface {
	Put(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}
