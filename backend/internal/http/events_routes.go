package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"vue-api/backend/internal/auth"
	"vue-api/backend/internal/events"
)

type EventRouteDeps struct {
	Users   auth.UserRepository
	Tokens  auth.TokenManager
	Tickets *events.TicketStore
	Broker  *events.Broker
}

func RegisterEventRoutes(router *echo.Echo, deps EventRouteDeps) {
	router.POST("/events/ticket", func(c echo.Context) error {
		user, err := authenticateRequest(c, AuthRouteDeps{
			Users:  deps.Users,
			Tokens: deps.Tokens,
		})
		if err != nil {
			return err
		}

		ticket, err := deps.Tickets.Issue(events.Subscriber{
			UserID:     user.ID,
			GlobalRole: user.GlobalRole,
		})
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"ticket": ticket,
		})
	})

	router.GET("/events", func(c echo.Context) error {
		ticket := c.QueryParam("ticket")
		subscriber, err := deps.Tickets.Consume(ticket)
		if errors.Is(err, events.ErrInvalidTicket) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid event ticket")
		}
		if err != nil {
			return err
		}

		res := c.Response()
		res.Header().Set(echo.HeaderContentType, "text/event-stream")
		res.Header().Set(echo.HeaderCacheControl, "no-cache")
		res.WriteHeader(http.StatusOK)

		if err := writeSSE(res, events.Event{Type: "connected"}); err != nil {
			return err
		}

		subscription := deps.Broker.Subscribe(c.Request().Context(), subscriber)
		defer subscription.Close()

		for {
			select {
			case <-c.Request().Context().Done():
				return nil
			case event, ok := <-subscription.Events():
				if !ok {
					return nil
				}
				if err := writeSSE(res, event); err != nil {
					return err
				}
			}
		}
	})
}

func writeSSE(res *echo.Response, event events.Event) error {
	if _, err := fmt.Fprintf(res, "event: %s\n", event.Type); err != nil {
		return err
	}

	if event.Data != nil {
		data, err := json.Marshal(event.Data)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(res, "data: %s\n", data); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(res, "\n"); err != nil {
		return err
	}
	res.Flush()

	return nil
}

func NewEventDeps(users auth.UserRepository, tokens auth.TokenManager) EventRouteDeps {
	return EventRouteDeps{
		Users:   users,
		Tokens:  tokens,
		Tickets: events.NewTicketStore(30 * time.Second),
		Broker:  events.NewBroker(),
	}
}
