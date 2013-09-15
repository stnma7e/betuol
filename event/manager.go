package event

import (
	"smig/common"
)

type EventListener func(evt *Event)

type Event struct {
	EventType string
}

type EventManager struct {
	eventlink 	chan *Event
	listenerMap map[string]EventListener
}

func MakeEventManager() *EventManager {
	return &EventManager {
		make(chan *Event),
		make(map[string]EventListener),
	}
}

func (em *EventManager) Tick(delta float64) {
	select {
	case evt := <-em.eventlink:
		listener, ok := em.listenerMap[evt.EventType]
		if !ok {
			common.Log.Error("no listener registered for %s", evt.EventType)
		}
		listener(evt)
	default:
	}
}

func (em *EventManager) RegisterListener(eventType string, listener EventListener) {
	em.listenerMap[eventType] = listener
}

func (em *EventManager) Send(evt *Event) {
	go func() {
		em.eventlink <- evt
	}()
}