package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the application.
// Values are read from environment variables with sensible defaults.
type Config struct {
	Host string
	Port string
	Env  string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	TestMode bool
}

// Load reads configuration from environment variables.
func Load() *Config {
	return &Config{
		Host: getEnv("HOST", "0.0.0.0"),
		Port: getEnv("PORT", "3000"),
		Env:  getEnv("APP_ENV", "development"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "autohaus"),
		DBUser:     getEnv("DB_USER", "autohaus"),
		DBPassword: getEnv("DB_PASSWORD", "p"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		TestMode: getEnv("TEST_MODE", "false") == "true",
	}
}

// Addr returns the full host:port address for the server to listen on.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// DSN returns a PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s search_path=autohaus",
		c.DBHost, c.DBPort, c.DBName, c.DBUser, c.DBPassword, c.DBSSLMode,
	)
}

// IsDevelopment returns true when running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
