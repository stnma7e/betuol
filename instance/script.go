package instance

import (
	//"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/event"
	//"github.com/stnma7e/betuol/math"
)

func (is *Instance) StartScript() error {
	goList, err := is.CreateFromMap("map1")
	if err != nil {
		return err
	}
	//is.em.Send(event.ChatEvent{goList[0], is.player, "Good god! You're finnally awake. I've been sitting here for hours. I thought you'd never wake up."})
	//is.em.Send(event.ChatEvent{goList[0], is.player, "Let's see if you can get up."})
	//is.qm.AddQuest(is.player, is.qm.FirstMoveQuest)
	//is.em.Send(event.ChatEvent{0, is.player, "Type 'player north' to move forward."})
	//is.em.Send(event.ChatEvent{0, is.player, "Type 'player south' to move backward."})
	//is.em.Send(event.ChatEvent{0, is.player, "Type 'player east' to move right."})
	//is.em.Send(event.ChatEvent{0, is.player, "Type 'player west' to move left."})

	//_, err = is.CreateObject("enemy", math.Vec3{00, 00, 05})
	//if err != nil {
	//common.LogErr.Println(err)
	//}

	// fade in
	is.em.Send(event.ChatEvent{goList[0], is.player, "*hushed* Hey, the jury is returning. They'll be giving the verdict to the judge."})
	is.em.Send(event.ChatEvent{goList[0], is.player, "Judge, the jury has reached a decision. Those in the jury have found the defendant guilty."})

	return nil
}
