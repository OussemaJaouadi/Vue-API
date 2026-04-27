package auth_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/auth"
)

func TestCreateUserAssignsUUIDV7ID(t *testing.T) {
	repo := auth.NewMemoryUserRepository()

	user, err := repo.CreateUser(context.Background(), auth.CreateUserParams{
		Email:        "user@example.com",
		Username:     "user",
		PasswordHash: "hash",
		GlobalRole:   auth.GlobalRoleUser,
	})
	require.NoError(t, err)

	uuidV7Pattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	require.Regexp(t, uuidV7Pattern, user.ID)
}

func TestCreateUserRequiresUniqueUsername(t *testing.T) {
	repo := auth.NewMemoryUserRepository()
	_, err := repo.CreateUser(context.Background(), auth.CreateUserParams{
		Email:        "one@example.com",
		Username:     "Owner",
		PasswordHash: "hash",
		GlobalRole:   auth.GlobalRoleUser,
	})
	require.NoError(t, err)

	_, err = repo.CreateUser(context.Background(), auth.CreateUserParams{
		Email:        "two@example.com",
		Username:     " owner ",
		PasswordHash: "hash",
		GlobalRole:   auth.GlobalRoleUser,
	})
	require.ErrorIs(t, err, auth.ErrUsernameAlreadyInUse)
}

func TestUpdateUsernameRequiresUniqueUsername(t *testing.T) {
	repo := auth.NewMemoryUserRepository()
	first, err := repo.CreateUser(context.Background(), auth.CreateUserParams{
		Email:        "one@example.com",
		Username:     "owner",
		PasswordHash: "hash",
		GlobalRole:   auth.GlobalRoleUser,
	})
	require.NoError(t, err)
	second, err := repo.CreateUser(context.Background(), auth.CreateUserParams{
		Email:        "two@example.com",
		Username:     "developer",
		PasswordHash: "hash",
		GlobalRole:   auth.GlobalRoleUser,
	})
	require.NoError(t, err)

	_, err = repo.UpdateUsername(context.Background(), second.ID, "OWNER")
	require.ErrorIs(t, err, auth.ErrUsernameAlreadyInUse)

	updated, err := repo.UpdateUsername(context.Background(), first.ID, " Lead-Dev ")
	require.NoError(t, err)
	require.Equal(t, "lead-dev", updated.Username)
}
