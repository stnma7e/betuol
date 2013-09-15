package ai

import (
	"fmt"
	"math/rand"
	
	"smig/component"
	"smig/component/character"
)

type AiComputer func(id component.GOiD, neighbors []component.GOiD, chars *character.CharacterManager)


func (am *AiManager) Attack(id1, id2 component.GOiD) float32 {
	attr1 := am.cm.GetCharacterAttributes(id1)
	attr2 := am.cm.GetCharacterAttributes(id2)
	hit := (rand.Float32() / 4) * attr1.Attributes[character.STRENGTH]
	attr2.Attributes[character.HEALTH] -= hit
	return hit
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
	fmt.Println(am.Attack(id, component.GOiD(i)))
}

func (am *AiManager) PlayerDecide(id component.GOiD, neighbors []component.GOiD, chars *character.CharacterManager) {
	ca := chars.GetCharacterAttributes(id)
	loc :=  chars.Scene.GetObjectLocation(id)
	fmt.Print(ca.Attributes[character.HEALTH], ca.Attributes[character.MANA], loc[:2], " --> ")

	var command string
	fmt.Scan(&command)
	character.ParsePlayerCommand(command, id, chars)
}