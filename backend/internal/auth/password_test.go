package auth_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/auth"
)

func TestPasswordHasherVerifiesCorrectPassword(t *testing.T) {
	hasher := auth.NewPasswordHasher(auth.PasswordHashParams{
		MemoryKB:    64 * 1024,
		Iterations:  1,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	})

	hash, err := hasher.Hash("correct horse battery staple")
	require.NoError(t, err)
	require.NotContains(t, hash, "correct horse battery staple")
	require.True(t, hasher.Verify("correct horse battery staple", hash))
	require.False(t, hasher.Verify("wrong password", hash))
}
