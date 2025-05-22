package config

import (
	"os"
)

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type AppConfig struct {
	Port     string
	Host     string
	Env      string
	LogLevel string
}

type Config struct {
	MySQL   MySQLConfig
	App     AppConfig
	AuthURL string
}

func NewConfig() *Config {
	return &Config{
		MySQL: MySQLConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "posts_db"),
		},
		App: AppConfig{
			Port:     getEnv("SERVER_PORT", "8081"),
			Host:     getEnv("SERVER_HOST", "localhost"),
			Env:      getEnv("APP_ENV", "development"),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		AuthURL: getEnv("AUTH_URL", "http://localhost:8080"),
	}
}

func LoadConfig() *Config {
	return NewConfig()
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
