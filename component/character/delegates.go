package character

import (
	"math/rand"
	"fmt"
        "strconv"

	"smig/event"
        "smig/common"
)

func chat(r, s, m string) {
  fmt.Printf("%s says to %s: %s\n", s, r, m)
}
func (cm *CharacterManager) HandleChat(evt event.Event) {
    cevt := evt.(event.ChatEvent)
    switch cevt.Reciever {
        case 0:
            chat("WORLD", strconv.Itoa(int(cevt.Sender)), cevt.Message)
        default:
            chat(strconv.Itoa(int(cevt.Reciever)), strconv.Itoa(int(cevt.Sender)), cevt.Message)
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
            cm.em.Send(event.KillEvent{ atevt.Char1, atevt.Char2 })
        }
	cm.UpdateId(atevt.Char2, attr2)
}