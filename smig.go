package main

import (
        "fmt"
        "time"

	"smig/instance"
	"smig/res"
	"smig/graphics"
	"smig/math"
        "smig/net"
        "smig/event"

	glfw "github.com/go-gl/glfw3"
)

const X, Y = 640, 480

func main() {
	returnlink := make(chan bool)
	rm := res.MakeResourceManager("/home/sam/go/src/smig/data/")

	target, eye, up := math.Vec3{0,0,0}, math.Vec3{0, 6, -12}, math.Vec3{0,1,0}
        cam := math.MakeFrustum(0.1, 100, 60, 1/1)
        cam.LookAt(target, eye, up)
        //mat := math.Mult4m4m(cam.LookAtMatrix(), cam.Projection())
        //fmt.Println(mat)
        //fmt.Println(mat.Inverse())
        //graphics.Trace(15,15, mat.Inverse())
        //return

        nm := net.MakeNetworkManager()
        err := nm.Connect("localhost:13572")
        if err != nil {
            fmt.Println(err)
        }
        nm.SendBytes([]byte("hello, world"))
        nm.Send(event.AttackEvent{1, 1})

	glg   := graphics.GlStart(X, Y, "smig", rm)
	gm    := graphics.MakeGraphicsManager(glg, rm)

	in := instance.MakeInstance(returnlink, rm, gm, nm)
	go in.Loop()


        oldtime := time.Now()
	for i := true; i; {
                secs := time.Since(oldtime).Seconds()
                oldtime = time.Now()
                fpsStr := fmt.Sprintf("%f", 100 / secs)
                spfStr := fmt.Sprintf("%f", secs / 100)

		x, y := gm.GetSize()
		cam := math.MakeFrustum(0.1, 100, 60, float32(y)/float32(x))

		for j := 0; j < 100 && i; j++ {
                        //nm.Tick()
			eye, target, up = gm.HandleInputs(eye, target, up)
			cam.LookAt(target, eye, up)

                        tm := in.GetSceneManagerSnapshot()
			gm.RenderAll(cam, tm)

                        gm.DrawString(10, 10, "fps: " + fpsStr)
                        gm.DrawString(10, 25, "spf: " + spfStr)

			i = gm.Tick()
			select {
			case <-returnlink:
				i = false
			default:
			}
		}
	}
	glfw.Terminate()
}

// loadobj player 0,0,0
// player attack 3
