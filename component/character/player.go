package character

import (
	"fmt"
	"strconv"

	"smig/common"
	"smig/component"
	"smig/event"
	"smig/math"
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
		fmt.Printf("\tInvalid player command, '%s'.\n", command)
	}
}

func PlayerLook(id component.GOiD, chars *CharacterManager) {
	loc := chars.sm.GetObjectLocation(id)
	ros := chars.attributeList[RANGEOFSIGHT][id]
	stk := chars.sm.GetObjectsInLocationRadius(loc, ros)
	numObj := stk.Size
	for i := 0; i < numObj; i++ {
		charId, err := stk.Dequeue()
		if err != nil {
			common.LogErr.Print(err)
		}

		if charId == int(id) || id == 0 {
			continue
		}

		ca := chars.GetCharacterAttributes(component.GOiD(charId))
		if ca.Description != "" {
			fmt.Println("\t", ca.Greet())
		}
	}
}

func PlayerMove(direction string, id component.GOiD, chars *CharacterManager) {
	transMat, err := chars.sm.GetTransformMatrix(id)
	if err != nil {
		common.LogErr.Println(err)
	}

	switch direction {
	case "north":
		chars.sm.SetLocationOverTime(id, math.Vec3{transMat[3], transMat[7], transMat[11] + 1}, 3)
	case "south":
		chars.sm.SetLocationOverTime(id, math.Vec3{transMat[3], transMat[7], transMat[11] - 1}, 3)
	case "east":
		chars.sm.SetLocationOverTime(id, math.Vec3{transMat[3] + 1, transMat[7], transMat[11]}, 3)
	case "west":
		chars.sm.SetLocationOverTime(id, math.Vec3{transMat[3] - 1, transMat[7], transMat[11]}, 3)
	}
	chars.sm.SetTransform(id, transMat)
}

func PlayerAttack(player, enemy component.GOiD, chars *CharacterManager) {
	chars.em.Send(event.AttackEvent{player, enemy})
}
