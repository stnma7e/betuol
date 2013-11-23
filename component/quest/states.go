package quest

import (
	"betuol/component"
	"betuol/event"
)

type QuestState func(id component.GOiD, evt event.Event)

var moveQuestTicker int

func (qm *QuestManager) FirstQuest(id component.GOiD, evt event.Event) {
	if evt.GetEventType() != "attack" {
		return
	}
	aevt := evt.(event.AttackEvent)
	qm.AddQuest(aevt.Char1, qm.KillQuest)
	qm.em.Send(event.QuestComplete{aevt.Char1, "FirstQuest"})
}

func (qm *QuestManager) KillQuest(id component.GOiD, evt event.Event) {
	if evt.GetEventType() != "kill" {
		return
	}
	kevt := evt.(event.KillEvent)
	if kevt.Killer != id {
		return
	}
	qm.em.Send(event.QuestComplete{kevt.Killer, "KillQuest"})
}

func (qm *QuestManager) FirstMoveQuest(id component.GOiD, evt event.Event) {
	if evt.GetEventType() != "characterMoved" {
		return
	}
	cevt := evt.(event.CharacterMoveEvent)
	if cevt.CharID != id {
		return
	}
	if moveQuestTicker < 1 {
		// the location of the player is changed once during gameobject creation and once during instance startup
		moveQuestTicker++
		return
	}
	qm.em.Send(event.QuestComplete{cevt.CharID, "FirstMoveQuest"})
	qm.quests[id] = nil
}
