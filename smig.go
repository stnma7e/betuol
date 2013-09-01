package main

import (
	"os"
	"flag"
	"runtime"
	"runtime/pprof"
	"log"

	"smig/game"
	"smig/graphics"
	"smig/math"
	"smig/res"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		runtime.SetCPUProfileRate(10000)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}


	returnlink := make(chan bool)
	gm := game.MakeGame("player", 10.0)
	go gm.Loop(returnlink)
	gm.CreateFromMap("map1")



	glg := graphics.GlStart(640, 480, "smig")
	rm := res.MakeResourceManager("/home/sam/go/data/")

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
	returnlink <- true
}
