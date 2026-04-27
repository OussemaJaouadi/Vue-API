package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/auth"
)

func TestTokenManagerValidatesAccessAndRefreshWithDifferentSecrets(t *testing.T) {
	manager := auth.NewTokenManager(auth.TokenConfig{
		AccessSecret:  "access-secret-at-least-32-bytes-long",
		RefreshSecret: "refresh-secret-at-least-32-bytes-long",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
		Issuer:        "vue-api",
	})

	accessToken, err := manager.IssueAccessToken("user_123", 4)
	require.NoError(t, err)

	refreshToken, err := manager.IssueRefreshToken("user_123", 4)
	require.NoError(t, err)

	accessClaims, err := manager.ValidateAccessToken(accessToken)
	require.NoError(t, err)
	require.Equal(t, "user_123", accessClaims.UserID)
	require.Equal(t, 4, accessClaims.TokenVersion)

	refreshClaims, err := manager.ValidateRefreshToken(refreshToken)
	require.NoError(t, err)
	require.Equal(t, "user_123", refreshClaims.UserID)
	require.Equal(t, 4, refreshClaims.TokenVersion)

	_, err = manager.ValidateAccessToken(refreshToken)
	require.ErrorIs(t, err, auth.ErrInvalidToken)

	_, err = manager.ValidateRefreshToken(accessToken)
	require.ErrorIs(t, err, auth.ErrInvalidToken)
}
