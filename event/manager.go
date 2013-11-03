package event

import (
	"smig/common"
)

type EventListener func(evt Event)

type Event interface {
	GetEventType() string
}

type EventManager struct {
	eventlink   chan Event
	listenerMap map[string]*common.Queue
}

func MakeEventManager() *EventManager {
	return &EventManager {
		make(chan Event),
		make(map[string]*common.Queue),
	}
}

func (em *EventManager) Tick(delta float64) {
    for {
	select {
	case evt := <-em.eventlink:
		listeners, ok := em.listenerMap[evt.GetEventType()]
		if !ok {
			common.LogErr.Printf("no listener registered for %s", evt.GetEventType())
		}
		listenersArray := listeners.Array()
		for i := range listenersArray {
			listenersArray[i].(EventListener)(evt)
		}
	default:
            return
	}
    }
}

func (em *EventManager) RegisterListener(eventType string, listener EventListener) {
        _, ok := em.listenerMap[eventType]
        if !ok {
            em.listenerMap[eventType] = &common.Queue{}
        }
        if em.listenerMap[eventType] == nil {
            em.listenerMap[eventType] = &common.Queue{}
        }
	em.listenerMap[eventType].Queue(listener)
}

func (em *EventManager) Send(evt Event) {
	go func() {
		em.eventlink <- evt
	}()
}
