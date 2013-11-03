package ai

import (
	"fmt"
	
	"smig/component"
	"smig/component/character"
	"smig/event"
)

type AiComputer func(id component.GOiD, neighbors []component.GOiD, chars *character.CharacterManager)

func (am *AiManager) PlayerDecide(id component.GOiD, neighbors []component.GOiD, chars *character.CharacterManager) {
	ca := chars.GetCharacterAttributes(id)
	loc :=  chars.Scene.GetObjectLocation(id)
	fmt.Print(ca.Attributes[character.HEALTH], ca.Attributes[character.MANA], loc[:2], " --> ")

	var command string
	fmt.Scan(&command)
	character.ParsePlayerCommand(command, id, chars)
}

func (am *AiManager) EnemyDecide(id component.GOiD, neighbors []component.GOiD, chars *character.CharacterManager) {
	fmt.Println("enemy: ", neighbors)
	this := chars.GetCharacterAttributes(id)
	var attr *character.CharacterAttributes
	var i int
	for i = 0; i < len(neighbors); i++ {
		attr = chars.GetCharacterAttributes(neighbors[i])
		if attr.Faction == this.Faction {
			continue
		}
		if attr.Attributes[character.HEALTH] > this.Attributes[character.HEALTH] {
			continue
		}
	}
	am.em.Send(event.AttackEvent{ id, component.GOiD(i) })
}

