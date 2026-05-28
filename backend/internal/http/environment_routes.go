package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"vue-api/backend/internal/auth"
	"vue-api/backend/internal/environment"
	"vue-api/backend/internal/workspace"
)

type EnvironmentRouteDeps struct {
	Environments environment.Repository
	Memberships  workspace.MembershipRepository
	Users        auth.UserRepository
	Tokens       auth.TokenManager
}

func RegisterEnvironmentRoutes(router *echo.Echo, deps EnvironmentRouteDeps) {
	g := router.Group("/v1/environments")
	g.Use(authMiddleware(deps.Users, deps.Tokens))

	g.GET("", func(c echo.Context) error {
		workspaceID := c.QueryParam("workspaceId")
		if workspaceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "workspaceId query parameter is required")
		}
		if !hasWorkspaceRole(c, deps.Memberships, workspaceID, auth.ScopedRoleTester) {
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized for this workspace")
		}

		envs, err := deps.Environments.ListEnvironments(c.Request().Context(), workspaceID)
		if err != nil {
			return err
		}

		result := make([]map[string]any, 0, len(envs))
		for _, env := range envs {
			vars, err := deps.Environments.ListVariables(c.Request().Context(), env.ID)
			if err != nil {
				return err
			}

			variableList := make([]map[string]any, 0, len(vars))
			for _, v := range vars {
				value := v.Value
				if v.Secret {
					value = "••••••••••••••••"
				}

				variableList = append(variableList, map[string]any{
					"id":     v.ID,
					"key":    v.Key,
					"value":  value,
					"secret": v.Secret,
				})
			}

			result = append(result, map[string]any{
				"id":         env.ID,
				"name":       env.Name,
				"visibility": env.Visibility,
				"variables":  variableList,
			})
		}

		return c.JSON(http.StatusOK, result)
	})

	g.POST("", func(c echo.Context) error {
		var req struct {
			WorkspaceID string `json:"workspaceId"`
			Name        string `json:"name"`
			Visibility  string `json:"visibility"`
		}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}
		if req.Name == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "name is required")
		}
		if req.WorkspaceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "workspaceId is required")
		}
		if !hasWorkspacePermission(c, deps.Memberships, req.WorkspaceID, auth.PermissionManageEnvironment) {
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized to manage environments")
		}

		env, err := deps.Environments.CreateEnvironment(c.Request().Context(), environment.CreateEnvironmentParams{
			WorkspaceID: req.WorkspaceID,
			Name:        req.Name,
			Visibility:  req.Visibility,
		})
		if errors.Is(err, environment.ErrEnvironmentNameTaken) {
			return echo.NewHTTPError(http.StatusConflict, "Environment name already exists")
		}
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"id":         env.ID,
			"name":       env.Name,
			"visibility": env.Visibility,
		})
	})

	g.PUT("/:id", func(c echo.Context) error {
		var req struct {
			WorkspaceID string  `json:"workspaceId"`
			Name        *string `json:"name"`
			Visibility  *string `json:"visibility"`
		}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}
		if req.WorkspaceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "workspaceId is required")
		}
		if !hasWorkspacePermission(c, deps.Memberships, req.WorkspaceID, auth.PermissionManageEnvironment) {
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized to manage environments")
		}

		env, err := deps.Environments.UpdateEnvironment(c.Request().Context(), c.Param("id"), environment.UpdateEnvironmentParams{
			Name:       req.Name,
			Visibility: req.Visibility,
		})
		if errors.Is(err, environment.ErrEnvironmentNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Environment not found")
		}
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{
			"id":         env.ID,
			"name":       env.Name,
			"visibility": env.Visibility,
		})
	})

	g.DELETE("/:id", func(c echo.Context) error {
		workspaceID := c.QueryParam("workspaceId")
		if workspaceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "workspaceId query parameter is required")
		}
		if !hasWorkspacePermission(c, deps.Memberships, workspaceID, auth.PermissionManageEnvironment) {
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized to manage environments")
		}

		if err := deps.Environments.DeleteEnvironment(c.Request().Context(), c.Param("id")); err != nil {
			if errors.Is(err, environment.ErrEnvironmentNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, "Environment not found")
			}
			return err
		}

		return c.NoContent(http.StatusNoContent)
	})

	g.POST("/:envId/variables", func(c echo.Context) error {
		var req struct {
			WorkspaceID string `json:"workspaceId"`
			Key         string `json:"key"`
			Value       string `json:"value"`
			Secret      bool   `json:"secret"`
		}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}
		if req.WorkspaceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "workspaceId is required")
		}
		if !hasWorkspacePermission(c, deps.Memberships, req.WorkspaceID, auth.PermissionManageEnvironment) {
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized to manage environments")
		}
		if req.Key == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "key is required")
		}

		v, err := deps.Environments.CreateVariable(c.Request().Context(), environment.CreateVariableParams{
			EnvironmentID: c.Param("envId"),
			Key:           req.Key,
			Value:         req.Value,
			Secret:        req.Secret,
		})
		if errors.Is(err, environment.ErrVariableKeyTaken) {
			return echo.NewHTTPError(http.StatusConflict, "Variable key already exists in this environment")
		}
		if errors.Is(err, environment.ErrEnvironmentNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Environment not found")
		}
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"id":     v.ID,
			"key":    v.Key,
			"value":  v.Value,
			"secret": v.Secret,
		})
	})

	g.PUT("/:envId/variables/:id", func(c echo.Context) error {
		var req struct {
			WorkspaceID string  `json:"workspaceId"`
			Key         *string `json:"key"`
			Value       *string `json:"value"`
			Secret      *bool   `json:"secret"`
		}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}
		if req.WorkspaceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "workspaceId is required")
		}
		if !hasWorkspacePermission(c, deps.Memberships, req.WorkspaceID, auth.PermissionManageEnvironment) {
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized to manage environments")
		}

		v, err := deps.Environments.UpdateVariable(c.Request().Context(), c.Param("id"), environment.UpdateVariableParams{
			Key:    req.Key,
			Value:  req.Value,
			Secret: req.Secret,
		})
		if errors.Is(err, environment.ErrVariableNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Variable not found")
		}
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{
			"id":     v.ID,
			"key":    v.Key,
			"value":  v.Value,
			"secret": v.Secret,
		})
	})

	g.DELETE("/:envId/variables/:id", func(c echo.Context) error {
		workspaceID := c.QueryParam("workspaceId")
		if workspaceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "workspaceId query parameter is required")
		}
		if !hasWorkspacePermission(c, deps.Memberships, workspaceID, auth.PermissionManageEnvironment) {
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized to manage environments")
		}

		if err := deps.Environments.DeleteVariable(c.Request().Context(), c.Param("id")); err != nil {
			if errors.Is(err, environment.ErrVariableNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, "Variable not found")
			}
			return err
		}

		return c.NoContent(http.StatusNoContent)
	})
}
