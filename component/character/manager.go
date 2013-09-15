package character

import (
	"fmt"
	"encoding/json"

	"smig/component"
	"smig/common"
)

type CharacterManager struct {
	attributeList 	[NUM_ATTRIBUTES][]float32
	descriptionList []string
	greetingList 	[]string
	factionList 	[]string

	movedlink 		chan component.GOiD

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
		}
	default:
	}
}

func (cm *CharacterManager) JsonCreate(index component.GOiD, data []byte) error {
	var comp struct {
		Health, Mana, Strength, Intelligence, RangeOfSight float32
		Description, Greeting, AiFunction, Faction string
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
		comp.Faction,
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
	cm.factionList[index] 	  = ca.Faction

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
	if cap(cm.factionList) - 1 < int(index) {
		tmp := cm.factionList
		cm.factionList = make([]string, index + RESIZESTEP)
		for i := range tmp {
			cm.factionList[i] = tmp[i]
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

func (cm *CharacterManager) Update(id component.GOiD, ca *CharacterAttributes) {
	for i := range cm.attributeList {
		cm.attributeList[i][id] = ca.Attributes[i]
	}
}

func (cm *CharacterManager) Died(id component.GOiD) {
	fmt.Println(id, "died: ", cm.attributeList[HEALTH][id])
	// send event of death
}