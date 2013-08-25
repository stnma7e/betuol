package ai

import (
	"time"

	"smig/component"
)

type aiEventManager struct {
	eventlink chan component.Event
	indexlink chan component.GOiD
	errorlink chan error
}

func MakeAiManager() *aiEventManager {
	am := aiManager{}
	ameg := &aiEventManager {
		make(chan component.Event),
		make(chan component.GOiD),
		make(chan error),
	}

	go am.HandleEvents(ameg)
	return ameg
}

func (am *aiManager) HandleEvents(ameg *aiEventManager) {
	for {
		evt := <-ameg.eventlink
		switch evt.EventType {
		case "create":
			ameg.errorlink <- nil
		}
		am.tick(1.0)
	}
}

func (ameg *aiEventManager) CreateComponent(index component.GOiD) error {
	ameg.eventlink <- component.Event {time.Now(), "create"}
	ameg.indexlink <- index
	return <-ameg.errorlink
}
func (ameg *aiEventManager) DeleteComponent(index component.GOiD) {
	ameg.eventlink <- component.Event {time.Now(), "delete"}
	ameg.indexlink <- index
}