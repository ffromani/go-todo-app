package config

import (
	"fmt"
	"strings"
)

// RedisConfig holds all the redis-related tunables
type RedisConfig struct {
	URL      string
	Password string
	Database int
}

// Config holds all the tunables
type Config struct {
	// Address is in the format `[host]:port`
	Address string
	Redis   RedisConfig
}

func (cfg Config) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "- address: %s\n", cfg.Address)
	fmt.Fprintf(&sb, "- redis:\n")
	fmt.Fprintf(&sb, "  - url:  %q\n", cfg.Redis.URL)
	fmt.Fprintf(&sb, "  - pass: %q\n", cfg.Redis.Password)
	fmt.Fprintf(&sb, "  - db:   %d\n", cfg.Redis.Database)
	return sb.String()
}

// Defaults return a Config initialized with the compiled-in defaults
func Defaults() Config {
	return Config{
		Address: "localhost:8181",
		Redis:   RedisConfig{},
	}
}
