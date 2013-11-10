package character

import (
	"fmt"
	"math/rand"

	"smig/component"
	"smig/common"
	"smig/event"
)

func chat(r, s component.GOiD, m string) {
	if r == 0 {
		fmt.Printf("%v says to WORLD: %s\n", s, m)
	} else {
		fmt.Printf("%v says to %v: %s\n", s, r, m)
	}
}
func (cm *CharacterManager) HandleChat(evt event.Event) {
	cevt := evt.(event.ChatEvent)
	switch cevt.Reciever {
	case 0:
		chat(0, cevt.Sender, cevt.Message)
	default:
		chat(cevt.Reciever, cevt.Sender, cevt.Message)
	}
}

func (cm *CharacterManager) HandleAttack(evt event.Event) {
	atevt := evt.(event.AttackEvent)
	attr1 := cm.GetCharacterAttributes(atevt.Char1)
	attr2 := cm.GetCharacterAttributes(atevt.Char2)

	if attr2.Attributes[HEALTH] <= 0 {
		common.LogErr.Print(atevt.Char2, " has a health below 0 during an attack.")
		return
	}
	hit := (rand.Float32() / 4) * attr1.Attributes[STRENGTH]
	fmt.Println(hit)
	attr2.Attributes[HEALTH] -= hit
	if attr2.Attributes[HEALTH] <= 0 {
		cm.em.Send(event.KillEvent{atevt.Char1, atevt.Char2})
	}
	cm.UpdateId(atevt.Char2, attr2)
}
