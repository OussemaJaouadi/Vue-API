package events_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/events"
)

func TestTicketStoreIssuesSingleUseTickets(t *testing.T) {
	store := events.NewTicketStore(time.Minute)

	ticket, err := store.Issue(events.Subscriber{UserID: "user-id", GlobalRole: "manager"})
	require.NoError(t, err)
	require.NotEmpty(t, ticket)

	subscriber, err := store.Consume(ticket)
	require.NoError(t, err)
	require.Equal(t, "user-id", subscriber.UserID)
	require.Equal(t, "manager", subscriber.GlobalRole)

	_, err = store.Consume(ticket)
	require.ErrorIs(t, err, events.ErrInvalidTicket)
}
