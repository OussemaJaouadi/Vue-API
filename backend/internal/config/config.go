package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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
	PasswordPepper           string
	JWTAccessSecret          string
	JWTRefreshSecret         string
	JWTAccessTTL             time.Duration
	JWTRefreshTTL            time.Duration
	RefreshCookieName        string
	RefreshCookieSecure      bool
	BootstrapManagerEnabled  bool
	BootstrapManagerEmail    string
	BootstrapManagerUsername string
	BootstrapManagerPassword string
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
			PasswordPepper:           os.Getenv("PASSWORD_PEPPER"),
			JWTAccessSecret:          getEnv("JWT_ACCESS_SECRET", "dev-access-secret-at-least-32-bytes"),
			JWTRefreshSecret:         getEnv("JWT_REFRESH_SECRET", "dev-refresh-secret-at-least-32-bytes"),
			JWTAccessTTL:             getDurationEnv("JWT_ACCESS_TTL", 15*time.Minute),
			JWTRefreshTTL:            getDurationEnv("JWT_REFRESH_TTL", 720*time.Hour),
			RefreshCookieName:        getEnv("REFRESH_COOKIE_NAME", "refresh_token"),
			RefreshCookieSecure:      getBoolEnv("REFRESH_COOKIE_SECURE", false),
			BootstrapManagerEnabled:  getBoolEnv("BOOTSTRAP_MANAGER_ENABLED", false),
			BootstrapManagerEmail:    os.Getenv("BOOTSTRAP_MANAGER_EMAIL"),
			BootstrapManagerUsername: os.Getenv("BOOTSTRAP_MANAGER_USERNAME"),
			BootstrapManagerPassword: os.Getenv("BOOTSTRAP_MANAGER_PASSWORD"),
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
	if cfg.App.Env == "production" && cfg.Auth.JWTAccessSecret == "dev-access-secret-at-least-32-bytes" {
		return errors.New("JWT_ACCESS_SECRET must be set to a production value")
	}
	if cfg.App.Env == "production" && cfg.Auth.JWTRefreshSecret == "dev-refresh-secret-at-least-32-bytes" {
		return errors.New("JWT_REFRESH_SECRET must be set to a production value")
	}
	if len(cfg.Auth.JWTAccessSecret) < 32 {
		return errors.New("JWT_ACCESS_SECRET must be at least 32 characters")
	}
	if len(cfg.Auth.JWTRefreshSecret) < 32 {
		return errors.New("JWT_REFRESH_SECRET must be at least 32 characters")
	}
	if cfg.Auth.JWTAccessTTL <= 0 {
		return errors.New("JWT_ACCESS_TTL must be positive")
	}
	if cfg.Auth.JWTRefreshTTL <= 0 {
		return errors.New("JWT_REFRESH_TTL must be positive")
	}
	if cfg.Auth.BootstrapManagerEnabled {
		if strings.TrimSpace(cfg.Auth.BootstrapManagerEmail) == "" {
			return errors.New("BOOTSTRAP_MANAGER_EMAIL is required when BOOTSTRAP_MANAGER_ENABLED is true")
		}
		if strings.TrimSpace(cfg.Auth.BootstrapManagerUsername) == "" {
			return errors.New("BOOTSTRAP_MANAGER_USERNAME is required when BOOTSTRAP_MANAGER_ENABLED is true")
		}
		if len(cfg.Auth.BootstrapManagerPassword) < 12 {
			return errors.New("BOOTSTRAP_MANAGER_PASSWORD must be at least 12 characters")
		}
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

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return duration
}

func getBoolEnv(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}
