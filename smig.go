package main

import (
	"smig/server"
)

func main() {

	returnlink := make(chan bool)
	sv := server.MakeServer(returnlink)
	go sv.Loop()

	<-returnlink
}
