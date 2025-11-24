package config

import (
	"flag"
	"os"
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

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		address = envRunAddr
	} else {
		flag.StringVar(&address, "a", "localhost:8080", "server address")
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		baseURL = envBaseURL
	} else {
		flag.StringVar(&baseURL, "b", "http://localhost:8080", "base URL for shortened links")
	}

	flag.Parse()

	return &Config{
		Address: address,
		BaseURL: baseURL,
	}
}
