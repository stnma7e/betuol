package character

import (
	"encoding/json"

	"smig/component"
	"smig/component/scene"
	"smig/common"
)

const (
	HEALTH 			= iota
	MANA   			= iota
	STRENGTH 		= iota
	INTELLIGENCE 	= iota

	NUM_ATTRIBUTES  = iota

	RESIZESTEP = 20
)

type AiComputer func(id component.GOiD)

type CharacterManager struct {
	attributeList 	[NUM_ATTRIBUTES][]float32
	descriptionList []string
	greetingList 	[]string

	aiList 			[]AiComputer

	movedlink chan Player
}

func (cm *CharacterManager) CreatePlayer(id component.GOiD, lookRange float32, sm *scene.SceneManager) {
	if cm.movedlink == nil {
		cm.movedlink = make(chan Player)
	}
	go StartPlayer(id, lookRange, cm.movedlink, sm, cm)
}

func (cm *CharacterManager) Tick(delta float64) {
	select {
	case pl := <-cm.movedlink:
		loc :=  pl.Scene.GetObjectLocation(pl.Id)
		stk := pl.Scene.GetObjectsInLocationRange(loc, pl.RangeOfSight + 20)
		numObj := stk.Size
		for i := 0; i < numObj; i++ {
			id,err := stk.Dequeue()
			if err != nil {
				common.Log.Warn(err)
			}
			if id == int(pl.Id) || id == 0 {
				continue
			}

			// cm.aiList[id](component.GOiD(id))
		}
	default:
		break
	}
}

func (cm *CharacterManager) JsonCreate(index component.GOiD, data []byte) error {
	var comp struct {
		Health, Mana, Strength, Intelligence float32
		Description, Greeting, AiFunction string
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

	return cm.CreateComponent(index, ca, comp.AiFunction)
}

func (cm *CharacterManager) CreateComponent(index component.GOiD, ca CharacterAttributes, aiFuncName string) error {
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
	if cap(cm.aiList) - 1 < int(index) {
		tmp := cm.aiList
		cm.aiList = make([]AiComputer, index + RESIZESTEP)
		for i := range tmp {
			cm.aiList[i] = tmp[i]
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