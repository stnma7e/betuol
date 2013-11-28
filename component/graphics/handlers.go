package graphics

import (
	"fmt"
	"runtime"
	"time"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/graphics"
	"github.com/stnma7e/betuol/math"

	"github.com/go-gl/gl"
	"github.com/go-gl/glfw3"
)

// GraphicsHandlerFunc is used by a GraphicsManager to launch rendering graphics handlers. Each handler gets its own thread and will respond to channel sends from the GraphicsManager.
type GraphicsHandlerFunc func(compslink chan *common.Vector, modellink chan graphics.ModelTransfer, deletelink chan component.GOiD, resizelink chan bool)

func (gm *GraphicsManager) RunGraphicsHandlerFunc(compslink chan *common.Vector, modellink chan graphics.ModelTransfer, deletelink chan component.GOiD, resizelink chan bool, ghf GraphicsHandlerFunc) {
	ghf(compslink, modellink, deletelink, resizelink)
	close(compslink)
	close(modellink)
	close(deletelink)
	close(resizelink)
}

// GlHandlerFunc starts an OpenGL rendering context and satisfies the GraphicsHandlerFunc function type.
func (gm *GraphicsManager) GlHandlerFunc(compslink chan *common.Vector, modellink chan graphics.ModelTransfer, deletelink chan component.GOiD, resizelink chan bool) {
	defer func() {
		r := recover()
		if r != nil {
			common.LogInfo.Printf("recovered: %s", r)
		}
	}()
	runtime.LockOSThread()
	glg, err := graphics.MakeGlGraphicsManager(640, 480, "betuol", gm.rm)
	if err != nil {
		common.LogErr.Println(err)
		return
	}

	target, eye, up := math.Vec3{0, 0, 0}, math.Vec3{0, 6, -12}, math.Vec3{0, 1, 0}
	gm.cam.LookAt(target, eye, up)

	oldTime := time.Now()
	comps := common.MakeVector()
	for i := true; i; {
		i = glg.Tick()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		eye, target, up = glg.HandleInputs0(eye, target, up)
		gm.cam.LookAt(target, eye, up)

		secs := time.Since(oldTime).Seconds()
		oldTime = time.Now()
		fpsStr := fmt.Sprintf("%f", 100/secs)
		spfStr := fmt.Sprintf("%f", secs/100)

		select {
		case comps = <-compslink:
		case model := <-modellink:
			gm.errorlink <- glg.LoadModel(model.Id, model.Gc)
		case id := <-deletelink:
			glg.DeleteModel(id)
		case <-resizelink:
			x, y := glg.GetSize()
			gm.cam = math.MakeFrustum(0.1, 100, 90, float32(y)/float32(x))
			gm.cam.LookAt(target, eye, up)
		}

		glg.Render(comps, gm.sm, gm.cam)

		glg.DrawString(10, 10, "fps: "+fpsStr)
		glg.DrawString(25, 25, "spf: "+spfStr)
	}
	glfw3.Terminate()
}

// TextHandlerFunc satisfies the GraphicsHandlerFunc function type, and runs a TextGraphicsHandler.
func (gm *GraphicsManager) TextHandlerFunc(compslink chan *common.Vector, modellink chan graphics.ModelTransfer, deletelink chan component.GOiD, resizelink chan bool) {
	tr := graphics.MakeTextGraphicsHandler()
	comps := common.MakeVector()
	for i := true; i; {
		select {
		case comps = <-compslink:
		case model := <-modellink:
			gm.errorlink <- tr.LoadModel(model.Id, model.Gc)
		case id := <-deletelink:
			tr.DeleteModel(id)
		case <-resizelink:
		}

		tr.Render(comps, gm.sm, gm.cam)
	}
}
