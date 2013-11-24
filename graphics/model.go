package graphics

import (
	"github.com/go-gl/gl"
	"unsafe"

	"github.com/stnma7e/betuol/math"
)

const (
	// buffers / attributes
	VERT  = iota
	NORM  = iota
	UV    = iota
	COLOR = iota
	INDEX = iota

	NUMBUFFERS = iota
	NUMATTR    = 3
)
const (
	//uniforms
	WORLD    = iota
	PROJ     = iota
	NORMALMV = iota

	DIRLIGHT   = iota
	LINTENSE   = iota
	LPOSITION  = iota
	AMBINTENSE = iota
	LMODELPOS  = iota

	SAMPLER = iota

	NUMUNIFORMS = iota
)

// Model contains data used by GlGraphicsManager to render and display a type of mesh.
type Model struct {
	buffers [NUMBUFFERS]gl.Buffer
	vao     gl.VertexArray
}

// MakeModel takes a slew of arrays to initialize a model with the necessary information to render.
func MakeModel(vertexList []math.Vec3, indexList []uint32, normalList []math.Vec3, texUvList []math.Vec2) Model {
	model := Model{}

	vao := gl.GenVertexArray()
	vao.Bind()
	model.vao = vao

	vertBuffer := gl.GenBuffer()
	vertBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexList)*int(unsafe.Sizeof(math.Vec3{})), &vertexList[0], gl.STATIC_DRAW)
	vertBuffer.Unbind(gl.ARRAY_BUFFER)
	model.buffers[VERT] = vertBuffer

	normBuffer := gl.GenBuffer()
	normBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(normalList)*int(unsafe.Sizeof(math.Vec3{})), &normalList[0], gl.STATIC_DRAW)
	normBuffer.Unbind(gl.ARRAY_BUFFER)
	model.buffers[NORM] = normBuffer

	uvBuffer := gl.GenBuffer()
	uvBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(texUvList)*int(unsafe.Sizeof(math.Vec2{})), &texUvList[0], gl.STATIC_DRAW)
	uvBuffer.Unbind(gl.ARRAY_BUFFER)
	model.buffers[UV] = uvBuffer

	indexBuffer := gl.GenBuffer()
	indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indexList)*int(unsafe.Sizeof(uint32(1))), &indexList[0], gl.STATIC_DRAW)

	indexBuffer.Unbind(gl.ELEMENT_ARRAY_BUFFER)
	model.buffers[INDEX] = indexBuffer

	return model
}

// Delete iterates through the OpenGL buffers and deletes each one and then it delete the VAO of the model.
func (mh *Model) Delete() {
	for i := range mh.buffers {
		mh.buffers[i].Delete()
	}
	mh.vao.Delete()
}
