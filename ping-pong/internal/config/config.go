package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Path string
}

func Load() *Config {
	_ = godotenv.Load()
	cfg := &Config{
		Port: getEnv("PORT", "8080"),
		Path: getEnv("LOG_FILE_PATH", "pong.log"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
