package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	apiv1 "github.com/ffromani/go-todo-app/api/v1"
	"github.com/ffromani/go-todo-app/ledger"
	"github.com/ffromani/go-todo-app/model"
	"github.com/ffromani/go-todo-app/store"
	"github.com/ffromani/go-todo-app/store/fake"
)

// solution:part1

func TestBacklogIndex(t *testing.T) {
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

	ld, err := ledger.New(fakeStore)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/backlog", strings.NewReader("{}"))
	rec := httptest.NewRecorder()

	hndl := New(ld)
	hndl.ServeHTTP(rec, req)

	res := rec.Body.String()
	assert.Equal(t, rec.Code, http.StatusOK, "response: %q", res)

	var resp apiv1.Response
	err = json.Unmarshal([]byte(res), &resp)
	assert.NoError(t, err)

	assert.Equal(t, resp.Status, apiv1.ResponseSuccess)
	assert.Nil(t, resp.Error)
	assert.NotNil(t, resp.Result)
	assert.Len(t, resp.Result.Items, count)
}
