package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	cfg := Defaults()
	assert.NotEmpty(t, cfg.Address, "Address is not set")
}

func TestString(t *testing.T) {
	addr := "localhost:8181"
	cfg := Config{
		Address: addr,
	}
	res := cfg.String()
	assert.NotEmpty(t, res, "String result is empty")
	assert.Contains(t, res, addr, "Expected address not in string result")
}
