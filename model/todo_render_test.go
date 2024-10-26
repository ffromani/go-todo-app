package model_test

// exercise

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	apiv1 "github.com/ffromani/go-todo-app/api/v1"
	"github.com/ffromani/go-todo-app/model"
)

var update = flag.Bool("update", false, "update .golden.json files")

func TestRender(t *testing.T) {

	tests := []struct {
		name     string
		toRender model.Todo
	}{
		{
			name: "non assigned",
			toRender: model.Todo{
				Title:       "todo2",
				Assignee:    "",
				Description: "second todo",
				Status:      apiv1.Pending,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rendered, err := tc.toRender.HTMLRow()
			if err != nil {
				t.Fatal("HTMLRow should not fail", err)
			}
			goldenFile := filepath.Join("testdata", strings.ReplaceAll(tc.name, " ", "-")+".golden")
			if *update {
				if err := os.WriteFile(goldenFile, rendered, os.ModePerm); err != nil {
					panic("write failed")
				}
			}
			expected, err := os.ReadFile(goldenFile)
			if err != nil {
				t.Errorf("failed to open golden file %s: %v", goldenFile, err)
			}
			if !bytes.Equal(expected, rendered) {
				t.Fail()
			}
		})
	}
}
