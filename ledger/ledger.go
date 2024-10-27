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

// Ledger represents a Todo object store
type Ledger struct {
	storer store.Storage
	blobs  map[store.ID]store.Blob
}

// Item binds a Todo object with its ID. Note that IDs are managed and owned by the Ledger.
type Item struct {
	ID   store.ID    `json:"id"`
	Todo *model.Todo `json:"todo,omitempty"`
}

// ToAPIv1 converts a Item on its API layer corresponding object
func (it Item) ToAPIv1() apiv1.Item {
	apiTodo := it.Todo.ToAPIv1()
	return apiv1.Item{
		ID:   apiv1.ID(it.ID),
		Todo: &apiTodo,
	}
}

// Items is a collection of Item
type Items []Item

// ToAPIv1 converts Items, a Item collection, on its API layer corresponding object
func (its Items) ToAPIv1() []apiv1.Item {
	apiItems := make([]apiv1.Item, 0, len(its))
	for _, it := range its {
		apiItems = append(apiItems, it.ToAPIv1())
	}
	return apiItems
}

// Wants is a filter used when processing a collection of Todo object to derive a new collection.
// Returns true if the given todo object should be included in the resulting collection.
type Wants func(todo model.Todo) bool

// New creates and initializes a new Ledger based on the given datastore and its contents.
// To initialize itself, a Ledger eagerly loads all the content of the datastore.
// Returns error if the initialization fails; in this case, the returned ledger instance must be ignored.
func New(storer store.Storage) (*Ledger, error) {
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

// Close deinitializes this ledger and closes the attached datastore.
func (ld *Ledger) Close() error {
	return ld.storer.Close()
}

// Filter returns all the known Item which matches the give Wants filter.
// On failure, the error value is not nil and the resulting collection
// must be ignored.
func (ld *Ledger) Filter(wants Wants) (Items, error) {
	var items []Item
	for id, blob := range ld.blobs {
		todo, err := model.DeserializeTodo(blob)
		if err != nil {
			return items, err
		}
		if !wants(todo) {
			continue
		}
		log.Printf("ledger: Filter: object %v included", id)
		items = append(items, Item{
			ID:   id,
			Todo: &todo,
		})
	}
	return items, nil
}

// Get returns a todo object from its id. On failure, error is not nil
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

// Set creates or updates Todo objects in the store.
func (ld *Ledger) Set(id store.ID, todo model.Todo) (rerr error) {
	blob, err := todo.Serialize()
	if err != nil {
		return err
	}

	if id == store.NullID {
		return errors.New("can't set null id")
	}

	log.Printf("ledger: Set: updating object %v (%d bytes)", id, len(blob))
	curBlob, found := ld.blobs[id]
	if !found {
		ld.blobs[id] = blob
		log.Printf("ledger: Set: created cache object %v (%d bytes)", id, len(blob))
		rerr = ld.storer.Create(id, blob)
		log.Printf("ledger: Set: created store object %v err=%v", id, err)
		return rerr
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
	return rerr
}

// Delete removes a Todo from the ledger. The ledger may recycle IDs of deleted objects.
// On failure, error is not nil.
func (ld *Ledger) Delete(id store.ID) error {
	err := ld.storer.Delete(id)
	if err != nil {
		return err
	}
	delete(ld.blobs, id)
	return nil
}
