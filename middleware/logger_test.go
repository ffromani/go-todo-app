package middleware

// exercise:part1

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	fake := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	stdout := os.Stdout
	defer func() {
		os.Stdout = stdout
	}()
	f, err := os.CreateTemp("", "midlogger")
	assert.NoError(t, err)
	defer func() {
		os.Remove(f.Name())
	}()
	os.Stdout = f

	req := httptest.NewRequest("GET", "/test", strings.NewReader("{}"))
	rec := httptest.NewRecorder()
	hndl := Logger(fake, "testlogger")
	hndl.ServeHTTP(rec, req)

	f.Sync()
	got, err := os.ReadFile(f.Name())
	assert.NoError(t, err)
	assert.Regexp(t, "middleware.*GET   	/test	testlogger.*", string(got))
}
