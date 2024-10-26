package ledger

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	apiv1 "github.com/ffromani/go-todo-app/api/v1"
	"github.com/ffromani/go-todo-app/model"
	"github.com/ffromani/go-todo-app/store"
	"github.com/ffromani/go-todo-app/store/fake"
)

// solution:part1

func TestNewLoad(t *testing.T) {
	fakeStore, err := fake.NewMem()
	assert.NoError(t, err)

	num := 0
	count := 5 // random number > 1
	fakeStore.Generate = func() (store.Item, bool, error) {
		done := (num == count)
		todo := model.Todo{
			Title:          fmt.Sprintf("fake todo %d", num),
			Description:    fmt.Sprintf("description %d", num),
			Status:         apiv1.Pending,
			LastUpdateTime: time.Now(),
		}
		data, err := todo.Serialize()
		if err != nil {
			return store.Item{}, true, err
		}
		item := store.Item{
			ID:   store.ID(fmt.Sprintf("%d", num)),
			Blob: store.Blob(data),
		}
		num += 1
		return item, done, nil
	}

	ld, err := New(fakeStore)
	assert.NoError(t, err)

	todo0, err := ld.Get("0")
	assert.NoError(t, err)
	assert.Equal(t, todo0.Title, "fake todo 0")
}
