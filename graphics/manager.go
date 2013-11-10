package graphics

import (
	"encoding/json"
	//"fmt"

	"smig/common"
	"smig/component"
	"smig/math"
	"smig/res"

	"github.com/go-gl/gl"
)

const (
	RESIZESTEP = 1
)

type GraphicsManager struct {
	window *GlGraphicsManager
	rm     *res.ResourceManager

	models []*Model

	modellink     chan ModelTransfer
	errorlink     chan error
	renderTypeMap map[string]Renderer
	renderMap     map[string]*common.Vector
}

type ModelTransfer struct {
	id    component.GOiD
	model *GraphicsComponent
}
type GraphicsComponent struct {
	ModelName, Mesh, MeshType, Renderer string
}

func MakeGraphicsManager(window *GlGraphicsManager, rm *res.ResourceManager) *GraphicsManager {
	gm := &GraphicsManager{
		window, rm,
		make([]*Model, 0),
		make(chan ModelTransfer),
		make(chan error),
		make(map[string]Renderer),
		make(map[string]*common.Vector),
	}

	gm.RegisterRenderer("fragmentLighting", MakeFragmentPointLightingRenderer(rm, window))
	for k, _ := range gm.renderTypeMap {
		gm.renderMap[k] = common.MakeVector()
	}

	return gm
}

func (gm *GraphicsManager) JsonCreate(id component.GOiD, compData []byte) error {
	obj := GraphicsComponent{}
	err := json.Unmarshal(compData, &obj)
	if err != nil {
		common.LogErr.Println(err)
	}
	// load model before sending it to the main thread
	// you can send a goroutine to load the non-GL information,
	// then use the main thread to do all the OpenGL related stuff

	gm.modellink <- ModelTransfer{id, &obj}

	return <-gm.errorlink
}

func (gm *GraphicsManager) CreateComponent(id component.GOiD, model *Model, renderer string) error {
	gm.resizeArrays(id)
	gm.models[id] = model
	rend, ok := gm.renderMap[renderer]
	if !ok {
		common.LogErr.Printf("no renderer for type '%s'\n", renderer)
	} else {
		rend.Push_back(id, 1, 1)
	}

	return nil
}

func (gm *GraphicsManager) resizeArrays(id component.GOiD) {
	if cap(gm.models)-1 < int(id) {
		newModels := gm.models
		gm.models = make([]*Model, id+RESIZESTEP)
		for i := range newModels {
			gm.models[i] = newModels[i]
		}
	}
}

func (gm *GraphicsManager) DeleteComponent(id component.GOiD) {
	if len(gm.models) > int(id) {
		gm.models[id] = nil
	}
}

func (gm *GraphicsManager) Tick() bool {
	if gm.window.Closing() {
		gm.window.close()
		return false
	}
	gm.window.Tick()
	for i := true; i; {
		select {
		case gc := <-gm.modellink:
			gm.errorlink <- gm.CreateComponent(gc.id, gm.window.LoadModel(gc.model, gm.rm), gc.model.Renderer)
		default:
			i = false
		}
	}

	return true
}

func (gm *GraphicsManager) RenderAll(camera *math.Frustum, tm component.SceneManager) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	//for i := 1; i < len(gm.models); i++ {
	//fmt.Println(gm.models[i])
	//transMat := tm.GetTransform4m(component.GOiD(i))
	//bounding := tm.GetBoundingSphere(component.GOiD(i))
	//if gm.models[i] == nil {
	//continue
	//}
	//if camera.ContainsSphere(bounding) > math.OUTSIDE {
	//gm.window.Render(*gm.models[i], transMat, camera.LookAtMatrix(), camera.Projection())
	//}
	//}

	//fmt.Println(camera.LookAtMatrix(), camera.Projection())

	for k, v := range gm.renderTypeMap {
		modelsOfRenderer := gm.renderMap[k].Array()
		for i := range modelsOfRenderer {
			id, ok := modelsOfRenderer[i].(component.GOiD)
			// get first component.GOiD in the list
			if !ok || gm.models[int(id)] == nil {
				// if modelsOfRenderer[i] is nil (and cannot be type asserted)
				// or if the manager's modelList's space is nil
				// then the component was destroyed and it can be erased from the renderer's list
				gm.renderMap[k].Erase(i)
				continue
			}
			transMat := tm.GetTransform4m(id)
			//gm.window.Render(*gm.models[int(id)], transMat, camera.LookAtMatrix(), camera.Projection())
			v.Render(*gm.models[int(id)], transMat, camera.LookAtMatrix(), camera.Projection())
		}
	}
}

func (gm *GraphicsManager) HandleInputs(eye, target, up math.Vec3) (math.Vec3, math.Vec3, math.Vec3) {
	return gm.window.HandleInputs(eye, target, up)
}

func (gm *GraphicsManager) DrawString(x, y float32, text string) {
	gm.window.DrawString(x, y, text)
}

func (gm *GraphicsManager) GetSize() (int, int) {
	return gm.window.GetSize()
}

func (gm *GraphicsManager) RegisterRenderer(rendType string, rend Renderer) {
	gm.renderTypeMap[rendType] = rend
}
