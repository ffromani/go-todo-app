package config

type Config struct {
	// Address is in the format `[host]:port`
	Address string `json:"address"`
}

func Defaults() Config {
	return Config{
		Address: "localhost:8181",
	}
}
