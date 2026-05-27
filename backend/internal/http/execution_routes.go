package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"vue-api/backend/internal/auth"
	"vue-api/backend/internal/execution"
)

type ExecutionRouteDeps struct {
	Execution execution.Service
	WS        *execution.WSManager
	Users     auth.UserRepository
	Tokens    auth.TokenManager
}

func RegisterExecutionRoutes(router *echo.Echo, deps ExecutionRouteDeps) {
	g := router.Group("/execute")
	g.Use(authMiddleware(deps.Users, deps.Tokens))

	g.POST("", func(c echo.Context) error {
		var req execution.Request
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		if req.URL == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "URL is required")
		}

		timeout := req.Timeout
		if timeout == 0 {
			timeout = 30 * time.Second
		}

		ctx, cancel := context.WithTimeout(c.Request().Context(), timeout)
		defer cancel()

		resp, err := deps.Execution.Execute(ctx, req)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, resp)
	})

	g.POST("/ws", func(c echo.Context) error {
		user := c.Get("user").(auth.User)
		var req execution.Request
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		id, err := deps.WS.Connect(c.Request().Context(), user.ID, req)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadGateway, err.Error())
		}

		return c.JSON(http.StatusCreated, map[string]string{"id": id})
	})

	g.POST("/:id/ws/send", func(c echo.Context) error {
		id := c.Param("id")
		var req struct {
			Payload string `json:"payload"`
		}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		if err := deps.WS.Send(id, req.Payload); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	})

	g.DELETE("/:id", func(c echo.Context) error {
		id := c.Param("id")
		if err := deps.WS.Close(id); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	})
}

func authMiddleware(users auth.UserRepository, tokens auth.TokenManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			token, found := strings.CutPrefix(authHeader, "Bearer ")
			if !found || token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Access token is required")
			}

			claims, err := tokens.ValidateAccessToken(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid access token")
			}

			user, err := users.FindUserByID(c.Request().Context(), claims.UserID)
			if err != nil || !user.Active || user.TokenVersion != claims.TokenVersion {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid access token")
			}

			c.Set("user", user)
			return next(c)
		}
	}
}
