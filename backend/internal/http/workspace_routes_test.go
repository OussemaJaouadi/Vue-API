package http_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/auth"
	apihttp "vue-api/backend/internal/http"
	"vue-api/backend/internal/workspace"
)

func workspaceTestDeps(t *testing.T) (*echo.Echo, string, string) {
	t.Helper()

	users := auth.NewMemoryUserRepository()
	tokens := testTokenManager()
	hasher := auth.NewPasswordHasher(auth.PasswordHashParams{
		MemoryKB:    1024,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	})

	workspaceRepo := workspace.NewMemoryWorkspaceRepository()
	membershipRepo := workspace.NewMemoryMembershipRepository()
	grantRepo := workspace.NewMemoryGrantRepository()

	router := echo.New()
	router.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler

	apihttp.RegisterAuthRoutes(router, apihttp.AuthRouteDeps{
		Users:             users,
		Passwords:         hasher,
		Tokens:            tokens,
		Events:            nil,
		RefreshCookieName: "refresh_token",
		RefreshCookieTTL:  24 * 60 * 60,
	})
	apihttp.RegisterUserRoutes(router, apihttp.UserRouteDeps{
		Users:  users,
		Tokens: tokens,
	})
	apihttp.RegisterWorkspaceRoutes(router, apihttp.WorkspaceRouteDeps{
		Workspaces:  workspaceRepo,
		Memberships: membershipRepo,
		Grants:      grantRepo,
		Users:       users,
		Tokens:      tokens,
	})

	// Register owner (becomes manager as first user)
	registerResp := performJSON(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "owner@example.com",
		"username": "owner",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusCreated, registerResp.Code)
	var registerBody map[string]any
	require.NoError(t, json.Unmarshal(registerResp.Body.Bytes(), &registerBody))
	ownerToken := registerBody["accessToken"].(string)
	require.NotEmpty(t, ownerToken)

	// Register a regular user for invite tests
	userResp := performJSON(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "alice@example.com",
		"username": "alice",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusCreated, userResp.Code)

	// Create a workspace for testing
	createResp := performJSON(router, http.MethodPost, "/v1/workspaces", map[string]string{
		"name": "Test Workspace",
	}, ownerToken)
	require.Equal(t, http.StatusCreated, createResp.Code)
	var createBody map[string]any
	require.NoError(t, json.Unmarshal(createResp.Body.Bytes(), &createBody))
	workspaceID := createBody["id"].(string)
	require.NotEmpty(t, workspaceID)

	return router, ownerToken, workspaceID
}

// --- Workspace CRUD ---

func TestWorkspaceRoutes_ListWorkspaces_NoAuth(t *testing.T) {
	router := echo.New()
	router.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler

	users := auth.NewMemoryUserRepository()
	tokens := testTokenManager()
	apihttp.RegisterWorkspaceRoutes(router, apihttp.WorkspaceRouteDeps{
		Workspaces:  workspace.NewMemoryWorkspaceRepository(),
		Memberships: workspace.NewMemoryMembershipRepository(),
		Grants:      workspace.NewMemoryGrantRepository(),
		Users:       users,
		Tokens:      tokens,
	})

	resp := performJSON(router, http.MethodGet, "/v1/workspaces", nil, "")
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestWorkspaceRoutes_ListWorkspaces_Empty(t *testing.T) {
	router, token, _ := workspaceTestDeps(t)

	// List workspaces
	resp := performJSON(router, http.MethodGet, "/v1/workspaces", nil, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body []map[string]any
	json.Unmarshal(resp.Body.Bytes(), &body)
	require.Len(t, body, 1)
	assert.Equal(t, "Test Workspace", body[0]["name"])
	assert.Equal(t, "admin", body[0]["role"])
}

func TestWorkspaceRoutes_CreateWorkspace_RequiresName(t *testing.T) {
	router, token, _ := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/workspaces", map[string]string{}, token)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestWorkspaceRoutes_CreateWorkspace_CreatesAndAutoAddsAdmin(t *testing.T) {
	router, token, _ := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/workspaces", map[string]string{
		"name": "My Second WS",
	}, token)
	require.Equal(t, http.StatusCreated, resp.Code)

	var body map[string]any
	json.Unmarshal(resp.Body.Bytes(), &body)
	assert.NotEmpty(t, body["id"])
	assert.Equal(t, "My Second WS", body["name"])

	// Verify it shows up in list
	listResp := performJSON(router, http.MethodGet, "/v1/workspaces", nil, token)
	var list []map[string]any
	json.Unmarshal(listResp.Body.Bytes(), &list)
	require.Len(t, list, 2)
}

func TestWorkspaceRoutes_CreateWorkspace_DuplicateName(t *testing.T) {
	router, token, _ := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/workspaces", map[string]string{
		"name": "Test Workspace",
	}, token)
	assert.Equal(t, http.StatusConflict, resp.Code)
}

