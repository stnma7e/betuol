package ai

import (
	"fmt"
	"math/rand"

	"smig/common"
	"smig/component"
	"smig/component/character"
	"smig/event"
)

//type AiComputer func(id component.GOiD, neighbors []component.GOiD, chars *character.CharacterManager)
type AiComputer func(id component.GOiD, eventlink chan event.Event)

func (am *AiManager) PlayerDecide(id component.GOiD) {
	ca := am.cm.GetCharacterAttributes(id)
	loc := am.tm.GetObjectLocation(id)
	fmt.Print(ca.Attributes[character.HEALTH], ca.Attributes[character.MANA], loc)

	fmt.Print(" --> ")
	var command string
	fmt.Scan(&command)
	character.ParsePlayerCommand(command, id, am.cm)
}

func (am *AiManager) EnemyDecide(id component.GOiD) {
	this := am.cm.GetCharacterAttributes(id)
	loc := am.tm.GetObjectLocation(id)
	idQueue := am.tm.GetObjectsInLocationRadius(loc, this.Attributes[character.RANGEOFSIGHT])
	size := idQueue.Size
	neighbors := make([]component.GOiD, size)
	for i := 0; i < size; i++ {
		val, err := idQueue.Dequeue()
		if err != nil {
			common.LogErr.Println("error: bad dequeue:", err)
			continue
		}
		neighbors[i] = component.GOiD(val)
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

func (am *AiManager) WanderDecide(id component.GOiD) {
	r := rand.Int31()
	switch r % 4 {
	case 0:
		loc := am.tm.GetObjectLocation(id)
		newLoc := loc
		newLoc[0] += 5
		am.tm.SetLocationOverTime(id, newLoc, 2.0)
	case 1:
		loc := am.tm.GetObjectLocation(id)
		newLoc := loc
		newLoc[0] += 5
		am.tm.SetLocationOverTime(id, newLoc, 2.0)
	case 2:
		loc := am.tm.GetObjectLocation(id)
		newLoc := loc
		newLoc[1] += 5
		am.tm.SetLocationOverTime(id, newLoc, 2.0)
	case 3:
		loc := am.tm.GetObjectLocation(id)
		newLoc := loc
		newLoc[2] += 5
		am.tm.SetLocationOverTime(id, newLoc, 2.0)
	default:
		common.LogErr.Print("this dont work shit")
	}
}

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
