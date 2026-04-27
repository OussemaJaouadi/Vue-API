package http_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/auth"
	"vue-api/backend/internal/events"
	apihttp "vue-api/backend/internal/http"
)

func TestEventsTicketRequiresAccessToken(t *testing.T) {
	router := echo.New()
	router.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler

	apihttp.RegisterEventRoutes(router, apihttp.EventRouteDeps{
		Users:   auth.NewMemoryUserRepository(),
		Tokens:  testTokenManager(),
		Tickets: events.NewTicketStore(time.Minute),
		Broker:  events.NewBroker(),
	})

	resp := performJSON(router, http.MethodPost, "/events/ticket", nil, "")
	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestEventsTicketCanBeIssuedForCurrentUser(t *testing.T) {
	router := echo.New()
	router.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler

	repo := auth.NewMemoryUserRepository()
	user, err := repo.CreateUser(t.Context(), auth.CreateUserParams{
		Email:        "manager@example.com",
		Username:     "manager",
		PasswordHash: "hash",
		GlobalRole:   auth.GlobalRoleManager,
	})
	require.NoError(t, err)

	tokens := testTokenManager()
	accessToken, err := tokens.IssueAccessToken(user.ID, user.TokenVersion)
	require.NoError(t, err)
	tickets := events.NewTicketStore(time.Minute)
	apihttp.RegisterEventRoutes(router, apihttp.EventRouteDeps{
		Users:   repo,
		Tokens:  tokens,
		Tickets: tickets,
		Broker:  events.NewBroker(),
	})

	ticketResp := performJSON(router, http.MethodPost, "/events/ticket", nil, accessToken)
	require.Equal(t, http.StatusOK, ticketResp.Code)

	var body map[string]string
	require.NoError(t, json.Unmarshal(ticketResp.Body.Bytes(), &body))
	require.NotEmpty(t, body["ticket"])

	subscriber, err := tickets.Consume(body["ticket"])
	require.NoError(t, err)
	require.Equal(t, user.ID, subscriber.UserID)
	require.Equal(t, auth.GlobalRoleManager, subscriber.GlobalRole)
}

func TestEventsStreamRejectsInvalidTicket(t *testing.T) {
	router := echo.New()
	router.HTTPErrorHandler = apihttp.CustomHTTPErrorHandler

	apihttp.RegisterEventRoutes(router, apihttp.EventRouteDeps{
		Users:   auth.NewMemoryUserRepository(),
		Tokens:  testTokenManager(),
		Tickets: events.NewTicketStore(time.Minute),
		Broker:  events.NewBroker(),
	})

	resp := performJSON(router, http.MethodGet, "/events?ticket=bad-ticket", nil, "")
	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func testTokenManager() auth.TokenManager {
	return auth.NewTokenManager(auth.TokenConfig{
		AccessSecret:  "access-secret-at-least-32-bytes-long",
		RefreshSecret: "refresh-secret-at-least-32-bytes-long",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
		Issuer:        "vue-api-test",
	})
}
