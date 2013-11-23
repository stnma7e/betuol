// Package quest handles the objective system of the game.
// It will register quest functions with GOiD's to allow an objective system to play out.
package quest

import (
	"betuol/common"
	"betuol/component"
	"betuol/event"
)

// QuestManager is the quest component manager that handles the creation, deletion, and updating of quest components.
type QuestManager struct {
	em     *event.EventManager
	quests map[component.GOiD]QuestState
}

// MakeQuestManager returns a pointer to a QuestManager
func MakeQuestManager(em *event.EventManager) *QuestManager {
	qm := &QuestManager{
		em,
		make(map[component.GOiD]QuestState),
	}

	return qm
}

// Tick updates character components based on elapsed time passed as an argument, delta.
func (qm *QuestManager) Tick(delta float64) {
}

// JsonCreate extracts creation data from a byte array of json text to pass to CreateComponent.
func (qm *QuestManager) JsonCreate(id component.GOiD, data []byte) error {
	return qm.CreateComponent(id)
}

// Uses extracted data from higher level component creation functions and quest a character component based on the id passed through.
func (qm *QuestManager) CreateComponent(id component.GOiD) error {
	return nil
}

// AddQuest is used to set the active quest of a quest component.
func (qm *QuestManager) AddQuest(id component.GOiD, state QuestState) {
	qm.quests[id] = state
}

// DeleteComponent implements the component.ComponentManager interface and deletes quest component data from the manager.
func (qm *QuestManager) DeleteComponent(id component.GOiD) {
	qm.quests[id] = nil
}

// HandleEvent is used to recieve game events from the event manager and distribute them to the active quest components.
func (qm *QuestManager) HandleEvent(evt event.Event) {
	for id, fun := range qm.quests {
		if fun == nil {
			continue
		}
		fun(id, evt)
	}
}

// QuestComplete is a helper function to print the status of a completed quest.
func (qm *QuestManager) QuestComplete(evt event.Event) {
	qevt := evt.(event.QuestComplete)
	common.LogInfo.Printf("%v completed the quest: %s.", qevt.Id, qevt.QuestName)
}
