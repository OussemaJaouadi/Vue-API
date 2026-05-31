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
	return collectionTestDepsWithRole(t, "admin")
}

func collectionTestDepsWithRole(t *testing.T, role string) (*echo.Echo, string, string) {
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

	managerResp := performJSON(router, http.MethodPost, "/auth/register", map[string]string{
		"email":    "manager@example.com",
		"username": "manager",
		"password": "correct horse battery staple",
	}, "")
	require.Equal(t, http.StatusCreated, managerResp.Code)

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
		Role:            role,
		CreatedByUserID: userID,
	})
	require.NoError(t, err)

	return router, accessToken, ws.ID
}

func registerCollectionRouteUser(t *testing.T, router *echo.Echo, email string, username string) string {
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

func TestCollectionRoutes_GetCollections_RequiresWorkspaceMembership(t *testing.T) {
	router, _, wsID := collectionTestDeps(t)
	otherToken := registerCollectionRouteUser(t, router, "other@example.com", "otheruser")

	resp := performJSON(router, http.MethodGet, "/v1/collections?workspaceId="+wsID, nil, otherToken)
	assert.Equal(t, http.StatusForbidden, resp.Code)
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

func TestCollectionRoutes_CreateFolder_RequiresManageCollectionsPermission(t *testing.T) {
	router, token, wsID := collectionTestDepsWithRole(t, "tester")

	resp := performJSON(router, http.MethodPost, "/v1/collections", map[string]string{
		"workspaceId": wsID,
		"name":        "Auth",
	}, token)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestCollectionRoutes_ImportWorkbenchExport(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/import", map[string]any{
		"workspaceId": wsID,
		"payload": map[string]any{
			"schema": "vue-api-workbench.collection.v1",
			"collections": []any{
				map[string]any{
					"name": "Authentication",
					"icon": "PhKey",
					"requests": []any{
						map[string]any{
							"method":       "POST",
							"name":         "Login",
							"path":         "/auth/login",
							"queryParams":  []any{map[string]any{"id": "q1", "key": "include_meta", "value": "true", "enabled": true}},
							"headers":      []any{map[string]any{"id": "h1", "key": "Content-Type", "value": "application/json", "enabled": true}},
							"body":         `{"login":"owner@example.com"}`,
							"bodyLanguage": "json",
							"authConfig":   map[string]any{"mode": "none"},
						},
						map[string]any{
							"method": "GET",
							"name":   "Current user",
							"path":   "/auth/me",
						},
					},
				},
			},
			"rootRequests": []any{
				map[string]any{
					"method": "GET",
					"name":   "Health",
					"path":   "/healthz",
				},
			},
		},
	}, token)
	require.Equal(t, http.StatusCreated, resp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Equal(t, "Workbench export", body["format"])
	assert.Equal(t, float64(1), body["collectionsCreated"])
	assert.Equal(t, float64(3), body["requestsCreated"])

	getResp := performJSON(router, http.MethodGet, "/v1/collections?workspaceId="+wsID, nil, token)
	require.Equal(t, http.StatusOK, getResp.Code)

	var collectionsBody map[string]any
	require.NoError(t, json.Unmarshal(getResp.Body.Bytes(), &collectionsBody))
	collections := collectionsBody["collections"].([]any)
	require.Len(t, collections, 1)
	imported := collections[0].(map[string]any)
	assert.Equal(t, "Authentication", imported["name"])
	assert.Equal(t, "PhKey", imported["icon"])

	requests := imported["requests"].([]any)
	require.Len(t, requests, 2)
	assert.Equal(t, "Login", requests[0].(map[string]any)["name"])
	assert.Equal(t, float64(0), requests[0].(map[string]any)["sortOrder"])
	assert.Equal(t, "Current user", requests[1].(map[string]any)["name"])
	assert.Equal(t, float64(1), requests[1].(map[string]any)["sortOrder"])

	rootRequests := collectionsBody["rootRequests"].([]any)
	require.Len(t, rootRequests, 1)
	assert.Equal(t, "Health", rootRequests[0].(map[string]any)["name"])
}

func TestCollectionRoutes_ImportWorkbenchExport_RenamesDuplicateCollection(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	createResp := performJSON(router, http.MethodPost, "/v1/collections", map[string]string{
		"workspaceId": wsID,
		"name":        "Authentication",
	}, token)
	require.Equal(t, http.StatusCreated, createResp.Code)

	resp := performJSON(router, http.MethodPost, "/v1/collections/import", map[string]any{
		"workspaceId": wsID,
		"payload": map[string]any{
			"schema": "vue-api-workbench.collection.v1",
			"collections": []any{
				map[string]any{
					"name":     "Authentication",
					"requests": []any{},
				},
			},
		},
	}, token)
	require.Equal(t, http.StatusCreated, resp.Code)

	getResp := performJSON(router, http.MethodGet, "/v1/collections?workspaceId="+wsID, nil, token)
	require.Equal(t, http.StatusOK, getResp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(getResp.Body.Bytes(), &body))
	collections := body["collections"].([]any)
	require.Len(t, collections, 2)
	assert.Equal(t, "Authentication", collections[0].(map[string]any)["name"])
	assert.Equal(t, "Authentication (import 2)", collections[1].(map[string]any)["name"])
}

func TestCollectionRoutes_ImportUnknownJSONPayload(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/import", map[string]any{
		"workspaceId": wsID,
		"payload": map[string]any{
			"not": "a collection format",
		},
	}, token)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}

