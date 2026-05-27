package http_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/auth"
	apihttp "vue-api/backend/internal/http"
)

func TestUserRoutes_ListUsers_NoAuth(t *testing.T) {
	router := echo.New()
	router.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler

	users := auth.NewMemoryUserRepository()
	tokens := testTokenManager()
	apihttp.RegisterUserRoutes(router, apihttp.UserRouteDeps{
		Users:  users,
		Tokens: tokens,
	})

	resp := performJSON(router, http.MethodGet, "/v1/users", nil, "")
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestUserRoutes_ListUsers_ReturnsOnlyRegistered(t *testing.T) {
	router, token := authTestDeps(t)

	resp := performJSON(router, http.MethodGet, "/v1/users", nil, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body []any
	json.Unmarshal(resp.Body.Bytes(), &body)
	require.Len(t, body, 1)
}

func TestUserRoutes_ListUsers_ReturnsRegisteredUsers(t *testing.T) {
	router, token := authTestDeps(t)

	secondResp := performJSON(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "alice@example.com",
		"username": "alice",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusCreated, secondResp.Code)

	resp := performJSON(router, http.MethodGet, "/v1/users", nil, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body []map[string]any
	json.Unmarshal(resp.Body.Bytes(), &body)
	require.Len(t, body, 2)

	found := false
	for _, u := range body {
		if u["email"] == "alice@example.com" {
			found = true
			assert.Equal(t, "alice", u["username"])
			assert.Equal(t, "user", u["role"])
			assert.NotEmpty(t, u["id"])
			assert.NotEmpty(t, u["active"])
			break
		}
	}
	assert.True(t, found, "expected alice in user list")
}

func authTestDeps(t *testing.T) (*echo.Echo, string) {
	t.Helper()

	users := auth.NewMemoryUserRepository()
	tokens := testTokenManager()
	hasher := auth.NewPasswordHasher(auth.PasswordHashParams{
		MemoryKB:    1024,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	})

	router := echo.New()
	router.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler

	apihttp.RegisterAuthRoutes(router, apihttp.AuthRouteDeps{
		Users:             users,
		Passwords:         hasher,
		Tokens:            tokens,
		Events:            nil,
		RefreshCookieName: "refresh_token",
		RefreshCookieTTL:  24 * 60 * 60,
	})
	apihttp.RegisterUserRoutes(router, apihttp.UserRouteDeps{
		Users:  users,
		Tokens: tokens,
	})

	registerResp := performJSON(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "owner@example.com",
		"username": "owner",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusCreated, registerResp.Code)

	var registerBody map[string]any
	require.NoError(t, json.Unmarshal(registerResp.Body.Bytes(), &registerBody))
	accessToken := registerBody["accessToken"].(string)
	require.NotEmpty(t, accessToken)

	return router, accessToken
}
