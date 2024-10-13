package config

import "flag"

func FromFlags(args ...string) (Config, error) {
	conf := Defaults()

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.StringVar(&conf.Address, "address", conf.Address, "address to listen to")

	err := flags.Parse(args)
	return conf, err
}
