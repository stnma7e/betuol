package graphics

import (
	"fmt"
	"unsafe"
	gomath "math"
	"time"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"

	"smig/math"
	"smig/common"
	"smig/res"
)

type GlGraphicsManager struct {
	window *glfw.Window
	modelMap map[string]*Model
}

func GlStart(sizeX, sizeY int, title string) *GlGraphicsManager {
	if !glfw.Init() {
		common.Log.Error("Can't init glfw")
	}

	glg := GlGraphicsManager{}
	glg.modelMap = make(map[string]*Model)

	window, err := glfw.CreateWindow(sizeX,sizeY,title,nil,nil)
	if err != nil {
		common.Log.Error(err)
	}
	glg.window = window
	glg.MakeContextCurrent()
	glg.window.SetKeyCallback(GlfwKeyCallback)

	gl.Enable(gl.CULL_FACE);
	gl.CullFace(gl.BACK);
	gl.FrontFace(gl.CCW);

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.DEPTH_CLAMP)
	gl.DepthFunc(gl.LESS)
	gl.DepthRange(0.0, 1.0)

	gl.ClearColor(0.0,0.6,0.6,0.0)
	gl.ClearDepth(1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Viewport(0,0,sizeX,sizeY)

	return &glg
}

func (glg *GlGraphicsManager) Tick() {
	glg.window.SwapBuffers()
	glfw.PollEvents()

	x,y := glg.window.GetSize()
	gl.Viewport(0,0,x,y)
}

func (glg *GlGraphicsManager) ShouldClose() {
	glg.window.SetShouldClose(true)
}
func (glg *GlGraphicsManager) Closing() bool {
	return glg.window.ShouldClose()
}
func (glg *GlGraphicsManager) close() {
	glg.window.Destroy()
}
func (glg *GlGraphicsManager) MakeContextCurrent() {
	glg.window.MakeContextCurrent()
}

func GlfwKeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape {
		window.SetShouldClose(true)
	}
}

func (glg *GlGraphicsManager) HandleInputs(eye math.Vec3) math.Vec3 {
	mouse_x, mouse_y   := glg.window.GetCursorPosition()
	window_x, window_y := glg.window.GetSize()

	horizontalAngle := 3.14159 + 0.001 * float64(window_x / 2) - mouse_x
	verticalAngle   := 0.001 * float64(window_y / 2) - mouse_y

	direction := math.Vec3{
		float32(gomath.Cos(verticalAngle) * gomath.Sin(horizontalAngle)),
		float32(gomath.Sin(verticalAngle)),
		float32(gomath.Cos(verticalAngle) * gomath.Cos(horizontalAngle)),
	}
	right := math.Vec3 {
		float32(gomath.Sin(horizontalAngle - 3.14159 / 2.0)),
		0,
		float32(gomath.Cos(horizontalAngle - 3.14159 / 2.0)),
	}
	fmt.Println(right)

	fmt.Println(eye)
	if glg.window.GetKey(glfw.KeyUp) == glfw.Press {
		eye = math.Add3v3v(eye, math.Mult3vf(direction, 0.1))
		fmt.Println("up")
	}	
	if glg.window.GetKey(glfw.KeyDown) == glfw.Press {
		eye = math.Sub3v3v(eye, math.Mult3vf(direction, 0.1))
		fmt.Println("down")
	}	
	if glg.window.GetKey(glfw.KeyLeft) == glfw.Press {
		eye = math.Add3v3v(eye, math.Mult3vf(right, 0.1))
		fmt.Println("left")
	}	
	if glg.window.GetKey(glfw.KeyRight) == glfw.Press {
		eye = math.Sub3v3v(eye, math.Mult3vf(right, 0.1))
		fmt.Println("right")
	}
	return eye
}

