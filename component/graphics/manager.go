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
	em *event.EventManager
	rm *res.ResourceManager
	sm component.SceneManager

	graphicsHandlersLink []chan *common.Vector
	modellink            []chan graphics.ModelTransfer
	deletelink           []chan component.GOiD
	resizelink           []chan bool
	errorlink            chan error
	cam                  *math.Frustum

	compList *common.Vector
}

// MakeGraphicsManager returns a pointer to a GraphicsManager
func MakeGraphicsManager(em *event.EventManager, rm *res.ResourceManager, sm component.SceneManager) *GraphicsManager {
	gm := &GraphicsManager{
		em,
		rm,
		sm,
		make([]chan *common.Vector, 2),
		make([]chan graphics.ModelTransfer, 2),
		make([]chan component.GOiD, 2),
		make([]chan bool, 2),
		make(chan error),
		math.MakeFrustum(0.1, 100, 90, 1/1),
		common.MakeVector(),
	}

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

	go gm.RunGraphicsHandlerFunc(gm.graphicsHandlersLink[0], gm.modellink[0], gm.deletelink[0], gm.resizelink[0], gm.GlHandlerFunc)
	go gm.RunGraphicsHandlerFunc(gm.graphicsHandlersLink[1], gm.modellink[1], gm.deletelink[1], gm.resizelink[1], gm.TextHandlerFunc)

	return gm
}

// Tick calls the Tick function of each GraphicsHandler in the manager's list.
// If any Tick functions return false, then GraphicsManager.Tick returns false.
func (gm *GraphicsManager) Tick(delta float64, sm component.SceneManager) {
	var i int
	defer func() {
		r := recover()
		if r != nil {
			common.LogInfo.Printf("recovered: %s", r)
		}
		gm.graphicsHandlersLink[i] = nil
		gm.resizelink[i] = nil
		gm.deletelink[i] = nil
	}()
	gm.sm = sm

	compsToSend := common.MakeVector()
	comps := gm.compList.Array()
	for i := range comps {
		if comps[i] == nil {
			continue
		}
		loc := gm.sm.GetObjectLocation(comps[i].(component.GOiD))
		if gm.cam.ContainsPoint(loc) {
			compsToSend.Insert(comps[i].(component.GOiD))
		}
	}

	for i = range gm.graphicsHandlersLink {
		if gm.graphicsHandlersLink[i] == nil {
			continue
		}
		gm.graphicsHandlersLink[i] <- compsToSend
	}

	for i = range gm.resizelink {
		if gm.resizelink[i] == nil {
			continue
		}
		gm.resizelink[i] <- true
	}

	return
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
