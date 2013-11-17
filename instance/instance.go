package instance

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"smig/common"
	"smig/component"
	"smig/component/ai"
	"smig/component/character"
	"smig/component/gofactory"
	"smig/component/quest"
	"smig/component/scenemanager"
	"smig/event"
	"smig/graphics"
	"smig/math"
	"smig/res"
)

type Instance struct {
	gof *gofactory.GameObjectFactory
	tm  *scenemanager.TransformManager
	cm  *character.CharacterManager
	rm  *res.ResourceManager
	em  *event.EventManager
	am  *ai.AiManager
	qm  *quest.QuestManager
	gm  *graphics.GraphicsManager

	tmSnapshot scenemanager.TransformManager

	returnlink  chan bool
	commandlink chan string

	player component.GOiD
}

func MakeInstance(returnlink chan bool, rm *res.ResourceManager, gm *graphics.GraphicsManager) *Instance {
	em := event.MakeEventManager()
	tm := scenemanager.MakeTransformManager(em)
	gof := gofactory.MakeGameObjectFactory(tm)
	cm := character.MakeCharacterManager(tm, em)
	is := &Instance{
		gof,
		tm,
		cm,
		rm,
		em,
		ai.MakeAiManager(tm, cm, em),
		quest.MakeQuestManager(em),
		gm,
		*tm,
		returnlink,
		make(chan string),
		0,
	}

	is.gof.Register("character", is.cm, is.cm.JsonCreate)
	is.gof.Register("ai", is.am, is.am.JsonCreate)
	is.gof.Register("graphics", is.gm, is.gm.JsonCreate)
	is.gof.Register("quest", is.qm, is.qm.JsonCreate)

	is.am.SetUpdateAiNearPlayer(false)
	//is.am.RegisterComputer("enemy", is.am.EnemyDecide)
	//is.am.RegisterComputer("player", is.am.PlayerDecide)
	is.am.RegisterComputer("player", is.am.PlayerAi)
	is.am.RegisterComputer("enemy", is.am.WanderAi)

	is.em.RegisterListener("attack", is.cm.HandleAttack)
	is.em.RegisterListener("death", is.gof.HandleDeath)
	is.em.RegisterListener("attack", is.qm.HandleEvent)
	is.em.RegisterListener("characterMoved", is.qm.HandleEvent)
	is.em.RegisterListener("kill", is.qm.HandleEvent)
	is.em.RegisterListener("questComplete", is.qm.QuestComplete)
	is.em.RegisterListener("chat", is.cm.HandleChat)

	rand.Seed(time.Now().UnixNano())

	return is
}

func (is *Instance) Loop() {
	defer is.Shutdown()
	oldTime := time.Now()
	ticks := time.NewTicker(time.Second / 60)

	is.player = is.CreateObject("player", math.Vec3{0, 0, 0})
	is.qm.AddQuest(is.player, is.qm.FirstQuest)
	is.tm.SetLocationOverTime(is.player, math.Vec3{3, 0, 0}, 1.5)

	is.StartScript()

	go func() {
		r := bufio.NewReaderSize(os.Stdin, 4*1024)
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				common.LogErr.Println(err)
			}
			s := string(line)
			if len(s) > 0 && s[len(s)-1] == '\n' {
				s = line[:len(s)-1]
			}
			if len(s) > 0 && s[len(s)-1] == '\r' {
				s = line[:len(s)-1]
			}
			if s != "" {
				is.commandlink <- s
				<-is.commandlink
				s = ""
			}
		}
	}()

	for numTicks := 0; ; numTicks++ {
		<-ticks.C

		newTime := time.Since(oldTime)
		oldTime = time.Now()
		secs := newTime.Seconds()

		// fmt.Println(newTime)

		//_, err := is.nm.RecieveBytes(100, 5)
		//if err != nil {
		//common.LogWarn.Print(err)
		//}
		//fmt.Println(data)

		is.ParseSysConsole()

		if numTicks > 60 {
			is.am.Tick(secs)
			numTicks = 0
		}

		is.em.Tick(secs)
		is.qm.Tick(secs)
		is.cm.Tick(secs)
		//is.pm.Tick(secs)
		is.tm.Tick(secs)
		is.tmSnapshot = *is.tm
	}
}

func (is *Instance) ParseSysConsole() {
	select {
	case command := <-is.commandlink:
		args := strings.SplitAfter(command, " ")
		for i := range args {
			if i == len(args)-1 {
				continue
			}
			args[i] = args[i][:len(args[i])-1]
		}
		fmt.Println(args)
		switch args[0] {
		case "exit":
			is.returnlink <- true
		case "loadmap":
			if !(len(args) >= 2) {
				common.LogErr.Print("not enough arguments to 'loadmap'")
				break
			}
			is.CreateFromMap(args[1])
		case "loadobj":
			if !(len(args) >= 3) {
				common.LogErr.Print("not enough arguments to 'loadobj'")
				break
			}
			//is.CreateObject(breed, location)
			is.CreateObject(args[1], decodeVec3String(args[2]))
		case "runai":
			if !(len(args) >= 2) {
				common.LogErr.Print("not enough arguments to 'runai'")
				break
			}
			arg, err := strconv.Atoi(args[1])
			if err != nil {
				common.LogErr.Println(err)
			}
			is.am.RunAi(component.GOiD(arg))
		case "player":
			is.am.RunAi(is.player)
		default:
			fmt.Println("\tInvalid command. Type \"help\" for choices.")
		}
		is.commandlink <- ""
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

func (is *Instance) CreateObject(objName string, loc math.Vec3) component.GOiD {
	components := is.rm.LoadGameObject(objName)
	id, err := is.gof.Create(components, loc)
	if err != nil {
		common.LogErr.Print(err)
	}
	// is.pm.AddForce(id, math.Vec3{0,0.5,0})
	common.LogInfo.Println("entity created, id:", id)

	return id
}

func (is *Instance) MoveObject(id component.GOiD, loc math.Vec3) {

}

func (is *Instance) GetSceneManagerSnapshot() component.SceneManager {
	snap := is.tmSnapshot
	// prevents the snapshot from changing during the rendering process
	return &snap
}

func (is *Instance) GetEventManager() *event.EventManager {
	return is.em
}

func decodeVec3String(vec3 string) math.Vec3 {
	strLoc := strings.Split(vec3, ",")
	f1, err := strconv.ParseFloat(strLoc[0], 32)
	f2, err := strconv.ParseFloat(strLoc[1], 32)
	f3, err := strconv.ParseFloat(strLoc[2], 32)
	if err != nil {
		common.LogErr.Println(err)
	}

	return math.Vec3{float32(f1), float32(f2), float32(f3)}
}
