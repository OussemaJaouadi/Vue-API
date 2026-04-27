package events

import (
	"context"
	"sync"
)

type Event struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

type Subscriber struct {
	UserID     string
	GlobalRole string
}

type Subscription interface {
	Events() <-chan Event
	Close()
}

type Broker struct {
	mu          sync.RWMutex
	subscribers map[*memorySubscription]Subscriber
}

type memorySubscription struct {
	broker *Broker
	events chan Event
	once   sync.Once
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[*memorySubscription]Subscriber),
	}
}

func (broker *Broker) Subscribe(ctx context.Context, subscriber Subscriber) Subscription {
	sub := &memorySubscription{
		broker: broker,
		events: make(chan Event, 16),
	}

	broker.mu.Lock()
	broker.subscribers[sub] = subscriber
	broker.mu.Unlock()

	go func() {
		<-ctx.Done()
		sub.Close()
	}()

	return sub
}

func (broker *Broker) PublishToManagers(event Event) {
	broker.publish(func(subscriber Subscriber) bool {
		return subscriber.GlobalRole == "manager"
	}, event)
}

func (broker *Broker) PublishToUser(userID string, event Event) {
	broker.publish(func(subscriber Subscriber) bool {
		return subscriber.UserID == userID
	}, event)
}

func (broker *Broker) publish(matches func(Subscriber) bool, event Event) {
	broker.mu.RLock()
	defer broker.mu.RUnlock()

	for sub, subscriber := range broker.subscribers {
		if !matches(subscriber) {
			continue
		}

		select {
		case sub.events <- event:
		default:
		}
	}
}

func (sub *memorySubscription) Events() <-chan Event {
	return sub.events
}

func (sub *memorySubscription) Close() {
	sub.once.Do(func() {
		sub.broker.mu.Lock()
		delete(sub.broker.subscribers, sub)
		sub.broker.mu.Unlock()
		close(sub.events)
	})
}
