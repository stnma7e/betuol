package character

import (
	"fmt"
	"strconv"

	"smig/component"
	"smig/event"
	"smig/common"
)

const (
	INTERACTING = 1 << iota
)

func ParsePlayerCommand(command string, id component.GOiD, chars *CharacterManager) {
	switch command {
	case "look":
		PlayerLook(id, chars)
	case "north":
		PlayerMove("north", id, chars)
	case "south":
		PlayerMove("south", id, chars)
	case "east":
		PlayerMove("east", id, chars)
	case "west":
		PlayerMove("west", id, chars)
	case "attack":
		var arg string
		fmt.Scan(&arg)
		enemy, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println("invalid enemy id")
			break
		}
		PlayerAttack(id, component.GOiD(enemy), chars)
	default:
		fmt.Println("\tInvalid command. Type \"help\" for choices.")
	}
}

func PlayerLook(id component.GOiD, chars *CharacterManager) {
	loc := chars.Scene.GetObjectLocation(id)
	ros := chars.attributeList[RANGEOFSIGHT][id]
	stk := chars.Scene.GetObjectsInLocationRadius(loc, ros)
	numObj := stk.Size
	for i := 0; i < numObj; i ++ {
		charId, err := stk.Dequeue()
		if err != nil {
			common.LogErr.Print(err)
		}

		if charId == int(id) || id == 0 {
			continue
		}

		ca := chars.GetCharacterAttributes(component.GOiD(charId))
                if ca.Description != "" {
		    fmt.Println("\t",ca.Greet())
                }
	}
}

func PlayerMove(direction string, id component.GOiD, chars *CharacterManager) {
	transMat := chars.Scene.GetTransformMatrix(id)

	switch direction {
	case "north":
		transMat[11]++
	case "south":
		transMat[11]--
	case "east":
		transMat[3]++
	case "west":
		transMat[3]-- }
	chars.Scene.SetTransform(id, transMat)
}

func PlayerAttack(player, enemy component.GOiD, chars *CharacterManager) {
	chars.em.Send(event.AttackEvent{ player, enemy })
}
