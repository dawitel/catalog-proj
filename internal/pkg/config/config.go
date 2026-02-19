package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/dawitel/product-catalog-service/internal/pkg/logger"
)

type GlobalConfig struct {
	SpannerProject      string
	SpannerInstance     string
	SpannerDatabase     string
	SpannerEmulatorHost string
	GrpcPort            string
	// MigrationsPath is the path to the initial schema SQL file (legacy single file).
	MigrationsPath string
	// MigrationsDir is the directory containing migration SQL files (run in lexicographic order).
	MigrationsDir string
}

// DatabasePath returns the Spanner database path: projects/{project}/instances/{instance}/databases/{database}.
func (c *GlobalConfig) DatabasePath() string {
	return "projects/" + c.SpannerProject + "/instances/" + c.SpannerInstance + "/databases/" + c.SpannerDatabase
}

func LoadFromEnv() GlobalConfig {
	if err := godotenv.Load(); err != nil {
		logger.Debug("no .env file loaded", "err", err)
	}
	return GlobalConfig{
		SpannerProject:      getEnv("SPANNER_PROJECT", "test-project"),
		SpannerInstance:     getEnv("SPANNER_INSTANCE", "test-instance"),
		SpannerDatabase:     getEnv("SPANNER_DATABASE", "product-catalog"),
		SpannerEmulatorHost: os.Getenv("SPANNER_EMULATOR_HOST"),
		GrpcPort:            grpcPortFromEnv(),
		MigrationsPath:      getEnv("MIGRATIONS_PATH", "migrations/001_initial_schema.sql"),
		MigrationsDir:       getEnv("MIGRATIONS_DIR", "migrations"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func grpcPortFromEnv() string {
	port := getEnv("GRPC_PORT", "8080")
	if _, err := strconv.Atoi(port); err != nil {
		return "8080"
	}
	return port
}
