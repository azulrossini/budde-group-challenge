package config

import (
	"fmt"
	"os"
)

type envVar string

const (
	envDatabaseURL envVar = "DATABASE_URL"
	envPort        envVar = "PORT"
)

const DefaultPort = "8080"

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() (Config, error) {
	databaseURL := os.Getenv(string(envDatabaseURL))
	if databaseURL == "" {
		return Config{}, fmt.Errorf("%s is required", envDatabaseURL)
	}

	port := os.Getenv(string(envPort))
	if port == "" {
		port = DefaultPort
	}

	return Config{DatabaseURL: databaseURL, Port: port}, nil
}
