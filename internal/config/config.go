package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Jwt      JWTConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	URL      string
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       string
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8000"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			URL:      getEnv("DB_URL", ""),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "aiki"),
			Password: getEnv("DB_PASSWORD", "aiki_password"),
			DBName:   getEnv("DB_NAME", "aiki_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Port:     getEnv("REDIS_PORT", "localhost"),
			Host:     getEnv("REDIS_Host", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Jwt: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "secret"),
			AccessExpiry:  parseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"), 15*time.Minute),
			RefreshExpiry: parseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"), 7*24*time.Hour),
		},
	}

	return cfg, nil
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Port, c.Host)
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func (c *DatabaseConfig) ConnectionUrl() string {
	return c.URL
}

func parseDuration(value string, duration time.Duration) time.Duration {
	d, err := time.ParseDuration(value)
	if err != nil {
		return duration
	}
	return d
}
