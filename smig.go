package main

import (
	"smig/instance"
	"smig/res"
)

func main() {
	returnlink := make(chan bool)
	rm := res.MakeResourceManager("/home/sam/go/src/smig/data/")
	in := instance.MakeInstance(returnlink, rm)
	go in.Loop()

	<-returnlink
}
