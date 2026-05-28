package http_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"vue-api/backend/internal/auth"
	apihttp "vue-api/backend/internal/http"
	gormstorage "vue-api/backend/internal/storage/gorm"
	"vue-api/backend/internal/workspace"
)

func environmentTestDeps(t *testing.T) (*echo.Echo, string, string) {
	return environmentTestDepsWithRole(t, "admin")
}

func environmentTestDepsWithRole(t *testing.T, role string) (*echo.Echo, string, string) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=private", t.Name())), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, gormstorage.Migrate(db))

	environments := gormstorage.NewEnvironmentRepository(db)
	users := auth.NewMemoryUserRepository()
	tokens := testTokenManager()
	hasher := auth.NewPasswordHasher(auth.PasswordHashParams{
		MemoryKB:    1024,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	})

	workspaceRepo := gormstorage.NewWorkspaceRepository(db)
	membershipRepo := gormstorage.NewMembershipRepository(db)

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
	apihttp.RegisterEnvironmentRoutes(router, apihttp.EnvironmentRouteDeps{
		Environments: environments,
		Memberships:  membershipRepo,
		Users:        users,
		Tokens:       tokens,
	})

	managerResp := performJSON(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "envmanager@example.com",
		"username": "envmanager",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusCreated, managerResp.Code)

	registerResp := performJSON(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "envtest@example.com",
		"username": "envtest",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusCreated, registerResp.Code)

	var registerBody map[string]any
	require.NoError(t, json.Unmarshal(registerResp.Body.Bytes(), &registerBody))
	accessToken := registerBody["accessToken"].(string)
	userID := registerBody["userId"].(string)

	ws, err := workspaceRepo.Create(nil, workspace.CreateWorkspaceParams{
		Name:            "Test Workspace",
		CreatedByUserID: userID,
	})
	require.NoError(t, err)

	_, err = membershipRepo.Create(nil, workspace.CreateMembershipParams{
		WorkspaceID:     ws.ID,
		UserID:          userID,
		Role:            role,
		CreatedByUserID: userID,
	})
	require.NoError(t, err)

	return router, accessToken, ws.ID
}

