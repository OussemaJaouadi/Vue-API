package auth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/auth"
)

func TestBootstrapManagerCreatesFirstUserWhenEnabled(t *testing.T) {
	repo := auth.NewMemoryUserRepository()
	hasher := testPasswordHasher()

	result, err := auth.BootstrapManager(context.Background(), repo, hasher, auth.BootstrapConfig{
		Enabled:  true,
		Email:    " Manager@Example.COM ",
		Username: "Manager",
		Password: "strong-password",
	})
	require.NoError(t, err)
	require.Equal(t, auth.BootstrapManagerSeeded, result.Status)
	require.Equal(t, "manager@example.com", result.Email)
	require.Equal(t, "manager", result.Username)
	require.NotEmpty(t, result.UserID)

	count, err := repo.CountUsers(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, count)

	user, err := repo.FindUserByEmail(context.Background(), "manager@example.com")
	require.NoError(t, err)
	require.Equal(t, auth.GlobalRoleManager, user.GlobalRole)
	require.Equal(t, "manager", user.Username)
	require.True(t, hasher.Verify("strong-password", user.PasswordHash))
}

func TestBootstrapManagerDoesNothingWhenDisabled(t *testing.T) {
	repo := auth.NewMemoryUserRepository()

	result, err := auth.BootstrapManager(context.Background(), repo, testPasswordHasher(), auth.BootstrapConfig{
		Enabled:  false,
		Email:    "manager@example.com",
		Username: "manager",
		Password: "strong-password",
	})
	require.NoError(t, err)
	require.Equal(t, auth.BootstrapManagerDisabled, result.Status)

	count, err := repo.CountUsers(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func TestBootstrapManagerDoesNotReseedExistingUsers(t *testing.T) {
	repo := auth.NewMemoryUserRepository()
	hasher := testPasswordHasher()
	hash, err := hasher.Hash("existing-password")
	require.NoError(t, err)
	_, err = repo.CreateUser(context.Background(), auth.CreateUserParams{
		Email:        "existing@example.com",
		Username:     "existing",
		PasswordHash: hash,
		GlobalRole:   auth.GlobalRoleUser,
	})
	require.NoError(t, err)

	result, err := auth.BootstrapManager(context.Background(), repo, hasher, auth.BootstrapConfig{
		Enabled:  true,
		Email:    "manager@example.com",
		Username: "manager",
		Password: "strong-password",
	})
	require.NoError(t, err)
	require.Equal(t, auth.BootstrapManagerSkippedExistingUsers, result.Status)

	count, err := repo.CountUsers(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, count)

	_, err = repo.FindUserByEmail(context.Background(), "manager@example.com")
	require.ErrorIs(t, err, auth.ErrUserNotFound)
}

func TestBootstrapManagerRequiresValidCredentialsWhenEnabled(t *testing.T) {
	repo := auth.NewMemoryUserRepository()

	_, err := auth.BootstrapManager(context.Background(), repo, testPasswordHasher(), auth.BootstrapConfig{
		Enabled:  true,
		Email:    "",
		Username: "manager",
		Password: "strong-password",
	})
	require.ErrorContains(t, err, "BOOTSTRAP_MANAGER_EMAIL")

	_, err = auth.BootstrapManager(context.Background(), repo, testPasswordHasher(), auth.BootstrapConfig{
		Enabled:  true,
		Email:    "manager@example.com",
		Username: "manager",
		Password: "short",
	})
	require.ErrorContains(t, err, "BOOTSTRAP_MANAGER_PASSWORD")

	_, err = auth.BootstrapManager(context.Background(), repo, testPasswordHasher(), auth.BootstrapConfig{
		Enabled:  true,
		Email:    "manager@example.com",
		Username: "",
		Password: "strong-password",
	})
	require.ErrorContains(t, err, "BOOTSTRAP_MANAGER_USERNAME")
}

func testPasswordHasher() auth.PasswordHasher {
	return auth.NewPasswordHasher(auth.PasswordHashParams{
		MemoryKB:    1024,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	})
}
