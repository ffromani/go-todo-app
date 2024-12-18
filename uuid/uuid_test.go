package uuid_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gotestbootcamp/go-todo-app/uuid"
)

func TestUUID(t *testing.T) {
	// TODO: change the code so it leverages this mock
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]string{"myuuid"})
	}))

	t.Cleanup(svr.Close)
	generator := uuid.New()
	uuid, err := generator.NewUUID()
	if err != nil {
		t.Fatalf("did not expect an error, got %v", err)
	}
	if uuid != "myuuid" {
		t.Fatalf("expecting myuuid, got %s", uuid)
	}
}
