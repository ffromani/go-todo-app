package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	apiv1 "github.com/ffromani/go-todo-app/api/v1"
)

var (
	ErrAlreadyAssigned = errors.New("todo already assigned")
	ErrNotAssigned     = errors.New("todo not assigned")
	ErrNotStarted      = errors.New("todo never started")
)

type Todo struct {
	Title          string
	Assignee       string
	Description    string
	Status         apiv1.Status
	LastUpdateTime time.Time
}

func (td Todo) ToAPIv1() apiv1.Todo {
	return apiv1.Todo{
		Title:          td.Title,
		Assignee:       td.Assignee,
		Description:    td.Description,
		Status:         td.Status,
		LastUpdateTime: td.LastUpdateTime,
	}
}

func (td Todo) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(td)
	return buf.Bytes(), err
}

func DeserializeTodo(data []byte) (Todo, error) {
	var todo Todo
	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&todo)
	return todo, err
}

func NewFromAPITodo(apiTodo apiv1.Todo) Todo {
	return Todo{
		Title:          apiTodo.Title,
		Description:    apiTodo.Description,
		Status:         apiv1.Pending,
		LastUpdateTime: time.Now(),
	}
}

func New(title string) Todo {
	return Todo{
		Title:          title,
		Status:         apiv1.Pending,
		LastUpdateTime: time.Now(),
	}
}

func (td Todo) Assign(assignee string) error {
	if td.Assignee != "" {
		return ErrAlreadyAssigned
	}
	td.Assignee = assignee
	td.Status = apiv1.Assigned
	td.LastUpdateTime = time.Now()
	return nil
}

func (td Todo) Begin() error {
	if td.Status != apiv1.Assigned {
		return ErrNotAssigned
	}
	td.Status = apiv1.Ongoing
	td.LastUpdateTime = time.Now()
	return nil
}

func (td Todo) Complete() error {
	if td.Status != apiv1.Ongoing {
		return ErrNotAssigned
	}
	td.Status = apiv1.Completed
	td.LastUpdateTime = time.Now()
	return nil
}

func (td Todo) Delete() {
	td.Status = apiv1.Deleted
	td.LastUpdateTime = time.Now()
}
