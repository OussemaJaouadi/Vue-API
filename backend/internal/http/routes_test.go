package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"vue-api/backend/internal/config"
	apihttp "vue-api/backend/internal/http"
)

func TestHealthz(t *testing.T) {
	// Setup
	e := echo.New()
	cfg := config.Config{
		App:      config.AppConfig{Env: "test"},
		Database: config.DatabaseConfig{Driver: "sqlite"},
	}
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	// Test handler
	apihttp.RegisterRoutes(e, cfg)

	// We need to find the handler since RegisterRoutes attaches it to the router
	e.ServeHTTP(rec, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "test", response["env"])
	assert.Equal(t, "sqlite", response["databaseDriver"])
}

func TestCustomErrorHandler(t *testing.T) {
	e := echo.New()
	e.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler

	req := httptest.NewRequest(http.MethodGet, "/not-found", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response apihttp.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Not Found", response.Error)
}

func TestCustomErrorHandlerDoesNotExposeInternalErrors(t *testing.T) {
	e := echo.New()
	e.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler
	e.GET("/internal-error", func(c echo.Context) error {
		return errors.New("database password leaked in stack")
	})

	req := httptest.NewRequest(http.MethodGet, "/internal-error", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response apihttp.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Internal Server Error", response.Error)
}
