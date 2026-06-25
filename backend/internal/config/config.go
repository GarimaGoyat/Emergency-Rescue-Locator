package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port              string
	GinMode           string
	DatabaseURL       string
	JWTSecret         string
	JWTExpiryHours    int
	CORSOrigins       []string
	RateLimitRequests int64
	RateLimitDuration time.Duration
}

func Load() *Config {
	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	rateLimitRequests, _ := strconv.ParseInt(getEnv("RATE_LIMIT_REQUESTS", "100"), 10, 64)
	rateLimitDuration, _ := time.ParseDuration(getEnv("RATE_LIMIT_DURATION", "1m"))

	corsOrigins := strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173"), ",")
	for i, origin := range corsOrigins {
		corsOrigins[i] = strings.TrimSpace(origin)
	}

	return &Config{
		Port:              getEnv("PORT", "8080"),
		GinMode:           getEnv("GIN_MODE", "debug"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		JWTSecret:         getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		JWTExpiryHours:    jwtExpiry,
		CORSOrigins:       corsOrigins,
		RateLimitRequests: rateLimitRequests,
		RateLimitDuration: rateLimitDuration,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
