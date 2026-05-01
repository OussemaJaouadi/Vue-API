package auth_test

import (
	"testing"

	"vue-api/backend/internal/auth"
	"vue-api/backend/internal/auth/authtest"
)

func TestMemoryUserRepositoryContract(t *testing.T) {
	authtest.RunUserRepositoryContract(t, func(t *testing.T) auth.UserRepository {
		return auth.NewMemoryUserRepository()
	})
}
