package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"vue-api/backend/internal/auth"
	"vue-api/backend/internal/collection"
	"vue-api/backend/internal/workspace"
)

type CollectionRouteDeps struct {
	Collections collection.Repository
	Memberships workspace.MembershipRepository
	Users       auth.UserRepository
	Tokens      auth.TokenManager
}

func RegisterCollectionRoutes(router *echo.Echo, deps CollectionRouteDeps) {
	g := router.Group("/v1/collections")
	g.Use(authMiddleware(deps.Users, deps.Tokens))

	g.GET("", func(c echo.Context) error {
		workspaceID := c.QueryParam("workspaceId")
		if workspaceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "workspaceId query parameter is required")
		}
		if !hasWorkspacePermission(c, deps.Memberships, workspaceID, auth.PermissionViewCollections) {
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized for this workspace")
		}

		folders, err := deps.Collections.ListFolders(c.Request().Context(), workspaceID)
		if err != nil {
			return err
		}

		result := make([]map[string]any, 0, len(folders))
		for _, folder := range folders {
			reqs, err := deps.Collections.ListRequests(c.Request().Context(), workspaceID, &folder.ID)
			if err != nil {
				return err
			}

			requestList := make([]map[string]any, 0, len(reqs))
			for _, req := range reqs {
				requestList = append(requestList, requestResponse(req))
			}

			result = append(result, map[string]any{
				"id":       folder.ID,
				"name":     folder.Name,
				"icon":     folder.Icon,
				"requests": requestList,
			})
		}

		rootReqs, err := deps.Collections.ListRootRequests(c.Request().Context(), workspaceID)
		if err != nil {
			return err
		}

		rootRequestList := make([]map[string]any, 0, len(rootReqs))
		for _, req := range rootReqs {
			rootRequestList = append(rootRequestList, requestResponse(req))
		}

		return c.JSON(http.StatusOK, map[string]any{
			"collections":  result,
			"rootRequests": rootRequestList,
		})
	})

	g.POST("", func(c echo.Context) error {
		var req struct {
			WorkspaceID string `json:"workspaceId"`
			Name        string `json:"name"`
			Icon        string `json:"icon"`
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
		if !hasWorkspacePermission(c, deps.Memberships, req.WorkspaceID, auth.PermissionManageCollections) {
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized to manage collections")
		}

		folder, err := deps.Collections.CreateFolder(c.Request().Context(), collection.CreateFolderParams{
			WorkspaceID: req.WorkspaceID,
			Name:        req.Name,
			Icon:        req.Icon,
		})
		if errors.Is(err, collection.ErrFolderNameTaken) {
			return echo.NewHTTPError(http.StatusConflict, "Collection name already exists")
		}
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"id":   folder.ID,
			"name": folder.Name,
			"icon": folder.Icon,
		})
	})

	g.PUT("/:id", func(c echo.Context) error {
		var req struct {
			WorkspaceID string  `json:"workspaceId"`
			Name        *string `json:"name"`
			Icon        *string `json:"icon"`
		}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}
		if req.WorkspaceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "workspaceId is required")
		}
		if !hasWorkspacePermission(c, deps.Memberships, req.WorkspaceID, auth.PermissionManageCollections) {
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized to manage collections")
		}

		folder, err := deps.Collections.UpdateFolder(c.Request().Context(), c.Param("id"), collection.UpdateFolderParams{
			Name: req.Name,
			Icon: req.Icon,
		})
		if errors.Is(err, collection.ErrFolderNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Collection not found")
		}
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{
			"id":   folder.ID,
			"name": folder.Name,
			"icon": folder.Icon,
		})
	})

	g.DELETE("/:id", func(c echo.Context) error {
		workspaceID := c.QueryParam("workspaceId")
		if workspaceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "workspaceId query parameter is required")
		}
		if !hasWorkspacePermission(c, deps.Memberships, workspaceID, auth.PermissionManageCollections) {
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized to manage collections")
		}

		if err := deps.Collections.DeleteFolder(c.Request().Context(), c.Param("id")); err != nil {
			if errors.Is(err, collection.ErrFolderNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, "Collection not found")
			}
			return err
		}

		return c.NoContent(http.StatusNoContent)
	})

	g.POST("/requests", func(c echo.Context) error {
		var req struct {
			WorkspaceID  string          `json:"workspaceId"`
			CollectionID *string         `json:"collectionId"`
			Method       string          `json:"method"`
			Name         string          `json:"name"`
			Path         string          `json:"path"`
			QueryParams  json.RawMessage `json:"queryParams"`
			Headers      json.RawMessage `json:"headers"`
			Body         string          `json:"body"`
			BodyLanguage string          `json:"bodyLanguage"`
			AuthConfig   json.RawMessage `json:"authConfig"`
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
		if !hasWorkspacePermission(c, deps.Memberships, req.WorkspaceID, auth.PermissionManageCollections) {
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized to manage collections")
		}

		method := req.Method
		if method == "" {
			method = "GET"
		}

		request, err := deps.Collections.CreateRequest(c.Request().Context(), collection.CreateRequestParams{
			CollectionID:    req.CollectionID,
			WorkspaceID:     req.WorkspaceID,
			Method:          method,
			Name:            req.Name,
			Path:            req.Path,
			QueryParamsJSON: rawJSONOrDefault(req.QueryParams, "[]"),
			HeadersJSON:     rawJSONOrDefault(req.Headers, "[]"),
			Body:            req.Body,
			BodyLanguage:    valueOrDefault(req.BodyLanguage, "json"),
			AuthConfigJSON:  rawJSONOrDefault(req.AuthConfig, "{}"),
		})
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, requestResponse(request))
	})

	g.PUT("/requests/:id", func(c echo.Context) error {
		var req struct {
			WorkspaceID  string           `json:"workspaceId"`
			Method       *string          `json:"method"`
			Name         *string          `json:"name"`
			Path         *string          `json:"path"`
			QueryParams  *json.RawMessage `json:"queryParams"`
			Headers      *json.RawMessage `json:"headers"`
			Body         *string          `json:"body"`
			BodyLanguage *string          `json:"bodyLanguage"`
			AuthConfig   *json.RawMessage `json:"authConfig"`
		}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}
		if req.WorkspaceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "workspaceId is required")
		}
		if !hasWorkspacePermission(c, deps.Memberships, req.WorkspaceID, auth.PermissionManageCollections) {
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized to manage collections")
		}

		params := collection.UpdateRequestParams{
			Method: req.Method,
			Name:   req.Name,
			Path:   req.Path,
		}
		if req.QueryParams != nil {
			value := rawJSONOrDefault(*req.QueryParams, "[]")
			params.QueryParamsJSON = &value
		}
		if req.Headers != nil {
			value := rawJSONOrDefault(*req.Headers, "[]")
			params.HeadersJSON = &value
		}
		if req.Body != nil {
			params.Body = req.Body
		}
		if req.BodyLanguage != nil {
			value := valueOrDefault(*req.BodyLanguage, "json")
			params.BodyLanguage = &value
		}
		if req.AuthConfig != nil {
			value := rawJSONOrDefault(*req.AuthConfig, "{}")
			params.AuthConfigJSON = &value
		}

		request, err := deps.Collections.UpdateRequest(c.Request().Context(), c.Param("id"), params)
		if errors.Is(err, collection.ErrRequestNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Request not found")
		}
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, requestResponse(request))
	})

	g.DELETE("/requests/:id", func(c echo.Context) error {
		workspaceID := c.QueryParam("workspaceId")
		if workspaceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "workspaceId query parameter is required")
		}
		if !hasWorkspacePermission(c, deps.Memberships, workspaceID, auth.PermissionManageCollections) {
			return echo.NewHTTPError(http.StatusForbidden, "Not authorized to manage collections")
		}

		if err := deps.Collections.DeleteRequest(c.Request().Context(), c.Param("id")); err != nil {
			if errors.Is(err, collection.ErrRequestNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, "Request not found")
			}
			return err
		}

		return c.NoContent(http.StatusNoContent)
	})
}

func requestResponse(req collection.Request) map[string]any {
	return map[string]any{
		"id":           req.ID,
		"method":       req.Method,
		"name":         req.Name,
		"path":         req.Path,
		"queryParams":  decodeJSONOrDefault(req.QueryParamsJSON, []any{}),
		"headers":      decodeJSONOrDefault(req.HeadersJSON, []any{}),
		"body":         req.Body,
		"bodyLanguage": valueOrDefault(req.BodyLanguage, "json"),
		"authConfig":   decodeJSONOrDefault(req.AuthConfigJSON, map[string]any{}),
	}
}

func rawJSONOrDefault(raw json.RawMessage, fallback string) string {
	if len(raw) == 0 || string(raw) == "null" {
		return fallback
	}

	return string(raw)
}

func decodeJSONOrDefault(raw string, fallback any) any {
	if raw == "" {
		return fallback
	}

	var decoded any
	if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
		return fallback
	}

	return decoded
}

func valueOrDefault(value string, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
