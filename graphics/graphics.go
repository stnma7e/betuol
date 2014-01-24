// Package graphics handles the visualization of a game world in various techniques.
package graphics

import (
	"github.com/stnma7e/betuol/component"
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
