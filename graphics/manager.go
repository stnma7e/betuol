package graphics

import (
	"encoding/json"
	//"fmt"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/math"
	"github.com/stnma7e/betuol/res"

	"github.com/go-gl/gl"
)

type modelTransfer struct {
	id component.GOiD
	gc graphicsComponent
}
type graphicsComponent struct {
	ModelName, Mesh, MeshType, Renderer, TextDescription string
}

// GraphicsManager is a component manager used to visualize the game onscreen.
// It uses multiple GraphicsHandlers to render the world in a variety of ways.
type GraphicsManager struct {
	rm *res.ResourceManager

	modellink chan modelTransfer
	errorlink chan error

	compList         *common.Vector
	graphicsHandlers *common.Vector
}

// MakeGraphicsManager returns a pointer to a GraphicsManager
func MakeGraphicsManager(window *GlGraphicsManager, rm *res.ResourceManager) *GraphicsManager {
	gm := &GraphicsManager{
		rm,
		make(chan modelTransfer),
		make(chan error),
		common.MakeVector(),
		common.MakeVector(),
	}

	gm.graphicsHandlers.Insert(window)
	gm.graphicsHandlers.Insert(MakeTextGraphicsHandler())

	return gm
}

// JsonCreate extracts creation data from a byte array of json text to pass to CreateComponent.
func (gm *GraphicsManager) JsonCreate(id component.GOiD, compData []byte) error {
	obj := graphicsComponent{}
	err := json.Unmarshal(compData, &obj)
	if err != nil {
		common.LogErr.Println(err)
	}
	// load model before sending it to the main thread
	// you can send a goroutine to load the non-GL information,
	// then use the main thread to do all the OpenGL related stuff

	gm.modellink <- modelTransfer{id, obj}

	return <-gm.errorlink
}

// Uses extracted data from higher level component creation functions and initializes a graphics component based on the id passed through.
// The function calls the LoadModel function of each GraphicsHandler in the manager's list.
func (gm *GraphicsManager) CreateComponent(id component.GOiD, gc graphicsComponent) error {
	graphicsHandlers := gm.graphicsHandlers.Array()
	for i := range graphicsHandlers {
		if err := graphicsHandlers[i].(GraphicsHandler).LoadModel(id, gc); err != nil {
			return err
		}
	}

	gm.compList.Insert(id)

	return nil
}

// DeleteComponent implements the component.ComponentManager interface and deletes graphics component data from the manager.
// The function calls the DeleteModel function of each GraphicsHandler in the manager's list.
func (gm *GraphicsManager) DeleteComponent(id component.GOiD) {
	comps := gm.compList.Array()
	graphicsHandlers := gm.graphicsHandlers.Array()
	for i := range comps {
		if comps[i] == id {
			gm.compList.Erase(i)
		}
	}
	for j := range graphicsHandlers {
		graphicsHandlers[j].(GraphicsHandler).DeleteModel(id)
	}
}

// Tick calls the Tick function of each GraphicsHandler in the manager's list.
// If any Tick functions return false, then GraphicsManager.Tick returns false.
func (gm *GraphicsManager) Tick() (ret bool) {
	for i := true; i; {
		select {
		case modelTrans := <-gm.modellink:
			gm.errorlink <- gm.CreateComponent(modelTrans.id, modelTrans.gc)
		default:
			i = false
		}
	}

	graphicsHandlers := gm.graphicsHandlers.Array()
	for i := range graphicsHandlers {
		ret = graphicsHandlers[i].(GraphicsHandler).Tick()
		if ret != true {
			common.LogInfo.Println("tick returning false from:", graphicsHandlers[i].(GraphicsHandler))
			return
		}
	}

	return
}

// RenderAll determines which objects are within view of the camera, and sends a list of those to each GraphicsHandler.
func (gm *GraphicsManager) RenderAll(camera *math.Frustum, sm component.SceneManager) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	compsToSend := common.MakeVector()
	comps := gm.compList.Array()
	graphicsHandlers := gm.graphicsHandlers.Array()
	for i := range comps {
		if comps[i] == nil {
			continue
		}
		loc := sm.GetObjectLocation(comps[i].(component.GOiD))
		if camera.ContainsPoint(loc) {
			compsToSend.Insert(comps[i].(component.GOiD))
		}
	}
	for i := range graphicsHandlers {
		graphicsHandlers[i].(GraphicsHandler).Render(gm.compList, sm, camera)
		// only send id's who will be visible
	}
}

// HandleInputs calls the HandleInputs function of each GraphicsHandler in the manager's list.
func (gm *GraphicsManager) HandleInputs(eye, target, up math.Vec3) (math.Vec3, math.Vec3, math.Vec3) {
	graphicsHandlers := gm.graphicsHandlers.Array()
	for i := range graphicsHandlers {
		/* inputs := */ graphicsHandlers[i].(GraphicsHandler).HandleInputs()
		// handle inputs
	}

	return graphicsHandlers[0].(*GlGraphicsManager).HandleInputs0(eye, target, up)
	//hack to keep the camera movement
}

// DrawString calls the DrawString function of each GraphicsHandler in the manager's list.
func (gm *GraphicsManager) DrawString(x, y float32, text string) {
	graphicsHandlers := gm.graphicsHandlers.Array()
	for i := range graphicsHandlers {
		graphicsHandlers[i].(GraphicsHandler).DrawString(x, y, text)
	}
}

// DrawString calls the DrawString function of each GraphicsHandler in the manager's list.
func (gm *GraphicsManager) GetSize() (int, int) {
	graphicsHandlers := gm.graphicsHandlers.Array()
	for i := range graphicsHandlers {
		return graphicsHandlers[i].(GraphicsHandler).GetSize()
	}

	return 0, 0
}

// RegisterGraphicsHandler addeds a GraphicsHandler to the manager's list to be included on all subsequent render and query function calls.
func (gm *GraphicsManager) RegisterGraphicsHandler(handler GraphicsHandler) {
	gm.graphicsHandlers.Insert(handler)
}
