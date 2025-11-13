package config

import (
	"flag"
	"log"
)

type Config struct {
	Address string
	BaseURL string
}

func MustLoad() *Config {
	var (
		address string
		baseURL string
	)

	flag.StringVar(&address, "a", "localhost:8080", "server address")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "base URL for shortened links")

	flag.Parse()

	if address == "" {
		log.Fatal("address not set")
	}
	if baseURL == "" {
		log.Fatal("base URL not set")
	}

	return &Config{
		Address: address,
		BaseURL: baseURL,
	}
}
