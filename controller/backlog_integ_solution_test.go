package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	apiv1 "github.com/gotestbootcamp/go-todo-app/api/v1"
	"github.com/gotestbootcamp/go-todo-app/ledger"
	"github.com/gotestbootcamp/go-todo-app/model"
	"github.com/gotestbootcamp/go-todo-app/store"
	"github.com/gotestbootcamp/go-todo-app/store/fake"
)

// solution:part1

func TestBacklogIndex(t *testing.T) {
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

	ld, err := ledger.New(fakeStore)
	if err != nil {
		t.Fatalf("ledger creation failed: %v", err)
	}

	req := httptest.NewRequest("GET", "/backlog", strings.NewReader("{}"))
	rec := httptest.NewRecorder()

	hndl := New(ld)
	hndl.ServeHTTP(rec, req)

	res := rec.Body.String()
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected response: code=%d body=%q", rec.Code, res)
	}

	var resp apiv1.Response
	err = json.Unmarshal([]byte(res), &resp)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	checkResponseIsSuccessWithItems(t, resp, count)
}

func checkResponseIsSuccessWithItems(t *testing.T, resp apiv1.Response, count int) {
	t.Helper()
	if resp.Status != apiv1.ResponseSuccess {
		t.Fatalf("response status code is not success: %v", resp.Status)
	}
	if resp.Error != nil {
		t.Fatalf("returned unexpected error: %v", resp.Error)
	}
	if resp.Result == nil {
		t.Fatalf("returned nil result")
	}
	if len(resp.Result.Items) != count {
		t.Fatalf("returned unexpected item count; got %d expected %d", len(resp.Result.Items), count)
	}
}
