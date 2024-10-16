package config

import (
	"fmt"
	"strings"
)

// Config holds all the tunables
type Config struct {
	// Address is in the format `[host]:port`
	Address string
}

func (cfg Config) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "- address: %s\n", cfg.Address)
	return sb.String()
}

// Defaults return a Config initialized with the compiled-in defaults
func Defaults() Config {
	return Config{
		Address: "localhost:8181",
	}
}
