package ai

import (
	"fmt"
	"math/rand"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/component/character"
	"github.com/stnma7e/betuol/event"
)

// Each ai component is tied to a AiComputer that manages the AI for the component.
type AiComputer func(cm *character.CharacterManager, sm component.SceneManager, em *event.EventManager, id component.GOiD, eventlink chan event.Event)

// AI function for the player.
// It allows players to be treated as normal characters, but they have a special AI function that checks for input.
func PlayerDecide(cm *character.CharacterManager, sm component.SceneManager, id component.GOiD) {
	ca := cm.GetCharacterAttributes(id)
	loc, err := sm.GetObjectLocation(id)
	if err != nil {
		common.LogErr.Printf("error in playerDecide, error: %s", err.Error())
	}
	fmt.Print(ca.Attributes[character.HEALTH], ca.Attributes[character.MANA], loc)

	fmt.Print(" --> ")
	var command string
	fmt.Scan(&command)
	character.ParsePlayerCommand(command, id, cm)
}

// AI function for a potential enemy.
// EnemyDecide looks for entities in close proximity and checks the health of these entities. If the health of the entity is lower than the health of the ai component, then the component initiates an attack with the character system.
// EnemyDecide returns a vector of event.Event. These are sent when the function returns to EnemyAi.
func EnemyDecide(cm *character.CharacterManager, sm component.SceneManager, id component.GOiD) (eventlist *common.Vector) {
	eventlist = common.MakeVector()
	this := cm.GetCharacterAttributes(id)
	loc, err := sm.GetObjectLocation(id)
	if err != nil {
		common.LogErr.Printf("error in enemyDecide, error: %s", err.Error())
	}
	idQueue := sm.GetObjectsInLocationRadius(loc, this.Attributes[character.RANGEOFSIGHT])
	size := idQueue.Size
	neighbors := make([]component.GOiD, size)
	for i := 0; i < size; i++ {
		val, err := idQueue.Dequeue()
		if err != nil {
			common.LogErr.Println("error: bad dequeue:", err)
			continue
		}
		neighbors[i] = val.(component.GOiD)
	}
	var attr *character.CharacterAttributes
	for i := 0; i < len(neighbors); i++ {
		if neighbors[i] == id {
			continue
		}
		attr = cm.GetCharacterAttributes(neighbors[i])
		if neighbors[i] == id ||
			attr.Attributes[character.HEALTH] >= this.Attributes[character.HEALTH] ||
			attr.Attributes[character.HEALTH] <= 0 {
			continue
		}
		eventlist.Insert(event.AttackEvent{id, component.GOiD(neighbors[i])})
		fmt.Println("enemy", id, "attacks", neighbors[i], neighbors)
		break
	}

	return
}

// AI function for random movement.
// WanderDecide wanders in random directions each time called.
func WanderDecide(cm *character.CharacterManager, sm component.SceneManager, id component.GOiD) {
	r := rand.Int31()
	loc, err := sm.GetObjectLocation(id)
	if err != nil {
		common.LogErr.Printf("error in wanderDecide, error: %s", err.Error())
	}
	switch r % 4 {
	case 0:
		loc[0] += 5
		sm.SetLocationOverTime(id, loc, 2.0)
	case 1:
		loc[0] += 5
		sm.SetLocationOverTime(id, loc, 2.0)
	case 2:
		loc[1] += 5
		sm.SetLocationOverTime(id, loc, 2.0)
	case 3:
		loc[2] += 5
		sm.SetLocationOverTime(id, loc, 2.0)
	default:
		common.LogErr.Print("this don't work shit")
	}
}

// EnemyAi implements the AiComputer function template.
// It handles the event response and registers for events to respond to.
func EnemyAi(cm *character.CharacterManager, sm component.SceneManager, em *event.EventManager, id component.GOiD, eventlink chan event.Event) {
	em.RegisterListeningChannel("attack", eventlink)
	for alive := true; alive; {
		evt := <-eventlink
		switch evt.GetEventType() {
		case "death":
			alive = false
		case "runAi":
			eventsToSend := EnemyDecide(cm, sm, id).Array()
			for i := range eventsToSend {
				em.Send(eventsToSend[i].(event.Event))
			}
			eventlink <- event.RunAiEvent{}
		case "attack":
			aevt := evt.(event.AttackEvent)
			if aevt.Char2 == id {
				em.Send(event.AttackEvent{id, aevt.Char1})
			}
		}
	}
}

// PlayerAi implements the AiComputer function template.
// It handles the event response and registers for events to respond to.
func PlayerAi(cm *character.CharacterManager, sm component.SceneManager, em *event.EventManager, id component.GOiD, eventlink chan event.Event) {
	for alive := true; alive; {
		evt := <-eventlink
		switch evt.GetEventType() {
		case "death":
			alive = false
		case "runAi":
			PlayerDecide(cm, sm, id)
			eventlink <- event.RunAiEvent{}
		}
	}
}

// WanderAi implements the AiComputer function template.
// It handles the event response and registers for events to respond to.
func WanderAi(cm *character.CharacterManager, sm component.SceneManager, em *event.EventManager, id component.GOiD, eventlink chan event.Event) {
	for alive := true; alive; {
		evt := <-eventlink
		switch evt.GetEventType() {
		case "death":
			alive = false
		case "runAi":
			WanderDecide(cm, sm, id)
			eventlink <- event.RunAiEvent{}
		}
	}
}
