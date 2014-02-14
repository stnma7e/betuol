package event

import (
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/math"
)

type AttackEvent struct {
	Char1, Char2 component.GOiD
}

func (at AttackEvent) GetType() string {
	return "attack"
}

type DeathEvent struct {
	Id component.GOiD
}

func (dt DeathEvent) GetType() string {
	return "death"
}

type KillEvent struct {
	Killer, Dead component.GOiD
}

func (ke KillEvent) GetType() string {
	return "kill"
}

type QuestComplete struct {
	Id        component.GOiD
	QuestName string
}

func (qt QuestComplete) GetType() string {
	return "questComplete"
}

type ChatEvent struct {
	Sender   component.GOiD
	Reciever component.GOiD
	Message  string
}

func (ce ChatEvent) GetType() string {
	return "chat"
}

type PlayerCreatedEvent struct {
	PlayerID component.GOiD
}

func (ce PlayerCreatedEvent) GetType() string {
	return "playerCreated"
}

type CharacterMoveEvent struct {
	CharID      component.GOiD
	NewLocation math.Vec3
}

func (cme CharacterMoveEvent) GetType() string {
	return "characterMoved"
}

type RunAiEvent struct {
}

func (rae RunAiEvent) GetType() string {
	return "runAi"
}
