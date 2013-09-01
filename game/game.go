package game

import (
	"time"

	"smig/component"
	"smig/component/physics"
	"smig/component/scene"
	"smig/component/character"
	"smig/math"
	"smig/res"
	"smig/common"
)

type Game struct {
	pl component.GOiD
	movedlink chan character.Player

	gof *component.GameObjectFactory
	sm  *scene.SceneManager
	pm  *physics.PhysicsManager
	cm  *character.CharacterManager
	rm  *res.ResourceManager
}

func MakeGame(playerName string, lookRange float32) *Game {
	sm := scene.MakeSceneManager()
	gm := &Game {
		0, make(chan character.Player),
		component.MakeGameObjectFactory(),
		sm, 
		physics.MakePhysicsManager(sm),
		&character.CharacterManager{}, 
		res.MakeResourceManager("/home/sam/go/data/"),
	}

	gm.gof.Register(component.SceneType, gm.sm, gm.sm.JsonCreate)
	gm.gof.Register("physics", gm.pm, gm.pm.JsonCreate)
	gm.gof.Register("character", gm.cm, gm.cm.JsonCreate)

	components := gm.rm.LoadGameObject("player")
	id, err := gm.gof.Create(components)
	if err != nil {
		common.Log.Error(err)
	}
	gm.cm.CreatePlayer(id, lookRange, gm.sm)

	gm.pl = id
	return gm
}

func (gm *Game) Loop(retrunlink chan bool) {
	oldTime := time.Now()
	ticks := time.NewTicker(time.Second / 60)
	for {
		select {
		case <-retrunlink:
			break
		default:
		}
		<-ticks.C

		newTime := time.Since(oldTime)
		secs := newTime.Seconds()

		// fmt.Println(newTime)

		gm.cm.Tick(secs)
		gm.pm.Tick(secs)
		gm.sm.Tick(secs)

		// for i := range list {
		// 	id := list[i]
		// 	trans,_ := tm.GetTransform(component.GOiD(id))
		// 	fmt.Println(id," ",trans.ToString())
		// }
		// fmt.Println()

		oldTime = time.Now()
	}
}

/*****************************************
*
* Component
*
*****************************************/

func (gm *Game) CreateFromMap(mapName string) {
	jmap := gm.rm.LoadJsonMap(mapName)
	gm.gof.CreateFromMap(&jmap)
}

/*****************************************
*
* Component/Transform
*
*****************************************/

func (gm *Game) GetObjectLocation(id component.GOiD) math.Vec3 {
	return gm.sm.GetObjectLocation(id)
}

func (gm *Game) GetObjectsInLocationRange(loc math.Vec3, lookRange float32) *common.IntQueue {
	return gm.sm.GetObjectsInLocationRange(loc, lookRange)
}

/*****************************************
*
* Component/Character
*
*****************************************/

func (gm *Game) GetCharacterAttributes(id component.GOiD) *character.CharacterAttributes {
	return gm.cm.GetCharacterAttributes(id)
}