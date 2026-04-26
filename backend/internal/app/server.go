package app

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"vue-api/backend/internal/config"
	apihttp "vue-api/backend/internal/http"
)

type Server struct {
	addr   string
	router *echo.Echo
}

func NewServer(cfg config.Config, logger *slog.Logger) *Server {
	router := echo.New()
	router.HideBanner = true
	router.HidePort = true

	router.Use(middleware.Recover())
	router.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler
	router.Use(requestLogger(logger))
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.CORS.AllowedOrigins,
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.PATCH,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderAuthorization,
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
		AllowCredentials: true,
	}))

	apihttp.RegisterRoutes(router, cfg)

	return &Server{
		addr:   cfg.HTTP.Addr,
		router: router,
	}
}

func (s *Server) Start() error {
	return s.router.Start(s.addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.router.Shutdown(ctx)
}

func requestLogger(logger *slog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			logger.Info("request completed",
				"method", c.Request().Method,
				"uri", values.URI,
				"status", values.Status,
				"error", values.Error,
			)
			return nil
		},
	})
}
