package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	StatsCacheTTL int
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// #nosec G101 -- local development fallback only; production must set DATABASE_URL.
		databaseURL = "postgres://postgres:postgres@localhost:5432/gotaskflow?sslmode=disable" // pragma: allowlist secret
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisDB := intEnv("REDIS_DB", 0)
	statsCacheTTL := intEnv("STATS_CACHE_TTL_SECONDS", 60)

	return Config{
		Port:          port,
		DatabaseURL:   databaseURL,
		RedisAddr:     redisAddr,
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
		StatsCacheTTL: statsCacheTTL,
	}
}

func intEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
