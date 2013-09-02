package server

import (
	"time"
	"fmt"
	"strings"
	"strconv"

	"smig/component"
	"smig/component/physics"
	"smig/component/character"
	"smig/res"
	"smig/common"
	"smig/math"
)

type Server struct {
	gof *component.GameObjectFactory
	sm  *component.SceneManager
	pm  *physics.PhysicsManager
	cm  *character.CharacterManager
	rm  *res.ResourceManager

	returnlink chan bool
}

func MakeServer(returnlink chan bool) *Server {
	sm 	:= component.MakeSceneManager()
	sv := &Server {
		component.MakeGameObjectFactory(sm),
		sm, 
		physics.MakePhysicsManager(sm),
		character.MakeCharacterManager(sm),
		res.MakeResourceManager("/home/sam/go/src/smig/data/"),
		returnlink,
	}

	sv.gof.Register("physics", sv.pm, sv.pm.JsonCreate)
	sv.gof.Register("character", sv.cm, sv.cm.JsonCreate)

	sv.cm.RegisterComputer("merchant", character.MerchantComputer)
	sv.cm.RegisterComputer("passive", character.PassiveComputer)
	sv.cm.RegisterComputer("player", character.PlayerComputer)

	return sv
}

func (sv *Server) Loop() {
	defer sv.Shutdown()
	oldTime := time.Now()
	ticks := time.NewTicker(time.Second / 60)
	for {
		<-ticks.C

		newTime := time.Since(oldTime)
		secs := newTime.Seconds()

		// fmt.Println(newTime)

		sv.ParseSysConsole()
		sv.cm.Tick(secs)
		sv.pm.Tick(secs)
		sv.sm.Tick(secs)

		// for i := range list {
		// 	id := list[i]
		// 	trans,_ := tm.GetTransform(component.GOiD(id))
		// 	fmt.Println(id," ",trans.ToString())
		// }
		// fmt.Println()

		oldTime = time.Now()
	}
}

func (sv *Server) ParseSysConsole() {
	fmt.Print("> ")
	var command string
	fmt.Scan(&command)
	switch command {
	case "exit":
		sv.returnlink <- true
	case "loadmap":
		var arg string
		fmt.Scan(&arg)
		sv.CreateFromMap(arg)
	case "loadobj":
		var breed, location string
		var radius float32
		fmt.Scan(&breed, &location, &radius)
		sv.CreateObject(breed, location, radius)
	case "runai":
		var arg component.GOiD
		fmt.Scan(&arg)
		sv.cm.RunAi(arg)
	default:
		fmt.Println("\tInvalid command. Type \"help\" for choices.")
	}
}

func (sv *Server) Shutdown() {

}

/*****************************************
*
* Component
*
*****************************************/

func (sv *Server) CreateFromMap(mapName string) {
	jmap := sv.rm.LoadJsonMap(mapName)
	sv.gof.CreateFromMap(&jmap)
}

func (sv *Server) CreateObject(objName, location string, radius float32) {
	components := sv.rm.LoadGameObject(objName)
	strLoc := strings.Split(location,",")
	f1, err := strconv.ParseFloat(strLoc[0], 32)
	f2, err := strconv.ParseFloat(strLoc[1], 32)
	f3, err := strconv.ParseFloat(strLoc[2], 32)
	if err != nil {
		fmt.Println(err)
	}
	id, err := sv.gof.Create(components, math.Vec3{float32(f1),float32(f2),float32(f3)}, radius)
	if err != nil {
		common.Log.Error(err)
	}
	fmt.Println("\tid:",id)
}