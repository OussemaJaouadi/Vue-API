package app

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"vue-api/backend/internal/auth"
	"vue-api/backend/internal/config"
	apihttp "vue-api/backend/internal/http"
)

type Server struct {
	addr   string
	router *echo.Echo
}

func NewServer(ctx context.Context, cfg config.Config, logger *slog.Logger) (*Server, error) {
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

	users := auth.NewMemoryUserRepository()
	passwords := auth.DefaultPasswordHasher()
	tokens := auth.NewTokenManager(auth.TokenConfig{
		AccessSecret:  cfg.Auth.JWTAccessSecret,
		RefreshSecret: cfg.Auth.JWTRefreshSecret,
		AccessTTL:     cfg.Auth.JWTAccessTTL,
		RefreshTTL:    cfg.Auth.JWTRefreshTTL,
		Issuer:        cfg.App.Name,
	})
	if err := auth.BootstrapManager(ctx, users, passwords, auth.BootstrapConfig{
		Enabled:  cfg.Auth.BootstrapManagerEnabled,
		Email:    cfg.Auth.BootstrapManagerEmail,
		Username: cfg.Auth.BootstrapManagerUsername,
		Password: cfg.Auth.BootstrapManagerPassword,
	}); err != nil {
		return nil, err
	}

	eventDeps := apihttp.NewEventDeps(users, tokens)

	apihttp.RegisterRoutes(router, cfg)
	apihttp.RegisterAuthRoutes(router, apihttp.AuthRouteDeps{
		Users:               users,
		Passwords:           passwords,
		Tokens:              tokens,
		Events:              eventDeps.Broker,
		RefreshCookieName:   cfg.Auth.RefreshCookieName,
		RefreshCookieTTL:    cfg.Auth.JWTRefreshTTL,
		RefreshCookieSecure: cfg.Auth.RefreshCookieSecure,
	})
	apihttp.RegisterEventRoutes(router, eventDeps)

	return &Server{
		addr:   cfg.HTTP.Addr,
		router: router,
	}, nil
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
