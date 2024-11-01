package store

import "fmt"

// ID is an opaque value which uniquely identifies a Todo. Can only be compared for equality
// Note: this incidentally is 1:1 with API objects, but this is an implementation
// detail rather than a requirement
type ID string
type Blob []byte

func (b Blob) MarshalBinary() ([]byte, error) {
	return b, nil
}

const (
	// NullID represents a invalid ID
	NullID ID = ""
)

type Storage interface {
	Close() error
	Create(ID, Blob) error
	LoadAll() ([]Item, error)
	Load(ID) (Blob, error)
	Save(ID, Blob) error
	Delete(ID) error
}

// Item binds a Todo with its ID identifier
// Note: this incidentally is 1:1 with API objects, but this is an implementation
// detail rather than a requirement
type Item struct {
	ID   ID
	Blob Blob
}

type ErrNotFound struct {
	ID ID
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("unknown id: %v", e.ID)
}

type ErrCorruptedContent struct {
	Name string
}

func (e ErrCorruptedContent) Error() string {
	return fmt.Sprintf("corrupted content: %q", e.Name)
}
