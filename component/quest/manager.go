package quest

import (
	"smig/common"
	"smig/component"
	"smig/event"
)

type QuestManager struct {
	em     *event.EventManager
	quests map[component.GOiD]QuestState
}

func MakeQuestManager(em *event.EventManager) *QuestManager {
	qm := &QuestManager{
		em,
		make(map[component.GOiD]QuestState),
	}

	return qm
}

func (qm *QuestManager) Tick(delta float64) {
}

func (qm *QuestManager) JsonCreate(id component.GOiD, data []byte) error {
	return qm.CreateComponent(id)
}

func (qm *QuestManager) CreateComponent(id component.GOiD) error {
	return nil
}

func (qm *QuestManager) AddQuest(id component.GOiD, state QuestState) {
	qm.quests[id] = state
}

func (qm *QuestManager) DeleteComponent(id component.GOiD) {
	qm.quests[id] = nil
}

func (qm *QuestManager) HandleEvent(evt event.Event) {
	for id, fun := range qm.quests {
		if fun == nil {
			continue
		}
		fun(id, evt)
	}
}

func (qm *QuestManager) QuestComplete(evt event.Event) {
	qevt := evt.(event.QuestComplete)
	common.LogInfo.Printf("%v completed the quest: %s.", qevt.Id, qevt.QuestName)
}