func TestCollectionRoutes_ImportOpenAPIJSON(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/import", map[string]any{
		"workspaceId": wsID,
		"payload": map[string]any{
			"openapi": "3.1.0",
			"info": map[string]any{
				"title": "Sample API",
			},
			"paths": map[string]any{
				"/auth/login": map[string]any{
					"post": map[string]any{
						"summary": "Login",
						"parameters": []any{
							map[string]any{"name": "include_meta", "in": "query"},
						},
						"requestBody": map[string]any{
							"content": map[string]any{
								"application/json": map[string]any{},
							},
						},
					},
				},
				"/auth/me": map[string]any{
					"get": map[string]any{
						"operationId": "currentUser",
					},
				},
			},
		},
	}, token)
	require.Equal(t, http.StatusCreated, resp.Code)

	var importBody map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &importBody))
	assert.Equal(t, "OpenAPI 3.1.0", importBody["format"])
	assert.Equal(t, float64(1), importBody["collectionsCreated"])
	assert.Equal(t, float64(2), importBody["requestsCreated"])

	getResp := performJSON(router, http.MethodGet, "/v1/collections?workspaceId="+wsID, nil, token)
	require.Equal(t, http.StatusOK, getResp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(getResp.Body.Bytes(), &body))
	collections := body["collections"].([]any)
	require.Len(t, collections, 1)
	imported := collections[0].(map[string]any)
	assert.Equal(t, "Sample API", imported["name"])

	requests := imported["requests"].([]any)
	require.Len(t, requests, 2)
	assert.Equal(t, "POST", requests[0].(map[string]any)["method"])
	assert.Equal(t, "Login", requests[0].(map[string]any)["name"])
	assert.Equal(t, "/auth/login", requests[0].(map[string]any)["path"])
	queryParams := requests[0].(map[string]any)["queryParams"].([]any)
	require.Len(t, queryParams, 1)
	assert.Equal(t, "include_meta", queryParams[0].(map[string]any)["key"])
	assert.Equal(t, "GET", requests[1].(map[string]any)["method"])
	assert.Equal(t, "currentUser", requests[1].(map[string]any)["name"])
}

