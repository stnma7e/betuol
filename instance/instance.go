package instance

import (
	"time"
	"math/rand"
	"fmt"
	"strings"
	"strconv"

	"smig/component"
	"smig/component/physics"
	"smig/component/character"
	"smig/component/ai"
	"smig/res"
	"smig/common"
	"smig/math"
	"smig/event"
)

type Instance struct {
	gof *component.GameObjectFactory
	sm  *component.SceneManager
	pm  *physics.PhysicsManager
	cm  *character.CharacterManager
	rm  *res.ResourceManager
	em  *event.EventManager
	am  *ai.AiManager

	returnlink  chan bool
	commandlink chan string
}

func MakeInstance(returnlink chan bool, rm *res.ResourceManager) *Instance {
	sm 	:= component.MakeSceneManager()
	gof := component.MakeGameObjectFactory(sm)
	em  := event.MakeEventManager()
	cm 	:= character.MakeCharacterManager(sm)
	is  := &Instance {
		gof,
		sm,
		physics.MakePhysicsManager(sm),
		cm,
		rm,
		em,
		ai.MakeAiManager(sm, cm),
		returnlink,
		make(chan string),
	}

	is.gof.Register("physics", is.pm, is.pm.JsonCreate)
	is.gof.Register("character", is.cm, is.cm.JsonCreate)
	is.gof.Register("ai", is.am, is.am.JsonCreate)

	is.am.RegisterComputer("enemy", is.am.EnemyDecide)
	is.am.RegisterComputer("player", is.am.PlayerDecide)

	is.em.RegisterListener("string", func(evt *event.Event) {
		fmt.Println(evt.EventType)
	})

	rand.Seed(time.Now().UnixNano())

	return is
}

func (is *Instance) Loop() {
	defer is.Shutdown()
	oldTime := time.Now()
	ticks := time.NewTicker(time.Second / 60)

	go func() {
		for {
			fmt.Print("> ")
			var command string
			fmt.Scan(&command)
			is.commandlink <- command
		}
	}()

	for {
		<-ticks.C

		newTime := time.Since(oldTime)
		secs := newTime.Seconds()

		// fmt.Println(newTime)

		is.ParseSysConsole()
		is.em.Tick(secs)
		is.am.Tick(secs)
		is.cm.Tick(secs)
		is.pm.Tick(secs)
		is.sm.Tick(secs)

		// for i := range list {
		// 	id := list[i]
		// 	trans,_ := tm.GetTransform(component.GOiD(id))
		// 	fmt.Println(id," ",trans.ToString())
		// }
		// fmt.Println()

		oldTime = time.Now()
	}
}

func (is *Instance) ParseSysConsole() {
	select {
	case command := <-is.commandlink:
		switch command {
		case "exit":
			is.returnlink <- true
		case "loadmap":
			var arg string
			fmt.Scan(&arg)
			is.CreateFromMap(arg)
		case "loadobj":
			var breed, location string
			var radius float32
			fmt.Scan(&breed, &location, &radius)
			is.CreateObject(breed, location, radius)
		case "runai":
			var arg component.GOiD
			fmt.Scan(&arg)
			is.am.RunAi(arg)
		default:
			fmt.Println("\tInvalid command. Type \"help\" for choices.")
		}
	default:
	}
	
}

func (is *Instance) Shutdown() {

}

/*****************************************
*
* Component
*
*****************************************/

func (is *Instance) CreateFromMap(mapName string) {
	jmap := is.rm.LoadJsonMap(mapName)
	is.gof.CreateFromMap(&jmap)
}

func (is *Instance) CreateObject(objName, location string, radius float32) {
	components := is.rm.LoadGameObject(objName)
	strLoc := strings.Split(location,",")
	f1, err := strconv.ParseFloat(strLoc[0], 32)
	f2, err := strconv.ParseFloat(strLoc[1], 32)
	f3, err := strconv.ParseFloat(strLoc[2], 32)
	if err != nil {
		fmt.Println(err)
	}
	id, err := is.gof.Create(components, math.Vec3{float32(f1),float32(f2),float32(f3)}, radius)
	if err != nil {
		common.Log.Error(err)
	}
	fmt.Println("\tid:",id)
}