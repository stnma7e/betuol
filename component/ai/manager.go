package ai

import (
	"fmt"
	"encoding/json"

	"smig/component"
        "smig/component/scenemanager"
	"smig/component/character"
	"smig/event"
        "smig/common"
)

const (
	IDLE_STATE	 = iota
	RUN_STATE	 = iota
	ATTACK_STATE = iota
)

type AiManager struct {
	computerMap map[component.GOiD]AiComputer
	computerTypeMap map[string]AiComputer

        players *common.Vector

	cm *character.CharacterManager
	tm *scenemanager.TransformManager
	em *event.EventManager
}

func MakeAiManager(tm *scenemanager.TransformManager, cm *character.CharacterManager, em *event.EventManager) *AiManager {
	am := AiManager{}
	am.computerMap      = make(map[component.GOiD]AiComputer)
	am.computerTypeMap  = make(map[string]AiComputer)
        am.players = common.MakeVector()
	am.cm = cm
	am.tm = tm
	am.em = em
	return &am
}

func (am *AiManager) Tick(delta float64) {
    players := am.players.Array()
    for i := range players {
        loc := am.tm.GetObjectLocation(players[i].(component.GOiD))
        charsInRadius := am.tm.GetObjectsInLocationRadius(loc, 5)
        chars := charsInRadius.Array()
        for j := 0; j < len(chars); j++ {
            if func() bool {
                              for k := range players {
                                   if chars[j] == int(players[k].(component.GOiD)) {
                                       return true
                                   }
                               }
                               return false
                           }() { //end func
                continue
            } else {
                am.RunAi(component.GOiD(chars[j]))
            }
        }
    }
}

func (am *AiManager) RunAi(id component.GOiD) {
	comp, ok := am.computerMap[id]
	if !ok || comp == nil {
		//common.LogWarn.Println("no computer for id:", id)
		return
	}
	attr := am.cm.GetCharacterAttributes(id)
	loc := am.tm.GetObjectLocation(id)
	idQueue := am.tm.GetObjectsInLocationRadius(loc, attr.Attributes[character.RANGEOFSIGHT])
	size := idQueue.Size
	neighbors := make([]component.GOiD, size)
	for i := 0; i < size; i++ {
		val, err := idQueue.Dequeue()
		if err != nil {
			common.LogErr.Println("error: bad dequeue:", err)
			continue
		}
		neighbors[i] = component.GOiD(val)
	}
	comp(id, neighbors, am.cm)
}

func (am *AiManager) JsonCreate(id component.GOiD, data []byte) error {
	var obj struct {
		Type string
	}
	json.Unmarshal(data, &obj)
	return am.CreateComponent(id, obj.Type)
}

func (am *AiManager) CreateComponent(id component.GOiD, computerType string) error {
	computer, ok := am.computerTypeMap[computerType]
	if !ok {
		return fmt.Errorf("bad ai type: %s", computerType)
	}
	am.computerMap[id] = computer
        if computerType == "player" {
            am.em.Send(event.PlayerCreatedEvent{ id })
        }

	return nil
}

func (am *AiManager) DeleteComponent(id component.GOiD) {
    _, ok := am.computerMap[id]; if ok {
        am.computerMap[id] = nil
    }
}

func (am *AiManager) RegisterComputer(aiType string, computer AiComputer) {
	am.computerTypeMap[aiType] = computer
}

func (am *AiManager) HandlePlayerCreated(evt event.Event) {
    if evt.GetEventType() != "playerCreated" {
        return
    }
    pcevt := evt.(event.PlayerCreatedEvent)
    am.players.Insert(pcevt.PlayerID)
}
