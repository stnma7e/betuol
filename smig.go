package main

import (
	"smig/instance"
	"smig/res"
	"smig/graphics"
	"smig/math"

	glfw "github.com/go-gl/glfw3"
)

const X, Y = 640, 480

func main() {
	returnlink := make(chan bool)
	modellink  := make(chan graphics.ModelTransfer)
	errorlink  := make(chan error)
	rm := res.MakeResourceManager("/home/sam/go/src/smig/data/")

	glg   := graphics.GlStart(X, Y, "smig")
	gm 	  := graphics.MakeGraphicsManager(glg, rm, modellink, errorlink)

	in := instance.MakeInstance(returnlink, rm, gm)
	sm := in.GetSceneManager()
	go in.Loop()

	cam := math.MakeFrustum(0.1, 100, 60, 1)
	target, eye, up := math.Vec3{}, math.Vec3{0, 6, -12}, math.Vec3{0,1,0}
	cam.LookAt(target, eye, up)
	for i := true; i; {
		cam.LookAt(target, eye, up)
		gm.RenderAll(cam, sm)

		i = gm.Tick()
		// eye = glg.HandleInputs(eye)
		select {
		case <-returnlink:
			i = false
		default:
		}

	}
	glfw.Terminate()
}

// loadobj player 0,0,0 1