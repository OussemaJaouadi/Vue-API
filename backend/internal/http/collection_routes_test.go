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

func collectionTestDeps(t *testing.T) (*echo.Echo, string, string) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=private", t.Name())), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, gormstorage.Migrate(db))

	collections := gormstorage.NewCollectionRepository(db)
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
	apihttp.RegisterCollectionRoutes(router, apihttp.CollectionRouteDeps{
		Collections: collections,
		Memberships: membershipRepo,
		Users:       users,
		Tokens:      tokens,
	})

	registerResp := performJSON(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "test@example.com",
		"username": "testuser",
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
		Role:            "admin",
		CreatedByUserID: userID,
	})
	require.NoError(t, err)

	return router, accessToken, ws.ID
}

func TestCollectionRoutes_GetCollections_NoAuth(t *testing.T) {
	router := echo.New()
	router.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler

	users := auth.NewMemoryUserRepository()
	tokens := testTokenManager()
	apihttp.RegisterCollectionRoutes(router, apihttp.CollectionRouteDeps{
		Collections: nil,
		Users:       users,
		Tokens:      tokens,
	})

	resp := performJSON(router, http.MethodGet, "/v1/collections?workspaceId=ws1", nil, "")
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestCollectionRoutes_GetCollections_MissingWorkspaceID(t *testing.T) {
	router, token, _ := collectionTestDeps(t)

	resp := performJSON(router, http.MethodGet, "/v1/collections", nil, token)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var body map[string]any
	json.Unmarshal(resp.Body.Bytes(), &body)
	assert.Contains(t, body["error"], "workspaceId")
}

func TestCollectionRoutes_GetCollections_Empty(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodGet, "/v1/collections?workspaceId="+wsID, nil, token)
	assert.Equal(t, http.StatusOK, resp.Code)

	var body map[string]any
	json.Unmarshal(resp.Body.Bytes(), &body)
	assert.Empty(t, body["collections"])
	assert.Empty(t, body["rootRequests"])
}

func TestCollectionRoutes_GetCollections_WithFoldersAndRequests(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	createResp := performJSON(router, http.MethodPost, "/v1/collections", map[string]string{
		"workspaceId": wsID,
		"name":        "Auth",
		"icon":        "PhKey",
	}, token)
	require.Equal(t, http.StatusCreated, createResp.Code)

	var createBody map[string]any
	json.Unmarshal(createResp.Body.Bytes(), &createBody)
	folderID := createBody["id"].(string)

	createResp = performJSON(router, http.MethodPost, "/v1/collections/requests", map[string]any{
		"workspaceId":  wsID,
		"collectionId": folderID,
		"method":       "POST",
		"name":         "Login",
		"path":         "/auth/login",
	}, token)
	require.Equal(t, http.StatusCreated, createResp.Code)

	createResp = performJSON(router, http.MethodPost, "/v1/collections/requests", map[string]any{
		"workspaceId": wsID,
		"method":      "GET",
		"name":        "Health",
		"path":        "/health",
	}, token)
	require.Equal(t, http.StatusCreated, createResp.Code)

	getResp := performJSON(router, http.MethodGet, "/v1/collections?workspaceId="+wsID, nil, token)
	require.Equal(t, http.StatusOK, getResp.Code)

	var body map[string]any
	json.Unmarshal(getResp.Body.Bytes(), &body)

	collections := body["collections"].([]any)
	require.Len(t, collections, 1)
	coll := collections[0].(map[string]any)
	assert.Equal(t, "Auth", coll["name"])
	assert.Equal(t, "PhKey", coll["icon"])

	requests := coll["requests"].([]any)
	require.Len(t, requests, 1)
	assert.Equal(t, "Login", requests[0].(map[string]any)["name"])

	rootReqs := body["rootRequests"].([]any)
	require.Len(t, rootReqs, 1)
	assert.Equal(t, "Health", rootReqs[0].(map[string]any)["name"])
}

func TestCollectionRoutes_CreateFolder_MissingFields(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections", map[string]string{
		"workspaceId": wsID,
	}, token)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	resp = performJSON(router, http.MethodPost, "/v1/collections", map[string]string{
		"name": "Test",
	}, token)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCollectionRoutes_CreateFolder_Duplicate(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	performJSON(router, http.MethodPost, "/v1/collections", map[string]string{
		"workspaceId": wsID,
		"name":        "Auth",
	}, token)

	resp := performJSON(router, http.MethodPost, "/v1/collections", map[string]string{
		"workspaceId": wsID,
		"name":        "Auth",
	}, token)
	assert.Equal(t, http.StatusConflict, resp.Code)
}