func TestWorkspaceRoutes_GetWorkspace_NotFound(t *testing.T) {
	router, token, _ := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodGet, "/v1/workspaces/nonexistent", nil, token)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestWorkspaceRoutes_GetWorkspace_Found(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodGet, "/v1/workspaces/"+wsID, nil, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body map[string]any
	json.Unmarshal(resp.Body.Bytes(), &body)
	assert.Equal(t, "Test Workspace", body["name"])
	assert.Equal(t, float64(1), body["memberCount"])
}

func TestWorkspaceRoutes_UpdateWorkspace_RequiresName(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodPut, "/v1/workspaces/"+wsID, map[string]string{}, token)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestWorkspaceRoutes_UpdateWorkspace_Renames(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodPut, "/v1/workspaces/"+wsID, map[string]string{
		"name": "Renamed",
	}, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body map[string]any
	json.Unmarshal(resp.Body.Bytes(), &body)
	assert.Equal(t, "Renamed", body["name"])
}

func TestWorkspaceRoutes_UpdateWorkspace_NotFound(t *testing.T) {
	router, token, _ := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodPut, "/v1/workspaces/nonexistent", map[string]string{
		"name": "Renamed",
	}, token)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestWorkspaceRoutes_DeleteWorkspace(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodDelete, "/v1/workspaces/"+wsID, nil, token)
	assert.Equal(t, http.StatusNoContent, resp.Code)

	getResp := performJSON(router, http.MethodGet, "/v1/workspaces/"+wsID, nil, token)
	assert.Equal(t, http.StatusNotFound, getResp.Code)

	listResp := performJSON(router, http.MethodGet, "/v1/workspaces", nil, token)
	var list []map[string]any
	require.NoError(t, json.Unmarshal(listResp.Body.Bytes(), &list))
	assert.Empty(t, list)
}

func TestWorkspaceRoutes_DeleteWorkspace_NotFound(t *testing.T) {
	router, token, _ := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodDelete, "/v1/workspaces/nonexistent", nil, token)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestWorkspaceRoutes_DeleteWorkspace_RequiresAdmin(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)
	inviteResp := performJSON(router, http.MethodPost, "/v1/workspaces/"+wsID+"/members", map[string]string{
		"email": "alice@example.com",
		"role":  "developer",
	}, token)
	require.Equal(t, http.StatusCreated, inviteResp.Code)

	loginResp := performJSON(router, http.MethodPost, "/auth/login", map[string]string{
		"login":    "alice@example.com",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusOK, loginResp.Code)
	var loginBody map[string]any
	require.NoError(t, json.Unmarshal(loginResp.Body.Bytes(), &loginBody))
	aliceToken := loginBody["accessToken"].(string)

	resp := performJSON(router, http.MethodDelete, "/v1/workspaces/"+wsID, nil, aliceToken)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

// --- Memberships ---

func TestWorkspaceRoutes_ListMembers_ReturnsMembers(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodGet, "/v1/workspaces/"+wsID+"/members", nil, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body []map[string]any
	json.Unmarshal(resp.Body.Bytes(), &body)
	require.Len(t, body, 1)
	assert.Equal(t, "owner", body[0]["username"])
	assert.Equal(t, "admin", body[0]["role"])
}

func TestWorkspaceRoutes_ListMembers_NonMemberForbidden(t *testing.T) {
	router, _, _ := workspaceTestDeps(t)

	// Register another user with no workspace membership
	otherResp := performJSON(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "bob@example.com",
		"username": "bob",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusCreated, otherResp.Code)
	var otherBody map[string]any
	json.Unmarshal(otherResp.Body.Bytes(), &otherBody)
	otherToken := otherBody["accessToken"].(string)

	resp := performJSON(router, http.MethodGet, "/v1/workspaces/test-ws/members", nil, otherToken)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestWorkspaceRoutes_InviteMember_ByEmail(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/workspaces/"+wsID+"/members", map[string]string{
		"email": "alice@example.com",
		"role":  "developer",
	}, token)
	require.Equal(t, http.StatusCreated, resp.Code)

	var body map[string]any
	json.Unmarshal(resp.Body.Bytes(), &body)
	assert.Equal(t, "alice", body["username"])
	assert.Equal(t, "developer", body["role"])

	// Verify member list has 2
	listResp := performJSON(router, http.MethodGet, "/v1/workspaces/"+wsID+"/members", nil, token)
	var list []map[string]any
	json.Unmarshal(listResp.Body.Bytes(), &list)
	assert.Len(t, list, 2)
}

func TestWorkspaceRoutes_InviteMember_AlreadyMember(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)

	// Invite alice first
	performJSON(router, http.MethodPost, "/v1/workspaces/"+wsID+"/members", map[string]string{
		"email": "alice@example.com",
		"role":  "developer",
	}, token)

	// Try inviting again
	resp := performJSON(router, http.MethodPost, "/v1/workspaces/"+wsID+"/members", map[string]string{
		"email": "alice@example.com",
		"role":  "admin",
	}, token)
	assert.Equal(t, http.StatusConflict, resp.Code)
}

