// Package instance implements a standalone structure that is capable of running a game.
package instance

import (
	"math/rand"
	"time"

	"github.com/stnma7e/betuol/component/ai"
	"github.com/stnma7e/betuol/component/character"
	"github.com/stnma7e/betuol/component/gofactory"
	"github.com/stnma7e/betuol/component/graphics"
	"github.com/stnma7e/betuol/component/physics"
	"github.com/stnma7e/betuol/component/quest"
	"github.com/stnma7e/betuol/component/scenemanager"
	"github.com/stnma7e/betuol/event"
	"github.com/stnma7e/betuol/res"
)

// The instance struct is used as an individual structure that is used to encapsulate the game state.
// Each instance has its own separate state from all other instances.
// It is the core piece of the game, and it represents one self-suficient portion of a game.
// Multiple instances could be bound together in a game to represent various levels or zones.
type Instance struct {
	gof *gofactory.GameObjectFactory
	tm  *scenemanager.TransformManager
	pm  *physics.PhysicsManager
	cm  *character.CharacterManager
	rm  *res.ResourceManager
	em  *event.EventManager
	am  *ai.AiManager
	qm  *quest.QuestManager
	gm  *graphics.GraphicsManager

	tmSnapshot  scenemanager.TransformManager
	commandlink chan string
	numTicks    int8
}

// MakeInstance returns a pointer to an Instance.
// A series of managers for various systems are created inside the function and a fully initialized Instance is returned.
func MakeInstance(rm *res.ResourceManager) *Instance {
	em := event.MakeEventManager()
	tm := scenemanager.MakeTransformManager(em)
	gof := gofactory.MakeGameObjectFactory(tm)
	cm := character.MakeCharacterManager(tm, em)
	is := &Instance{
		gof,
		tm,
		physics.MakePhysicsManager(tm),
		cm,
		rm,
		em,
		ai.MakeAiManager(tm, cm, em),
		quest.MakeQuestManager(em),
		graphics.MakeGraphicsManager(em, rm, tm),
		*tm,
		make(chan string),
		0,
	}

	is.gof.Register("character", is.cm, is.cm.JsonCreate)
	is.gof.Register("ai", is.am, is.am.JsonCreate)
	is.gof.Register("graphics", is.gm, is.gm.JsonCreate)
	is.gof.Register("quest", is.qm, is.qm.JsonCreate)
	is.gof.Register("physics", is.pm, is.pm.JsonCreate)

	is.am.SetUpdateAiNearPlayer(false)
	is.am.RegisterComputer("player", ai.PlayerAi)
	is.am.RegisterComputer("enemy", ai.EnemyAi)
	is.am.RegisterComputer("wander", ai.WanderAi)

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

// Loop is launched as a goroutine and updates on its own.
// This function does some initialization, but then jumps into an infinite loop.
// When the loop breaks, a bool will be sent along the returnlink that the instance was created with.
// The returnlink is expected to be periodically checked for activity because this indicates an exit of the main loop of the instance.
func (is *Instance) Tick(secs float64) {
	// fmt.Println(newTime)

	//_, err := is.nm.RecieveBytes(100, 5)
	//if err != nil {
	//common.LogWarn.Print(err)
	//}
	//fmt.Println(data)

	if is.numTicks > 60 {
		is.am.Tick(secs)
		is.numTicks = 0
	}

	is.em.Tick(secs)
	is.qm.Tick(secs)
	is.cm.Tick(secs)
	is.pm.Tick(secs)
	is.tm.Tick(secs)

	is.gm.Tick(secs, &is.tmSnapshot)
	is.tmSnapshot = *is.tm
	is.numTicks++
}
