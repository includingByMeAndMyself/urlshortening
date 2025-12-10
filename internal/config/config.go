package config

import (
	"flag"
	"os"
	"path/filepath"
)

type Config struct {
	Address         string
	BaseURL         string
	FileStoragePath string
}

func Load() *Config {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	addressFlag := flag.String("a", "localhost:8080", "server address")
	baseURLFlag := flag.String("b", "http://localhost:8080", "base URL for shortened links")

	defaultStoragePath := filepath.Join(os.TempDir(), "shortener-db.json")
	fileStorageFlag := flag.String("f", defaultStoragePath, "path to file storage")

	flag.Parse()

	config := &Config{
		Address:         *addressFlag,
		BaseURL:         *baseURLFlag,
		FileStoragePath: *fileStorageFlag,
	}

	if envRunAddr, ok := os.LookupEnv("SERVER_ADDRESS"); ok && envRunAddr != "" {
		config.Address = envRunAddr
	}

	if envBaseURL, ok := os.LookupEnv("BASE_URL"); ok && envBaseURL != "" {
		config.BaseURL = envBaseURL
	}

	if envFileStorage, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok && envFileStorage != "" {
		config.FileStoragePath = envFileStorage
	}

	return config
}
