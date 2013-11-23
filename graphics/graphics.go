// Package graphics handles the visualization of a game world in various techniques.
package graphics

import (
	"betuol/common"
	"betuol/component"
	"betuol/math"
)

type Inputs struct {
}

// GraphicsHandler represents an interface that is used to render the game world regardless of output media.
type GraphicsHandler interface {
	Render(ids *common.Vector, sm component.SceneManager, cam *math.Frustum)
	LoadModel(id component.GOiD, gc graphicsComponent) error
	DeleteModel(id component.GOiD)
	Tick() bool
	HandleInputs() Inputs
	DrawString(x, y float32, text string)
	GetSize() (int, int)
}
