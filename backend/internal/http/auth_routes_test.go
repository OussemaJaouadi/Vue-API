package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/auth"
	"vue-api/backend/internal/events"
	apihttp "vue-api/backend/internal/http"
)

func TestAuthRoutesRegisterLoginRefreshAndCurrentUser(t *testing.T) {
	router := echo.New()
	router.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler

	repo := auth.NewMemoryUserRepository()
	hasher := auth.NewPasswordHasher(auth.PasswordHashParams{
		MemoryKB:    8 * 1024,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	})
	tokens := auth.NewTokenManager(auth.TokenConfig{
		AccessSecret:  "access-secret-at-least-32-bytes-long",
		RefreshSecret: "refresh-secret-at-least-32-bytes-long",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
		Issuer:        "vue-api-test",
	})

	apihttp.RegisterAuthRoutes(router, apihttp.AuthRouteDeps{
		Users:             repo,
		Passwords:         hasher,
		Tokens:            tokens,
		Events:            events.NewBroker(),
		RefreshCookieName: "refresh_token",
		RefreshCookieTTL:  24 * time.Hour,
	})

	registerResp := performJSON(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "owner@example.com",
		"username": "Owner",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusCreated, registerResp.Code)

	var registerBody map[string]any
	require.NoError(t, json.Unmarshal(registerResp.Body.Bytes(), &registerBody))
	require.NotEmpty(t, registerBody["accessToken"])
	require.Equal(t, "owner", registerBody["username"])
	require.Equal(t, "manager", registerBody["globalRole"])

	loginResp := performJSON(router, http.MethodPost, "/auth/login", map[string]string{
		"login":    "owner@example.com",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusOK, loginResp.Code)

	var loginBody map[string]any
	require.NoError(t, json.Unmarshal(loginResp.Body.Bytes(), &loginBody))
	accessToken := loginBody["accessToken"].(string)
	require.NotEmpty(t, accessToken)

	usernameLoginResp := performJSON(router, http.MethodPost, "/auth/login", map[string]string{
		"login":    "owner",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusOK, usernameLoginResp.Code)

	refreshCookie := findCookie(loginResp.Result().Cookies(), "refresh_token")
	require.NotNil(t, refreshCookie)
	require.True(t, refreshCookie.HttpOnly)

	currentResp := performJSON(router, http.MethodGet, "/auth/me", nil, accessToken)
	require.Equal(t, http.StatusOK, currentResp.Code)

	var currentBody map[string]any
	require.NoError(t, json.Unmarshal(currentResp.Body.Bytes(), &currentBody))
	require.Equal(t, "owner@example.com", currentBody["email"])
	require.Equal(t, "owner", currentBody["username"])
	require.Equal(t, "manager", currentBody["globalRole"])

	refreshReq := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	refreshReq.AddCookie(refreshCookie)
	refreshRec := httptest.NewRecorder()
	router.ServeHTTP(refreshRec, refreshReq)
	require.Equal(t, http.StatusOK, refreshRec.Code)

	var refreshBody map[string]any
	require.NoError(t, json.Unmarshal(refreshRec.Body.Bytes(), &refreshBody))
	require.NotEmpty(t, refreshBody["accessToken"])
}

func TestAuthRoutesPublishRegisteredUserEventsToManagers(t *testing.T) {
	router := echo.New()
	router.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler

	repo := auth.NewMemoryUserRepository()
	manager, err := repo.CreateUser(t.Context(), auth.CreateUserParams{
		Email:        "manager@example.com",
		Username:     "manager",
		PasswordHash: "hash",
		GlobalRole:   auth.GlobalRoleManager,
	})
	require.NoError(t, err)

	broker := events.NewBroker()
	subscription := broker.Subscribe(t.Context(), events.Subscriber{
		UserID:     manager.ID,
		GlobalRole: manager.GlobalRole,
	})
	defer subscription.Close()

	apihttp.RegisterAuthRoutes(router, apihttp.AuthRouteDeps{
		Users: repo,
		Passwords: auth.NewPasswordHasher(auth.PasswordHashParams{
			MemoryKB:    1024,
			Iterations:  1,
			Parallelism: 1,
			SaltLength:  16,
			KeyLength:   32,
		}),
		Tokens:            testTokenManager(),
		Events:            broker,
		RefreshCookieName: "refresh_token",
		RefreshCookieTTL:  24 * time.Hour,
	})

	registerResp := performJSON(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "pending@example.com",
		"username": "pending",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusCreated, registerResp.Code)

	select {
	case event := <-subscription.Events():
		require.Equal(t, "user.registered", event.Type)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for user.registered event")
	}
}

func TestAuthRoutesRejectDuplicateUsername(t *testing.T) {
	router := echo.New()
	router.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler

	repo := auth.NewMemoryUserRepository()
	apihttp.RegisterAuthRoutes(router, apihttp.AuthRouteDeps{
		Users: repo,
		Passwords: auth.NewPasswordHasher(auth.PasswordHashParams{
			MemoryKB:    1024,
			Iterations:  1,
			Parallelism: 1,
			SaltLength:  16,
			KeyLength:   32,
		}),
		Tokens: auth.NewTokenManager(auth.TokenConfig{
			AccessSecret:  "access-secret-at-least-32-bytes-long",
			RefreshSecret: "refresh-secret-at-least-32-bytes-long",
			AccessTTL:     15 * time.Minute,
			RefreshTTL:    24 * time.Hour,
			Issuer:        "vue-api-test",
		}),
		RefreshCookieName: "refresh_token",
		RefreshCookieTTL:  24 * time.Hour,
	})

	firstResp := performJSON(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "one@example.com",
		"username": "Owner",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusCreated, firstResp.Code)

	secondResp := performJSON(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "two@example.com",
		"username": "owner",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusConflict, secondResp.Code)
}

func performJSON(router *echo.Echo, method string, path string, payload any, bearer string) *httptest.ResponseRecorder {
	var body bytes.Buffer
	if payload != nil {
		_ = json.NewEncoder(&body).Encode(payload)
	}

	req := httptest.NewRequest(method, path, &body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if bearer != "" {
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+bearer)
	}
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	return rec
}

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}

	return nil
}
