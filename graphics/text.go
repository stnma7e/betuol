package graphics

import (
	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/math"
)

// TextGraphicsHandler implements the GraphicsHandler interface, but instead of rendering to a graphics context, TextGraphicsHandler instead outputs textual descriptions of each model.
type TextGraphicsHandler struct {
	compList   []string
	lastIdList common.Vector
}

// MakeTextGraphicsHandler returns a pointer to a TextGraphicsHandler.
func MakeTextGraphicsHandler() *TextGraphicsHandler {
	tgh := &TextGraphicsHandler{
		make([]string, 0),
		*common.MakeVector(),
	}
	return tgh
}

// Tick returns the status of the text window.
func (tgh *TextGraphicsHandler) Tick() bool {
	return true
}

// Render implements the Renderer interface and outputs text based on the GOiD's in the list passed as an argument.
// The function will only output a new text description if the model has newly come into the scene.
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

// LoadModel implements the GraphicsHandler interface and adds data used to render the components later.
func (tgh *TextGraphicsHandler) LoadModel(id component.GOiD, gc graphicsComponent) error {
	tgh.resizeArrays(id)
	tgh.compList[id] = gc.TextDescription

	return nil
}

// DeleteModel implements the GraphicsHandler interface and deletes the data used for rendering.
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

// HandleInputs implements the GraphicsHandler interface and returns the inputs recieved since the last query.
func (tgh *TextGraphicsHandler) HandleInputs() Inputs {
	return Inputs{}
}

// DrawString implements the GraphicsHandler interface and outputs a the string passed in as an arguement.
// The x, y coordinates are ignored.
func (tgh *TextGraphicsHandler) DrawString(x, y float32, text string) {

}

// GetSize implements the GraphicsHandler interface, but returns 0, 0 always because the text window has no size.
func (tgh *TextGraphicsHandler) GetSize() (int, int) {
	return 0, 0
}
