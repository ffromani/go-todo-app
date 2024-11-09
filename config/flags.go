package config

import (
	"flag"
	"fmt"
	"os"
)

// FromFlags creates a Config object out of the command line args
// If succesfull, returns the resulting Config; otherwise returns
// a zero-valued Config and the error describing the failure.
func FromFlags(args ...string) (Config, error) {
	conf := Defaults()

	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flags.StringVar(&conf.Address, "url", conf.Address, "url to listen to")
	flags.StringVar(&conf.Redis.URL, "redis-url", conf.Redis.URL, "redis URL")
	flags.StringVar(&conf.Redis.Password, "redis-password", conf.Redis.Password, "redis password")
	flags.IntVar(&conf.Redis.Database, "redis-database", conf.Redis.Database, "redis database index")
	flags.StringVar(&conf.FSDir.Path, "fsdir-path", conf.FSDir.Path, "filesystem directory store path")

	flags.Usage = func() {
		w := flags.Output()
		fmt.Fprintf(w, "Usage of %s:\n", os.Args[0])
		flags.PrintDefaults()
		os.Exit(0)
	}

	err := flags.Parse(args)
	return conf, err
}
