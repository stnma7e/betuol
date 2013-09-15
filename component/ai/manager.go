package ai

import (
	"fmt"
	"encoding/json"

	"smig/component"
	"smig/component/character"
)

const (
	IDLE_STATE 	 = iota
	RUN_STATE 	 = iota
	ATTACK_STATE = iota
)

type AiManager struct {
	stateList []int

	computerMap map[component.GOiD]AiComputer
	computerTypeMap map[string]AiComputer

	cm *character.CharacterManager
	sm *component.SceneManager
}

func MakeAiManager(sm *component.SceneManager, cm *character.CharacterManager) *AiManager {
	am := AiManager{}
	am.computerMap  	= make(map[component.GOiD]AiComputer)
	am.computerTypeMap  = make(map[string]AiComputer)
	am.cm = cm
	am.sm = sm
	return &am
}

func (am *AiManager) Tick(delta float64) {

}

func (am *AiManager) RunAi(id component.GOiD) {
	comp, ok := am.computerMap[id]
	if !ok {
		fmt.Println("no computer for id:", id)
		return
	}
	attr := am.cm.GetCharacterAttributes(id)
	loc := am.sm.GetObjectLocation(id)
	idQueue := am.sm.GetObjectsInLocationRange(loc, attr.Attributes[character.RANGEOFSIGHT])
	size := idQueue.Size
	neighbors := make([]component.GOiD, size)
	for i := 0; i < size; i++ {
		val, err := idQueue.Dequeue()
		if err != nil {
			fmt.Println("error: bad dequeue:")
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

	return nil
}

func (am *AiManager) DeleteComponent(id component.GOiD) {

}

func (am *AiManager) RegisterComputer(aiType string, computer AiComputer) {
	am.computerTypeMap[aiType] = computer
}