func TestWorkspaceRoutes_InviteMember_UserNotFound(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/workspaces/"+wsID+"/members", map[string]string{
		"email": "nobody@example.com",
		"role":  "developer",
	}, token)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestWorkspaceRoutes_InviteMember_RequiresValidRole(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/workspaces/"+wsID+"/members", map[string]string{
		"email": "alice@example.com",
		"role":  "invalid",
	}, token)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestWorkspaceRoutes_UpdateMemberRole(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)

	// Invite alice
	inviteResp := performJSON(router, http.MethodPost, "/v1/workspaces/"+wsID+"/members", map[string]string{
		"email": "alice@example.com",
		"role":  "developer",
	}, token)
	require.Equal(t, http.StatusCreated, inviteResp.Code)
	var inviteBody map[string]any
	json.Unmarshal(inviteResp.Body.Bytes(), &inviteBody)
	aliceID := inviteBody["userId"].(string)
	require.NotEmpty(t, aliceID)

	// Update role using userId
	resp := performJSON(router, http.MethodPut, "/v1/workspaces/"+wsID+"/members/"+aliceID, map[string]string{
		"role": "admin",
	}, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var updated map[string]any
	json.Unmarshal(resp.Body.Bytes(), &updated)
	assert.Equal(t, "admin", updated["role"])
}

func TestWorkspaceRoutes_KickMember(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)

	// Invite alice
	performJSON(router, http.MethodPost, "/v1/workspaces/"+wsID+"/members", map[string]string{
		"email": "alice@example.com",
		"role":  "developer",
	}, token)

	// Find alice's userId
	membersResp := performJSON(router, http.MethodGet, "/v1/workspaces/"+wsID+"/members", nil, token)
	var members []map[string]any
	json.Unmarshal(membersResp.Body.Bytes(), &members)

	var aliceID string
	for _, m := range members {
		if m["username"] == "alice" {
			aliceID = m["userId"].(string)
			break
		}
	}
	require.NotEmpty(t, aliceID)

	// Kick alice
	resp := performJSON(router, http.MethodDelete, "/v1/workspaces/"+wsID+"/members/"+aliceID, nil, token)
	assert.Equal(t, http.StatusNoContent, resp.Code)

	// Verify only owner remains
	listResp := performJSON(router, http.MethodGet, "/v1/workspaces/"+wsID+"/members", nil, token)
	var list []map[string]any
	json.Unmarshal(listResp.Body.Bytes(), &list)
	assert.Len(t, list, 1)
}

func TestWorkspaceRoutes_KickSelf_Forbidden(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)

	// Find owner's userId
	membersResp := performJSON(router, http.MethodGet, "/v1/workspaces/"+wsID+"/members", nil, token)
	var members []map[string]any
	json.Unmarshal(membersResp.Body.Bytes(), &members)

	ownerID := members[0]["userId"].(string)

	resp := performJSON(router, http.MethodDelete, "/v1/workspaces/"+wsID+"/members/"+ownerID, nil, token)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// --- Grants ---

func TestWorkspaceRoutes_GetGrants_Empty(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)

	resp := performJSON(router, http.MethodGet, "/v1/workspaces/"+wsID+"/members/anyuser/grants", nil, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body map[string]any
	json.Unmarshal(resp.Body.Bytes(), &body)
	assert.NotNil(t, body["collections"])
	assert.NotNil(t, body["environments"])
	assert.NotNil(t, body["secrets"])
}

func TestWorkspaceRoutes_SetGrants(t *testing.T) {
	router, token, wsID := workspaceTestDeps(t)

	// Find owner's userId
	membersResp := performJSON(router, http.MethodGet, "/v1/workspaces/"+wsID+"/members", nil, token)
	var members []map[string]any
	json.Unmarshal(membersResp.Body.Bytes(), &members)
	ownerID := members[0]["userId"].(string)

	// Set grants
	setResp := performJSON(router, http.MethodPut, "/v1/workspaces/"+wsID+"/members/"+ownerID+"/grants", map[string]any{
		"grants": []map[string]string{
			{"resourceType": "collection", "resourceId": "c1", "accessLevel": "read"},
			{"resourceType": "environment", "resourceId": "e1", "accessLevel": "admin"},
		},
	}, token)
	assert.Equal(t, http.StatusNoContent, setResp.Code)

	// Get grants
	getResp := performJSON(router, http.MethodGet, "/v1/workspaces/"+wsID+"/members/"+ownerID+"/grants", nil, token)
	require.Equal(t, http.StatusOK, getResp.Code)

	var body map[string]any
	json.Unmarshal(getResp.Body.Bytes(), &body)

	cols := body["collections"].(map[string]any)
	envs := body["environments"].(map[string]any)
	assert.Equal(t, "read", cols["c1"])
	assert.Equal(t, "admin", envs["e1"])
}
