package character

import (
	"fmt"

	"smig/component/transform"
	"smig/component"
	"smig/math"
	"smig/common"
)

type Player struct {
	Id component.GOiD
	scene *transform.SceneManager
	chars *CharacterManager
}

func StartPlayer(id component.GOiD, tm *transform.SceneManager, cm *CharacterManager) {
	pl := Player{}
	pl.Id = id
	pl.scene = tm
	pl.chars = cm

	pl.ProcessCommandLine()
}

func (pl *Player) ProcessCommandLine() {
	for {
		var com string
		ca := pl.chars.GetCharacterAttributes(pl.Id)
		locMat, err :=  pl.scene.GetTransform(pl.Id)
		if err != nil {
			common.Log.Error(err)
		}
		loc := math.Mult(math.Vec3{}, locMat)
		fmt.Print(ca.Attributes[HEALTH], ca.Attributes[MANA], loc[:2], " --> ")
		fmt.Scan(&com)
		pl.ParseCommand(com)
	}
}

func (pl *Player) ParseCommand(com string) {
	switch com {
	case "help":

	case "look":
		pl.Look()
	case "north":
		pl.Move("north")
	case "south":
		pl.Move("south")
	case "east":
		pl.Move("east")
	case "west":
		pl.Move("west")
	default:
		Println("Invalid command. Type help for choices.")
	}
}

func (pl *Player) Look() {
	locMat, err :=  pl.scene.GetTransform(pl.Id)
	if err != nil {
		common.Log.Error(err)
	}
	loc := math.Mult(math.Vec3{}, locMat)
	stk := pl.scene.GetObjectsInLocationRange(loc, 10.0)
	numObj := stk.Size
	for i := 0; i < numObj; i ++ {
		id, err := stk.Dequeue()
		if err != nil {
			common.Log.Warn(err)
		}

		if id == int(pl.Id) || id == 0 {
			continue
		}

		ca := pl.chars.GetCharacterAttributes(component.GOiD(id))
		Println(ca.Greet())
	}
}

func (pl *Player) Move(direction string) {
	transMat, err := pl.scene.GetTransform(pl.Id)
	if err != nil {
		common.Log.Error(err)
	}

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
	pl.scene.Transform(pl.Id, transMat)
}


func Println(s string) {
	fmt.Printf("\t%s\n", s)
}