package memory

var Store = make(map[string]string)

func Put(key string, value string) error {
	Store[key] = value
	return nil
}

func Get(key string) (string, error) {
	value, ok := Store[key]
	if !ok {
		return "", ErrorNoSuchKey
	}
	return value, nil
}

func Delete(key string) error {
	delete(Store, key)
	return nil
}
