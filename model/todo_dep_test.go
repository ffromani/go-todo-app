package model_test

import (
	"testing"
	"time"

	apiv1 "github.com/ffromani/go-todo-app/api/v1"
	"github.com/ffromani/go-todo-app/model"
)

func TestNewTodo(t *testing.T) {
	newTodo := model.New("foo")

	updateTime, err := time.Parse("2006-Jan-02", "2014-Feb-04")
	if err != nil {
		panic(err)
	}
	toCompare := model.Todo{
		Title:          "foo",
		Status:         apiv1.Pending,
		LastUpdateTime: updateTime,
	}

	if newTodo != toCompare {
		t.Fatalf("expecting %v, got %v", toCompare, newTodo)
	}
}
