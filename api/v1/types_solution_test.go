package v1

// extra:part1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	//	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

// this is not sufficient to ensure full coverage, yet:
// ok  	github.com/gotestbootcamp/go-todo-app/api/v1	0.009s	coverage: 100.0% of statements
func TestJSONRoundtrip(t *testing.T) {
	type testCase struct {
		name string
		todo Todo
	}

	for _, tcase := range []testCase{
		{
			name: "empty",
			todo: Todo{},
		},
		/*
			{
				name: "full init",
				todo: Todo{
					Title:          "test todo 1",
					Assignee:       "Jane Doe",
					Description:    "testing data with all fields set",
					Status:         Assigned,
					LastUpdateTime: time.Now(), // intentionally changing every run
				},
			},
		*/
	} {
		t.Run(tcase.name, func(t *testing.T) {
			data, err := tcase.todo.ToJSON()
			assert.NoError(t, err)
			todo2, err := NewTodoFromJSON(data)
			assert.NoError(t, err)

			//			options := []cmp.Option{
			//				cmpopts.IgnoreFields(Todo{}, "LastUpdateTime"),
			//			}
			//			delta := cmp.Diff(tcase.todo, todo2, options...)
			delta := cmp.Diff(tcase.todo, todo2)
			assert.Empty(t, delta, "roundtripped object differs from expected object")
		})
	}
}
