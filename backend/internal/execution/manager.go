package execution

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"vue-api/backend/internal/events"
)

type WSManager struct {
	mu          sync.RWMutex
	executions  map[string]*WSExecution
	eventBroker *events.Broker
}

type WSExecution struct {
	ID        string
	Conn      *websocket.Conn
	Cancel    context.CancelFunc
	UserID    string
	Target    string
	CreatedAt time.Time
}

func NewWSManager(broker *events.Broker) *WSManager {
	return &WSManager{
		executions:  make(map[string]*WSExecution),
		eventBroker: broker,
	}
}

func (m *WSManager) Connect(ctx context.Context, userID string, req Request) (string, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
		TLSClientConfig:  nil, // Optional: handle custom certs
	}

	header := http.Header{}
	for _, h := range req.Headers {
		if h.Enabled && h.Key != "" {
			header.Add(h.Key, h.Value)
		}
	}

	conn, _, err := dialer.DialContext(ctx, req.URL, header)
	if err != nil {
		return "", err
	}

	executionID := fmt.Sprintf("ws_%d", time.Now().UnixNano())
	
	execCtx, cancel := context.WithCancel(context.Background())
	exec := &WSExecution{
		ID:        executionID,
		Conn:      conn,
		Cancel:    cancel,
		UserID:    userID,
		Target:    req.URL,
		CreatedAt: time.Now(),
	}

	m.mu.Lock()
	m.executions[executionID] = exec
	m.mu.Unlock()

	go m.handleRead(execCtx, exec)

	return executionID, nil
}

func (m *WSManager) Send(executionID string, payload string) error {
	m.mu.RLock()
	exec, ok := m.executions[executionID]
	m.mu.RUnlock()

	if !ok {
		return fmt.Errorf("execution not found: %s", executionID)
	}

	err := exec.Conn.WriteMessage(websocket.TextMessage, []byte(payload))
	if err != nil {
		return err
	}

	if m.eventBroker != nil {
		m.eventBroker.PublishToUser(exec.UserID, events.Event{
			Type: "ws.message.out",
			Data: map[string]any{
				"executionId": executionID,
				"payload":     payload,
				"timestamp":   time.Now().Format(time.RFC3339),
			},
		})
	}

	return nil
}

func (m *WSManager) Close(executionID string) error {
	m.mu.Lock()
	exec, ok := m.executions[executionID]
	if ok {
		delete(m.executions, executionID)
	}
	m.mu.Unlock()

	if !ok {
		return fmt.Errorf("execution not found: %s", executionID)
	}

	exec.Cancel()
	return exec.Conn.Close()
}

func (m *WSManager) handleRead(ctx context.Context, exec *WSExecution) {
	defer func() {
		m.Close(exec.ID)
		if m.eventBroker != nil {
			m.eventBroker.PublishToUser(exec.UserID, events.Event{
				Type: "ws.closed",
				Data: map[string]any{
					"executionId": exec.ID,
					"timestamp":   time.Now().Format(time.RFC3339),
				},
			})
		}
	}()

	if m.eventBroker != nil {
		m.eventBroker.PublishToUser(exec.UserID, events.Event{
			Type: "ws.connected",
			Data: map[string]any{
				"executionId": exec.ID,
				"target":      exec.Target,
				"timestamp":   time.Now().Format(time.RFC3339),
			},
		})
	}

	for {
		messageType, payload, err := exec.Conn.ReadMessage()
		if err != nil {
			if m.eventBroker != nil {
				m.eventBroker.PublishToUser(exec.UserID, events.Event{
					Type: "ws.error",
					Data: map[string]any{
						"executionId": exec.ID,
						"error":       err.Error(),
						"timestamp":   time.Now().Format(time.RFC3339),
					},
				})
			}
			return
		}

		if m.eventBroker != nil {
			var eventType string
			switch messageType {
			case websocket.TextMessage:
				eventType = "ws.message.in"
			case websocket.BinaryMessage:
				eventType = "ws.message.in.binary"
			default:
				continue
			}

			m.eventBroker.PublishToUser(exec.UserID, events.Event{
				Type: eventType,
				Data: map[string]any{
					"executionId": exec.ID,
					"payload":     string(payload),
					"sizeBytes":   len(payload),
					"timestamp":   time.Now().Format(time.RFC3339),
				},
			})
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
