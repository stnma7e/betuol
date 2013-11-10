package event

import (
	"time"

	"smig/common"
)

type EventListener func(evt Event)

type Event interface {
	GetEventType() string
}
type EventMessage struct {
	time time.Time
	evt Event
}

type EventManager struct {
	eventlink   chan EventMessage
	listenerMap map[string]*common.Queue
	eventList   [2]*common.Vector
	changeQueue bool
}

func MakeEventManager() *EventManager {
	em := &EventManager{
		make(chan EventMessage),
		make(map[string]*common.Queue),
		[2]*common.Vector{common.MakeVector(), common.MakeVector()},
		false,
	}
	go em.Sort()
	return em
}

func (em *EventManager) Sort() {
	// need to implement some sorting scheme to order events by the time that they were created
	for {
		evtMsg := <-em.eventlink
		if em.changeQueue {
			em.eventList[0].Insert(evtMsg.evt)
		} else {
			em.eventList[1].Insert(evtMsg.evt)
		}
	}
}

func (em *EventManager) Tick(delta float64) {
	var events []interface{}
	if em.changeQueue {
		em.changeQueue = false
		events = em.eventList[0].Array()
		em.eventList[0].Empty()
	} else {
		em.changeQueue = true
		events = em.eventList[1].Array()
		em.eventList[1].Empty()
	}
	for i := range events {
		evt := events[i].(Event)
		listeners, ok := em.listenerMap[evt.GetEventType()]
		if !ok {
			common.LogErr.Printf("no listener registered for %s", evt.GetEventType())
			break
		}
		listenersArray := listeners.Array()
		for i := range listenersArray {
			listenersArray[i].(EventListener)(evt)
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
		em.eventlink <- EventMessage{time.Now(), evt}
	}()
}
