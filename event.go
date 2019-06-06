package event

import (
	"fmt"
	"sync"
)

type Event interface {
	Name() string
	//	correlationId() string
}

type subscription struct {
	name   string
	filter func(Event) bool
	c      chan Event
}

type bus struct {
	subscribers map[string]*subscription
	mutex       sync.Mutex
}

func NewEventBus() *bus {
	return &bus{
		subscribers: make(map[string]*subscription),
	}
}

func (e *bus) Publish(evt Event) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, v := range e.subscribers {
		if v.filter == nil || v.filter(evt) {
			go func(v *subscription) {
				v.c <- evt
			}(v)
		}
	}
}

func (e *bus) Subscribe(name string, filter subscriptionFilter) (chan Event, error) {
	_, subscriptionExists := e.subscribers[name]

	if subscriptionExists {
		return nil, fmt.Errorf("subscription %s already exists", name)
	}

	e.mutex.Lock()
	e.subscribers[name] = &subscription{
		name:   name,
		filter: filter,
		c:      make(chan Event),
	}
	e.mutex.Unlock()

	return e.subscribers[name].c, nil
}

func (e *bus) Unsubscribe(name string) bool {
	sub, ok := e.subscribers[name]
	if ok {
		close(sub.c)
		for range sub.c {
		}
	}
	// TODO Drain subscription channel?

	e.mutex.Lock()
	delete(e.subscribers, name)
	e.mutex.Unlock()

	return ok
}
