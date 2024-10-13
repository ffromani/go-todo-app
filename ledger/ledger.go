package ledger

import (
	"errors"
	"log"

	apiv1 "github.com/ffromani/go-todo-app/api/v1"
	"github.com/ffromani/go-todo-app/model"
	"github.com/ffromani/go-todo-app/store"
)

var (
	ErrNotFound = errors.New("object not found")
)

type BlobStorer interface {
	Close() error
	Create(data store.Blob) (store.ID, error)
	LoadAll() ([]store.Item, error)
	Load(id store.ID) (store.Blob, error)
	Save(id store.ID, blob store.Blob) error
	Delete(id store.ID) error
}

type Ledger struct {
	storer BlobStorer
	blobs  map[store.ID]store.Blob
}

type Item struct {
	ID   store.ID    `json:"id"`
	Todo *model.Todo `json:"todo,omitempty"`
}

func (it Item) ToAPIv1() apiv1.Item {
	apiTodo := it.Todo.ToAPIv1()
	return apiv1.Item{
		ID:   apiv1.ID(it.ID),
		Todo: &apiTodo,
	}
}

type Items []Item

func (its Items) ToAPIv1() []apiv1.Item {
	apiItems := make([]apiv1.Item, 0, len(its))
	for _, it := range its {
		apiItems = append(apiItems, it.ToAPIv1())
	}
	return apiItems
}

type Wants func(todo model.Todo) bool

func New(storer BlobStorer) (*Ledger, error) {
	items, err := storer.LoadAll()
	if err != nil {
		return nil, err
	}
	ld := Ledger{
		storer: storer,
		blobs:  make(map[store.ID]store.Blob, len(items)),
	}
	for _, item := range items {
		ld.blobs[item.ID] = item.Blob
	}
	return &ld, nil

}

func (ld *Ledger) Close() error {
	return ld.storer.Close()
}

func (ld *Ledger) Filter(wants Wants) (Items, error) {
	var items []Item
	for id, blob := range ld.blobs {
		log.Printf("ledger: Filter: object %v (%d bytes)", id, len(blob))
		todo, err := model.DeserializeTodo(blob)
		if err != nil {
			return items, err
		}
		if !wants(todo) {
			continue
		}
		log.Printf("ledger: Filter: object %v is wanted", id)
		items = append(items, Item{
			ID:   id,
			Todo: &todo,
		})
	}
	return items, nil
}

func (ld *Ledger) Get(id store.ID) (model.Todo, error) {
	blob, ok := ld.blobs[id]
	if !ok {
		return model.Todo{}, store.ErrNotFound{ID: id}
	}
	todo, err := model.DeserializeTodo(blob)
	if err != nil {
		return model.Todo{}, err
	}
	return todo, nil
}

// Set creates or updates Todo objects in the store
func (ld *Ledger) Set(id store.ID, todo model.Todo) (blobID store.ID, rerr error) {
	blob, err := todo.Serialize()
	if err != nil {
		return store.NullID, err
	}

	if id == store.NullID {
		log.Printf("ledger: Set: create new object (%d bytes)", len(blob))
		todoID, err := ld.storer.Create(blob)
		if err != nil {
			return store.NullID, err
		}
		ld.blobs[todoID] = blob
		log.Printf("ledger: Set: created new object %v (%d bytes)", todoID, len(blob))
		return todoID, nil
	}

	log.Printf("ledger: Set: updaring object %v (%v)", id, blob)
	curBlob, found := ld.blobs[id]
	if !found {
		return store.NullID, ErrNotFound
	}
	// rollback
	defer func() {
		if rerr == nil {
			return
		}
		log.Printf("ledger: Set: rollbacking object %v", id)
		ld.blobs[id] = curBlob
	}()
	ld.blobs[id] = blob
	log.Printf("ledger: Set: updated cache object %v (%d bytes)", id, len(blob))
	rerr = ld.storer.Save(id, blob)
	log.Printf("ledger: Set: updated store object %v err=%v", id, err)
	return id, rerr
}

func (ld *Ledger) Delete(id store.ID) error {
	err := ld.storer.Delete(id)
	if err != nil {
		return err
	}
	delete(ld.blobs, id)
	return nil
}
