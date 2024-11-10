package ledger

import (
	"fmt"
	"testing"
	"time"

	apiv1 "github.com/gotestbootcamp/go-todo-app/api/v1"
	"github.com/gotestbootcamp/go-todo-app/model"
	"github.com/gotestbootcamp/go-todo-app/store"
	"github.com/gotestbootcamp/go-todo-app/store/fake"
)

// solution:part1

func TestNewLoad(t *testing.T) {
	fakeStore, err := fake.NewMem()
	if err != nil {
		t.Fatalf("store creation failed: %v", err)
	}

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
	if err != nil {
		t.Fatalf("ledger creation failed: %v", err)
	}

	todo0, err := ld.Get("0")
	if err != nil {
		t.Fatalf("ledger get failed: %v", err)
	}
	if todo0.Title != "fake todo 0" {
		t.Fatalf("got unexpected data: %v", todo0)
	}
}
