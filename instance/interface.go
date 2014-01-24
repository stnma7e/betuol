package instance

import (
	"fmt"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/component/quest"
	"github.com/stnma7e/betuol/event"
	"github.com/stnma7e/betuol/math"
)

type IInstance interface {
	CreateFromMap(mapName string) ([]component.GOiD, error)
	CreateObject(objName string, loc math.Vec3) (component.GOiD, error)
	MoveGameObject(id component.GOiD, newLoc math.Vec3)
	AddQuest(id component.GOiD, questName string)
	RunAi(id component.GOiD)
	RenderFromPerspective(id component.GOiD)

	GetSceneManagerSnapshot() component.SceneManager
	GetEventManager() *event.EventManager

	Tick(dtime float64)
}

// CreateFromMap is a helper function to create a map from a filename string.
func (is *Instance) CreateFromMap(mapName string) ([]component.GOiD, error) {
	jmap, err := is.rm.LoadJsonMap(mapName)
	if err != nil {
		return []component.GOiD{}, err
	}
	return is.gof.CreateFromMap(&jmap)
}

// CreateObject is a helper function to create a GameObject with a starting location.
func (is *Instance) CreateObject(objName string, loc math.Vec3) (component.GOiD, error) {
	components, err := is.rm.LoadGameObject(objName)
	if err != nil {
		return 0, err
	}
	id, err := is.gof.Create(components, loc)
	if err != nil {
		return 0, fmt.Errorf("gameobject %s creation failed, error: %s", objName, err.Error())
	}

	// is.pm.AddForce(id, math.Vec3{0,0.5,0})
	common.LogInfo.Println("entity created, id:", id)

	return id, nil
}

// GetSceneManagerSnapshot is used by the graphics manager to request a snapshot of the location data during the current frame. The graphics manager uses this to depict the scene for the current frame without data corruption.
func (is *Instance) GetSceneManagerSnapshot() component.SceneManager {
	snap := is.tmSnapshot
	// prevents the snapshot from changing during the rendering process
	return &snap
}

// GetEventManager returns a pointer to the instance's event manager.
func (is *Instance) GetEventManager() *event.EventManager {
	return is.em
}

// AddQuest adds a quest using the quest manager.
func (is *Instance) AddQuest(id component.GOiD, questName string) {
	var QuestNames = map[string]quest.QuestState{
		"attack":    is.qm.AttackQuest,
		"kill":      is.qm.KillQuest,
		"firstMove": is.qm.FirstMoveQuest,
	}
	// this needs to be redone with the quest overhaul. Shit fix.

	is.qm.AddQuest(id, QuestNames[questName])
}

// MoveGameObject moves the GameObject in the transform manager.
func (is *Instance) MoveGameObject(id component.GOiD, newLoc math.Vec3) {
	is.tm.SetLocation(id, newLoc)
}

// RunAi calls the RunAi function on the ai manager.
func (is *Instance) RunAi(id component.GOiD) {
	is.am.RunAi(id)
}

func (is *Instance) RenderFromPerspective(id component.GOiD) {
	compsToSend, errs := is.gm.RenderAllFromPerspective(id, is.tm)
	if errs != nil {
		errArray := errs.Array()
		if errArray != nil && len(errArray) > 0 {
			for i := range errArray {
				common.LogErr.Print(errArray[i].(error))
			}
		}
	}
	is.gm.ForceRender(compsToSend)
}
