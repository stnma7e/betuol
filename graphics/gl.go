package graphics

import (
	//"unsafe"
	gomath "math"
	"os"
	"time"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/gltext"

	"smig/common"
	"smig/component"
	"smig/math"
	"smig/res"
)

const NOPROGRAM gl.Program = 0
const NOVAO gl.VertexArray = 0

var g_tex gl.Texture //= GlLoadTexture("/home/sam/downloads/tower/tower_diffuse.png")

type GlGraphicsManager struct {
	rm *res.ResourceManager

	window   *glfw.Window
	modelMap map[string]*Model

	loadedFont *gltext.Font
	program    gl.Program

	modelList []*Model

	renderMap   map[component.GOiD]Renderer
	renderTypes map[string]Renderer
}

func GlStart(sizeX, sizeY int, title string, rm *res.ResourceManager) *GlGraphicsManager {
	if !glfw.Init() {
		common.LogErr.Fatal("GLFW init failed.")
	}

	glg := GlGraphicsManager{}
	glg.modelMap = make(map[string]*Model)

	window, err := glfw.CreateWindow(sizeX, sizeY, title, nil, nil)
	if err != nil {
		common.LogErr.Fatal(err)
	}
	glg.window = window
	glg.MakeContextCurrent()
	glg.window.SetKeyCallback(GlfwKeyCallback)

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.DEPTH_CLAMP)
	gl.DepthFunc(gl.LESS)
	gl.DepthRange(0.0, 1.0)

	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ClearDepth(1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Viewport(0, 0, sizeX, sizeY)

	g_tex = GlLoadTexture("/home/sam/downloads/tower/tower_diffuse.png")

	glg.renderMap = make(map[component.GOiD]Renderer)
	glg.renderTypes = make(map[string]Renderer)
	glg.renderTypes["fragmentLighting"] = MakeFragmentPointLightingRenderer(rm, &glg)

	return &glg
}

func (glg *GlGraphicsManager) Tick() bool {
	glg.SwapBuffers()
	glfw.PollEvents()

	x, y := glg.window.GetSize()
	gl.Viewport(0, 0, x, y)

	return !glg.Closing()
}

func (glg *GlGraphicsManager) SwapBuffers() {
	glg.window.SwapBuffers()
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
func (glg *GlGraphicsManager) GetSize() (int, int) {
	return glg.window.GetSize()
}

func GlfwKeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape {
		window.SetShouldClose(true)
	}
}

func (glg *GlGraphicsManager) HandleInputs() Inputs {
	return Inputs{}
}

func (glg *GlGraphicsManager) HandleInputs0(eye, target, up math.Vec3) (math.Vec3, math.Vec3, math.Vec3) {
	const STEPSIZE = 0.1
	const MOUSESPEED = 0.0001
	mouse_x, mouse_y := glg.window.GetCursorPosition()
	win_w, win_h := glg.window.GetSize()

	horizontalAngle := 3.14159 + (MOUSESPEED * (float64(win_w/2) - mouse_x))
	verticalAngle := MOUSESPEED * (float64(win_h/2) - mouse_y)
	direction := math.Vec3{
		float32(gomath.Cos(verticalAngle) * gomath.Sin(horizontalAngle)),
		float32(gomath.Sin(verticalAngle)),
		float32(gomath.Cos(verticalAngle) * gomath.Cos(horizontalAngle)),
	}
	right := math.Vec3{
		float32(gomath.Sin(horizontalAngle - (3.14159 / 2.0))),
		0,
		float32(gomath.Cos(horizontalAngle - (3.14159 / 2.0))),
	}
	up = math.Cross3v3v(right, direction)

	left := math.Cross3v3v(target, up)
	left = math.Mult3vf(math.Normalize3v(left), STEPSIZE)

	if glg.window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		if glg.window.GetKey(glfw.KeyUp) == glfw.Press {
			eye = math.Add3v3v(eye, math.Vec3{0, 0.1, 0})
		}
		if glg.window.GetKey(glfw.KeyDown) == glfw.Press {
			eye = math.Sub3v3v(eye, math.Vec3{0, 0.1, 0})
		}
	} else {
		if glg.window.GetKey(glfw.KeyUp) == glfw.Press {
			eye = math.Add3v3v(eye, math.Vec3{0, 0, 0.1})
		}
		if glg.window.GetKey(glfw.KeyDown) == glfw.Press {
			eye = math.Sub3v3v(eye, math.Vec3{0, 0, 0.1})
		}
		if glg.window.GetKey(glfw.KeyLeft) == glfw.Press {
			eye = math.Add3v3v(eye, math.Vec3{0.1, 0, 0})
		}
		if glg.window.GetKey(glfw.KeyRight) == glfw.Press {
			eye = math.Sub3v3v(eye, math.Vec3{0.1, 0, 0})
		}
	}
	return eye, target, up
}

