package authtest

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/auth"
)

type UserRepositoryFactory func(t *testing.T) auth.UserRepository

func RunUserRepositoryContract(t *testing.T, factory UserRepositoryFactory) {
	t.Helper()

	t.Run("creates and finds users", func(t *testing.T) {
		repo := factory(t)

		user, err := repo.CreateUser(context.Background(), auth.CreateUserParams{
			Email:        " Owner@Example.COM ",
			Username:     " Owner ",
			PasswordHash: "hash",
			GlobalRole:   auth.GlobalRoleManager,
		})
		require.NoError(t, err)

		uuidV7Pattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
		require.Regexp(t, uuidV7Pattern, user.ID)
		require.Equal(t, "owner@example.com", user.Email)
		require.Equal(t, "owner", user.Username)
		require.Equal(t, "hash", user.PasswordHash)
		require.Equal(t, auth.GlobalRoleManager, user.GlobalRole)
		require.Equal(t, 1, user.TokenVersion)
		require.True(t, user.Active)
		require.False(t, user.CreatedAt.IsZero())
		require.False(t, user.UpdatedAt.IsZero())

		count, err := repo.CountUsers(context.Background())
		require.NoError(t, err)
		require.Equal(t, 1, count)

		byEmail, err := repo.FindUserByEmail(context.Background(), "OWNER@example.com")
		require.NoError(t, err)
		require.Equal(t, user.ID, byEmail.ID)

		byID, err := repo.FindUserByID(context.Background(), user.ID)
		require.NoError(t, err)
		require.Equal(t, user.Email, byID.Email)
	})

	t.Run("rejects duplicate email and username", func(t *testing.T) {
		repo := factory(t)
		_, err := repo.CreateUser(context.Background(), auth.CreateUserParams{
			Email:        "one@example.com",
			Username:     "owner",
			PasswordHash: "hash",
			GlobalRole:   auth.GlobalRoleUser,
		})
		require.NoError(t, err)

		_, err = repo.CreateUser(context.Background(), auth.CreateUserParams{
			Email:        " ONE@example.com ",
			Username:     "other",
			PasswordHash: "hash",
			GlobalRole:   auth.GlobalRoleUser,
		})
		require.ErrorIs(t, err, auth.ErrEmailAlreadyInUse)

		_, err = repo.CreateUser(context.Background(), auth.CreateUserParams{
			Email:        "two@example.com",
			Username:     " Owner ",
			PasswordHash: "hash",
			GlobalRole:   auth.GlobalRoleUser,
		})
		require.ErrorIs(t, err, auth.ErrUsernameAlreadyInUse)
	})

	t.Run("updates username", func(t *testing.T) {
		repo := factory(t)
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
	})

	t.Run("returns not found errors", func(t *testing.T) {
		repo := factory(t)

		_, err := repo.FindUserByID(context.Background(), "missing")
		require.ErrorIs(t, err, auth.ErrUserNotFound)

		_, err = repo.FindUserByEmail(context.Background(), "missing@example.com")
		require.ErrorIs(t, err, auth.ErrUserNotFound)

		_, err = repo.UpdateUsername(context.Background(), "missing", "new-name")
		require.ErrorIs(t, err, auth.ErrUserNotFound)
	})
}
