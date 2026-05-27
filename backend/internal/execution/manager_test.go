package execution_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/events"
	"vue-api/backend/internal/execution"
)

func newTestWSServer(t *testing.T) *httptest.Server {
	t.Helper()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			err = conn.WriteMessage(mt, []byte("echo: "+string(msg)))
			if err != nil {
				return
			}
		}
	}))

	return ts
}

func TestNewWSManager_CreatesEmpty(t *testing.T) {
	broker := events.NewBroker()
	mgr := execution.NewWSManager(broker)
	require.NotNil(t, mgr)
}

func TestConnect_InvalidURL(t *testing.T) {
	broker := events.NewBroker()
	mgr := execution.NewWSManager(broker)

	_, err := mgr.Connect(context.Background(), "user1", execution.Request{
		Method: execution.MethodGet,
		URL:    "://invalid",
	})
	require.Error(t, err)
}

func TestConnect_ConnectionRefused(t *testing.T) {
	broker := events.NewBroker()
	mgr := execution.NewWSManager(broker)

	_, err := mgr.Connect(context.Background(), "user1", execution.Request{
		Method: execution.MethodGet,
		URL:    "ws://127.0.0.1:1",
	})
	require.Error(t, err)
}

func TestSend_NonExistentExecution(t *testing.T) {
	broker := events.NewBroker()
	mgr := execution.NewWSManager(broker)

	err := mgr.Send("nonexistent", "hello")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestClose_NonExistentExecution(t *testing.T) {
	broker := events.NewBroker()
	mgr := execution.NewWSManager(broker)

	err := mgr.Close("nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestConnect_Send_Close_Lifecycle(t *testing.T) {
	ts := newTestWSServer(t)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	broker := events.NewBroker()
	mgr := execution.NewWSManager(broker)

	id, err := mgr.Connect(context.Background(), "user1", execution.Request{
		Method: execution.MethodGet,
		URL:    wsURL,
		Headers: []execution.Header{
			{Key: "Origin", Value: ts.URL, Enabled: true},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, id)
	assert.True(t, strings.HasPrefix(id, "ws_"))

	err = mgr.Send(id, "ping")
	require.NoError(t, err)

	err = mgr.Close(id)
	require.NoError(t, err)
}

func TestConnect_WithBrokerConnectedEvent(t *testing.T) {
	ts := newTestWSServer(t)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	broker := events.NewBroker()
	sub := broker.Subscribe(context.Background(), events.Subscriber{UserID: "user1"})
	defer sub.Close()

	mgr := execution.NewWSManager(broker)

	id, err := mgr.Connect(context.Background(), "user1", execution.Request{
		Method: execution.MethodGet,
		URL:    wsURL,
	})
	require.NoError(t, err)

	select {
	case evt := <-sub.Events():
		assert.Equal(t, "ws.connected", evt.Type)
		data, ok := evt.Data.(map[string]any)
		require.True(t, ok)
		assert.Equal(t, id, data["executionId"])
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for ws.connected event")
	}

	mgr.Close(id)
}

func TestClose_PublishesClosedEvent(t *testing.T) {
	ts := newTestWSServer(t)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	broker := events.NewBroker()
	sub := broker.Subscribe(context.Background(), events.Subscriber{UserID: "user2"})
	defer sub.Close()

	mgr := execution.NewWSManager(broker)

	id, err := mgr.Connect(context.Background(), "user2", execution.Request{
		Method: execution.MethodGet,
		URL:    wsURL,
	})
	require.NoError(t, err)

	drainEvents(sub)

	mgr.Close(id)

	select {
	case evt := <-sub.Events():
		assert.Equal(t, "ws.closed", evt.Type)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for ws.closed event")
	}
}

func TestSend_PublishesOutgoingMessageEvent(t *testing.T) {
	ts := newTestWSServer(t)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	broker := events.NewBroker()
	sub := broker.Subscribe(context.Background(), events.Subscriber{UserID: "user3"})
	defer sub.Close()

	mgr := execution.NewWSManager(broker)

	id, err := mgr.Connect(context.Background(), "user3", execution.Request{
		Method: execution.MethodGet,
		URL:    wsURL,
	})
	require.NoError(t, err)

	drainEvents(sub)

	mgr.Send(id, "hello")

	select {
	case evt := <-sub.Events():
		assert.Equal(t, "ws.message.out", evt.Type)
		data, ok := evt.Data.(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "hello", data["payload"])
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for ws.message.out event")
	}

	mgr.Close(id)
}

func TestDoubleClose_ReturnsError(t *testing.T) {
	ts := newTestWSServer(t)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	broker := events.NewBroker()
	mgr := execution.NewWSManager(broker)

	id, err := mgr.Connect(context.Background(), "user1", execution.Request{
		Method: execution.MethodGet,
		URL:    wsURL,
	})
	require.NoError(t, err)

	err = mgr.Close(id)
	require.NoError(t, err)

	err = mgr.Close(id)
	require.Error(t, err)
}

func TestConnect_SetsTargetAndUserID(t *testing.T) {
	ts := newTestWSServer(t)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	broker := events.NewBroker()
	mgr := execution.NewWSManager(broker)

	id, err := mgr.Connect(context.Background(), "user42", execution.Request{
		Method: execution.MethodGet,
		URL:    wsURL,
	})
	require.NoError(t, err)

	err = mgr.Close(id)
	require.NoError(t, err)
}

func TestManager_SupportsMultipleConnections(t *testing.T) {
	ts := newTestWSServer(t)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	broker := events.NewBroker()
	mgr := execution.NewWSManager(broker)

	id1, err := mgr.Connect(context.Background(), "user1", execution.Request{
		Method: execution.MethodGet,
		URL:    wsURL,
	})
	require.NoError(t, err)

	id2, err := mgr.Connect(context.Background(), "user2", execution.Request{
		Method: execution.MethodGet,
		URL:    wsURL,
	})
	require.NoError(t, err)

	assert.NotEqual(t, id1, id2)

	err = mgr.Send(id1, "to-first")
	require.NoError(t, err)

	err = mgr.Send(id2, "to-second")
	require.NoError(t, err)

	err = mgr.Close(id1)
	require.NoError(t, err)

	err = mgr.Close(id2)
	require.NoError(t, err)
}

func TestConnect_WithBrokerNil(t *testing.T) {
	ts := newTestWSServer(t)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	mgr := execution.NewWSManager(nil)

	id, err := mgr.Connect(context.Background(), "user1", execution.Request{
		Method: execution.MethodGet,
		URL:    wsURL,
	})
	require.NoError(t, err)
	require.NotEmpty(t, id)

	err = mgr.Send(id, "test")
	require.NoError(t, err)

	err = mgr.Close(id)
	require.NoError(t, err)

	err = mgr.Close(id)
	require.Error(t, err)
}

func drainEvents(sub events.Subscription) {
	for {
		select {
		case <-sub.Events():
		default:
			return
		}
	}
}
