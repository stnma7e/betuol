package graphics

import (
	"fmt"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"

	"smig/math"
)

type GlGraphics struct {
	window *glfw.Window
}

func GlStart(sizeX, sizeY int, title string) *GlGraphics {
	if !glfw.Init() {
		panic("Can't init glfw")
	}

	glg := GlGraphics{}
	window, err := glfw.CreateWindow(sizeX,sizeY,title,nil,nil)
	if err != nil {
		panic(err)
	}
	glg.window = window
	glg.MakeContextCurrent()
	glg.window.SetKeyCallback(GlfwKeyCallback)

	gl.ClearColor(0.0,0.0,0.0,0.0)
	gl.ClearDepth(1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Viewport(0,0,sizeX,sizeY);

	return &glg
}

func (glg *GlGraphics) Run() bool {
	if glg.window.ShouldClose() {
		glg.close()
		return false
	}
	glg.window.SwapBuffers()
	glfw.PollEvents()

	return true
}

func (glg *GlGraphics) ShouldClose() {
	glg.window.SetShouldClose(true)
}
func (glg *GlGraphics) Closing() bool {
	return glg.window.ShouldClose()
}
func (glg *GlGraphics) close() {
	glg.window.Destroy()
}
func (glg *GlGraphics) MakeContextCurrent() {
	glg.window.MakeContextCurrent()
}

func GlfwKeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape {
		window.SetShouldClose(true)
	}
}

func (glg *GlGraphics) Render(mh *Mesh, transMat *math.Mat4x4) {
	if mh.program == 0 {
		return
	}

	mh.vao.Bind()
	mh.program.Use()

	var transArray [16]float32 = [16]float32(*transMat)
	mh.uniforms[MVP].UniformMatrix4f(false, &transArray)

	mh.buffers[VERT].Bind(gl.ARRAY_BUFFER)
	mh.attribArray[0].EnableArray()
	mh.attribArray[0].AttribPointer(3, gl.FLOAT, false, 0, nil)

	mh.buffers[INDEX].Bind(gl.ELEMENT_ARRAY_BUFFER)

	// gl.DrawArrays(gl.TRIANGLES,0,3)
	gl.DrawElements(gl.TRIANGLES,3,gl.UNSIGNED_INT,nil)

	mh.buffers[VERT].Unbind(gl.ARRAY_BUFFER)
	mh.buffers[INDEX].Unbind(gl.ELEMENT_ARRAY_BUFFER)
	mh.attribArray[0].DisableArray()
}

//*****************************
// TOOLS
//*****************************

func (glg *GlGraphics) MakeShader(shadType gl.GLenum, shadStr string) gl.Shader {
	gl.Init()
	shader := gl.CreateShader(shadType)
	shader.Source(shadStr)
	shader.Compile()
	if ok := shader.Get(gl.COMPILE_STATUS); ok < 1 {
		log := shader.GetInfoLog()
		shader.Delete()
		panic(fmt.Sprintf("failed to compile shader. type: %e\n\t%s", shadType, log))
	}

	return shader
}

func (glg *GlGraphics) LinkProgram(program gl.Program, shaderList []gl.Shader) []gl.AttribLocation {
	for i := range shaderList {
		program.AttachShader(shaderList[i])
	}

	attribArray := make([]gl.AttribLocation,1)
	for i := range attribArray {
		attribArray[i] = gl.AttribLocation(i)
	}

	program.BindAttribLocation(0,"position")
	// program.BindAttribLocation(1,"normal")
	// program.BindAttribLocation(2,"color")
	// program.BindAttribLocation(3,"uv")

	program.Link()

	if ok := program.Get(gl.LINK_STATUS); ok < 1 {
		log := program.GetInfoLog()
		panic(fmt.Sprintf("failed to link program\n\t%v", log))
	}
	for i := range shaderList {
		program.DetachShader(shaderList[i])
		shaderList[i].Delete()
	}

	return attribArray
}