func TestCollectionRoutes_ImportSwaggerJSON(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/import", map[string]any{
		"workspaceId": wsID,
		"payload": map[string]any{
			"swagger":  "2.0",
			"basePath": "/v1",
			"info": map[string]any{
				"title": "Legacy API",
			},
			"paths": map[string]any{
				"/orders": map[string]any{
					"post": map[string]any{
						"operationId": "createOrder",
						"consumes":    []any{"application/json"},
						"parameters": []any{
							map[string]any{"name": "trace", "in": "query"},
							map[string]any{"name": "payload", "in": "body"},
						},
					},
				},
			},
		},
	}, token)
	require.Equal(t, http.StatusCreated, resp.Code)

	var importBody map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &importBody))
	assert.Equal(t, "Swagger 2.0", importBody["format"])
	assert.Equal(t, float64(1), importBody["collectionsCreated"])
	assert.Equal(t, float64(1), importBody["requestsCreated"])

	getResp := performJSON(router, http.MethodGet, "/v1/collections?workspaceId="+wsID, nil, token)
	require.Equal(t, http.StatusOK, getResp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(getResp.Body.Bytes(), &body))
	collections := body["collections"].([]any)
	require.Len(t, collections, 1)
	imported := collections[0].(map[string]any)
	assert.Equal(t, "Legacy API", imported["name"])

	requests := imported["requests"].([]any)
	require.Len(t, requests, 1)
	request := requests[0].(map[string]any)
	assert.Equal(t, "POST", request["method"])
	assert.Equal(t, "createOrder", request["name"])
	assert.Equal(t, "/v1/orders", request["path"])
	assert.Equal(t, "{}", request["body"])

	queryParams := request["queryParams"].([]any)
	require.Len(t, queryParams, 1)
	assert.Equal(t, "trace", queryParams[0].(map[string]any)["key"])
}

func TestCollectionRoutes_ImportOpenAPIYAML(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/import", map[string]any{
		"workspaceId": wsID,
		"fileName":    "openapi.yaml",
		"content": `
openapi: 3.1.0
info:
  title: YAML API
paths:
  /healthz:
    get:
      summary: Health check
`,
	}, token)
	require.Equal(t, http.StatusCreated, resp.Code)

	var importBody map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &importBody))
	assert.Equal(t, "OpenAPI 3.1.0", importBody["format"])
	assert.Equal(t, float64(1), importBody["collectionsCreated"])
	assert.Equal(t, float64(1), importBody["requestsCreated"])

	getResp := performJSON(router, http.MethodGet, "/v1/collections?workspaceId="+wsID, nil, token)
	require.Equal(t, http.StatusOK, getResp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(getResp.Body.Bytes(), &body))
	collections := body["collections"].([]any)
	require.Len(t, collections, 1)
	assert.Equal(t, "YAML API", collections[0].(map[string]any)["name"])
	requests := collections[0].(map[string]any)["requests"].([]any)
	require.Len(t, requests, 1)
	assert.Equal(t, "Health check", requests[0].(map[string]any)["name"])
}

func TestCollectionRoutes_ImportPostmanCollection(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/import", map[string]any{
		"workspaceId": wsID,
		"payload": map[string]any{
			"info": map[string]any{
				"_postman_id": "postman-id",
				"name":        "Postman API",
			},
			"item": []any{
				map[string]any{
					"name": "Authentication",
					"item": []any{
						map[string]any{
							"name": "Login",
							"request": map[string]any{
								"method": "POST",
								"header": []any{
									map[string]any{"key": "Content-Type", "value": "application/json"},
								},
								"body": map[string]any{"mode": "raw", "raw": `{"login":"owner@example.com"}`},
								"url": map[string]any{
									"raw": "https://api.example.com/auth/login?include_meta=true",
									"query": []any{
										map[string]any{"key": "include_meta", "value": "true"},
									},
								},
								"auth": map[string]any{
									"type": "bearer",
									"bearer": []any{
										map[string]any{"key": "token", "value": "{{accessToken}}"},
									},
								},
							},
						},
					},
				},
				map[string]any{
					"name": "Health",
					"request": map[string]any{
						"method": "GET",
						"url":    "https://api.example.com/healthz",
					},
				},
			},
		},
	}, token)
	require.Equal(t, http.StatusCreated, resp.Code)

	var importBody map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &importBody))
	assert.Equal(t, "Postman collection", importBody["format"])
	assert.Equal(t, float64(1), importBody["collectionsCreated"])
	assert.Equal(t, float64(2), importBody["requestsCreated"])

	getResp := performJSON(router, http.MethodGet, "/v1/collections?workspaceId="+wsID, nil, token)
	require.Equal(t, http.StatusOK, getResp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(getResp.Body.Bytes(), &body))
	collections := body["collections"].([]any)
	require.Len(t, collections, 1)
	assert.Equal(t, "Authentication", collections[0].(map[string]any)["name"])
	requests := collections[0].(map[string]any)["requests"].([]any)
	require.Len(t, requests, 1)
	login := requests[0].(map[string]any)
	assert.Equal(t, "POST", login["method"])
	assert.Equal(t, "/auth/login", login["path"])
	assert.Equal(t, `{"login":"owner@example.com"}`, login["body"])
	assert.Equal(t, "bearer", login["authConfig"].(map[string]any)["mode"])
	assert.Equal(t, "{{accessToken}}", login["authConfig"].(map[string]any)["bearerToken"])

	queryParams := login["queryParams"].([]any)
	require.Len(t, queryParams, 1)
	assert.Equal(t, "include_meta", queryParams[0].(map[string]any)["key"])

	rootRequests := body["rootRequests"].([]any)
	require.Len(t, rootRequests, 1)
	assert.Equal(t, "Health", rootRequests[0].(map[string]any)["name"])
	assert.Equal(t, "/healthz", rootRequests[0].(map[string]any)["path"])
}