func (glg *GlGraphicsManager) resizeArrays(id component.GOiD) {
	const RESIZESTEP = 1
	if cap(glg.modelList)-1 < int(id) {
		newCompList := make([]*Model, id+RESIZESTEP)
		for i := range glg.modelList {
			newCompList[i] = glg.modelList[i]
		}
		glg.modelList = newCompList
	}
}

//*****************************
// TOOLS
//*****************************

func (glg *GlGraphicsManager) LoadModel(id component.GOiD, comp GraphicsComponent) error {
	oldTime := time.Now()
	glg.resizeArrays(id)

	modelPtr, ok := glg.modelMap[comp.ModelName]
	if ok {
		glg.modelList[id] = modelPtr
		glg.renderMap[id] = glg.renderTypes[comp.Renderer]
		return nil
	} else {
		common.LogInfo.Println("mesh not yet loaded:", comp.Mesh)
	}

	var vertsVector, indiciesVector, normsVector, uvVector *common.Vector
	//var boundingRadius float32
	switch comp.MeshType {
	case "wavefront":
		vertsVector, indiciesVector, normsVector, uvVector /*boundingRadius*/, _ = glg.rm.LoadModelWavefront(comp.Mesh)
	}
	//fmt.Println(boundingRadius)

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
	uv := uvVector.Array()
	uvs := make([]math.Vec2, len(uv))
	for i := range uvs {
		uvs[i] = uv[i].(math.Vec2)
	}
	if len(uvs) == 0 {
		uvs = make([]math.Vec2, 1)
	}

	tex := make([]gl.Texture, 1)
	tex[0] = g_tex

	//fmt.Println(verts, inx)
	model := MakeModel(verts, indicies, norms, uvs, tex)
	glg.modelMap[comp.ModelName] = &model

	glg.modelList[id] = &model
	glg.renderMap[id] = glg.renderTypes[comp.Renderer]

	common.LogInfo.Println("model loaded", time.Since(oldTime))
	return nil
}

func (glg *GlGraphicsManager) Render(ids *common.Vector, sm component.SceneManager, cam *math.Frustum) {
	comps := ids.Array()
	for i := range comps {
		id := comps[i].(component.GOiD)
		loc := sm.GetTransform4m(id)
		glg.renderMap[id].Render(*glg.modelList[id], loc, cam.LookAtMatrix(), cam.Projection())
	}
}

func LoadShader(shadType gl.GLenum, shadStr string) gl.Shader {
	gl.Init()
	shader := gl.CreateShader(shadType)
	shader.Source(shadStr)
	shader.Compile()
	if shader.Get(gl.COMPILE_STATUS) < 1 {
		infoLog := shader.GetInfoLog()
		shader.Delete()
		common.LogErr.Fatalf("failed to compile shader type: %e\n\t%s", shadType, infoLog)
	}

	return shader
}

func LinkProgram(program gl.Program, shaderList []gl.Shader) [NUMATTR]gl.AttribLocation {
	for i := range shaderList {
		program.AttachShader(shaderList[i])
	}

	attribArray := [NUMATTR]gl.AttribLocation{}
	for i := range attribArray {
		attribArray[i] = gl.AttribLocation(i)
	}

	program.Link()

	if program.Get(gl.LINK_STATUS) < 1 {
		infoLog := program.GetInfoLog()
		common.LogErr.Printf("failed to link program\n\t%v", infoLog)
	}
	for i := range shaderList {
		program.DetachShader(shaderList[i])
		shaderList[i].Delete()
	}

	return attribArray
}

func (glg *GlGraphicsManager) LoadFont() {
	fd, err := os.Open("/home/sam/go/src/smig/data/AkashiMF.ttf")
	if err != nil {
		common.LogErr.Print(err)
	}
	defer fd.Close()

	glg.loadedFont, err = gltext.LoadTruetype(fd, 18, 32, 127, gltext.LeftToRight)
	if err != nil {
		common.LogErr.Print(err)
	}
}

func (glg *GlGraphicsManager) DrawString(x, y float32, text string) {
	const sample = "0 1 2 3 4 5 6 7 8 9 A B C D E F"
	if glg.loadedFont == nil {
		common.LogWarn.Print("font not yet loaded")
		glg.LoadFont()
	}
	_, h := glg.loadedFont.GlyphBounds()
	y = y + float32(h)
	sw, sh := glg.loadedFont.Metrics(sample)
	gl.Color4f(0.1, 0.1, 0.1, 0.7)
	gl.Rectf(x, y, x+float32(sw), y+float32(sh))

	// Render the string.
	gl.Color4f(1, 1, 1, 1)
	err := glg.loadedFont.Printf(x, y, text)
	if err != nil {
		//common.LogErr.Print(err)
		// spams error messages
		// NEED TO FIX LATER
	}
}
