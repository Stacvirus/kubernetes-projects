package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	NatsURL      string
	NatsUser     string
	NatsPassword string
	NatsSubject  string
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		Port:         getString("PORT", "8050"),
		NatsURL:      getString("NATS_URL", "nats://localhost:4222"),
		NatsUser:     getString("NATS_USER", "your-user"),
		NatsPassword: getString("NATS_PASSWORD", "your-password"),
		NatsSubject:  getString("NATS_SUBJECT", "updates"),
	}
}

func getString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if raw, ok := os.LookupEnv(key); ok {
		value, err := strconv.Atoi(raw)
		if err != nil {
			log.Printf("⚠️ Invalid value for %s=%s, using fallback %d", key, raw, fallback)
			return fallback
		}
		return value
	}
	return fallback
}