func TestCollectionRoutes_ImportRequiresManageCollectionsPermission(t *testing.T) {
	router, token, wsID := collectionTestDepsWithRole(t, "tester")

	resp := performJSON(router, http.MethodPost, "/v1/collections/import", map[string]any{
		"workspaceId": wsID,
		"payload": map[string]any{
			"schema":      "vue-api-workbench.collection.v1",
			"collections": []any{},
		},
	}, token)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestCollectionRoutes_PreviewWorkbenchImport(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/import/preview", map[string]any{
		"workspaceId": wsID,
		"fileName":    "collections.json",
		"content": `{
			"schema": "vue-api-workbench.collection.v1",
			"collections": [
				{"name": "Authentication", "requests": [{"name": "Login"}, {"name": "Me"}]}
			],
			"rootRequests": [{"name": "Health"}]
		}`,
	}, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Equal(t, "collections.json", body["fileName"])
	assert.Equal(t, "Workbench export", body["format"])
	assert.Equal(t, "ready", body["status"])
	assert.Equal(t, "1 collections / 3 requests", body["summary"])
}

func TestCollectionRoutes_PreviewOpenAPIImport(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/import/preview", map[string]any{
		"workspaceId": wsID,
		"fileName":    "openapi.json",
		"content":     `{"openapi":"3.1.0","paths":{"/healthz":{"get":{}}}}`,
	}, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Equal(t, "OpenAPI 3.1.0", body["format"])
	assert.Equal(t, "ready", body["status"])
	assert.Equal(t, "1 paths detected", body["summary"])
}

func TestCollectionRoutes_PreviewSwaggerImport(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/import/preview", map[string]any{
		"workspaceId": wsID,
		"fileName":    "swagger.json",
		"content":     `{"swagger":"2.0","paths":{"/orders":{"post":{}}}}`,
	}, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Equal(t, "Swagger 2.0", body["format"])
	assert.Equal(t, "ready", body["status"])
	assert.Equal(t, "1 paths detected", body["summary"])
}

func TestCollectionRoutes_PreviewOpenAPIYAMLImport(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/import/preview", map[string]any{
		"workspaceId": wsID,
		"fileName":    "openapi.yaml",
		"content": `
openapi: 3.1.0
paths:
  /healthz:
    get: {}
`,
	}, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Equal(t, "OpenAPI 3.1.0", body["format"])
	assert.Equal(t, "ready", body["status"])
	assert.Equal(t, "1 paths detected", body["summary"])
}

func TestCollectionRoutes_PreviewPostmanImport(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/import/preview", map[string]any{
		"workspaceId": wsID,
		"fileName":    "postman.json",
		"content":     `{"info":{"_postman_id":"id","name":"Postman API"},"item":[{"name":"Health","request":{"method":"GET","url":"https://api.example.com/healthz"}}]}`,
	}, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Equal(t, "Postman collection", body["format"])
	assert.Equal(t, "ready", body["status"])
	assert.Equal(t, "1 top-level items detected", body["summary"])
}

