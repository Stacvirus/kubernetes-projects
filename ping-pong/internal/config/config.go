package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	Path         string
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func Load() *Config {
	_ = godotenv.Load()
	cfg := &Config{
		Port:         getString("PORT", "8080"),
		Path:         getString("LOG_FILE_PATH", "pong.log"),
		Addr:         getString("PG_ADDR", "postgres://stac:password@localhost:5436/pingpong?sslmode=disable"),
		MaxOpenConns: getInt("PG_MAX_OPEN_CONNS", 30),
		MaxIdleConns: getInt("PG_MAX_IDLE_CONNS", 30),
		MaxIdleTime:  getString("PG_MAX_IDLE_TIME", "15m"),
	}
	return cfg
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
