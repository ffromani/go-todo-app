package fake

import (
	"github.com/ffromani/go-todo-app/store"
)

type Mem struct {
	Blobs        map[store.ID]store.Blob
	LastObjectID store.ID
	Error        error
	Generate     func() (store.Item, bool, error)
}

func NewMem() (*Mem, error) {
	return &Mem{
		Blobs:        make(map[store.ID]store.Blob),
		LastObjectID: 0,
		Generate: func() (store.Item, bool, error) {
			return store.Item{}, true, nil
		},
	}, nil
}

func (mm *Mem) Close() error {
	return mm.Error
}

func (mm *Mem) Create(data store.Blob) (store.ID, error) {
	if mm.Error != nil {
		return store.NullID, mm.Error
	}
	objectID := mm.LastObjectID + 1
	mm.Blobs[objectID] = data
	mm.LastObjectID = objectID
	return objectID, nil

}

func (mm *Mem) LoadAll() ([]store.Item, error) {
	if mm.Error != nil {
		return nil, mm.Error
	}
	var items []store.Item
	for {
		item, done, err := mm.Generate()
		if err != nil {
			return items, err
		}
		if done {
			break
		}
		items = append(items, item)
	}
	return items, nil
}

func (mm *Mem) Load(id store.ID) (store.Blob, error) {
	if mm.Error != nil {
		return nil, mm.Error
	}
	blob, ok := mm.Blobs[id]
	if !ok {
		return nil, store.ErrNotFound{ID: id}
	}
	return blob, nil
}

func (mm *Mem) Save(id store.ID, blob store.Blob) error {
	if mm.Error != nil {
		return mm.Error
	}
	mm.Blobs[id] = blob
	return nil
}

func (mm *Mem) Delete(id store.ID) error {
	if mm.Error != nil {
		return mm.Error
	}
	delete(mm.Blobs, id)
	return nil
}
