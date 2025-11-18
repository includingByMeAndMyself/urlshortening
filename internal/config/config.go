package config

import (
	"flag"
	"log"
)

type Config struct {
	Address string
}

func MustLoad() *Config {
	var address string
	flag.StringVar(&address, "a", "localhost:8080", "server address")
	flag.Parse()

	if address == "" {
		log.Fatal("address not set")
	}

	return &Config{Address: address}
}
