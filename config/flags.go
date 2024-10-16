package config

import "flag"

// FromFlags creates a Config object out of the command line args
// If succesfull, returns the resulting Config; otherwise returns
// a zero-valued Config and the error describing the failure.
func FromFlags(args ...string) (Config, error) {
	conf := Defaults()

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.StringVar(&conf.Address, "address", conf.Address, "address to listen to")

	err := flags.Parse(args)
	return conf, err
}
