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
}

// Load reads configuration from environment variables.
func Load() *Config {
	return &Config{
		Host: getEnv("HOST", "0.0.0.0"),
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("APP_ENV", "development"),
	}
}

// Addr returns the full host:port address for the server to listen on.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
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
