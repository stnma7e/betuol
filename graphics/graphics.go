package graphics

import (
	"smig/common"
	"smig/component"
	"smig/math"
	"smig/res"
)

type WindowManager interface {
	Closing() bool
	ShouldClose()
	close()
	Tick()
	MakeContextCurrent()
	Render(mh Model, transMat, camMat, projectMat math.Mat4x4)
	LoadModel(comp *GraphicsComponent, rm *res.ResourceManager) *Model
	HandleInputs(eye, target, up math.Vec3) (math.Vec3, math.Vec3, math.Vec3)
	DrawString(x, y float32, text string)
	GetSize() (int, int)
	SwapBuffers()
}

type Inputs struct {
}

type GraphicsHandler interface {
	Render(ids *common.Vector, sm component.SceneManager, cam *math.Frustum)
	LoadModel(id component.GOiD, gc GraphicsComponent) error
	DeleteModel(id component.GOiD)
	Tick() bool
	HandleInputs() Inputs
	DrawString(x, y float32, text string)
	GetSize() (int, int)
}
