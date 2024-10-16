package v1

import (
	"bytes"
	"encoding/json"
	"time"
)

// Status represent the status of a Todo
// Transition between statuses are enforced by the Todo object methods.
type Status string

const (
	// Pending means a Todo is in the common backlog
	Pending Status = "pending"
	// Assigned means a Todo has got an assignee, and work has thus begun
	Assigned Status = "assigned"
	// Completed means a Todo has been completed by its assignee and is no longer active
	Completed Status = "completed"
	// Deleted means a Todo has been deleted, regardless of its previous state, and it is no longer relevant
	Deleted Status = "deleted"
)

// ID is an opaque value which uniquely identifies a Todo. Can only be compared for equality
type ID int64

// Todo represent a todo item managed by the system
type Todo struct {
	// Title is a short summary of the todo
	Title string `json:"title"`
	// Assignee is the identifier of the agent working on the Todo
	Assignee string `json:"assignee,omitempty"`
	// Description is a longer description of the todo
	Description string `json:"description,omitempty"`
	// Status is the current processing status of the todo
	Status Status `json:"status"`
	// LastUpdateTime records the last time a todo was modified in any way in the system
	LastUpdateTime time.Time `json:"updated"`
}

// ToJSON returns a bytestream JSON encoding of the Todo; if succesfull, err is nil;
// otherwise contains the encoding error.
func (td Todo) ToJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(td)
	return buf.Bytes(), err
}

// NewTodoFromJSON creates a new todo from a bytestream, which is expected to contain
// a valid JSON encoding (e.g. from ToJSON). If succesfull, returns the decoded todo
// and error is nil; otherwise, returns a zero-valued Todo and an error describing the
// decoding failure
func NewTodoFromJSON(data []byte) (Todo, error) {
	var todo Todo
	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&todo)
	return todo, err
}

// Item binds a Todo with its ID identifier
type Item struct {
	// ID is the ID which identifies the Todo processed by the operation
	ID ID `json:"id"`
	// Todo is the todo object processed by the operation. If the object
	// is implicit and can be unanbiguosly inferred, can be omitted
	Todo *Todo `json:"todo,omitempty"`
}

// Errors give informations about a processing error
type Error struct {
	// Processing error code. If positive, it is a HTTP status code
	Code int `json:"code"`
	// Optional human friendly description of the error
	Text string `json:"text,omitempty"`
}

// Result represent the status of a succesfull processing.
type Result struct {
	// Items includes the updated objects as returned by the operation.
	// Can be empty in succesfull operations (e.g. a query produced no values)
	Items []Item `json:"items,omitempty"`
	// Optional human friendly description of the operation
	Text string `json:"text,omitempty"`
}

// ResponseStatus represent the overlal status (e.g. success/error) of a operation
type ResponseStatus string

const (
	ResponseSuccess ResponseStatus = "success"
	ResponseError   ResponseStatus = "error"
)

// Response represents the outcome of a operation
type Response struct {
	Status ResponseStatus `json:"status"`
	Error  *Error         `json:"error,omitempty"`
	Result *Result        `json:"result,omitempty"`
}