func TestCollectionRoutes_UpdateFolder(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	createResp := performJSON(router, http.MethodPost, "/v1/collections", map[string]string{
		"workspaceId": wsID,
		"name":        "Old",
	}, token)
	require.Equal(t, http.StatusCreated, createResp.Code)

	var createBody map[string]any
	json.Unmarshal(createResp.Body.Bytes(), &createBody)
	folderID := createBody["id"].(string)

	updateResp := performJSON(router, http.MethodPut, "/v1/collections/"+folderID, map[string]string{
		"workspaceId": wsID,
		"name":        "Renamed",
		"icon":        "PhLock",
	}, token)
	require.Equal(t, http.StatusOK, updateResp.Code)

	var updateBody map[string]any
	json.Unmarshal(updateResp.Body.Bytes(), &updateBody)
	assert.Equal(t, "Renamed", updateBody["name"])
	assert.Equal(t, "PhLock", updateBody["icon"])
}

func TestCollectionRoutes_UpdateFolder_NotFound(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPut, "/v1/collections/nonexistent", map[string]string{
		"workspaceId": wsID,
		"name":        "Nope",
	}, token)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestCollectionRoutes_DeleteFolder(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	createResp := performJSON(router, http.MethodPost, "/v1/collections", map[string]string{
		"workspaceId": wsID,
		"name":        "Temp",
	}, token)
	require.Equal(t, http.StatusCreated, createResp.Code)

	var createBody map[string]any
	json.Unmarshal(createResp.Body.Bytes(), &createBody)
	folderID := createBody["id"].(string)

	delResp := performJSON(router, http.MethodDelete, "/v1/collections/"+folderID+"?workspaceId="+wsID, nil, token)
	assert.Equal(t, http.StatusNoContent, delResp.Code)
}

func TestCollectionRoutes_DeleteFolder_NotFound(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodDelete, "/v1/collections/nonexistent?workspaceId="+wsID, nil, token)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestCollectionRoutes_CreateRequest_MissingFields(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/requests", map[string]string{
		"workspaceId": wsID,
	}, token)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	resp = performJSON(router, http.MethodPost, "/v1/collections/requests", map[string]string{
		"name": "Test",
	}, token)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCollectionRoutes_CreateRequest_DefaultsMethod(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/requests", map[string]string{
		"workspaceId": wsID,
		"name":        "Default",
		"path":        "/",
	}, token)
	require.Equal(t, http.StatusCreated, resp.Code)

	var body map[string]any
	json.Unmarshal(resp.Body.Bytes(), &body)
	assert.Equal(t, "GET", body["method"])
}

func TestCollectionRoutes_UpdateRequest(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	createResp := performJSON(router, http.MethodPost, "/v1/collections/requests", map[string]string{
		"workspaceId": wsID,
		"method":      "GET",
		"name":        "Old",
		"path":        "/old",
	}, token)
	require.Equal(t, http.StatusCreated, createResp.Code)

	var createBody map[string]any
	json.Unmarshal(createResp.Body.Bytes(), &createBody)
	reqID := createBody["id"].(string)

	updateResp := performJSON(router, http.MethodPut, "/v1/collections/requests/"+reqID, map[string]string{
		"workspaceId": wsID,
		"method":      "POST",
		"name":        "Updated",
		"path":        "/new",
	}, token)
	require.Equal(t, http.StatusOK, updateResp.Code)

	var updateBody map[string]any
	json.Unmarshal(updateResp.Body.Bytes(), &updateBody)
	assert.Equal(t, "POST", updateBody["method"])
	assert.Equal(t, "Updated", updateBody["name"])
	assert.Equal(t, "/new", updateBody["path"])
}

func TestCollectionRoutes_UpdateRequest_NotFound(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPut, "/v1/collections/requests/nonexistent", map[string]string{
		"workspaceId": wsID,
		"name":        "Nope",
	}, token)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestCollectionRoutes_DeleteRequest(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	createResp := performJSON(router, http.MethodPost, "/v1/collections/requests", map[string]string{
		"workspaceId": wsID,
		"name":        "Temp",
		"path":        "/temp",
	}, token)
	require.Equal(t, http.StatusCreated, createResp.Code)

	var createBody map[string]any
	json.Unmarshal(createResp.Body.Bytes(), &createBody)
	reqID := createBody["id"].(string)

	delResp := performJSON(router, http.MethodDelete, "/v1/collections/requests/"+reqID+"?workspaceId="+wsID, nil, token)
	assert.Equal(t, http.StatusNoContent, delResp.Code)
}

func TestCollectionRoutes_DeleteRequest_NotFound(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodDelete, "/v1/collections/requests/nonexistent?workspaceId="+wsID, nil, token)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}
