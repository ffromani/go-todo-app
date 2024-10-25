package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

// this is not sufficient to ensure full coverage, yet:
// ok  	github.com/ffromani/go-todo-app/config	0.008s	coverage: 100.0% of statements
func TestFromFlags(t *testing.T) {
	type testCase struct {
		name        string
		args        []string
		expectedCfg Config
		expectedErr bool
	}

	for _, tcase := range []testCase{
		{
			name:        "no flags", // but argv[0] is always set
			args:        []string{"test"},
			expectedCfg: Defaults(),
			expectedErr: false,
		},
	} {
		t.Run(tcase.name, func(t *testing.T) {
			cfg, err := FromFlags(tcase.args...)
			gotErr := (err != nil)
			assert.Equal(t, gotErr, tcase.expectedErr, "error not expected value")
			delta := cmp.Diff(tcase.expectedCfg, cfg)
			assert.Empty(t, delta, "got object differs from expected object")
		})
	}
}
