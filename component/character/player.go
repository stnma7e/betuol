package character

import (
	"fmt"

	"smig/component"
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
	case "interact":
		var arg string
		fmt.Scan(&arg)
		fmt.Println("\t",arg)
	default:
		fmt.Println("\tInvalid command. Type \"help\" for choices.")
	}
}

func PlayerLook(id component.GOiD, chars *CharacterManager) {
	loc := chars.Scene.GetObjectLocation(id)
	ros := chars.attributeList[RANGEOFSIGHT][id]
	stk := chars.Scene.GetObjectsInLocationRange(loc, ros)
	numObj := stk.Size
	for i := 0; i < numObj; i ++ {
		charId, err := stk.Dequeue()
		if err != nil {
			common.Log.Warn(err)
		}

		if charId == int(id) || id == 0 {
			continue
		}

		ca := chars.GetCharacterAttributes(component.GOiD(charId))
		fmt.Println("\t",ca.Greet())
	}
}

func PlayerMove(direction string, id component.GOiD, chars *CharacterManager) {
	transMat := chars.Scene.GetTransformMatrix(id)

	switch direction {
	case "north":
		transMat[7]++
	case "south":
		transMat[7]--
	case "east":
		transMat[3]++
	case "west":
		transMat[3]--
	}
	chars.Scene.Transform(id, transMat)
}