func (glg *GlGraphicsManager) Render(mh Model, transMat, camMat, projectMat math.Mat4x4) {
	if mh.program == 0 {
		return
	}

	mh.vao.Bind()
	mh.program.Use()

	lightDirection := math.Vec3{0.866, 0.5, 0.0}
	lightDirCameraSpace := math.Mult4m3v(camMat, lightDirection)
	lightPosition := math.Vec3{40,20,20}
	lightPosCameraSpace := math.Normalize3v(math.Mult4m3v(camMat, lightPosition))

	var transArray [16]float32 = [16]float32(math.Mult4m4m(camMat, transMat))
	var projection [16]float32 = [16]float32(projectMat)
	var normArray  [9]float32
	for i := range normArray {
		normArray[i] = transArray[i]
	}

	mh.uniforms[WORLD].UniformMatrix4f(false, &transArray)
	mh.uniforms[PROJ].UniformMatrix4f(false, &projection)
	mh.uniforms[NORMALMV].UniformMatrix3f(false, &normArray)
	mh.uniforms[DIRLIGHT].Uniform3fv(1, lightDirCameraSpace[:])
	mh.uniforms[LINTENSE].Uniform4f(0.8, 0.8, 0.8, 1.0)
	mh.uniforms[LPOSITION].Uniform4fv(1, lightPosCameraSpace[:])
	mh.uniforms[AMBINTENSE].Uniform4f(0.2, 0.2, 0.2, 0.2)

	mh.buffers[VERT].Bind(gl.ARRAY_BUFFER)
	mh.attribArray[VERT].EnableArray()
	mh.attribArray[VERT].AttribPointer(3, gl.FLOAT, false, 0, nil)
	mh.buffers[VERT].Unbind(gl.ARRAY_BUFFER)

	mh.buffers[NORM].Bind(gl.ARRAY_BUFFER)
	mh.attribArray[NORM].EnableArray()
	mh.attribArray[NORM].AttribPointer(3, gl.FLOAT, false, 0, nil)
	mh.buffers[NORM].Unbind(gl.ARRAY_BUFFER)

	mh.buffers[INDEX].Bind(gl.ELEMENT_ARRAY_BUFFER)

	// gl.DrawArrays(gl.TRIANGLES,0,3)
	sizeptr := gl.GetBufferParameteriv(gl.ELEMENT_ARRAY_BUFFER, gl.BUFFER_SIZE)
	size := sizeptr / int32(unsafe.Sizeof(float32(1)))
	gl.DrawElements(gl.TRIANGLES, int(size), gl.UNSIGNED_INT, nil)

	mh.buffers[INDEX].Unbind(gl.ELEMENT_ARRAY_BUFFER)
	mh.attribArray[VERT].DisableArray()
	mh.attribArray[NORM].DisableArray()
}

//*****************************
// TOOLS
//*****************************

func (glg *GlGraphicsManager) LoadModel(comp *GraphicsComponent, rm *res.ResourceManager) Model {
	oldTime := time.Now()

	modelPtr, ok := glg.modelMap[comp.ModelName]
	if ok {
		return *modelPtr
	} else {
		fmt.Println("mesh not yet loaded: ", comp.Mesh)
	}

	vertStr := string(rm.GetFileContents("graphics/shader/" + comp.Vertex + ".vert"))
	fragStr := string(rm.GetFileContents("graphics/shader/" + comp.Fragment + ".frag"))
	shaders := make([]gl.Shader, 2)
	shaders[0] = glg.MakeShader(gl.VERTEX_SHADER, vertStr)
	shaders[1] = glg.MakeShader(gl.FRAGMENT_SHADER, fragStr)
	program := gl.CreateProgram()
	attribArray := glg.LinkProgram(program, shaders)

	var vertsVector, indiciesVector, normsVector, texVector *common.Vector
	switch comp.MeshType {
		case "wavefront":
			vertsVector, indiciesVector, normsVector, texVector = rm.LoadModelWavefront(comp.Mesh)
	}

	vts := vertsVector.Array()
	verts := make([]math.Vec3, len(vts))
	for i := range verts {
		verts[i] = vts[i].(math.Vec3)
	}
	inx := indiciesVector.Array()
	indicies := make([]uint32, len(inx))
	for i := range indicies {
		indicies[i] = inx[i].(uint32)
	}
	nms := normsVector.Array()
	norms := make([]math.Vec3, len(nms))
	for i := range norms {
		norms[i] = nms[i].(math.Vec3)
	}
	uv := texVector.Array()
	uvs := make([]math.Vec2, len(uv))
	for i := range uvs {
		uvs[i] = uv[i].(math.Vec2)
	}
	if len(uvs) == 0 {
		uvs = make([]math.Vec2, 1)
	}
	
	tex := make([]gl.Texture, 1)

	model := MakeModel(program, attribArray, verts, indicies, norms, uvs, tex)
	glg.modelMap[comp.ModelName] = &model

	fmt.Println("model loaded", time.Since(oldTime))
	return model
}

func (glg *GlGraphicsManager) MakeShader(shadType gl.GLenum, shadStr string) gl.Shader {
	gl.Init()
	shader := gl.CreateShader(shadType)
	shader.Source(shadStr)
	shader.Compile()
	if ok := shader.Get(gl.COMPILE_STATUS); ok < 1 {
		log := shader.GetInfoLog()
		shader.Delete()
		common.Log.Error(fmt.Sprintf("failed to compile shader. type: %e\n\t%s", shadType, log))
	}

	return shader
}

func (glg *GlGraphicsManager) LinkProgram(program gl.Program, shaderList []gl.Shader) [NUMATTR]gl.AttribLocation {
	for i := range shaderList {
		program.AttachShader(shaderList[i])
	}

	attribArray := [NUMATTR]gl.AttribLocation{}
	for i := range attribArray {
		attribArray[i] = gl.AttribLocation(i)
	}

	program.BindAttribLocation(VERT,"position")
	program.BindAttribLocation(NORM,"normal")
	// program.BindAttribLocation(COLOR,"color")
	// program.BindAttribLocation(UV,"uv")

	program.Link()

	if ok := program.Get(gl.LINK_STATUS); ok < 1 {
		log := program.GetInfoLog()
		common.Log.Error(fmt.Sprintf("failed to link program\n\t%v", log))
	}
	for i := range shaderList {
		program.DetachShader(shaderList[i])
		shaderList[i].Delete()
	}

	return attribArray
}