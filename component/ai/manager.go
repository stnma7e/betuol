package ai

import (
	"encoding/json"
	"fmt"

	"smig/common"
	"smig/component"
	"smig/component/character"
	"smig/component/scenemanager"
	"smig/event"
)

type AiManager struct {
	computerTypeMap map[string]AiComputer
	compList        []chan event.Event
	players         *common.Vector

	cm *character.CharacterManager
	tm *scenemanager.TransformManager
	em *event.EventManager

	aiTicker func(delta float64)
}

func MakeAiManager(tm *scenemanager.TransformManager, cm *character.CharacterManager, em *event.EventManager) *AiManager {
	am := AiManager{
		make(map[string]AiComputer),
		nil,
		common.MakeVector(),
		cm,
		tm,
		em,
		nil,
	}

	am.aiTicker = am.UpdateAiNearPlayer
	return &am
}

func (am *AiManager) Tick(delta float64) {
	am.aiTicker(delta)
}

func (am *AiManager) UpdateAi(delta float64) {
	players := am.players.Array()
	for i := range am.compList {
		if func() bool {
			for j := range players {
				if i == int(players[j].(component.GOiD)) || am.compList[i] == nil {
					return true
				}
			}
			return false
		}() { //end func
			continue
		} else {
			am.RunAi(component.GOiD(i))
		}
	}
}

func (am *AiManager) UpdateAiNearPlayer(delta float64) {
	players := am.players.Array()
	for i := range players {
		loc := am.tm.GetObjectLocation(players[i].(component.GOiD))
		charsInRadius := am.tm.GetObjectsInLocationRadius(loc, 5)
		chars := charsInRadius.Array()
		for j := 0; j < len(chars); j++ {
			if func() bool {
				for k := range players {
					if chars[j] == int(players[k].(component.GOiD)) || am.compList[chars[j]] == nil {
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
	if len(am.compList) < int(id) || am.compList[id] == nil {
		common.LogErr.Printf("no ai routine for id %v", id)
		return
	}
	am.compList[id] <- event.RunAiEvent{}
	<-am.compList[id]
}

func (am *AiManager) JsonCreate(id component.GOiD, data []byte) error {
	var obj struct {
		Type string
	}
	json.Unmarshal(data, &obj)
	return am.CreateComponent(id, obj.Type)
}

func (am *AiManager) CreateComponent(id component.GOiD, computerType string) error {
	am.resizeArray(id)
	computer, ok := am.computerTypeMap[computerType]
	if !ok {
		return fmt.Errorf("unregistered ai type: %s", computerType)
	}

	am.compList[id] = make(chan event.Event)
	go computer(id, am.compList[id])

	if computerType == "player" {
		am.players.Insert(id)
		am.em.Send(event.PlayerCreatedEvent{id})
	}

	return nil
}

func (am *AiManager) DeleteComponent(id component.GOiD) {
	if len(am.compList) <= int(id) {
		return
	}
	if am.compList[id] != nil {
		am.compList[id] <- event.DeathEvent{id}
		am.compList[id] = nil
	}
}

func (am *AiManager) resizeArray(index component.GOiD) {
	const RESIZESTEP = 5
	if cap(am.compList)-1 < int(index) {
		newCompList := make([]chan event.Event, index+RESIZESTEP)
		for i := range am.compList {
			newCompList[i] = am.compList[i]
		}
		am.compList = newCompList
	}
}

func (am *AiManager) RegisterComputer(aiType string, computer AiComputer) {
	am.computerTypeMap[aiType] = computer
}

func (am *AiManager) SetUpdateAiNearPlayer(yes bool) {
	if yes {
		am.aiTicker = am.UpdateAiNearPlayer
	} else {
		am.aiTicker = am.UpdateAi
	}
}
