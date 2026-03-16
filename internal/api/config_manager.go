package api

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	DB PostgresConfig
}

type PostgresConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DB: PostgresConfig{
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PWD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Database: os.Getenv("POSTGRES_DB"),
		},
	}

	return cfg, nil
}
