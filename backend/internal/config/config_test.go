package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/config"
)

func TestLoadAuthConfigFromEnvironment(t *testing.T) {
	t.Setenv("CONFIG_LOAD_DOTENV", "false")
	t.Setenv("JWT_ACCESS_SECRET", "access-secret-at-least-32-bytes-long")
	t.Setenv("JWT_REFRESH_SECRET", "refresh-secret-at-least-32-bytes-long")
	t.Setenv("JWT_ACCESS_TTL", "10m")
	t.Setenv("JWT_REFRESH_TTL", "168h")
	t.Setenv("REFRESH_COOKIE_NAME", "refresh")
	t.Setenv("REFRESH_COOKIE_SECURE", "true")

	cfg, err := config.Load()
	require.NoError(t, err)

	require.Equal(t, "access-secret-at-least-32-bytes-long", cfg.Auth.JWTAccessSecret)
	require.Equal(t, "refresh-secret-at-least-32-bytes-long", cfg.Auth.JWTRefreshSecret)
	require.Equal(t, 10*time.Minute, cfg.Auth.JWTAccessTTL)
	require.Equal(t, 168*time.Hour, cfg.Auth.JWTRefreshTTL)
	require.Equal(t, "refresh", cfg.Auth.RefreshCookieName)
	require.True(t, cfg.Auth.RefreshCookieSecure)
}

func TestLoadBootstrapManagerConfigFromEnvironment(t *testing.T) {
	t.Setenv("CONFIG_LOAD_DOTENV", "false")
	t.Setenv("BOOTSTRAP_MANAGER_ENABLED", "true")
	t.Setenv("BOOTSTRAP_MANAGER_EMAIL", "manager@example.com")
	t.Setenv("BOOTSTRAP_MANAGER_USERNAME", "manager")
	t.Setenv("BOOTSTRAP_MANAGER_PASSWORD", "strong-password")

	cfg, err := config.Load()
	require.NoError(t, err)

	require.True(t, cfg.Auth.BootstrapManagerEnabled)
	require.Equal(t, "manager@example.com", cfg.Auth.BootstrapManagerEmail)
	require.Equal(t, "manager", cfg.Auth.BootstrapManagerUsername)
	require.Equal(t, "strong-password", cfg.Auth.BootstrapManagerPassword)
}

func TestBootstrapManagerRequiresCredentialsWhenEnabled(t *testing.T) {
	t.Setenv("CONFIG_LOAD_DOTENV", "false")
	t.Setenv("BOOTSTRAP_MANAGER_ENABLED", "true")
	t.Setenv("BOOTSTRAP_MANAGER_EMAIL", "manager@example.com")
	t.Setenv("BOOTSTRAP_MANAGER_USERNAME", "manager")
	t.Setenv("BOOTSTRAP_MANAGER_PASSWORD", "short")

	_, err := config.Load()
	require.ErrorContains(t, err, "BOOTSTRAP_MANAGER_PASSWORD")
}

func TestLoadDotenvFilesWithLocalOverride(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte("APP_NAME=Root App\nBOOTSTRAP_MANAGER_ENABLED=false\n"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "backend"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "backend", ".env"), []byte("BOOTSTRAP_MANAGER_ENABLED=true\nBOOTSTRAP_MANAGER_EMAIL=manager@example.com\nBOOTSTRAP_MANAGER_USERNAME=manager\nBOOTSTRAP_MANAGER_PASSWORD=strong-password\n"), 0o644))

	originalWD, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(filepath.Join(dir, "backend")))
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(originalWD))
	})

	cfg, err := config.Load()
	require.NoError(t, err)

	require.Equal(t, "Root App", cfg.App.Name)
	require.True(t, cfg.Auth.BootstrapManagerEnabled)
	require.Equal(t, "manager@example.com", cfg.Auth.BootstrapManagerEmail)
}

func TestProductionRequiresJWTSecrets(t *testing.T) {
	t.Setenv("CONFIG_LOAD_DOTENV", "false")
	t.Setenv("APP_ENV", "production")

	_, err := config.Load()
	require.ErrorContains(t, err, "JWT_ACCESS_SECRET")
}
