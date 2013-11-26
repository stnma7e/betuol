package main

import (
	//"fmt"

	"github.com/stnma7e/betuol/instance"
	//"github.com/stnma7e/betuol/math"
	"github.com/stnma7e/betuol/res"
)

func main() {
	rm := res.MakeResourceManager("./data/")

	//target, eye, up := math.Vec3{0, 0, 0}, math.Vec3{0, 6, -12}, math.Vec3{0, 1, 0}
	//cam := math.MakeFrustum(0.1, 100, 90, 1/1)
	//cam.LookAt(target, eye, up)
	//mat := math.Mult4m4m(cam.LookAtMatrix(), cam.Projection())
	//fmt.Println(mat)
	//fmt.Println(mat.Inverse())
	//graphics.Trace(15,15, mat.Inverse())
	//return

	in := instance.MakeInstance(rm)
	in.Loop()
}

// loadobj player 0,0,0
// player attack 3
