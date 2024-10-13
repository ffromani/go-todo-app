package store

type Redis struct{}

func NewRedis() (*Redis, error) {
	return nil, nil
}

func (rd *Redis) Close() error {
	return nil
}

func (rd *Redis) Create(data Blob) (ID, error) {
	return NullID, nil
}

func (rd *Redis) LoadAll() ([]Item, error) {
	return nil, nil
}

func (rd *Redis) Load(id ID) (Blob, error) {
	return nil, nil
}

func (rd *Redis) Save(id ID, blob Blob) error {
	return nil
}

func (rd *Redis) Delete(id ID) error {
	return nil
}
