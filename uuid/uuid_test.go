package uuid

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUUID(t *testing.T) {
	oldUrl := uuidURL
	defer func() {
		uuidURL = oldUrl
	}()

	tests := []struct {
		name       string
		uuid       string
		shouldFail bool
	}{
		{
			name: "short",
			uuid: "123",
		},
		{
			name: "long",
			uuid: "1234567890",
		},
		{
			name:       "should error",
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := 0

			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called++
				if tt.shouldFail {
					http.Error(w, "error", http.StatusInternalServerError)
					return
				}
				json.NewEncoder(w).Encode([]string{tt.uuid})
			}))

			t.Cleanup(svr.Close)
			uuidURL = svr.URL

			generator := New()
			uuid, err := generator.NewUUID()
			if called != 1 {
				t.Fatalf("expecting one call to the server, got %d", called)
			}
			if err == nil && tt.shouldFail {
				t.Fatalf("expecting an error, got success")
			}
			if err != nil && !tt.shouldFail {
				t.Fatalf("did not expect an error, got %v", err)
			}
			if !tt.shouldFail && uuid != tt.uuid {
				t.Fatalf("expecting myuuid, got %s", uuid)
			}
		})
	}
}
