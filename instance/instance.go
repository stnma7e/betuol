package instance

import (
	"time"
	"math/rand"
	"fmt"
	"strings"
	"strconv"

	"smig/component"
	"smig/component/gofactory"
	"smig/component/physics"
	"smig/component/character"
	"smig/component/ai"
        "smig/component/quest"
	"smig/res"
	"smig/math"
	"smig/event"
	"smig/graphics"
        "smig/net"
	"smig/common"
)

type Instance struct {
	gof *gofactory.GameObjectFactory
	tm  *component.TransformManager
	pm  *physics.PhysicsManager
	cm  *character.CharacterManager
	rm  *res.ResourceManager
	em  *event.EventManager
	am  *ai.AiManager
        qm  *quest.QuestManager
	gm  *graphics.GraphicsManager
        nm  *net.NetworkManager

	returnlink  chan bool
	commandlink chan string

	player component.GOiD
}

func MakeInstance(returnlink chan bool, rm *res.ResourceManager, gm *graphics.GraphicsManager, nm *net.NetworkManager) *Instance {
	tm  := component.MakeTransformManager()
	gof := gofactory.MakeGameObjectFactory(tm)
	em  := event.MakeEventManager()
	cm  := character.MakeCharacterManager(tm, em)
	is  := &Instance {
		gof,
		tm,
		physics.MakePhysicsManager(tm),
		cm,
		rm,
		em,
		ai.MakeAiManager(tm, cm, em),
                quest.MakeQuestManager(),
		gm,
                nm,
		returnlink,
		make(chan string),
		0,
	}

	is.gof.Register("physics", is.pm, is.pm.JsonCreate)
	is.gof.Register("character", is.cm, is.cm.JsonCreate)
	is.gof.Register("ai", is.am, is.am.JsonCreate)
	is.gof.Register("graphics", is.gm, is.gm.JsonCreate)
        is.gof.Register("quest", is.qm, is.qm.JsonCreate)

	is.am.RegisterComputer("enemy", is.am.EnemyDecide)
	is.am.RegisterComputer("player", is.am.PlayerDecide)

	is.em.RegisterListener("attack", is.cm.HandleAttack)
	is.em.RegisterListener("death", is.gof.HandleDeath)
	is.em.RegisterListener("attack", is.qm.HandleEvent)
	is.em.RegisterListener("kill", is.qm.HandleEvent)
        is.em.RegisterListener("chat", is.cm.HandleChat)
        is.em.RegisterListener("playerCreated", is.am.HandlePlayerCreated)

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

	is.player = is.CreateObject("player", "0,0,0")
        is.qm.AddQuest(is.player, is.qm.FirstQuest)
        is.tm.SetLocation(is.player, math.Vec3{10,10,10})

        is.StartScript()


        for numTicks := 0; ; {
		<-ticks.C

		newTime := time.Since(oldTime)
		oldTime = time.Now()
		secs := newTime.Seconds()

		// fmt.Println(newTime)

                _, err := is.nm.RecieveBytes(100, 5)
                if err != nil {
                    //common.LogWarn.Print(err)
                }
                //fmt.Println(data)

		is.ParseSysConsole()
		is.em.Tick(secs)
                is.qm.Tick(secs)
		is.cm.Tick(secs)
                //is.pm.Tick(secs)
		is.tm.Tick(secs)

                if numTicks % 10 == 0 {
                    is.am.Tick(secs)
                }
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
			fmt.Scan(&breed, &location)
			is.CreateObject(breed, location)
		case "runai":
			var arg component.GOiD
			fmt.Scan(&arg)
			is.am.RunAi(arg)
		case "player":
			is.am.RunAi(is.player)
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

func (is *Instance) CreateFromMap(mapName string) []component.GOiD {
	jmap := is.rm.LoadJsonMap(mapName)
        return is.gof.CreateFromMap(&jmap)
}

func (is *Instance) CreateObject(objName, location string) component.GOiD {
	components := is.rm.LoadGameObject(objName)
	strLoc := strings.Split(location,",")
	f1, err := strconv.ParseFloat(strLoc[0], 32)
	f2, err := strconv.ParseFloat(strLoc[1], 32)
	f3, err := strconv.ParseFloat(strLoc[2], 32)
	if err != nil {
		fmt.Println(err)
	}
	id, err := is.gof.Create(components, math.Vec3{float32(f1),float32(f2),float32(f3)})
	if err != nil {
		common.LogErr.Print(err)
	}
	// is.pm.AddForce(id, math.Vec3{0,0.5,0})
	fmt.Println("\tid:",id)

	return id
}

func (is *Instance) GetSceneManager() component.SceneManager {
	return is.tm
}

func (is *Instance) GetEventManager() *event.EventManager {
    return is.em
}
