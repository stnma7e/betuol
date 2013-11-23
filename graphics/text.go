package graphics

import (
	"betuol/common"
	"betuol/component"
	"betuol/math"
)

type TextGraphicsHandler struct {
	compList   []string
	lastIdList common.Vector
}

func MakeTextGraphicsHandler() *TextGraphicsHandler {
	tgh := &TextGraphicsHandler{
		make([]string, 0),
		*common.MakeVector(),
	}
	return tgh
}

func (tgh *TextGraphicsHandler) Tick() bool {
	return true
}

func (tgh *TextGraphicsHandler) Render(ids *common.Vector, sm component.SceneManager, cam *math.Frustum) {
	diff := tgh.lastIdList.Difference(ids)

	if ids.Length == tgh.lastIdList.Length && diff.Length < 1 {
		//common.LogInfo.Println(ids)
		return
	}

	//common.LogInfo.Println(ids)
	comps := diff.Array()
	for i := range comps {
		id := comps[i].(int)
		if id == 0 {
			continue
		}
		common.LogInfo.Printf("%d, \"%s\"\n", id, tgh.compList[id])
	}
	tgh.lastIdList = *ids
}

func (tgh *TextGraphicsHandler) LoadModel(id component.GOiD, gc GraphicsComponent) error {
	tgh.resizeArrays(id)
	tgh.compList[id] = gc.TextDescription

	return nil
}

func (tgh *TextGraphicsHandler) DeleteModel(id component.GOiD) {
	tgh.compList[id] = ""
}

func (tgh *TextGraphicsHandler) resizeArrays(id component.GOiD) {
	const RESIZESTEP = 1
	if cap(tgh.compList)-1 < int(id) {
		newCompList := make([]string, id+RESIZESTEP)
		for i := range tgh.compList {
			newCompList[i] = tgh.compList[i]
		}
		tgh.compList = newCompList
	}
}

func (tgh *TextGraphicsHandler) HandleInputs() Inputs {
	return Inputs{}
}

func (tgh *TextGraphicsHandler) DrawString(x, y float32, text string) {

}

func (tgh *TextGraphicsHandler) GetSize() (int, int) {
	return 0, 0
}
