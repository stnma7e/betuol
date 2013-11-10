package instance

import (
	"smig/event"
)

func (is *Instance) StartScript() {
	goList := is.CreateFromMap("map1")
	is.em.Send(event.ChatEvent{goList[0], is.player, "Good god! You're finnally awake. I've been sitting here for hours. I thought you'd never wake up."})
	is.em.Send(event.ChatEvent{goList[0], is.player, "Let's see if you can get up."})
	is.qm.AddQuest(is.player, is.qm.FirstMoveQuest)
	is.em.Send(event.ChatEvent{goList[0], is.player, "Type 'player north' to move forward."})
	is.em.Send(event.ChatEvent{goList[0], is.player, "Type 'player south' to move backward."})
	is.em.Send(event.ChatEvent{goList[0], is.player, "Type 'player east' to move right."})
	is.em.Send(event.ChatEvent{goList[0], is.player, "Type 'player west' to move left."})
}
