package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseDSN    string
	Listen         string
	NameServiceURL string
}

func New() Config {
	_ = godotenv.Load()

	return Config{
		DatabaseDSN:    getEnv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		Listen:         getEnv("LISTEN", ":8080"),
		NameServiceURL: getEnv("NAME_SERVICE_URL", ""),
	}
}

func (c Config) IsNameServiceEnabled() bool {
	return c.NameServiceURL != ""
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
