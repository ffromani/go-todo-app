package middleware

// solution:part1

import (
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	fake := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	stdout := os.Stdout
	defer func() {
		os.Stdout = stdout
	}()
	f, err := os.CreateTemp("", "midlogger")
	if err != nil {
		t.Fatalf("temp file creation failed: %v", err)
	}
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
	if err != nil {
		t.Fatalf("file read failed: %v", err)
	}
	ok, err := regexp.Match(".*middleware.*GET   	/test	testlogger.*", got)
	if err != nil {
		t.Fatalf("regexp match failed: %v", err)
	}
	if !ok {
		t.Fatalf("logging does not match the expectation: %q", string(got))
	}
}
