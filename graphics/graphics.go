package graphics

import (
	"smig/math"
)

type GraphicsManager interface {
	Closing() bool
	ShouldClose()
	MakeContextCurrent()
	Render(mh *Mesh, transMat *math.Mat4x4)
}