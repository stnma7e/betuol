package graphics

import (
	"encoding/json"
	"fmt"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/event"
	"github.com/stnma7e/betuol/graphics"
	"github.com/stnma7e/betuol/math"
	"github.com/stnma7e/betuol/res"
)

// GraphicsManager is a component manager used to visualize the game onscreen.
// It uses multiple GraphicsHandlers to render the world in a variety of ways.
type GraphicsManager struct {
	em                *event.EventManager
	rm                *res.ResourceManager
	sm                component.SceneManager
	justForcedARender bool

	graphicsHandlersLink []chan *common.Vector
	modellink            []chan graphics.ModelTransfer
	deletelink           []chan component.GOiD
	resizelink           []chan bool
	errorlink            chan error
	cam                  *math.Frustum

	compList *common.Vector
}

// MakeGraphicsManager returns a pointer to a GraphicsManager.
func MakeGraphicsManager(em *event.EventManager, rm *res.ResourceManager, sm component.SceneManager) *GraphicsManager {
	gm := &GraphicsManager{
		em,
		rm,
		sm,
		false,
		make([]chan *common.Vector, 1),
		make([]chan graphics.ModelTransfer, 1),
		make([]chan component.GOiD, 1),
		make([]chan bool, 1),
		make(chan error),
		math.MakeFrustum(0.1, 100, 90, 1/1),
		common.MakeVector(),
	}
	target, eye, up := math.Vec3{0, 0, 0}, math.Vec3{0, 6, -12}, math.Vec3{0, 1, 0}
	gm.cam.LookAt(target, eye, up)

	for i := range gm.graphicsHandlersLink {
		gm.graphicsHandlersLink[i] = make(chan *common.Vector)
	}
	for i := range gm.modellink {
		gm.modellink[i] = make(chan graphics.ModelTransfer)
	}
	for i := range gm.deletelink {
		gm.deletelink[i] = make(chan component.GOiD)
	}
	for i := range gm.resizelink {
		gm.resizelink[i] = make(chan bool)
	}

	go gm.RunGraphicsHandlerFunc(gm.graphicsHandlersLink[0], gm.modellink[0], gm.deletelink[0], gm.resizelink[0], gm.TextHandlerFunc)

	return gm
}

func (gm *GraphicsManager) handleClosedGraphicsHandler(indexOfClosedHandler int) {
	r := recover()
	if r != nil {
		common.LogErr.Printf("a graphics handler might have closed. deleting the handler now. recovered: %s", r)
		gm.graphicsHandlersLink[indexOfClosedHandler] = nil
		gm.resizelink[indexOfClosedHandler] = nil
		gm.deletelink[indexOfClosedHandler] = nil
	}
}

// Tick calls the Tick function of each GraphicsHandler in the manager's list.
// If any Tick functions return false, then GraphicsManager.Tick returns false.
func (gm *GraphicsManager) Tick(delta float64, sm component.SceneManager) {
	gm.sm = sm

	if gm.justForcedARender {
		gm.justForcedARender = false
		return
	}

	compsToSend, errs := gm.RenderAll(sm)
	if errs != nil {
		errArray := errs.Array()
		if errArray != nil && len(errArray) > 0 {
			for i := range errArray {
				common.LogErr.Print(errArray[i].(error))
			}
		}
	}
	gm.Render(compsToSend)
}

// ForceRender sends a resize message to all of the handlers, signaling a redraw.
func (gm *GraphicsManager) ForceRender(compsToSend *common.Vector) {
	handlerIndex := 0
	defer gm.handleClosedGraphicsHandler(handlerIndex)

	gm.Render(compsToSend)
	for handlerIndex = range gm.resizelink {
		if gm.resizelink[handlerIndex] == nil {
			continue
		}
		gm.resizelink[handlerIndex] <- true
	}

	gm.justForcedARender = true

}

