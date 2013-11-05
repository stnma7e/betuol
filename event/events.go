package event

import (
	"smig/component"
)

type AttackEvent struct {
    Char1,Char2 component.GOiD
}
func (at AttackEvent) GetEventType() string {
    return "attack"
}

type DeathEvent struct {
    Id component.GOiD
}
func (dt DeathEvent) GetEventType() string {
    return "death"
}

type KillEvent struct {
    Killer, Dead component.GOiD
}
func (ke KillEvent) GetEventType() string {
    return "kill"
}

type QuestComplete struct {
    Id component.GOiD
    QuestName string
}
func (qt QuestComplete) GetEventType() string {
    return "questComplete"
}

type ChatEvent struct {
    Sender component.GOiD
    Reciever component.GOiD
    Message string
}
func (ce ChatEvent) GetEventType() string {
    return "chat"
}

type PlayerCreatedEvent struct {
    PlayerID component.GOiD
}
func (ce PlayerCreatedEvent) GetEventType() string {
    return "playerCreated"
}