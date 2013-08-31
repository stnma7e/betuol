package character

import (
	"encoding/json"

	"smig/component"
)

const (
	HEALTH 			= iota
	MANA   			= iota
	STRENGTH 		= iota
	INTELLIGENCE 	= iota

	NUM_ATTRIBUTES  = iota

	RESIZESTEP = 20
)

type CharacterManager struct {
	attributeList 	[NUM_ATTRIBUTES][]float32
	descriptionList []string
	greetingList 	[]string
}

func (cm *CharacterManager) JsonCreate(index component.GOiD, data []byte) error {
	var comp struct {
		Health, Mana, Strength, Intelligence float32
		Description, Greeting string
	}
	json.Unmarshal(data, &comp)

	ca := CharacterAttributes {
		[NUM_ATTRIBUTES]float32 {
			comp.Health,
			comp.Mana,
			comp.Strength,
			comp.Intelligence,
		},
		comp.Description,
		comp.Greeting,
	}

	return cm.CreateComponent(index, ca)
}

func (cm *CharacterManager) CreateComponent(index component.GOiD, ca CharacterAttributes) error {
	cm.resizeLists(index)
	for i := range ca.Attributes {
		cm.attributeList[i][index] = ca.Attributes[i]
	}

	cm.descriptionList[index] = ca.Description
	cm.greetingList[index]	  = ca.Greeting

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
}

func (cm *CharacterManager) DeleteComponent(index component.GOiD) {
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