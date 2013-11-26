// Package graphics handles the visualization of a game world in various techniques.
package graphics

import (
	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/math"
)

type Inputs struct {
}

// ModelTransfer is used to send creation information to the renderers over a channel.
type ModelTransfer struct {
	Id component.GOiD
	Gc GraphicsComponent
}

// GraphicsComponent represents data extracted from a graphics component file.
type GraphicsComponent struct {
	ModelName, Mesh, MeshType, Renderer, TextDescription string
}

// GraphicsManager is a component manager used to visualize the game onscreen.
// GraphicsHandler represents an interface that is used to render the game world regardless of output media.
type GraphicsHandler interface {
	Render(ids *common.Vector, sm component.SceneManager, cam *math.Frustum)
	LoadModel(id component.GOiD, gc GraphicsComponent) error
	DeleteModel(id component.GOiD)
	Tick() bool
	HandleInputs() Inputs
	DrawString(x, y float32, text string)
	GetSize() (int, int)
}
