package http

import (
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
				requestList = append(requestList, map[string]any{
					"id":     req.ID,
					"method": req.Method,
					"name":   req.Name,
					"path":   req.Path,
				})
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
			rootRequestList = append(rootRequestList, map[string]any{
				"id":     req.ID,
				"method": req.Method,
				"name":   req.Name,
				"path":   req.Path,
			})
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
		if !hasWorkspaceRole(c, deps.Memberships, req.WorkspaceID, "tester") {
			return echo.NewHTTPError(http.StatusForbidden, "Not a member of this workspace")
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
		if !hasWorkspaceRole(c, deps.Memberships, workspaceID, "tester") {
			return echo.NewHTTPError(http.StatusForbidden, "Not a member of this workspace")
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
			WorkspaceID  string  `json:"workspaceId"`
			CollectionID *string `json:"collectionId"`
			Method       string  `json:"method"`
			Name         string  `json:"name"`
			Path         string  `json:"path"`
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

		method := req.Method
		if method == "" {
			method = "GET"
		}

		request, err := deps.Collections.CreateRequest(c.Request().Context(), collection.CreateRequestParams{
			CollectionID: req.CollectionID,
			WorkspaceID:  req.WorkspaceID,
			Method:       method,
			Name:         req.Name,
			Path:         req.Path,
		})
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"id":     request.ID,
			"method": request.Method,
			"name":   request.Name,
			"path":   request.Path,
		})
	})

	g.PUT("/requests/:id", func(c echo.Context) error {
		var req struct {
			WorkspaceID string  `json:"workspaceId"`
			Method      *string `json:"method"`
			Name        *string `json:"name"`
			Path        *string `json:"path"`
		}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}
		if req.WorkspaceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "workspaceId is required")
		}
		if !hasWorkspaceRole(c, deps.Memberships, req.WorkspaceID, "tester") {
			return echo.NewHTTPError(http.StatusForbidden, "Not a member of this workspace")
		}

		request, err := deps.Collections.UpdateRequest(c.Request().Context(), c.Param("id"), collection.UpdateRequestParams{
			Method: req.Method,
			Name:   req.Name,
			Path:   req.Path,
		})
		if errors.Is(err, collection.ErrRequestNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Request not found")
		}
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{
			"id":     request.ID,
			"method": request.Method,
			"name":   request.Name,
			"path":   request.Path,
		})
	})

	g.DELETE("/requests/:id", func(c echo.Context) error {
		workspaceID := c.QueryParam("workspaceId")
		if workspaceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "workspaceId query parameter is required")
		}
		if !hasWorkspaceRole(c, deps.Memberships, workspaceID, "tester") {
			return echo.NewHTTPError(http.StatusForbidden, "Not a member of this workspace")
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
