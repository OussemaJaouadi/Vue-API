package config

import "os"

type Config struct {
	Addr string
	Env  string
}

func Load() Config {
	return Config{
		Addr: getEnv("BACKEND_ADDR", ":8080"),
		Env:  getEnv("APP_ENV", "development"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
