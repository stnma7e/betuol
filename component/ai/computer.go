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
type AiComputer func(id component.GOiD, eventlink chan event.Event)

// AI function for the player.
// It allows players to be treated as normal characters, but they have a special AI function that checks for input.
func (am *AiManager) PlayerDecide(id component.GOiD) {
	ca := am.cm.GetCharacterAttributes(id)
	loc, err := am.tm.GetObjectLocation(id)
	if err != nil {
		common.LogErr.Printf("error in playerDecide, error: %s", err.Error())
	}
	fmt.Print(ca.Attributes[character.HEALTH], ca.Attributes[character.MANA], loc)

	fmt.Print(" --> ")
	var command string
	fmt.Scan(&command)
	character.ParsePlayerCommand(command, id, am.cm)
}

// AI function for a potential enemy.
// EnemyDecide looks for entities in close proximity and checks the health of these entities. If the health of the entity is lower than the health of the ai component, then the component initiates an attack with the character system.
func (am *AiManager) EnemyDecide(id component.GOiD) {
	this := am.cm.GetCharacterAttributes(id)
	loc, err := am.tm.GetObjectLocation(id)
	if err != nil {
		common.LogErr.Printf("error in enemyDecide, error: %s", err.Error())
	}
	idQueue := am.tm.GetObjectsInLocationRadius(loc, this.Attributes[character.RANGEOFSIGHT])
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
		attr = am.cm.GetCharacterAttributes(neighbors[i])
		if neighbors[i] == id ||
			attr.Attributes[character.HEALTH] >= this.Attributes[character.HEALTH] ||
			attr.Attributes[character.HEALTH] <= 0 {
			continue
		}
		am.em.Send(event.AttackEvent{id, component.GOiD(neighbors[i])})
		fmt.Println("enemy", id, "attacks", neighbors[i], neighbors)
		break
	}
}

// AI function for random movement.
// WanderDecide wanders in random directions each time called.
func (am *AiManager) WanderDecide(id component.GOiD) {
	r := rand.Int31()
	loc, err := am.tm.GetObjectLocation(id)
	if err != nil {
		common.LogErr.Printf("error in wanderDecide, error: %s", err.Error())
	}
	switch r % 4 {
	case 0:
		loc[0] += 5
		am.tm.SetLocationOverTime(id, loc, 2.0)
	case 1:
		loc[0] += 5
		am.tm.SetLocationOverTime(id, loc, 2.0)
	case 2:
		loc[1] += 5
		am.tm.SetLocationOverTime(id, loc, 2.0)
	case 3:
		loc[2] += 5
		am.tm.SetLocationOverTime(id, loc, 2.0)
	default:
		common.LogErr.Print("this don't work shit")
	}
}

// EnemyAi implements the AiComputer function template.
// It handles the event response and registers for events to respond to.
func (am *AiManager) EnemyAi(id component.GOiD, eventlink chan event.Event) {
	am.em.RegisterListeningChannel("attack", eventlink)
	for alive := true; alive; {
		evt := <-eventlink
		switch evt.GetEventType() {
		case "death":
			alive = false
		case "runAi":
			am.EnemyDecide(id)
			eventlink <- event.RunAiEvent{}
		case "attack":
			aevt := evt.(event.AttackEvent)
			if aevt.Char2 == id {
				am.em.Send(event.AttackEvent{id, aevt.Char1})
			}
		}
	}
}

// PlayerAi implements the AiComputer function template.
// It handles the event response and registers for events to respond to.
func (am *AiManager) PlayerAi(id component.GOiD, eventlink chan event.Event) {
	for alive := true; alive; {
		evt := <-eventlink
		switch evt.GetEventType() {
		case "death":
			alive = false
		case "runAi":
			am.PlayerDecide(id)
			eventlink <- event.RunAiEvent{}
		}
	}
}

// WanderAi implements the AiComputer function template.
// It handles the event response and registers for events to respond to.
func (am *AiManager) WanderAi(id component.GOiD, eventlink chan event.Event) {
	for alive := true; alive; {
		evt := <-eventlink
		switch evt.GetEventType() {
		case "death":
			alive = false
		case "runAi":
			am.WanderDecide(id)
			eventlink <- event.RunAiEvent{}
		}
	}
}
