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
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	addressFlag := flag.String("a", "localhost:8080", "server address")
	baseURLFlag := flag.String("b", "http://localhost:8080", "base URL for shortened links")

	flag.Parse()

	config := &Config{
		Address: *addressFlag,
		BaseURL: *baseURLFlag,
	}

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		config.Address = envRunAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		config.BaseURL = envBaseURL
	}

	return config
}
