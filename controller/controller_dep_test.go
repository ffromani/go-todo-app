package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	apiv1 "github.com/gotestbootcamp/go-todo-app/api/v1"
	"github.com/gotestbootcamp/go-todo-app/controller"
	"github.com/gotestbootcamp/go-todo-app/ledger"
	"github.com/gotestbootcamp/go-todo-app/model"
	"github.com/gotestbootcamp/go-todo-app/store/fake"
)

type mockUUID struct {
	uuid      string
	shouldErr bool
}

func (m mockUUID) NewUUID() (string, error) {
	if m.shouldErr {
		return "", fmt.Errorf("error")
	}
	return m.uuid, nil
}

var _ controller.UUIDer = mockUUID{}

func TestTodoCreate(t *testing.T) {
	uuid := mockUUID{uuid: "1234"}
	handler := controller.New(memoryStorage(), uuid)

	t.Run("test post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/todos", bodyFromTodo(model.Todo{Title: "foo", Assignee: "fede", Description: "todo"}))
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		apiRes := &apiv1.Response{}
		err := json.NewDecoder(res.Body).Decode(&apiRes)
		if err != nil {
			t.Fatalf("expected error to be nil got %v", err)
		}

		if len(apiRes.Result.Items) != 1 {
			t.Fatalf("expecting one item back")
		}

		if apiRes.Result.Items[0].ID != "1234" {
			t.Fatalf("expecting id 1234 got %s", apiRes.Result.Items[0].ID)
		}
	})

	t.Run("test post should fail when uuid fails", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/todos", bodyFromTodo(model.Todo{Title: "foo", Assignee: "fede", Description: "todo"}))
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		apiRes := &apiv1.Response{}
		err := json.NewDecoder(res.Body).Decode(&apiRes)
		if err != nil {
			t.Fatalf("expected error to be nil got %v", err)
		}

		if len(apiRes.Result.Items) != 1 {
			t.Fatalf("expecting one item back")
		}

		if apiRes.Result.Items[0].ID != "1234" {
			t.Fatalf("expecting id 1234 got %s", apiRes.Result.Items[0].ID)
		}
	})
}

func memoryStorage() *ledger.Ledger {
	st, err := fake.NewMem()
	if err != nil {
		panic("failed to initialize the memory storage")
	}
	ldg, err := ledger.New(st)
	if err != nil {
		panic("failed to initialize the ledger")
	}

	return ldg
}
func bodyFromTodo(t model.Todo) io.Reader {
	serialized, err := t.Serialize()
	if err != nil {
		panic("")
	}

	res := bytes.NewReader(serialized)
	return res
}
