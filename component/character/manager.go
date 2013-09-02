package character

import (
	"fmt"
	"encoding/json"

	"smig/component"
	"smig/common"
)

const (
	HEALTH 			= iota
	MANA   			= iota
	STRENGTH 		= iota
	INTELLIGENCE 	= iota
	RANGEOFSIGHT 	= iota

	NUM_ATTRIBUTES  = iota

	RESIZESTEP = 20
)

type CharacterManager struct {
	attributeList 	[NUM_ATTRIBUTES][]float32
	descriptionList []string
	greetingList 	[]string
	aiList 			[]AiComputer

	movedlink 		chan component.GOiD

	aiFunctionName	map[string]AiComputer
	Scene 			*component.SceneManager
}

func MakeCharacterManager(sm *component.SceneManager) *CharacterManager {
	cm := CharacterManager{}
	cm.movedlink = make(chan component.GOiD)
	cm.Scene 	 = sm

	return &cm
}

func (cm *CharacterManager) Tick(delta float64) {
	select {
	case id := <-cm.movedlink:
		loc :=  cm.Scene.GetObjectLocation(id)
		stk := cm.Scene.GetObjectsInLocationRange(loc, cm.attributeList[RANGEOFSIGHT][id] + 20)
		numObj := stk.Size
		for i := 0; i < numObj; i++ {
			charId,err := stk.Dequeue()
			if err != nil {
				common.Log.Warn(err)
			}
			if charId == int(id) || id == 0 {
				continue
			}

			cm.RunAi(component.GOiD(charId))
		}
	default:
	}
}

func (cm *CharacterManager) JsonCreate(index component.GOiD, data []byte) error {
	var comp struct {
		Health, Mana, Strength, Intelligence, RangeOfSight float32
		Description, Greeting, AiFunction string
	}
	json.Unmarshal(data, &comp)

	ca := CharacterAttributes {
		[NUM_ATTRIBUTES]float32 {
			comp.Health,
			comp.Mana,
			comp.Strength,
			comp.Intelligence,
			comp.RangeOfSight,
		},
		comp.Description,
		comp.Greeting,
	}

	return cm.CreateComponent(index, ca, comp.AiFunction)
}

func (cm *CharacterManager) CreateComponent(index component.GOiD, ca CharacterAttributes, aiFuncName string) error {
	cm.resizeLists(index)
	for i := range ca.Attributes {
		cm.attributeList[i][index] = ca.Attributes[i]
	}

	cm.descriptionList[index] = ca.Description
	cm.greetingList[index]	  = ca.Greeting
	var err error
	cm.aiList[index], err	  = cm.GetComputer(aiFuncName)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (cm *CharacterManager) resizeLists(index component.GOiD) {
	for i := range cm.attributeList {
		if cap(cm.attributeList[i]) - 1 < int(index) {
			tmp := cm.attributeList[i]
			cm.attributeList[i] = make([]float32, index + RESIZESTEP)
			for j := range tmp {
				cm.attributeList[i][j] = tmp[j]
			}
		}
	}

	if cap(cm.descriptionList) - 1 < int(index) {
		tmp := cm.descriptionList
		cm.descriptionList = make([]string, index + RESIZESTEP)
		for i := range tmp {
			cm.descriptionList[i] = tmp[i]
		}
	}
	if cap(cm.greetingList) - 1 < int(index) {
		tmp := cm.greetingList
		cm.greetingList = make([]string, index + RESIZESTEP)
		for i := range tmp {
			cm.greetingList[i] = tmp[i]
		}
	}
	if cap(cm.aiList) - 1 < int(index) {
		tmp := cm.aiList
		cm.aiList = make([]AiComputer, index + RESIZESTEP)
		for i := range tmp {
			cm.aiList[i] = tmp[i]
		}
	}
}

func (cm *CharacterManager) DeleteComponent(index component.GOiD) {
	for i := range cm.attributeList {
		cm.attributeList[i][index] = -1
	}
	cm.descriptionList[index] = ""
	cm.greetingList[index] = ""
}

func (cm *CharacterManager) GetCharacterAttributes(index component.GOiD) *CharacterAttributes {
	ca := &CharacterAttributes{}
	for i := range ca.Attributes {
		ca.Attributes[i] = cm.attributeList[i][index]
	}
	ca.Description = cm.descriptionList[index]
	ca.Greeting    = cm.greetingList[index]

	return ca
}