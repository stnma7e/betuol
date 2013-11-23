// Package event implements an event management system.
package event

import (
	"time"

	"betuol/common"
)

type EventListener func(evt Event)

// Event is an interface to allow many different structs to be used as events as long as they specify their type.
type Event interface {
	GetEventType() string
}

// EventMessage is a helper struct for keeping track of the time that an event was added to the queue.
type EventMessage struct {
	time time.Time
	evt  Event
}

// EventManager handles incoming events and stores them until they are dispersed according to registered listeners.
type EventManager struct {
	eventlink         chan EventMessage
	listenerMap       map[string]*common.Vector
	listeningChannels map[string]*common.Vector
	eventList         [2]*common.Vector
	changeQueue       bool
}

// MakeEventManager returns a pointer to an EventManager.
func MakeEventManager() *EventManager {
	em := &EventManager{
		make(chan EventMessage),
		make(map[string]*common.Vector),
		make(map[string]*common.Vector),
		[2]*common.Vector{common.MakeVector(), common.MakeVector()},
		false,
	}
	go em.Sort()
	return em
}

// Sort will recieve events and sort them based on time of arrival.
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

// Tick is used to dispatch events to listeners that registered to an event type.
// Events are dispatched to both function listeners and to channels during the same tick.
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
			common.LogWarn.Printf("no listener registered for %s", evt.GetEventType())
		} else {
			listenersArray := listeners.Array()
			for j := range listenersArray {
				listenersArray[j].(EventListener)(evt)
			}
		}

		channels, ok := em.listeningChannels[evt.GetEventType()]
		if !ok {
			//common.LogWarn.Printf("no channel registered for %s", evt.GetEventType())
		} else {
			channelsArray := channels.Array()
			for j := range channelsArray {
				channelsArray[j].(chan Event) <- evt
			}
		}
	}
}

// RegisterListener registers a listening function to be called every time an event of type eventType is processed.
func (em *EventManager) RegisterListener(eventType string, listener EventListener) {
	_, ok := em.listenerMap[eventType]
	if !ok {
		em.listenerMap[eventType] = common.MakeVector()
	}
	if em.listenerMap[eventType] == nil {
		em.listenerMap[eventType] = common.MakeVector()
	}
	em.listenerMap[eventType].Insert(listener)
}

// RegisterListeningChannel registers a listening channel to be sent the event every time an event of type eventType is processed.
func (em *EventManager) RegisterListeningChannel(eventType string, eventlink chan Event) {
	_, ok := em.listeningChannels[eventType]
	if !ok {
		em.listeningChannels[eventType] = common.MakeVector()
	}
	if em.listeningChannels[eventType] == nil {
		em.listeningChannels[eventType] = common.MakeVector()
	}
	em.listeningChannels[eventType].Insert(eventlink)
}

// Send will add an event passed as an argument to the event queue to be processed.
func (em *EventManager) Send(evt Event) {
	go func() {
		em.eventlink <- EventMessage{time.Now(), evt}
	}()
}
