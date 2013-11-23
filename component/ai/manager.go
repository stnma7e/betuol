// Package ai manages ai components and implements ai functions used by GameObjects.
package ai

import (
	"encoding/json"
	"fmt"

	"betuol/common"
	"betuol/component"
	"betuol/component/character"
	"betuol/component/scenemanager"
	"betuol/event"
)

// AiManager is the ai component manager that handles the creation, deletion, and updating of ai components.
type AiManager struct {
	computerTypeMap map[string]AiComputer
	compList        []chan event.Event
	players         *common.Vector

	cm *character.CharacterManager
	tm *scenemanager.TransformManager
	em *event.EventManager

	aiTicker func(delta float64)
}

// MakeAiManager returns a pointer to an AiManager.
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

// Tick is called to update the ai components based on the difference in time between the last update.
// delta is used to specify the elapsed time since the last update.
func (am *AiManager) Tick(delta float64) {
	am.aiTicker(delta)
}

// UpdateAi implements an update sequence for updating ai components.
// This function can be registered with an AiManager to be called when AiManager.Tick() is called.
// This function will update all ai components created by the AiManager.
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

// UpdateAiNearPlayer implements an update sequence for updating ai components.
// This function can be registered with an AiManager to be called when AiManager.Tick() is called.
// This function uses a list of players maintained by the AiManager, and the function will update ai components within a certain proximity to any player. If the GameObject who owns the ai component is not within a certain radius of a player, then its ai component will not be updated.
func (am *AiManager) UpdateAiNearPlayer(delta float64) {
	players := am.players.Array()
	for i := range players {
		loc := am.tm.GetObjectLocation(players[i].(component.GOiD))
		charsInRadius := am.tm.GetObjectsInLocationRadius(loc, 5)
		chars := charsInRadius.Array()
		for j := 0; j < len(chars); j++ {
			if func() bool {
				for k := range players {
					if chars[j].(component.GOiD) == players[k].(component.GOiD) ||
						am.compList[int(chars[j].(component.GOiD))] == nil {
						return true
					}
				}
				return false
			}() { //end func
				continue
			} else {
				am.RunAi(chars[j].(component.GOiD))
			}
		}
	}
}

// RunAi can be called on demand to immediately run the ai function of an ai component.
// This will run an ai component's computer function before, after, or during a call to AiManager.Tick() which updates all components according to the update function used.
func (am *AiManager) RunAi(id component.GOiD) {
	if len(am.compList) < int(id) || am.compList[id] == nil {
		common.LogErr.Printf("no ai routine for id %v", id)
		return
	}
	am.compList[id] <- event.RunAiEvent{}
	<-am.compList[id]
}

// JsonCreate will use a byte array of json creation data passed to it in order to initialize an ai component for the given GOiD.
func (am *AiManager) JsonCreate(id component.GOiD, data []byte) error {
	var obj struct {
		Type string
	}
	json.Unmarshal(data, &obj)
	return am.CreateComponent(id, obj.Type)
}

// CreateComponent does the low level initialization and is called from higher level creation functions.
// Higher level creation functions extract the type of ai computer to be used for the component from a data source and pass the low level information to CreateComponent in order to do the real work of initialization.
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

// DeleteComponent implements the component.ComponentManager interface and removes the ai component from the entity component system.
func (am *AiManager) DeleteComponent(id component.GOiD) {
	if len(am.compList) <= int(id) {
		return
	}
	if am.compList[id] != nil {
		am.compList[id] <- event.DeathEvent{id}
		am.compList[id] = nil
	}
}

// resizeArray is a helper function to resize the array of components to accomodate a new component.
// If the GOiD of the new component is larger than the size of the array, then resizeArrays will grow the array and copy data over in order to fit the new component.
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

// RegisterComputer will register an ai computer type to be used for component initialization.
func (am *AiManager) RegisterComputer(aiType string, computer AiComputer) {
	am.computerTypeMap[aiType] = computer
}

// SetUpdateAiNearPlayer is a helper function to set the update function of AiManager.Tick().
// If true is passed to the function, AiManager will use the UpdateAiNearPlayer function to update ai components.
// If false if passed then UpdateAi is used instead.
func (am *AiManager) SetUpdateAiNearPlayer(yes bool) {
	if yes {
		am.aiTicker = am.UpdateAiNearPlayer
	} else {
		am.aiTicker = am.UpdateAi
	}
}
