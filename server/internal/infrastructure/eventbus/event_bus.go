package eventbus

import (
	appdto "dalkak/pkg/dto/app"
	"sync"
)

type Event struct {
	Type         string
	UserInfo		 *appdto.UserInfo
	Payload      interface{}
	TxID         string
	ResponseChan chan<- appdto.Response
}

type EventHandler func(Event)

type EventBus struct {
	listeners map[string][]EventHandler
	lock      sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		listeners: make(map[string][]EventHandler),
	}
}

func (bus *EventBus) Subscribe(eventType string, handler EventHandler) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	bus.listeners[eventType] = append(bus.listeners[eventType], handler)
}

func (bus *EventBus) Publish(event Event) {
	bus.lock.RLock()
	defer bus.lock.RUnlock()
	if handlers, found := bus.listeners[event.Type]; found {
		for _, handler := range handlers {
			handler(event)
		}
	}
}
