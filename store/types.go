package store

import "fmt"

type ID int64
type Blob []byte

const (
	NullID ID = 0
)

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
