package main

import (
	"fmt"
	"time"
	"os"
	"flag"
	"runtime"
	"runtime/pprof"
	"log"

	"smig/component"
	"smig/component/transform"
	"smig/res"
	// "smig/component/ai"
	"smig/component/physics"
	"smig/graphics"
	"smig/math"
	// "smig/common"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

type Obj struct {
	id component.GOiD
	sp math.Sphere
}

var rm  *res.ResourceManager
var gof *component.GameObjectFactory
var tm  *transform.SceneManager
var pm  *physics.PhysicsManager
var running bool

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
	if err != nil {
		log.Fatal(err)
	}
		runtime.SetCPUProfileRate(100)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}


	returnlink := make(chan bool)

	gof = component.MakeGameObjectFactory()
	tm  = transform.MakeSceneManager()
	pm  = physics.MakePhysicsManager(tm)
	rm := res.MakeResourceManager("/home/sam/go/data/")
	// am := ai.MakeAiManager()
	gof.Register(component.SceneType, tm, tm.JsonCreate)
	gof.Register("physics", pm, pm.JsonCreate)
	// gof.Register("ai", am)

	jmap := rm.LoadJsonMap("map1")
	fmt.Println(jmap)

	glg := graphics.GlStart(640, 480, "smig")

	vertStr := string(rm.GetFileContents("graphics/shader/" + "MatrixTransform.vert"))
	fragStr := string(rm.GetFileContents("graphics/shader/" + "White.frag"))
	vert := glg.MakeShader(gl.VERTEX_SHADER, vertStr)
	frag := glg.MakeShader(gl.FRAGMENT_SHADER, fragStr)
	shaders := make([]gl.Shader, 2)
	shaders[0] = vert
	shaders[1] = frag

	program := gl.CreateProgram()
	attribArray := glg.LinkProgram(program, shaders)

	mat := math.MakePerspectiveMatrix(1, 10, 60, 1)
	mat.MakeIdentity()

	verts := make([]math.Vec3, 3)
	verts[0] = math.Vec3{0.75, 0.75, 0.0}
	verts[1] = math.Vec3{0.75, -0.75, 0.0}
	verts[2] = math.Vec3{-0.75, -0.75, 0.0}
	indicies := make([]uint32, 3)
	indicies[0] = 0
	indicies[1] = 1
	indicies[2] = 2
	uv := make([]math.Vec2, 1)
	normals := make([]math.Vec3, 1)
	tex := make([]gl.Texture, 1)
	mesh := graphics.MakeMesh(program, attribArray, verts, indicies, uv, normals, tex)

	go func() {
		running = true

		deerStr := rm.GetFileContents("map/gameobject/deer/obj.json")
		components := rm.LoadGameObject(deerStr, "deer")
		fmt.Println(components)

		list := gof.CreateFromMap(&jmap)
		// list := [5000]byte{}

		// fmt.Println(list)

		var oldTime time.Time
		for i := range list {
			// list[i], err := gof.Create(components)
			// if err != nil {
			// 	panic(err)
			// }
			pm.AddForce(list[i], &math.Vec3{-1, -1, -1})
			pm.AddForce(list[i], &math.Vec3{3, 0, -2})
		}

		oldTime = time.Now()
		for running {
			time.Sleep(250 * time.Millisecond)

			newTime := time.Since(oldTime)
			secs := newTime.Seconds()
			fmt.Println(newTime)

			pm.Tick(secs)
			tm.Tick(secs)
			// tm.Render(glg)

			for i := range list {
				id := list[i]
				trans,_ := tm.GetTransform(component.GOiD(id))
				fmt.Println(id," ",trans.ToString())
			}
			fmt.Println()

			oldTime = oldTime.Add(newTime)
		}

		// const amt int = 1000
		// var list  [amt]Obj
		// var bools [amt]bool

		// frust := math.MakeFrustum(1,10,60,1)
		// frust.LookAt(&math.Vec3{0,0,0}, &math.Vec3{0,0,1}, &math.Vec3{0,0,0})

		// for i := range bools {
		// 	tm.Tick(1.0)
		// 	pm.Tick(1.0)
		// 	bools[i] = false
		// 	list[i].id,_ = gof.Create(compList)
		// 	list[i].sp   = math.Sphere{math.Vec3{0,0,float32(-i)}, 1}
		// 	is := frust.IsSphereInside(&list[i].sp)
		// 	// fmt.Println(is)
		// 	switch is {
		// 	case 1:
		// 		bools[i] = true
		// 	case 2:
		// 		bools[i] = true
		// 	}
		// }
		// for i := range bools {
		// 	if bools[i] != true {
		// 		fmt.Println(i)
		// 	}
		// }

		// vec2  := math.Vec3{0,0,-10}
		// frust = math.MakeFrustum(1,100,60,1)
		// frust.LookAt(&math.Vec3{0,0,0}, &math.Vec3{0,0,1}, &math.Vec3{0,1,0})
		// fmt.Println(frust.IsPointInside(&vec2))
		// sp := math.Sphere{math.Vec3{0,0,0}, 2}
		// fmt.Println(frust.IsSphereInside(&sp))

		returnlink <- true
	}()

	for i := true; i; {
		glg.Render(mesh, &mat)
		i = glg.Run()
	}
	glfw.Terminate()
	running = false

	<-returnlink
}
