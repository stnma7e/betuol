package instance

import (
    "smig/event"
)

func (is *Instance) StartScript() {
    goList := is.CreateFromMap("map1")
    is.em.Send(event.ChatEvent{ goList[0], 1, "Good god! You're finnally awake. I've been sitting here for hours. I thought you'd never wake up." })
}