func TestCollectionRoutes_PreviewInvalidYAMLImport(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/import/preview", map[string]any{
		"workspaceId": wsID,
		"fileName":    "broken.yaml",
		"content":     `openapi: [`,
	}, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Equal(t, "Invalid YAML", body["format"])
	assert.Equal(t, "error", body["status"])
}

func TestCollectionRoutes_PreviewInvalidJSONImport(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodPost, "/v1/collections/import/preview", map[string]any{
		"workspaceId": wsID,
		"fileName":    "broken.json",
		"content":     `{`,
	}, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Equal(t, "Invalid JSON", body["format"])
	assert.Equal(t, "error", body["status"])
}

func TestCollectionRoutes_PreviewImportRequiresWorkspaceMembership(t *testing.T) {
	router, _, wsID := collectionTestDeps(t)
	otherToken := registerCollectionRouteUser(t, router, "preview-other@example.com", "previewother")

	resp := performJSON(router, http.MethodPost, "/v1/collections/import/preview", map[string]any{
		"workspaceId": wsID,
		"fileName":    "collections.json",
		"content":     `{}`,
	}, otherToken)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestCollectionRoutes_ExportWorkbenchData(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	folderResp := performJSON(router, http.MethodPost, "/v1/collections", map[string]string{
		"workspaceId": wsID,
		"name":        "Authentication",
		"icon":        "PhKey",
	}, token)
	require.Equal(t, http.StatusCreated, folderResp.Code)

	var folderBody map[string]any
	require.NoError(t, json.Unmarshal(folderResp.Body.Bytes(), &folderBody))
	folderID := folderBody["id"].(string)

	requestResp := performJSON(router, http.MethodPost, "/v1/collections/requests", map[string]any{
		"workspaceId":  wsID,
		"collectionId": folderID,
		"method":       "POST",
		"name":         "Login",
		"path":         "/auth/login",
		"queryParams":  []any{map[string]any{"id": "q1", "key": "include_meta", "value": "true", "enabled": true}},
		"headers":      []any{map[string]any{"id": "h1", "key": "Content-Type", "value": "application/json", "enabled": true}},
		"body":         `{"login":"owner@example.com"}`,
		"bodyLanguage": "json",
		"authConfig":   map[string]any{"mode": "bearer", "bearerToken": "{{accessToken}}"},
	}, token)
	require.Equal(t, http.StatusCreated, requestResp.Code)

	rootResp := performJSON(router, http.MethodPost, "/v1/collections/requests", map[string]any{
		"workspaceId": wsID,
		"method":      "GET",
		"name":        "Health",
		"path":        "/healthz",
	}, token)
	require.Equal(t, http.StatusCreated, rootResp.Code)

	resp := performJSON(router, http.MethodGet, "/v1/collections/export?workspaceId="+wsID, nil, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Equal(t, "vue-api-workbench.collection.v1", body["schema"])
	assert.NotEmpty(t, body["exportedAt"])

	collections := body["collections"].([]any)
	require.Len(t, collections, 1)
	exportedCollection := collections[0].(map[string]any)
	assert.Equal(t, "Authentication", exportedCollection["name"])
	assert.Equal(t, "PhKey", exportedCollection["icon"])

	requests := exportedCollection["requests"].([]any)
	require.Len(t, requests, 1)
	exportedRequest := requests[0].(map[string]any)
	assert.Equal(t, "POST", exportedRequest["method"])
	assert.Equal(t, "Login", exportedRequest["name"])
	assert.Equal(t, "/auth/login", exportedRequest["path"])
	assert.Equal(t, `{"login":"owner@example.com"}`, exportedRequest["body"])
	assert.Equal(t, "json", exportedRequest["bodyLanguage"])
	assert.Equal(t, "bearer", exportedRequest["authConfig"].(map[string]any)["mode"])

	rootRequests := body["rootRequests"].([]any)
	require.Len(t, rootRequests, 1)
	assert.Equal(t, "Health", rootRequests[0].(map[string]any)["name"])
}

func TestCollectionRoutes_ExportSelectedCollection(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	authFolderResp := performJSON(router, http.MethodPost, "/v1/collections", map[string]string{
		"workspaceId": wsID,
		"name":        "Authentication",
		"icon":        "PhKey",
	}, token)
	require.Equal(t, http.StatusCreated, authFolderResp.Code)

	var authFolderBody map[string]any
	require.NoError(t, json.Unmarshal(authFolderResp.Body.Bytes(), &authFolderBody))
	authFolderID := authFolderBody["id"].(string)

	realtimeFolderResp := performJSON(router, http.MethodPost, "/v1/collections", map[string]string{
		"workspaceId": wsID,
		"name":        "Realtime",
		"icon":        "PhBroadcast",
	}, token)
	require.Equal(t, http.StatusCreated, realtimeFolderResp.Code)

	var realtimeFolderBody map[string]any
	require.NoError(t, json.Unmarshal(realtimeFolderResp.Body.Bytes(), &realtimeFolderBody))
	realtimeFolderID := realtimeFolderBody["id"].(string)

	authRequestResp := performJSON(router, http.MethodPost, "/v1/collections/requests", map[string]any{
		"workspaceId":  wsID,
		"collectionId": authFolderID,
		"method":       "POST",
		"name":         "Login",
		"path":         "/auth/login",
	}, token)
	require.Equal(t, http.StatusCreated, authRequestResp.Code)

	realtimeRequestResp := performJSON(router, http.MethodPost, "/v1/collections/requests", map[string]any{
		"workspaceId":  wsID,
		"collectionId": realtimeFolderID,
		"method":       "GET",
		"name":         "Event stream",
		"path":         "/events",
	}, token)
	require.Equal(t, http.StatusCreated, realtimeRequestResp.Code)

	rootResp := performJSON(router, http.MethodPost, "/v1/collections/requests", map[string]any{
		"workspaceId": wsID,
		"method":      "GET",
		"name":        "Health",
		"path":        "/healthz",
	}, token)
	require.Equal(t, http.StatusCreated, rootResp.Code)

	resp := performJSON(router, http.MethodGet, "/v1/collections/export?workspaceId="+wsID+"&collectionId="+authFolderID, nil, token)
	require.Equal(t, http.StatusOK, resp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))

	collections := body["collections"].([]any)
	require.Len(t, collections, 1)
	exportedCollection := collections[0].(map[string]any)
	assert.Equal(t, "Authentication", exportedCollection["name"])

	requests := exportedCollection["requests"].([]any)
	require.Len(t, requests, 1)
	assert.Equal(t, "Login", requests[0].(map[string]any)["name"])

	rootRequests := body["rootRequests"].([]any)
	assert.Empty(t, rootRequests)
}

func TestCollectionRoutes_ExportSelectedCollectionNotFound(t *testing.T) {
	router, token, wsID := collectionTestDeps(t)

	resp := performJSON(router, http.MethodGet, "/v1/collections/export?workspaceId="+wsID+"&collectionId=missing-collection", nil, token)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestCollectionRoutes_ExportRequiresWorkspaceMembership(t *testing.T) {
	router, _, wsID := collectionTestDeps(t)
	otherToken := registerCollectionRouteUser(t, router, "export-other@example.com", "exportother")

	resp := performJSON(router, http.MethodGet, "/v1/collections/export?workspaceId="+wsID, nil, otherToken)
	assert.Equal(t, http.StatusForbidden, resp.Code)
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

func TestCollectionRoutes_UpdateFolder_RequiresManageCollectionsPermission(t *testing.T) {
	router, token, wsID := collectionTestDepsWithRole(t, "tester")

	resp := performJSON(router, http.MethodPut, "/v1/collections/some-folder", map[string]string{
		"workspaceId": wsID,
		"name":        "Renamed",
	}, token)
	assert.Equal(t, http.StatusForbidden, resp.Code)
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

func TestCollectionRoutes_CreateRequest_RequiresManageCollectionsPermission(t *testing.T) {
	router, token, wsID := collectionTestDepsWithRole(t, "tester")

	resp := performJSON(router, http.MethodPost, "/v1/collections/requests", map[string]string{
		"workspaceId": wsID,
		"name":        "Blocked",
		"path":        "/blocked",
	}, token)
	assert.Equal(t, http.StatusForbidden, resp.Code)
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