// Render sends a new list of components to be rendered to the graphics handlers.
func (gm *GraphicsManager) Render(compsToSend *common.Vector) {
	handlerIndex := 0
	defer gm.handleClosedGraphicsHandler(handlerIndex)

	//common.LogInfo.Println(compsToSend)
	for handlerIndex = range gm.graphicsHandlersLink {
		if gm.graphicsHandlersLink[handlerIndex] == nil {
			continue
		}
		gm.graphicsHandlersLink[handlerIndex] <- compsToSend
	}
}

// RenderAllFromPerspective returns a list of all the game objects with graphics components that can be seen from the perspective of a single game object, id.
func (gm *GraphicsManager) RenderAllFromPerspective(id component.GOiD, sm component.SceneManager) (*common.Vector, *common.Vector) {
	errs := common.MakeVector()
	compsToSend := common.MakeVector()
	comps := gm.compList.Array()

	perspLoc, err := sm.GetObjectLocation(id)
	if err != nil {
		errs.Insert(fmt.Errorf("requesting location from scene manager failed in perspective render, error %s", err.Error()))
		return nil, errs
	}
	compsNearPerspective := sm.GetObjectsInLocationRadius(perspLoc, 5.0).Array()

	for i := range comps {
		if comps[i] == nil {
			continue
		}

		if comps[i].(component.GOiD) == id || comps[i].(component.GOiD) == 0 {
			continue
		}

		for j := range compsNearPerspective {
			if comps[i].(component.GOiD) == compsNearPerspective[j].(component.GOiD) {
				compsToSend.Insert(comps[i].(component.GOiD))
			}
		}
	}

	return compsToSend, errs
}

// RenderAll returns a list of all of the game objects with graphics components.
func (gm *GraphicsManager) RenderAll(sm component.SceneManager) (*common.Vector, *common.Vector) {
	errs := common.MakeVector()
	compsToSend := common.MakeVector()
	comps := gm.compList.Array()

	for i := range comps {
		if comps[i] == nil {
			continue
		}
		compsToSend.Insert(comps[i].(component.GOiD))
	}

	return compsToSend, errs
}

// JsonCreate extracts creation data from a byte array of json text to pass to CreateComponent.
func (gm *GraphicsManager) JsonCreate(id component.GOiD, compData []byte) error {
	obj := graphics.GraphicsComponent{}
	err := json.Unmarshal(compData, &obj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal graphics component, error: %s", err.Error())
	}

	gm.CreateComponent(id, obj)

	return nil
}

// Uses extracted data from higher level component creation functions and initializes a graphics component based on the id passed through.
// The function calls the LoadModel function of each GraphicsHandler in the manager's list.
func (gm *GraphicsManager) CreateComponent(id component.GOiD, gc graphics.GraphicsComponent) error {
	for i := range gm.modellink {
		gm.modellink[i] <- graphics.ModelTransfer{id, gc}
		err := <-gm.errorlink
		if err != nil {
			return fmt.Errorf("failed to create model with GraphicsHandler #%d, error: %s", i, err.Error())
		}
	}

	gm.compList.Insert(id)

	return fmt.Errorf("failed to create model with GraphicsHandler #%d, error: %s", id, "heop")
	return nil
}

// DeleteComponent implements the component.ComponentManager interface and deletes graphics component data from the manager.
// The function calls the DeleteModel function of each GraphicsHandler in the manager's list.
func (gm *GraphicsManager) DeleteComponent(id component.GOiD) {
	comps := gm.compList.Array()
	for i := range comps {
		if comps[i] == id {
			gm.compList.Erase(i)
		}
	}

	for i := range gm.modellink {
		if gm.deletelink[i] == nil {
			continue
		}
		gm.deletelink[i] <- id
	}
}

// RegisterGraphicsHandler addeds a GraphicsHandlerFunc to the manager's list to be included on all subsequent render and query function calls.
func (gm *GraphicsManager) RegisterGraphicsHandler(handler GraphicsHandlerFunc) {
	// resize arrays
	// make(chan) for all the channels
	// launch go routine with gm.RunGraphicsHandlerFunc()
}
