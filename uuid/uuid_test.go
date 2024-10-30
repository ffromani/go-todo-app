package uuid_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ffromani/go-todo-app/uuid"
)

// exercise
func TestUUID(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]string{"uuid"})
	}))

	t.Cleanup(svr.Close)
	generator := uuid.New()
	uuid, err := generator.NewUUID()
	if err != nil {
		t.Fatalf("did not expect an error, got %v", err)
	}
	fmt.Printf("TODO find a way to verify the uuid", uuid)
}
