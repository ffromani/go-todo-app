package model_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	apiv1 "github.com/gotestbootcamp/go-todo-app/api/v1"
	"github.com/gotestbootcamp/go-todo-app/model"
)

func TestMerge(t *testing.T) {
	// Setup code
	toMerge := model.Todo{
		Title:       "todo1",
		Assignee:    "",
		Description: "first todo",
		Status:      apiv1.Pending,
	}

	tests := []struct {
		name      string
		toMerge   model.Todo
		expected  model.Todo
		shouldErr bool
	}{
		{
			name: "non assigned",
			toMerge: model.Todo{
				Title:       "todo2",
				Assignee:    "",
				Description: "second todo",
				Status:      apiv1.Pending,
			},
			expected: model.Todo{
				Title:       "todo1-todo2",
				Assignee:    "",
				Description: "first todo-second todo",
				Status:      apiv1.Pending,
			},
			shouldErr: false,
		},
		{
			name: "assigned",
			toMerge: model.Todo{
				Title:       "todoassigned",
				Assignee:    "fede",
				Description: "second todo assigned",
				Status:      apiv1.Assigned,
			},
			expected: model.Todo{
				Title:       "todo1-todoassigned",
				Assignee:    "fede",
				Description: "first todo-second todo assigned",
				Status:      apiv1.Assigned,
			},
			shouldErr: false,
		},
		{
			name: "should fail when completed",
			toMerge: model.Todo{
				Title:       "todoassigned",
				Assignee:    "fede",
				Description: "second todo assigned",
				Status:      apiv1.Completed,
			},
			expected:  model.Todo{},
			shouldErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := model.Merge(toMerge, tc.toMerge)
			if err != nil && !tc.shouldErr {
				t.Fatal("merged failed while should not have", err)
			}
			if err == nil && tc.shouldErr {
				t.Fatal("merged failed while should have")
			}
			if res != tc.expected {
				t.Fatalf("expecting %v, got %v", spew.Sdump(res), spew.Sdump(tc.expected))
			}
		})
	}
}
