package character

import (
	"encoding/json"

	"betuol/component"
	"betuol/component/scenemanager"
	"betuol/event"
)

// CharacterManager is the character component manager that handles the creation, deletion, and updating of character components.
type CharacterManager struct {
	attributeList   [NUM_ATTRIBUTES][]float32
	descriptionList []string
	greetingList    []string

	movedlink chan component.GOiD

	sm *scenemanager.TransformManager
	em *event.EventManager
}

// MakeCharacterManager returns a pointer to a CharacterManager.
func MakeCharacterManager(tm *scenemanager.TransformManager, em *event.EventManager) *CharacterManager {
	cm := CharacterManager{}
	cm.movedlink = make(chan component.GOiD)
	cm.sm = tm
	cm.em = em

	return &cm
}

// Tick updates character components based on elapsed time passed as an argument, delta.
func (cm *CharacterManager) Tick(delta float64) {
}

// JsonCreate extracts creation data from a byte array of json text to pass to CreateComponent.
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

// Uses extracted data from higher level component creation functions and initializes a character component based on the id passed through.
func (cm *CharacterManager) CreateComponent(index component.GOiD, ca CharacterAttributes) error {
	cm.resizeArrays(index)
	for i := range ca.Attributes {
		cm.attributeList[i][index] = ca.Attributes[i]
	}

	cm.descriptionList[index] = ca.Description
	cm.greetingList[index] = ca.Greeting

	return nil
}

// resizeArray is a helper function to resize the array of components to accomodate a new component.
// If the GOiD of the new component is larger than the size of the array, then resizeArrays will grow the array and copy data over in order to fit the new component.
func (cm *CharacterManager) resizeArrays(index component.GOiD) {
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

// DeleteComponent implements the component.ComponentManager interface and deletes character component data from the manager.
func (cm *CharacterManager) DeleteComponent(index component.GOiD) {
	for i := range cm.attributeList {
		cm.attributeList[i][index] = -1
	}
	cm.descriptionList[index] = ""
	cm.greetingList[index] = ""
}

// GetCharacterAttributes is a helper function that will compile attribute data of a character component and return a pointer to a structure with the desired information.
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

// UpdateId updates a character component's attribute information based on the data in the CharacterAttributes structure passed in as an argument.
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
