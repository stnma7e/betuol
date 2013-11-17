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
	rm *res.ResourceManager

	modellink chan ModelTransfer
	errorlink chan error

	compList         *common.Vector
	graphicsHandlers *common.Vector
}

type ModelTransfer struct {
	id component.GOiD
	gc GraphicsComponent
}
type GraphicsComponent struct {
	ModelName, Mesh, MeshType, Renderer string
}

func MakeGraphicsManager(window *GlGraphicsManager, rm *res.ResourceManager) *GraphicsManager {
	gm := &GraphicsManager{
		rm,
		make(chan ModelTransfer),
		make(chan error),
		common.MakeVector(),
		common.MakeVector(),
	}

	gm.graphicsHandlers.Insert(window)

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

	gm.modellink <- ModelTransfer{id, obj}

	return <-gm.errorlink
}

func (gm *GraphicsManager) CreateComponent(id component.GOiD, gc GraphicsComponent) error {
	graphicsHandlers := gm.graphicsHandlers.Array()
	for i := range graphicsHandlers {
		if err := graphicsHandlers[i].(GraphicsHandler).LoadModel(id, gc); err != nil {
			return err
		}
	}

	gm.compList.Insert(id)

	return nil
}

func (gm *GraphicsManager) DeleteComponent(id component.GOiD) {
}

func (gm *GraphicsManager) Tick() (ret bool) {
	graphicsHandlers := gm.graphicsHandlers.Array()
	for i := range graphicsHandlers {
		ret = graphicsHandlers[i].(GraphicsHandler).Tick()
	}
	for i := true; i; {
		select {
		case modelTrans := <-gm.modellink:
			gm.errorlink <- gm.CreateComponent(modelTrans.id, modelTrans.gc)
		default:
			i = false
		}
	}

	return
}

func (gm *GraphicsManager) RenderAll(camera *math.Frustum, sm component.SceneManager) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	graphicsHandlers := gm.graphicsHandlers.Array()
	for i := range graphicsHandlers {
		graphicsHandlers[i].(GraphicsHandler).Render(gm.compList, sm, camera)
		// only send id's who will be visible
	}
}

func (gm *GraphicsManager) HandleInputs(eye, target, up math.Vec3) (math.Vec3, math.Vec3, math.Vec3) {
	graphicsHandlers := gm.graphicsHandlers.Array()
	for i := range graphicsHandlers {
		/* inputs := */ graphicsHandlers[i].(GraphicsHandler).HandleInputs()
		// handle inputs
	}

	return graphicsHandlers[0].(*GlGraphicsManager).HandleInputs0(eye, target, up)
	//hack to keep the camera movement
}

func (gm *GraphicsManager) DrawString(x, y float32, text string) {
	graphicsHandlers := gm.graphicsHandlers.Array()
	for i := range graphicsHandlers {
		graphicsHandlers[i].(GraphicsHandler).DrawString(x, y, text)
	}
}

func (gm *GraphicsManager) GetSize() (int, int) {
	graphicsHandlers := gm.graphicsHandlers.Array()
	for i := range graphicsHandlers {
		return graphicsHandlers[i].(GraphicsHandler).GetSize()
	}

	return 0, 0
}

func (gm *GraphicsManager) RegisterGraphicsHandler(handler GraphicsHandler) {
	gm.graphicsHandlers.Insert(handler)
}
