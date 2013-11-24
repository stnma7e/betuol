package character

import (
	"fmt"
	"math/rand"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/event"
)

func chat(r, s component.GOiD, m string) {
	if r == 0 {
		fmt.Printf("%v says to WORLD: %s\n", s, m)
	} else {
		fmt.Printf("%v says to %v: %s\n", s, r, m)
	}
}

// HandleChat responds to event.ChatEvents's and prints a message based on the data in the event.
func (cm *CharacterManager) HandleChat(evt event.Event) {
	cevt := evt.(event.ChatEvent)
	switch cevt.Reciever {
	case 0:
		chat(0, cevt.Sender, cevt.Message)
	default:
		chat(cevt.Reciever, cevt.Sender, cevt.Message)
	}
}

// HandleAttack responds to event.AttackEvents that are launched by other component manager systems and updates character data based on outcomes of random hits.
func (cm *CharacterManager) HandleAttack(evt event.Event) {
	atevt := evt.(event.AttackEvent)
	attr1 := cm.GetCharacterAttributes(atevt.Char1)
	attr2 := cm.GetCharacterAttributes(atevt.Char2)

	if attr2.Attributes[HEALTH] <= 0 {
		common.LogErr.Print(atevt.Char2, " had a health below 0 when it was attacked.")
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
