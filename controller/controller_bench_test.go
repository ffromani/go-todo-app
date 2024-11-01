package controller

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ffromani/go-todo-app/model"
)

// exercise

func BenchmarkTodoFromRequest(b *testing.B) {
	todo := model.New("todo")
	_ = todo.Assign("fede")
	_ = todo.Describe("hello")

	serialized, err := todo.Serialize()
	if err != nil {
		panic("")
	}

	body := bytes.NewReader(serialized)
	for i := 0; i < b.N; i++ {
		_, _ = body.Seek(0, io.SeekStart)
		req := httptest.NewRequest(http.MethodGet, "/foo", body)
		_, code, err := todoFromRequest(req)
		// _, code, err := todoFromRequestReader(req)
		if err != nil {
			b.Fatal("error", err)
		}
		if code != 0 {
			b.Fatal("error code", code)
		}
	}
}
