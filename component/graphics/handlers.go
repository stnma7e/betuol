package graphics

import (
	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/graphics"
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
