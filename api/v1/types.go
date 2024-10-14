package v1

import (
	"bytes"
	"encoding/json"
	"time"
)

type Status string

const (
	Pending   Status = "pending"
	Assigned  Status = "assigned"
	Completed Status = "completed"
	Deleted   Status = "deleted"
)

type ID int64

type Todo struct {
	Title          string    `json:"title"`
	Assignee       string    `json:"assignee,omitempty"`
	Description    string    `json:"description,omitempty"`
	Status         Status    `json:"status"`
	LastUpdateTime time.Time `json:"updated"`
}

func (td Todo) ToJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(td)
	return buf.Bytes(), err
}

func NewTodoFromJSON(data []byte) (Todo, error) {
	var todo Todo
	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&todo)
	return todo, err
}

type Item struct {
	ID   ID    `json:"id"`
	Todo *Todo `json:"todo,omitempty"`
}

type Error struct {
	Code int    `json:"code"`
	Text string `json:"text,omitempty"`
}

type Result struct {
	Items []Item `json:"items,omitempty"`
	Text  string `json:"text,omitempty"`
}

type ResponseStatus string

const (
	ResponseSuccess ResponseStatus = "success"
	ResponseError   ResponseStatus = "error"
)

type Response struct {
	Status ResponseStatus `json:"status"`
	Error  *Error         `json:"error,omitempty"`
	Result *Result        `json:"result,omitempty"`
}
