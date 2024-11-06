package model_test

// exercise

import (
	"testing"

	apiv1 "github.com/gotestbootcamp/go-todo-app/api/v1"
	"github.com/gotestbootcamp/go-todo-app/model"
)

func TestMerge(t *testing.T) {
	first := model.Todo{
		Title:       "todo1",
		Assignee:    "",
		Description: "first todo",
		Status:      apiv1.Pending,
	}

	second := model.Todo{
		Title:       "todo2",
		Assignee:    "",
		Description: "second todo",
		Status:      apiv1.Pending,
	}

	res, err := model.Merge(first, second)
	if err != nil {
		t.Fatal("merged failed", err)
	}
	expected := model.Todo{
		Title:       "todo1-todo2",
		Description: "first todo-second todo",
		Status:      apiv1.Pending,
	}
	if res != expected {
		t.Fatal("merged failed", err)
	}
}
