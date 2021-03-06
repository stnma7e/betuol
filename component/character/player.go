package character

import (
	"fmt"
	"strconv"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/event"
	"github.com/stnma7e/betuol/math"
)

// ParsePlayerCommand is called by the ai component manager to parse a command input into the ai manager.
// This function parses the command and responds accordingly.
func ParsePlayerCommand(command string, id component.GOiD, chars *CharacterManager) {
	switch command {
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

// PlayerMove moves the player in one of four cardinal directions within the game.
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

// PlayerAttack sends an attack event that will be parsed by the charater manager to afflict damage to another character component.
func PlayerAttack(player, enemy component.GOiD, chars *CharacterManager) {
	chars.em.Send(event.AttackEvent{player, enemy})
}
