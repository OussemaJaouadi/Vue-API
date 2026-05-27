package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vue-api/backend/internal/auth"
)

type UserRouteDeps struct {
	Users  auth.UserRepository
	Tokens auth.TokenManager
}

func RegisterUserRoutes(router *echo.Echo, deps UserRouteDeps) {
	g := router.Group("/v1/users")
	g.Use(authMiddleware(deps.Users, deps.Tokens))

	g.GET("", func(c echo.Context) error {
		users, err := deps.Users.ListUsers(c.Request().Context())
		if err != nil {
			return err
		}

		result := make([]map[string]any, 0, len(users))
		for _, u := range users {
			result = append(result, map[string]any{
				"id":       u.ID,
				"email":    u.Email,
				"username": u.Username,
				"role":     u.GlobalRole,
				"active":   u.Active,
			})
		}

		return c.JSON(http.StatusOK, result)
	})
}
