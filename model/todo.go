package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"time"

	apiv1 "github.com/gotestbootcamp/go-todo-app/api/v1"
)

var (
	ErrAlreadyAssigned = errors.New("todo already assigned")
	ErrNotAssigned     = errors.New("todo not assigned")
	ErrFinalized       = errors.New("todo finalized")
)

// Todo represent a todo item managed by the system.
// Note: this incidentally is 1:1 with API objects, but this is an implementation
// detail rather than a requirement
type Todo struct {
	// Title is a short summary of the todo
	Title string
	// Assignee is the identifier of the agent working on the Todo
	Assignee string
	// Description is a longer description of the todo
	Description string
	// Status is the current processing status of the todo
	Status apiv1.Status
	// LastUpdateTime records the last time a todo was modified in any way in the system
	LastUpdateTime time.Time
}

// ToAPIv1 converts the object into the corresponding API layer object
func (td Todo) ToAPIv1() apiv1.Todo {
	return apiv1.Todo{
		Title:          td.Title,
		Assignee:       td.Assignee,
		Description:    td.Description,
		Status:         td.Status,
		LastUpdateTime: td.LastUpdateTime,
	}
}

// Serialize encodes the object in its canonical bytestream representation.
// If succesfull, returns the representation; otherwise the representation
// must be ignored, and the error will describe the failure.
func (td Todo) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(td)
	return buf.Bytes(), err
}

// Serialize decodes the object from its canonical bytestream representation.
// If succesfull, returns the decode object; otherwise returns a zero valued
// object, and the error will describe the failure.
func DeserializeTodo(data []byte) (Todo, error) {
	var todo Todo
	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&todo)
	return todo, err
}

// / NewFromAPIv1 creates a new object from its corresponding API layer object
func NewFromAPIv1(apiTodo apiv1.Todo) Todo {
	return Todo{
		Title:          apiTodo.Title,
		Description:    apiTodo.Description,
		Status:         apiv1.Pending,
		LastUpdateTime: time.Now(),
	}
}

var now = time.Now

// New creates a new Todo with the given title and with sane defaults
func New(title string) Todo {
	return Todo{
		Title:          title,
		Status:         apiv1.Pending,
		LastUpdateTime: now(),
	}
}

// IsOngoing returns true if the todo object is processable.
// In turn, an object is processable if not in a final state.
// An object in final state is terminated and can't be manipulated anymore
// (hence the "final").
func (td Todo) IsOngoing() bool {
	return td.Status == apiv1.Pending || td.Status == apiv1.Assigned
}

func (t Todo) HTMLRow() ([]byte, error) {
	const tt = `<tr>
    <td>{{ .Title }}</td>
    <td>{{ .Assignee }}</td>
    <td>{{ .Description }}</td>
    <td>{{ .Status }}</td>
    <td>{{ printTime .LastUpdateTime }}</td>
  </tr>`
	templ, err := template.New("tablerow").Funcs(
		template.FuncMap{
			"printTime": func(t time.Time) string {
				return t.Format(time.RFC3339)
			},
		}).Parse(tt)
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer
	err = templ.Execute(&buff, t)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// Describe changes the description of an object.
// This method is idempotent: the description can be changed any number of time
// while the object is processable. Returns error if the description update
// fails.
func (td *Todo) Describe(description string) error {
	if !td.IsOngoing() {
		return ErrFinalized
	}
	td.Description = description
	td.LastUpdateTime = time.Now()
	return nil
}

// Assign grants an assignee to a todo. Assignation can only be done once,
// e.g. Todos can't be reassigned once set. Returns error if the assignation fails.
func (td *Todo) Assign(assignee string) error {
	if !td.IsOngoing() {
		return ErrFinalized
	}
	if td.Assignee != "" {
		return ErrAlreadyAssigned
	}
	td.Assignee = assignee
	td.Status = apiv1.Assigned
	td.LastUpdateTime = time.Now()
	return nil
}

// Complete marks a todo as completed, which is a final state. Hence, a todo can be only completed once.
// Returns error if the completion fails.
func (td *Todo) Complete() error {
	if td.Status != apiv1.Assigned {
		return ErrNotAssigned
	}
	td.Status = apiv1.Completed
	td.LastUpdateTime = time.Now()
	return nil
}

// Delete marks a todo as deleted, which is a final state. Hence, a todo can be only deleted once.
// Note this is a soft-deletion. This method will (and must) not actually remove the Todo from the system.
// Returns error if the completion fails.
func (td *Todo) Delete() error {
	if !td.IsOngoing() {
		return ErrFinalized
	}
	td.Status = apiv1.Deleted
	td.LastUpdateTime = time.Now()
	return nil
}

// Merge takes two todo items, merges them into a new Todo item..
func Merge(td1, td2 Todo) (Todo, error) {
	if !td1.IsOngoing() || !td2.IsOngoing() {
		return Todo{}, ErrFinalized
	}
	if td1.Assignee != "" && td2.Assignee != "" && td1.Assignee != td2.Assignee {
		return Todo{}, errors.New("can't merge items with different assignees")
	}
	assignee := td1.Assignee
	if assignee == "" {
		assignee = td2.Assignee
	}

	status := apiv1.Pending
	if td1.Status == apiv1.Assigned || td2.Status == apiv1.Assigned {
		status = apiv1.Assigned
	}
	lastUpdateTime := td1.LastUpdateTime
	if lastUpdateTime.Before(td2.LastUpdateTime) {
		lastUpdateTime = td2.LastUpdateTime
	}

	res := Todo{
		Title:          fmt.Sprintf("%s-%s", td1.Title, td2.Title),
		Description:    fmt.Sprintf("%s-%s", td1.Description, td2.Description),
		Assignee:       assignee,
		Status:         status,
		LastUpdateTime: lastUpdateTime,
	}
	return res, nil
}
