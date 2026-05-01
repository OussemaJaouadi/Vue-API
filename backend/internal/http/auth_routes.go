package http

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"vue-api/backend/internal/auth"
	"vue-api/backend/internal/events"
)

type AuthRouteDeps struct {
	Users               auth.UserRepository
	Passwords           auth.PasswordHasher
	Tokens              auth.TokenManager
	Events              *events.Broker
	RefreshCookieName   string
	RefreshCookieTTL    time.Duration
	RefreshCookieSecure bool
}

type authRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	AccessToken string `json:"accessToken"`
	UserID      string `json:"userId"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	GlobalRole  string `json:"globalRole"`
}

func RegisterAuthRoutes(router *echo.Echo, deps AuthRouteDeps) {
	if deps.RefreshCookieName == "" {
		deps.RefreshCookieName = "refresh_token"
	}

	router.POST("/auth/register", func(c echo.Context) error {
		var req authRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		email := auth.NormalizeEmail(req.Email)
		username := auth.NormalizeUsername(req.Username)
		if email == "" || username == "" || len(req.Password) < 12 {
			return echo.NewHTTPError(http.StatusBadRequest, "Email, username, and password are required")
		}

		existingUsers, err := deps.Users.CountUsers(c.Request().Context())
		if err != nil {
			return err
		}

		globalRole := auth.GlobalRoleUser
		if existingUsers == 0 {
			globalRole = auth.GlobalRoleManager
		}

		hash, err := deps.Passwords.Hash(req.Password)
		if err != nil {
			return err
		}

		user, err := deps.Users.CreateUser(c.Request().Context(), auth.CreateUserParams{
			Email:        email,
			Username:     username,
			PasswordHash: hash,
			GlobalRole:   globalRole,
		})
		if errors.Is(err, auth.ErrEmailAlreadyInUse) {
			return echo.NewHTTPError(http.StatusConflict, "Email already in use")
		}
		if errors.Is(err, auth.ErrUsernameAlreadyInUse) {
			return echo.NewHTTPError(http.StatusConflict, "Username already in use")
		}
		if err != nil {
			return err
		}

		if deps.Events != nil && user.GlobalRole != auth.GlobalRoleManager {
			deps.Events.PublishToManagers(events.Event{
				Type: "user.registered",
				Data: map[string]string{
					"userId":   user.ID,
					"email":    user.Email,
					"username": user.Username,
				},
			})
		}

		return writeAuthResponse(c, deps, user, http.StatusCreated)
	})

	router.POST("/auth/login", func(c echo.Context) error {
		var req authRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		user, err := findLoginUser(c, deps.Users, req.Login)
		if errors.Is(err, auth.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid login or password")
		}
		if err != nil {
			return err
		}
		if !user.Active || !deps.Passwords.Verify(req.Password, user.PasswordHash) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid login or password")
		}

		return writeAuthResponse(c, deps, user, http.StatusOK)
	})

	router.POST("/auth/refresh", func(c echo.Context) error {
		cookie, err := c.Cookie(deps.RefreshCookieName)
		if err != nil || cookie.Value == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Refresh token is required")
		}

		claims, err := deps.Tokens.ValidateRefreshToken(cookie.Value)
		if errors.Is(err, auth.ErrExpiredToken) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Refresh token expired")
		}
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
		}

		user, err := deps.Users.FindUserByID(c.Request().Context(), claims.UserID)
		if errors.Is(err, auth.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
		}
		if err != nil {
			return err
		}
		if !user.Active || user.TokenVersion != claims.TokenVersion {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
		}

		accessToken, err := deps.Tokens.IssueAccessToken(user.ID, user.TokenVersion)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, authResponse{
			AccessToken: accessToken,
			UserID:      user.ID,
			Email:       user.Email,
			Username:    user.Username,
			GlobalRole:  user.GlobalRole,
		})
	})

	router.POST("/auth/logout", func(c echo.Context) error {
		clearRefreshCookie(c, deps)
		return c.NoContent(http.StatusNoContent)
	})

	router.GET("/auth/me", func(c echo.Context) error {
		user, err := authenticateRequest(c, deps)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"userId":     user.ID,
			"email":      user.Email,
			"username":   user.Username,
			"globalRole": user.GlobalRole,
		})
	})
}

func findLoginUser(c echo.Context, users auth.UserRepository, login string) (auth.User, error) {
	ctx := c.Request().Context()
	if !strings.Contains(login, "@") {
		if user, err := users.FindUserByUsername(ctx, login); err == nil {
			return user, nil
		} else if !errors.Is(err, auth.ErrUserNotFound) {
			return auth.User{}, err
		}

		return users.FindUserByEmail(ctx, login)
	}

	if user, err := users.FindUserByEmail(ctx, login); err == nil {
		return user, nil
	} else if !errors.Is(err, auth.ErrUserNotFound) {
		return auth.User{}, err
	}

	return users.FindUserByUsername(ctx, login)
}

func writeAuthResponse(c echo.Context, deps AuthRouteDeps, user auth.User, status int) error {
	accessToken, err := deps.Tokens.IssueAccessToken(user.ID, user.TokenVersion)
	if err != nil {
		return err
	}

	refreshToken, err := deps.Tokens.IssueRefreshToken(user.ID, user.TokenVersion)
	if err != nil {
		return err
	}

	setRefreshCookie(c, deps, refreshToken)

	return c.JSON(status, authResponse{
		AccessToken: accessToken,
		UserID:      user.ID,
		Email:       user.Email,
		Username:    user.Username,
		GlobalRole:  user.GlobalRole,
	})
}

func authenticateRequest(c echo.Context, deps AuthRouteDeps) (auth.User, error) {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	token, found := strings.CutPrefix(authHeader, "Bearer ")
	if !found || token == "" {
		return auth.User{}, echo.NewHTTPError(http.StatusUnauthorized, "Access token is required")
	}

	claims, err := deps.Tokens.ValidateAccessToken(token)
	if errors.Is(err, auth.ErrExpiredToken) {
		return auth.User{}, echo.NewHTTPError(http.StatusUnauthorized, "Access token expired")
	}
	if err != nil {
		return auth.User{}, echo.NewHTTPError(http.StatusUnauthorized, "Invalid access token")
	}

	user, err := deps.Users.FindUserByID(c.Request().Context(), claims.UserID)
	if errors.Is(err, auth.ErrUserNotFound) {
		return auth.User{}, echo.NewHTTPError(http.StatusUnauthorized, "Invalid access token")
	}
	if err != nil {
		return auth.User{}, err
	}
	if !user.Active || user.TokenVersion != claims.TokenVersion {
		return auth.User{}, echo.NewHTTPError(http.StatusUnauthorized, "Invalid access token")
	}

	return user, nil
}

func setRefreshCookie(c echo.Context, deps AuthRouteDeps, value string) {
	c.SetCookie(&http.Cookie{
		Name:     deps.RefreshCookieName,
		Value:    value,
		Path:     "/",
		MaxAge:   int(deps.RefreshCookieTTL.Seconds()),
		Expires:  time.Now().UTC().Add(deps.RefreshCookieTTL),
		HttpOnly: true,
		Secure:   deps.RefreshCookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearRefreshCookie(c echo.Context, deps AuthRouteDeps) {
	c.SetCookie(&http.Cookie{
		Name:     deps.RefreshCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0).UTC(),
		HttpOnly: true,
		Secure:   deps.RefreshCookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}
