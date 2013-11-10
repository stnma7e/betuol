package character

import (
	"encoding/json"

	"smig/common"
	"smig/component"
	"smig/component/scenemanager"
	"smig/event"
)

type CharacterManager struct {
	attributeList   [NUM_ATTRIBUTES][]float32
	descriptionList []string
	greetingList    []string

	movedlink chan component.GOiD

	sm *scenemanager.TransformManager
	em *event.EventManager
}

func MakeCharacterManager(tm *scenemanager.TransformManager, em *event.EventManager) *CharacterManager {
	cm := CharacterManager{}
	cm.movedlink = make(chan component.GOiD)
	cm.sm = tm
	cm.em = em

	return &cm
}

func (cm *CharacterManager) Tick(delta float64) {
	select {
	case id := <-cm.movedlink:
		loc := cm.sm.GetObjectLocation(id)
		stk := cm.sm.GetObjectsInLocationRadius(loc, cm.attributeList[RANGEOFSIGHT][id]+20)
		numObj := stk.Size
		for i := 0; i < numObj; i++ {
			charId, err := stk.Dequeue()
			if err != nil {
				common.LogWarn.Print(err)
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
		Description, Greeting                              string
	}
	err := json.Unmarshal(data, &comp)
	if err != nil {
		return err
	}

	ca := CharacterAttributes{
		[NUM_ATTRIBUTES]float32{
			comp.Health,
			comp.Mana,
			comp.Strength,
			comp.Intelligence,
			comp.RangeOfSight,
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
	cm.greetingList[index] = ca.Greeting

	return nil
}

func (cm *CharacterManager) resizeLists(index component.GOiD) {
	for i := range cm.attributeList {
		if cap(cm.attributeList[i])-1 < int(index) {
			tmp := cm.attributeList[i]
			cm.attributeList[i] = make([]float32, index+RESIZESTEP)
			for j := range tmp {
				cm.attributeList[i][j] = tmp[j]
			}
		}
	}

	if cap(cm.descriptionList)-1 < int(index) {
		tmp := cm.descriptionList
		cm.descriptionList = make([]string, index+RESIZESTEP)
		for i := range tmp {
			cm.descriptionList[i] = tmp[i]
		}
	}
	if cap(cm.greetingList)-1 < int(index) {
		tmp := cm.greetingList
		cm.greetingList = make([]string, index+RESIZESTEP)
		for i := range tmp {
			cm.greetingList[i] = tmp[i]
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
	if index == 0 {
		return ca
	}

	for i := range ca.Attributes {
		ca.Attributes[i] = cm.attributeList[i][index]
	}
	ca.Description = cm.descriptionList[index]
	ca.Greeting = cm.greetingList[index]

	return ca
}

func (cm *CharacterManager) UpdateId(id component.GOiD, ca *CharacterAttributes) {
	if id == 0 {
		return
	}

	if ca.Attributes[HEALTH] <= 0 {
		cm.em.Send(event.DeathEvent{id})
	}
	for i := range cm.attributeList {
		cm.attributeList[i][id] = ca.Attributes[i]
	}
}
