package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"vue-api/backend/internal/auth"
	"vue-api/backend/internal/config"
	"vue-api/backend/internal/execution"
	apihttp "vue-api/backend/internal/http"
	gormstorage "vue-api/backend/internal/storage/gorm"
	"vue-api/backend/internal/workspace"
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

	db, err := gormstorage.Open(cfg.Database)
	if err != nil {
		return nil, err
	}
	if err := gormstorage.VerifySchema(db); err != nil {
		return nil, err
	}
	logger.Info("database ready", "driver", cfg.Database.Driver)

	users := gormstorage.NewUserRepository(db)
	passwords := auth.DefaultPasswordHasher()
	tokens := auth.NewTokenManager(auth.TokenConfig{
		AccessSecret:  cfg.Auth.JWTAccessSecret,
		RefreshSecret: cfg.Auth.JWTRefreshSecret,
		AccessTTL:     cfg.Auth.JWTAccessTTL,
		RefreshTTL:    cfg.Auth.JWTRefreshTTL,
		Issuer:        cfg.App.Name,
	})

	bootstrapResult, err := auth.BootstrapManager(ctx, users, passwords, auth.BootstrapConfig{
		Enabled:  cfg.Auth.BootstrapManagerEnabled,
		Email:    cfg.Auth.BootstrapManagerEmail,
		Username: cfg.Auth.BootstrapManagerUsername,
		Password: cfg.Auth.BootstrapManagerPassword,
	})
	if err != nil {
		return nil, err
	}
	logBootstrapResult(logger, bootstrapResult)

	collections := gormstorage.NewCollectionRepository(db)
	environments := gormstorage.NewEnvironmentRepository(db)
	workspaceRepo := gormstorage.NewWorkspaceRepository(db)
	membershipRepo := gormstorage.NewMembershipRepository(db)
	grantRepo := gormstorage.NewGrantRepository(db)

	// Auto-create workspace for bootstrap manager
	if bootstrapResult.Status == auth.BootstrapManagerSeeded {
		ws, err := workspaceRepo.Create(ctx, workspace.CreateWorkspaceParams{
			Name:            bootstrapResult.Username + "'s Workspace",
			CreatedByUserID: bootstrapResult.UserID,
		})
		if err != nil {
			return nil, fmt.Errorf("auto-create workspace: %w", err)
		}
		_, err = membershipRepo.Create(ctx, workspace.CreateMembershipParams{
			WorkspaceID:     ws.ID,
			UserID:          bootstrapResult.UserID,
			Role:            "admin",
			CreatedByUserID: bootstrapResult.UserID,
		})
		if err != nil {
			return nil, fmt.Errorf("auto-add admin membership: %w", err)
		}
		logger.Info("workspace auto-created for bootstrap manager",
			"workspaceId", ws.ID,
			"workspaceName", ws.Name,
		)
	}

	eventDeps := apihttp.NewEventDeps(users, tokens)

	executionService := execution.NewService()
	wsManager := execution.NewWSManager(eventDeps.Broker)

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
	apihttp.RegisterUserRoutes(router, apihttp.UserRouteDeps{
		Users:  users,
		Tokens: tokens,
	})
	apihttp.RegisterCollectionRoutes(router, apihttp.CollectionRouteDeps{
		Collections: collections,
		Memberships: membershipRepo,
		Users:       users,
		Tokens:      tokens,
	})
	apihttp.RegisterEnvironmentRoutes(router, apihttp.EnvironmentRouteDeps{
		Environments: environments,
		Memberships:  membershipRepo,
		Users:        users,
		Tokens:       tokens,
	})
	apihttp.RegisterExecutionRoutes(router, apihttp.ExecutionRouteDeps{
		Execution: executionService,
		WS:        wsManager,
		Users:     users,
		Tokens:    tokens,
	})
	apihttp.RegisterWorkspaceRoutes(router, apihttp.WorkspaceRouteDeps{
		Workspaces:  workspaceRepo,
		Memberships: membershipRepo,
		Grants:      grantRepo,
		Users:       users,
		Tokens:      tokens,
	})

	return &Server{
		addr:   cfg.HTTP.Addr,
		router: router,
	}, nil
}

func logBootstrapResult(logger *slog.Logger, result auth.BootstrapResult) {
	switch result.Status {
	case auth.BootstrapManagerSeeded:
		logger.Info("bootstrap manager seeded",
			"userId", result.UserID,
			"email", result.Email,
			"username", result.Username,
		)
	case auth.BootstrapManagerSkippedExistingUsers:
		logger.Info("bootstrap manager skipped", "reason", "users already exist")
	case auth.BootstrapManagerDisabled:
		logger.Info("bootstrap manager disabled")
	default:
		logger.Info("bootstrap manager status unknown", "status", result.Status)
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
				"uri", sanitizeRequestURI(values.URI),
				"status", values.Status,
				"error", values.Error,
			)
			return nil
		},
	})
}

func sanitizeRequestURI(uri string) string {
	parsed, err := url.ParseRequestURI(uri)
	if err != nil {
		return uri
	}

	query := parsed.Query()
	for _, key := range []string{"access_token", "refresh_token", "token", "ticket"} {
		if query.Has(key) {
			query.Set(key, "[redacted]")
		}
	}
	parsed.RawQuery = strings.ReplaceAll(query.Encode(), "%5Bredacted%5D", "[redacted]")

	return parsed.RequestURI()
}
