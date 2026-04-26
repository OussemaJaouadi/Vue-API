package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DatabaseConfig
	Auth     AuthConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Env  string
	Name string
}

type HTTPConfig struct {
	Addr        string
	PublicURL   string
	FrontendURL string
}

type DatabaseConfig struct {
	Driver    string
	URL       string
	AuthToken string
}

type AuthConfig struct {
	SessionSecret  string
	PasswordPepper string
}

type CORSConfig struct {
	AllowedOrigins []string
}

func Load() (Config, error) {
	cfg := Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Name: getEnv("APP_NAME", "Vue API Workbench"),
		},
		HTTP: HTTPConfig{
			Addr:        getEnv("BACKEND_ADDR", ":8080"),
			PublicURL:   getEnv("BACKEND_PUBLIC_URL", "http://localhost:8080"),
			FrontendURL: getEnv("FRONTEND_ORIGIN", "http://localhost:3000"),
		},
		Database: DatabaseConfig{
			Driver:    getEnv("DATABASE_DRIVER", "sqlite"),
			URL:       getEnv("DATABASE_URL", "file:../my.local.db"),
			AuthToken: os.Getenv("DATABASE_AUTH_TOKEN"),
		},
		Auth: AuthConfig{
			SessionSecret:  getEnv("SESSION_SECRET", "dev-only-change-me"),
			PasswordPepper: os.Getenv("PASSWORD_PEPPER"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getCSVEnv("CORS_ALLOWED_ORIGINS", []string{
				"http://localhost:3000",
				"http://127.0.0.1:3000",
			}),
		},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg Config) Validate() error {
	if cfg.Database.Driver == "" {
		return errors.New("DATABASE_DRIVER is required")
	}
	if cfg.Database.URL == "" {
		return errors.New("DATABASE_URL is required")
	}
	if cfg.App.Env == "production" && cfg.Auth.SessionSecret == "dev-only-change-me" {
		return errors.New("SESSION_SECRET must be set to a production value")
	}
	if cfg.Database.Driver != "sqlite" && cfg.Database.Driver != "libsql" {
		return fmt.Errorf("unsupported DATABASE_DRIVER %q", cfg.Database.Driver)
	}

	return nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getCSVEnv(key string, fallback []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			values = append(values, item)
		}
	}
	if len(values) == 0 {
		return fallback
	}

	return values
}
