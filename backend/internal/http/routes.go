package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vue-api/backend/internal/config"
)

func RegisterRoutes(router *echo.Echo, cfg config.Config) {
	router.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":         "ok",
			"env":            cfg.App.Env,
			"databaseDriver": cfg.Database.Driver,
		})
	})
}
