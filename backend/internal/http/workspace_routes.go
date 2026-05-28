package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"vue-api/backend/internal/auth"
	"vue-api/backend/internal/workspace"
)

func hasWorkspaceRole(c echo.Context, memberships workspace.MembershipRepository, workspaceID string, minRole string) bool {
	user := c.Get("user").(auth.User)
	m, err := memberships.FindByUserAndWorkspace(c.Request().Context(), user.ID, workspaceID)
	if err != nil {
		return false
	}
	return workspaceRoleAtLeast(m.Role, minRole)
}

func workspaceRoleAtLeast(role, minRole string) bool {
	order := map[string]int{"tester": 1, "developer": 2, "admin": 3}
	return order[role] >= order[minRole]
}

func hasWorkspacePermission(c echo.Context, memberships workspace.MembershipRepository, workspaceID string, permission string) bool {
	user := c.Get("user").(auth.User)
	m, err := memberships.FindByUserAndWorkspace(c.Request().Context(), user.ID, workspaceID)
	if err != nil {
		return false
	}

	return auth.Can(user.GlobalRole, m.Role, permission)
}

type WorkspaceRouteDeps struct {
	Workspaces  workspace.WorkspaceRepository
	Memberships workspace.MembershipRepository
	Grants      workspace.GrantRepository
	Users       auth.UserRepository
	Tokens      auth.TokenManager
}

type createWorkspaceRequest struct {
	Name string `json:"name"`
}

type updateWorkspaceRequest struct {
	Name string `json:"name"`
}

type inviteMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type updateMemberRoleRequest struct {
	Role string `json:"role"`
}

type setGrantsRequest struct {
	Grants []workspace.GrantInput `json:"grants"`
}

