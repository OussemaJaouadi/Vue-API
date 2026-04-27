package events

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

var ErrInvalidTicket = errors.New("invalid event ticket")

type TicketStore struct {
	mu      sync.Mutex
	ttl     time.Duration
	tickets map[string]ticketEntry
	now     func() time.Time
}

type ticketEntry struct {
	subscriber Subscriber
	expiresAt  time.Time
}

func NewTicketStore(ttl time.Duration) *TicketStore {
	return &TicketStore{
		ttl:     ttl,
		tickets: make(map[string]ticketEntry),
		now:     time.Now,
	}
}

func (store *TicketStore) Issue(subscriber Subscriber) (string, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}

	ticket := base64.RawURLEncoding.EncodeToString(raw[:])

	store.mu.Lock()
	defer store.mu.Unlock()
	store.tickets[ticket] = ticketEntry{
		subscriber: subscriber,
		expiresAt:  store.now().UTC().Add(store.ttl),
	}

	return ticket, nil
}

func (store *TicketStore) Consume(ticket string) (Subscriber, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	entry, exists := store.tickets[ticket]
	if !exists {
		return Subscriber{}, ErrInvalidTicket
	}
	delete(store.tickets, ticket)

	if !entry.expiresAt.After(store.now().UTC()) {
		return Subscriber{}, ErrInvalidTicket
	}

	return entry.subscriber, nil
}
