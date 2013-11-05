package quest

import (
    "smig/component"
    "smig/event"
)

type QuestState func(id component.GOiD, evt event.Event)

func (qm *QuestManager) FirstQuest(id component.GOiD, evt event.Event) {
    if evt.GetEventType()!= "attack" {
        return
    }
    aevt := evt.(event.AttackEvent)
    qm.AddQuest(aevt.Char1, qm.KillQuest)
    qm.em.Send(event.QuestComplete{ aevt.Char1, "FirstQuest" })
}

func (qm *QuestManager) KillQuest(id component.GOiD, evt event.Event) {
    if evt.GetEventType() != "kill" {
        return
    }
    kevt := evt.(event.KillEvent)
    if kevt.Killer != id {
        return
    }
    qm.em.Send(event.QuestComplete{ kevt.Killer, "KillQuest" })
}