func RegisterWorkspaceRoutes(router *echo.Echo, deps WorkspaceRouteDeps) {
	g := router.Group("/v1/workspaces")
	g.Use(authMiddleware(deps.Users, deps.Tokens))

	g.GET("", func(c echo.Context) error {
		user := c.Get("user").(auth.User)
		workspaces, err := deps.Workspaces.ListByUser(c.Request().Context(), user.ID)
		if err != nil {
			return err
		}

		result := make([]map[string]any, 0, len(workspaces))
		for _, ws := range workspaces {
			count, _ := deps.Memberships.CountByWorkspace(c.Request().Context(), ws.ID)
			m, err := deps.Memberships.FindByUserAndWorkspace(c.Request().Context(), user.ID, ws.ID)
			role := ""
			if err == nil {
				role = m.Role
			}
			result = append(result, map[string]any{
				"id":          ws.ID,
				"name":        ws.Name,
				"role":        role,
				"memberCount": count,
			})
		}
		return c.JSON(http.StatusOK, result)
	})

	g.POST("", func(c echo.Context) error {
		user := c.Get("user").(auth.User)
		if user.GlobalRole != auth.GlobalRoleManager {
			return echo.NewHTTPError(http.StatusForbidden, "Only managers can create workspaces")
		}

		var req createWorkspaceRequest
		if err := c.Bind(&req); err != nil || req.Name == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "name is required")
		}

		ws, err := deps.Workspaces.Create(c.Request().Context(), workspace.CreateWorkspaceParams{
			Name:            req.Name,
			CreatedByUserID: user.ID,
		})
		if errors.Is(err, workspace.ErrWorkspaceNameTaken) {
			return echo.NewHTTPError(http.StatusConflict, "Workspace name already taken")
		}
		if err != nil {
			return err
		}

		deps.Memberships.Create(c.Request().Context(), workspace.CreateMembershipParams{
			WorkspaceID:     ws.ID,
			UserID:          user.ID,
			Role:            "admin",
			CreatedByUserID: user.ID,
		})

		return c.JSON(http.StatusCreated, map[string]any{
			"id":   ws.ID,
			"name": ws.Name,
		})
	})

	g.GET("/:id", func(c echo.Context) error {
		user := c.Get("user").(auth.User)
		ws, err := deps.Workspaces.FindByID(c.Request().Context(), c.Param("id"))
		if errors.Is(err, workspace.ErrWorkspaceNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Workspace not found")
		}
		if err != nil {
			return err
		}

		_, err = deps.Memberships.FindByUserAndWorkspace(c.Request().Context(), user.ID, ws.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "Not a member of this workspace")
		}

		count, _ := deps.Memberships.CountByWorkspace(c.Request().Context(), ws.ID)
		return c.JSON(http.StatusOK, map[string]any{
			"id":          ws.ID,
			"name":        ws.Name,
			"memberCount": count,
		})
	})

	g.PUT("/:id", func(c echo.Context) error {
		wsID := c.Param("id")
		if _, err := deps.Workspaces.FindByID(c.Request().Context(), wsID); errors.Is(err, workspace.ErrWorkspaceNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Workspace not found")
		}

		if !hasWorkspaceRole(c, deps.Memberships, wsID, "admin") {
			return echo.NewHTTPError(http.StatusForbidden, "Only workspace admins can update workspaces")
		}

		var req updateWorkspaceRequest
		if err := c.Bind(&req); err != nil || req.Name == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "name is required")
		}

		ws, err := deps.Workspaces.Update(c.Request().Context(), wsID, workspace.UpdateWorkspaceParams{Name: &req.Name})
		if errors.Is(err, workspace.ErrWorkspaceNameTaken) {
			return echo.NewHTTPError(http.StatusConflict, "Workspace name already taken")
		}
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{"id": ws.ID, "name": ws.Name})
	})

	g.DELETE("/:id", func(c echo.Context) error {
		wsID := c.Param("id")
		if _, err := deps.Workspaces.FindByID(c.Request().Context(), wsID); errors.Is(err, workspace.ErrWorkspaceNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Workspace not found")
		} else if err != nil {
			return err
		}

		if !hasWorkspacePermission(c, deps.Memberships, wsID, auth.PermissionDeleteProject) {
			return echo.NewHTTPError(http.StatusForbidden, "Only workspace admins can delete workspaces")
		}

		if err := deps.Workspaces.Delete(c.Request().Context(), wsID); errors.Is(err, workspace.ErrWorkspaceNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Workspace not found")
		} else if err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	})

	m := g.Group("/:id/members")

	m.GET("", func(c echo.Context) error {
		workspaceID := c.Param("id")
		user := c.Get("user").(auth.User)
		if _, err := deps.Memberships.FindByUserAndWorkspace(c.Request().Context(), user.ID, workspaceID); err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "Not a member of this workspace")
		}

		members, err := deps.Memberships.ListByWorkspace(c.Request().Context(), workspaceID)
		if err != nil {
			return err
		}

		result := make([]map[string]any, 0, len(members))
		for _, m := range members {
			u, err := deps.Users.FindUserByID(c.Request().Context(), m.UserID)
			username := ""
			email := ""
			if err == nil {
				username = u.Username
				email = u.Email
			}
			result = append(result, map[string]any{
				"id":       m.ID,
				"userId":   m.UserID,
				"username": username,
				"email":    email,
				"role":     m.Role,
			})
		}
		return c.JSON(http.StatusOK, result)
	})

	m.POST("", func(c echo.Context) error {
		workspaceID := c.Param("id")
		if !hasWorkspaceRole(c, deps.Memberships, workspaceID, "admin") {
			return echo.NewHTTPError(http.StatusForbidden, "Only workspace admins can invite members")
		}

		var req inviteMemberRequest
		if err := c.Bind(&req); err != nil || req.Email == "" || !auth.ValidScopedRole(req.Role) {
			return echo.NewHTTPError(http.StatusBadRequest, "email and valid role (admin/developer/tester) are required")
		}

		user := c.Get("user").(auth.User)
		targetUser, err := deps.Users.FindUserByEmail(c.Request().Context(), req.Email)
		if errors.Is(err, auth.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		if err != nil {
			return err
		}

		m, err := deps.Memberships.Create(c.Request().Context(), workspace.CreateMembershipParams{
			WorkspaceID:     workspaceID,
			UserID:          targetUser.ID,
			Role:            req.Role,
			CreatedByUserID: user.ID,
		})
		if errors.Is(err, workspace.ErrAlreadyMember) {
			return echo.NewHTTPError(http.StatusConflict, "User is already a member of this workspace")
		}
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"id":       m.ID,
			"userId":   targetUser.ID,
			"username": targetUser.Username,
			"email":    targetUser.Email,
			"role":     m.Role,
		})
	})

	m.PUT("/:userId", func(c echo.Context) error {
		workspaceID := c.Param("id")
		if !hasWorkspaceRole(c, deps.Memberships, workspaceID, "admin") {
			return echo.NewHTTPError(http.StatusForbidden, "Only workspace admins can update roles")
		}

		var req updateMemberRoleRequest
		if err := c.Bind(&req); err != nil || !auth.ValidScopedRole(req.Role) {
			return echo.NewHTTPError(http.StatusBadRequest, "valid role (admin/developer/tester) is required")
		}

		members, err := deps.Memberships.ListByWorkspace(c.Request().Context(), workspaceID)
		if err != nil {
			return err
		}
		var targetMembership *workspace.Membership
		for _, m := range members {
			if m.UserID == c.Param("userId") {
				targetMembership = &m
				break
			}
		}
		if targetMembership == nil {
			return echo.NewHTTPError(http.StatusNotFound, "Membership not found")
		}

		updated, err := deps.Memberships.UpdateRole(c.Request().Context(), targetMembership.ID, workspace.UpdateMembershipParams{Role: &req.Role})
		if err != nil {
			return err
		}

		u, _ := deps.Users.FindUserByID(c.Request().Context(), updated.UserID)
		return c.JSON(http.StatusOK, map[string]any{
			"id":       updated.ID,
			"userId":   updated.UserID,
			"username": u.Username,
			"email":    u.Email,
			"role":     updated.Role,
		})
	})

	m.DELETE("/:userId", func(c echo.Context) error {
		workspaceID := c.Param("id")
		if !hasWorkspaceRole(c, deps.Memberships, workspaceID, "admin") {
			return echo.NewHTTPError(http.StatusForbidden, "Only workspace admins can remove members")
		}

		user := c.Get("user").(auth.User)
		members, err := deps.Memberships.ListByWorkspace(c.Request().Context(), workspaceID)
		if err != nil {
			return err
		}

		var targetMembership *workspace.Membership
		for _, m := range members {
			if m.UserID == c.Param("userId") {
				targetMembership = &m
				break
			}
		}
		if targetMembership == nil {
			return echo.NewHTTPError(http.StatusNotFound, "Membership not found")
		}

		if targetMembership.UserID == user.ID {
			return echo.NewHTTPError(http.StatusBadRequest, "Cannot remove yourself")
		}

		if err := deps.Memberships.Delete(c.Request().Context(), targetMembership.ID); err != nil {
			return err
		}
		return c.NoContent(http.StatusNoContent)
	})

	gr := m.Group("/:userId/grants")

	gr.GET("", func(c echo.Context) error {
		workspaceID := c.Param("id")
		userID := c.Param("userId")

		user := c.Get("user").(auth.User)
		if _, err := deps.Memberships.FindByUserAndWorkspace(c.Request().Context(), user.ID, workspaceID); err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "Not a member of this workspace")
		}

		grants, err := deps.Grants.ListByUserAndWorkspace(c.Request().Context(), userID, workspaceID)
		if err != nil {
			return err
		}

		collections := make(map[string]string)
		environments := make(map[string]string)
		secrets := make(map[string]string)
		for _, g := range grants {
			switch g.ResourceType {
			case "collection":
				collections[g.ResourceID] = g.AccessLevel
			case "environment":
				environments[g.ResourceID] = g.AccessLevel
			case "secret":
				secrets[g.ResourceID] = g.AccessLevel
			}
		}

		return c.JSON(http.StatusOK, map[string]any{
			"collections":  collections,
			"environments": environments,
			"secrets":      secrets,
		})
	})

	gr.PUT("", func(c echo.Context) error {
		workspaceID := c.Param("id")
		if !hasWorkspaceRole(c, deps.Memberships, workspaceID, "admin") {
			return echo.NewHTTPError(http.StatusForbidden, "Only workspace admins can set grants")
		}

		var req setGrantsRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		userID := c.Param("userId")
		_, err := deps.Grants.Set(c.Request().Context(), userID, workspaceID, req.Grants)
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	})
}
