package character

import (
	"fmt"

	"smig/component/scene"
	"smig/component"
	"smig/common"
)

const (
	INTERACTING = 1 << iota
)

type Player struct {
	Id component.GOiD
	Scene *scene.SceneManager
	Chars *CharacterManager
	RangeOfSight float32
	movedlink chan Player

	stateMask int
}

func StartPlayer(id component.GOiD, lookRange float32, movedlink chan Player, sm *scene.SceneManager, cm *CharacterManager) {
	pl := Player{}
	pl.Id = id
	pl.RangeOfSight = lookRange
	pl.movedlink = movedlink
	pl.Scene = sm
	pl.Chars = cm

	pl.ProcessCommandLine()
}

func (pl *Player) ProcessCommandLine() {
	for {
		var com string
		ca := pl.Chars.GetCharacterAttributes(pl.Id)
		loc :=  pl.Scene.GetObjectLocation(pl.Id)
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
	case "interact":
		var arg string
		fmt.Scan(&arg)
		Println(arg)
		pl.stateMask = INTERACTING
	default:
		Println("Invalid command. Type \"help\" for choices.")
	}
}

func (pl *Player) Look() {
	loc :=  pl.Scene.GetObjectLocation(pl.Id)
	stk := pl.Scene.GetObjectsInLocationRange(loc, pl.RangeOfSight)
	numObj := stk.Size
	for i := 0; i < numObj; i ++ {
		id, err := stk.Dequeue()
		if err != nil {
			common.Log.Warn(err)
		}

		if id == int(pl.Id) || id == 0 {
			continue
		}

		ca := pl.Chars.GetCharacterAttributes(component.GOiD(id))
		Println(ca.Greet())
	}
}

func (pl *Player) Move(direction string) {
	transMat := *pl.Scene.GetTransformPointer(pl.Id)

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
	pl.Scene.Transform(pl.Id, &transMat)

	pl.movedlink <- *pl
}


func Println(s string) {
	fmt.Printf("\t%s\n", s)
}