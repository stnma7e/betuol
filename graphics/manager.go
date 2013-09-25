package graphics

import (
	"encoding/json"

	"smig/math"
	"smig/component"
	"smig/res"

	"github.com/go-gl/gl"
)

const (
	RESIZESTEP = 1
)


type GraphicsManager struct {
	window WindowManager
	rm *res.ResourceManager

	models []Model

	modellink chan ModelTransfer
	errorlink chan error
}

type ModelTransfer struct {
	id component.GOiD
	model *GraphicsComponent
}
type GraphicsComponent struct {
	ModelName, Mesh, MeshType, Vertex, Fragment string
}

func MakeGraphicsManager(window WindowManager, rm *res.ResourceManager, modellink chan ModelTransfer, errorlink chan error) *GraphicsManager {
	return &GraphicsManager{
		window, rm,
		make([]Model, 0),
		modellink,
		errorlink,
	}
}

func (gm *GraphicsManager) JsonCreate(id component.GOiD, compData []byte) error {
	obj := GraphicsComponent{}
	json.Unmarshal(compData, &obj)

	gm.modellink <- ModelTransfer{ id, &obj }

	return <-gm.errorlink
}

func (gm *GraphicsManager) CreateComponent(id component.GOiD, model Model) error {
	gm.resizeArrays(id)
	gm.models[id] = model

	return nil
}

func (gm *GraphicsManager) resizeArrays(id component.GOiD) {
	if cap(gm.models) - 1 < int(id) {
		newModels := gm.models
		gm.models = make([]Model, id + RESIZESTEP)
		for i := range newModels {
			gm.models[i] = newModels[i]
		}
	}
}

func (gm *GraphicsManager) DeleteComponent(id component.GOiD) {
	gm.models[id] = Model{}
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
			gm.errorlink <- gm.CreateComponent(gc.id, gm.window.LoadModel(gc.model, gm.rm))
		default:
			i = false
		}
	}

	return true
}

func (gm *GraphicsManager) RenderAll(camera *math.Frustum, sm *component.SceneManager) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for i := 1; i < len(gm.models); i++ {
		transMat := sm.GetTransformMatrix(component.GOiD(i))
		bounding := sm.GetBoundingSphere(component.GOiD(i))
		if camera.ContainsSphere(bounding) > math.OUTSIDE {
			gm.window.Render(gm.models[i], transMat, camera.LookAtMatrix(), camera.Projection())
		}
	}
}