func registerEnvironmentRouteUser(t *testing.T, router *echo.Echo, email string, username string) string {
	t.Helper()

	resp := performJSON(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    email,
		"username": username,
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusCreated, resp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	return body["accessToken"].(string)
}

func TestEnvironmentRoutes_GetEnvironments_NoAuth(t *testing.T) {
	router := echo.New()
	router.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler

	users := auth.NewMemoryUserRepository()
	tokens := testTokenManager()
	apihttp.RegisterEnvironmentRoutes(router, apihttp.EnvironmentRouteDeps{
		Environments: nil,
		Users:        users,
		Tokens:       tokens,
	})

	resp := performJSON(router, http.MethodGet, "/v1/environments?workspaceId=ws1", nil, "")
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestEnvironmentRoutes_GetEnvironments_MissingWorkspaceID(t *testing.T) {
	router, token, _ := environmentTestDeps(t)

	resp := performJSON(router, http.MethodGet, "/v1/environments", nil, token)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var body map[string]any
	json.Unmarshal(resp.Body.Bytes(), &body)
	assert.Contains(t, body["error"].(string), "workspaceId")
}

func TestEnvironmentRoutes_GetEnvironments_Empty(t *testing.T) {
	router, token, wsID := environmentTestDeps(t)

	resp := performJSON(router, http.MethodGet, "/v1/environments?workspaceId="+wsID, nil, token)
	assert.Equal(t, http.StatusOK, resp.Code)

	var body []any
	json.Unmarshal(resp.Body.Bytes(), &body)
	assert.Empty(t, body)
}

func TestEnvironmentRoutes_GetEnvironments_RequiresWorkspaceMembership(t *testing.T) {
	router, _, wsID := environmentTestDeps(t)
	otherToken := registerEnvironmentRouteUser(t, router, "envother@example.com", "envother")

	resp := performJSON(router, http.MethodGet, "/v1/environments?workspaceId="+wsID, nil, otherToken)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestEnvironmentRoutes_GetEnvironments_WithData(t *testing.T) {
	router, token, wsID := environmentTestDeps(t)

	envResp := performJSON(router, http.MethodPost, "/v1/environments", map[string]string{
		"workspaceId": wsID,
		"name":        "Production",
		"visibility":  "project",
	}, token)
	require.Equal(t, http.StatusCreated, envResp.Code)

	var createBody map[string]any
	json.Unmarshal(envResp.Body.Bytes(), &createBody)
	envID := createBody["id"].(string)

	performJSON(router, http.MethodPost, "/v1/environments/"+envID+"/variables", map[string]any{
		"workspaceId": wsID,
		"key":         "DB_URL",
		"value":       "localhost",
		"secret":      false,
	}, token)

	performJSON(router, http.MethodPost, "/v1/environments/"+envID+"/variables", map[string]any{
		"workspaceId": wsID,
		"key":         "PASSWORD",
		"value":       "hunter2",
		"secret":      true,
	}, token)

	getResp := performJSON(router, http.MethodGet, "/v1/environments?workspaceId="+wsID, nil, token)
	require.Equal(t, http.StatusOK, getResp.Code)

	var body []any
	json.Unmarshal(getResp.Body.Bytes(), &body)
	require.Len(t, body, 1)

	env := body[0].(map[string]any)
	assert.Equal(t, "Production", env["name"])
	assert.Equal(t, "project", env["visibility"])

	vars := env["variables"].([]any)
	require.Len(t, vars, 2)

	var dbURLVar, passwordVar map[string]any
	for _, v := range vars {
		vmap := v.(map[string]any)
		if vmap["key"] == "DB_URL" {
			dbURLVar = vmap
		} else {
			passwordVar = vmap
		}
	}

	assert.Equal(t, "localhost", dbURLVar["value"])
	assert.Equal(t, false, dbURLVar["secret"])
	assert.Equal(t, "••••••••••••••••", passwordVar["value"])
	assert.Equal(t, true, passwordVar["secret"])
}

func TestEnvironmentRoutes_CreateEnvironment_MissingFields(t *testing.T) {
	router, token, wsID := environmentTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/environments", map[string]string{
		"workspaceId": wsID,
	}, token)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	resp = performJSON(router, http.MethodPost, "/v1/environments", map[string]string{
		"name": "Test",
	}, token)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestEnvironmentRoutes_CreateEnvironment_Duplicate(t *testing.T) {
	router, token, wsID := environmentTestDeps(t)

	performJSON(router, http.MethodPost, "/v1/environments", map[string]string{
		"workspaceId": wsID,
		"name":        "Prod",
	}, token)

	resp := performJSON(router, http.MethodPost, "/v1/environments", map[string]string{
		"workspaceId": wsID,
		"name":        "Prod",
	}, token)
	assert.Equal(t, http.StatusConflict, resp.Code)
}

func TestEnvironmentRoutes_CreateEnvironment_RequiresManageEnvironmentPermission(t *testing.T) {
	router, token, wsID := environmentTestDepsWithRole(t, "tester")

	resp := performJSON(router, http.MethodPost, "/v1/environments", map[string]string{
		"workspaceId": wsID,
		"name":        "Prod",
	}, token)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestEnvironmentRoutes_UpdateEnvironment(t *testing.T) {
	router, token, wsID := environmentTestDeps(t)

	createResp := performJSON(router, http.MethodPost, "/v1/environments", map[string]string{
		"workspaceId": wsID,
		"name":        "Old",
	}, token)
	require.Equal(t, http.StatusCreated, createResp.Code)

	var createBody map[string]any
	json.Unmarshal(createResp.Body.Bytes(), &createBody)
	envID := createBody["id"].(string)

	updateResp := performJSON(router, http.MethodPut, "/v1/environments/"+envID, map[string]string{
		"workspaceId": wsID,
		"name":        "Renamed",
		"visibility":  "personal",
	}, token)
	require.Equal(t, http.StatusOK, updateResp.Code)

	var updateBody map[string]any
	json.Unmarshal(updateResp.Body.Bytes(), &updateBody)
	assert.Equal(t, "Renamed", updateBody["name"])
	assert.Equal(t, "personal", updateBody["visibility"])
}

func TestEnvironmentRoutes_UpdateEnvironment_RequiresManageEnvironmentPermission(t *testing.T) {
	router, token, wsID := environmentTestDepsWithRole(t, "tester")

	resp := performJSON(router, http.MethodPut, "/v1/environments/some-env", map[string]string{
		"workspaceId": wsID,
		"name":        "Renamed",
	}, token)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestEnvironmentRoutes_UpdateEnvironment_NotFound(t *testing.T) {
	router, token, wsID := environmentTestDeps(t)

	resp := performJSON(router, http.MethodPut, "/v1/environments/nonexistent", map[string]string{
		"workspaceId": wsID,
		"name":        "Nope",
	}, token)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestEnvironmentRoutes_DeleteEnvironment(t *testing.T) {
	router, token, wsID := environmentTestDeps(t)

	createResp := performJSON(router, http.MethodPost, "/v1/environments", map[string]string{
		"workspaceId": wsID,
		"name":        "Temp",
	}, token)
	require.Equal(t, http.StatusCreated, createResp.Code)

	var createBody map[string]any
	json.Unmarshal(createResp.Body.Bytes(), &createBody)
	envID := createBody["id"].(string)

	delResp := performJSON(router, http.MethodDelete, "/v1/environments/"+envID+"?workspaceId="+wsID, nil, token)
	assert.Equal(t, http.StatusNoContent, delResp.Code)
}

func TestEnvironmentRoutes_DeleteEnvironment_NotFound(t *testing.T) {
	router, token, wsID := environmentTestDeps(t)

	resp := performJSON(router, http.MethodDelete, "/v1/environments/nonexistent?workspaceId="+wsID, nil, token)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestEnvironmentRoutes_CreateVariable_NoKey(t *testing.T) {
	router, token, wsID := environmentTestDeps(t)

	envResp := performJSON(router, http.MethodPost, "/v1/environments", map[string]string{
		"workspaceId": wsID,
		"name":        "Test",
	}, token)
	require.Equal(t, http.StatusCreated, envResp.Code)

	var envBody map[string]any
	json.Unmarshal(envResp.Body.Bytes(), &envBody)
	envID := envBody["id"].(string)

	resp := performJSON(router, http.MethodPost, "/v1/environments/"+envID+"/variables", map[string]any{
		"workspaceId": wsID,
		"value":       "val",
	}, token)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestEnvironmentRoutes_CreateVariable_DuplicateKey(t *testing.T) {
	router, token, wsID := environmentTestDeps(t)

	envResp := performJSON(router, http.MethodPost, "/v1/environments", map[string]string{
		"workspaceId": wsID,
		"name":        "Test",
	}, token)
	require.Equal(t, http.StatusCreated, envResp.Code)

	var envBody map[string]any
	json.Unmarshal(envResp.Body.Bytes(), &envBody)
	envID := envBody["id"].(string)
	require.NotEmpty(t, envID)

	performJSON(router, http.MethodPost, "/v1/environments/"+envID+"/variables", map[string]any{
		"workspaceId": wsID,
		"key":         "DB_URL",
		"value":       "first",
	}, token)

	resp := performJSON(router, http.MethodPost, "/v1/environments/"+envID+"/variables", map[string]any{
		"workspaceId": wsID,
		"key":         "DB_URL",
		"value":       "second",
	}, token)
	assert.Equal(t, http.StatusConflict, resp.Code)
}

func TestEnvironmentRoutes_UpdateVariable(t *testing.T) {
	router, token, wsID := environmentTestDeps(t)

	envResp := performJSON(router, http.MethodPost, "/v1/environments", map[string]string{
		"workspaceId": wsID,
		"name":        "Test",
	}, token)
	require.Equal(t, http.StatusCreated, envResp.Code)

	var envBody map[string]any
	json.Unmarshal(envResp.Body.Bytes(), &envBody)
	envID := envBody["id"].(string)

	createResp := performJSON(router, http.MethodPost, "/v1/environments/"+envID+"/variables", map[string]any{
		"workspaceId": wsID,
		"key":         "OLD",
		"value":       "old",
	}, token)
	require.Equal(t, http.StatusCreated, createResp.Code)

	var createBody map[string]any
	json.Unmarshal(createResp.Body.Bytes(), &createBody)
	varID := createBody["id"].(string)

	updateResp := performJSON(router, http.MethodPut, "/v1/environments/"+envID+"/variables/"+varID, map[string]any{
		"workspaceId": wsID,
		"key":         "NEW",
		"value":       "new",
		"secret":      true,
	}, token)
	require.Equal(t, http.StatusOK, updateResp.Code)

	var updateBody map[string]any
	json.Unmarshal(updateResp.Body.Bytes(), &updateBody)
	assert.Equal(t, "NEW", updateBody["key"])
	assert.Equal(t, "new", updateBody["value"])
	assert.Equal(t, true, updateBody["secret"])
}

func TestEnvironmentRoutes_UpdateVariable_NotFound(t *testing.T) {
	router, token, wsID := environmentTestDeps(t)

	resp := performJSON(router, http.MethodPut, "/v1/environments/ws1/variables/nonexistent", map[string]string{
		"workspaceId": wsID,
		"value":       "nope",
	}, token)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestEnvironmentRoutes_DeleteVariable(t *testing.T) {
	router, token, wsID := environmentTestDeps(t)

	envResp := performJSON(router, http.MethodPost, "/v1/environments", map[string]string{
		"workspaceId": wsID,
		"name":        "Test",
	}, token)
	require.Equal(t, http.StatusCreated, envResp.Code)

	var envBody map[string]any
	json.Unmarshal(envResp.Body.Bytes(), &envBody)
	envID := envBody["id"].(string)

	createResp := performJSON(router, http.MethodPost, "/v1/environments/"+envID+"/variables", map[string]any{
		"workspaceId": wsID,
		"key":         "TEMP",
		"value":       "x",
	}, token)
	require.Equal(t, http.StatusCreated, createResp.Code)

	var createBody map[string]any
	json.Unmarshal(createResp.Body.Bytes(), &createBody)
	varID := createBody["id"].(string)

	delResp := performJSON(router, http.MethodDelete, "/v1/environments/"+envID+"/variables/"+varID+"?workspaceId="+wsID, nil, token)
	assert.Equal(t, http.StatusNoContent, delResp.Code)
}

func TestEnvironmentRoutes_DeleteVariable_NotFound(t *testing.T) {
	router, token, wsID := environmentTestDeps(t)

	resp := performJSON(router, http.MethodDelete, "/v1/environments/ws1/variables/nonexistent?workspaceId="+wsID, nil, token)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}
