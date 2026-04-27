package events_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/auth"
	"vue-api/backend/internal/events"
)

func TestBrokerPublishesRegisteredUsersToManagersOnly(t *testing.T) {
	broker := events.NewBroker()
	managerSub := broker.Subscribe(context.Background(), events.Subscriber{
		UserID:     "manager-id",
		GlobalRole: auth.GlobalRoleManager,
	})
	userSub := broker.Subscribe(context.Background(), events.Subscriber{
		UserID:     "user-id",
		GlobalRole: auth.GlobalRoleUser,
	})
	defer managerSub.Close()
	defer userSub.Close()

	broker.PublishToManagers(events.Event{
		Type: "user.registered",
		Data: map[string]string{"userId": "new-user-id"},
	})

	require.Equal(t, "user.registered", receiveEvent(t, managerSub).Type)
	requireNoEvent(t, userSub)
}

func TestBrokerPublishesMembershipEventsToOneUser(t *testing.T) {
	broker := events.NewBroker()
	targetSub := broker.Subscribe(context.Background(), events.Subscriber{UserID: "target-id"})
	otherSub := broker.Subscribe(context.Background(), events.Subscriber{UserID: "other-id"})
	defer targetSub.Close()
	defer otherSub.Close()

	broker.PublishToUser("target-id", events.Event{
		Type: "membership.created",
		Data: map[string]string{"workspaceId": "workspace-id"},
	})

	require.Equal(t, "membership.created", receiveEvent(t, targetSub).Type)
	requireNoEvent(t, otherSub)
}

func receiveEvent(t *testing.T, sub events.Subscription) events.Event {
	t.Helper()

	select {
	case event := <-sub.Events():
		return event
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event")
	}

	return events.Event{}
}

func requireNoEvent(t *testing.T, sub events.Subscription) {
	t.Helper()

	select {
	case event := <-sub.Events():
		t.Fatalf("unexpected event: %#v", event)
	case <-time.After(20 * time.Millisecond):
	}
}
