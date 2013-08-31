package main

import (
	// "fmt"
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
	"smig/component/character"
	"smig/graphics"
	"smig/math"
	// "smig/common"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

var rm  *res.ResourceManager
var gof *component.GameObjectFactory
var tm  *transform.SceneManager
var pm  *physics.PhysicsManager
var cm  *character.CharacterManager
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
	cm  = &character.CharacterManager{}
	pm  = physics.MakePhysicsManager(tm)
	rm := res.MakeResourceManager("/home/sam/go/data/")
	// am := ai.MakeAiManager()
	gof.Register(component.SceneType, tm, tm.JsonCreate)
	gof.Register("physics", pm, pm.JsonCreate)
	gof.Register("character", cm, cm.JsonCreate)
	// gof.Register("ai", am)

	go func() {
		running = true

		components := rm.LoadGameObject("player")
		id, _ := gof.Create(components)
		go character.StartPlayer(id, tm, cm)

		jmap := rm.LoadJsonMap("map1")
		gof.CreateFromMap(&jmap)

		// list := [5000]byte{}

		// fmt.Println(list)

		// for i := range list {
		// 	// list[i], err := gof.Create(components)
		// 	// if err != nil {
		// 	// 	common.Log.LogError(err)
		// 	// }
		// 	pm.AddForce(list[i], &math.Vec3{-1, -1, -1})
		// 	pm.AddForce(list[i], &math.Vec3{3, 0, -2})
		// }

		oldTime := time.Now()
		for running {
			time.Sleep(500 * time.Millisecond)

			newTime := time.Since(oldTime)
			secs := newTime.Seconds()
			// fmt.Println(newTime)

			pm.Tick(secs)
			tm.Tick(secs)

			// for i := range list {
			// 	id := list[i]
			// 	trans,_ := tm.GetTransform(component.GOiD(id))
			// 	fmt.Println(id," ",trans.ToString())
			// }
			// fmt.Println()

			oldTime = oldTime.Add(newTime)
		}


		returnlink <- true
	}()

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


	for i := true; i; {
		glg.Render(mesh, &mat)
		i = glg.Run()
	}
	glfw.Terminate()
	running = false

	<-returnlink
}
