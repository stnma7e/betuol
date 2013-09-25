package graphics

import (
	"unsafe"
	"github.com/go-gl/gl"

	"smig/math"
)

const (
	// buffers / attributes
	VERT	= iota
	NORM 	= iota
	COLOR 	= iota
	UV 		= iota
	INDEX 	= iota

	NUMBUFFERS = iota
	NUMATTR = 2
)
const (
	//uniforms
	WORLD  		= iota
	PROJ 		= iota
	NORMALMV	= iota

	DIRLIGHT 	= iota
	LINTENSE 	= iota
	LPOSITION   = iota
	AMBINTENSE  = iota

	NUMUNIFORMS = iota
)

type Model struct {
	buffers		[NUMBUFFERS]gl.Buffer
	uniforms 	[NUMUNIFORMS]gl.UniformLocation
	attribArray [NUMATTR]gl.AttribLocation
	program 	gl.Program
	vao 		gl.VertexArray
}

func MakeModel(program gl.Program,
			  attribArray [NUMATTR]gl.AttribLocation,
			  vertexList  []math.Vec3,
			  indexList   []uint32,
			  normalList  []math.Vec3,
			  texUvList   []math.Vec2,
			  textureList []gl.Texture) Model {

	model := Model{}
	model.program = program

	vao := gl.GenVertexArray()
	vao.Bind()
	model.vao = vao

	vertBuffer := gl.GenBuffer()
	vertBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexList) * int(unsafe.Sizeof(math.Vec3{})), &vertexList[0], gl.STATIC_DRAW)
	vertBuffer.Unbind(gl.ARRAY_BUFFER)
	model.buffers[VERT] = vertBuffer

	normBuffer := gl.GenBuffer()
	normBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(normalList) * int(unsafe.Sizeof(math.Vec3{})), &normalList[0], gl.STATIC_DRAW)
	normBuffer.Unbind(gl.ARRAY_BUFFER)
	model.buffers[NORM] = normBuffer

	indexBuffer := gl.GenBuffer()
	indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indexList) * int(unsafe.Sizeof(uint32(1))), &indexList[0], gl.STATIC_DRAW)
	indexBuffer.Unbind(gl.ELEMENT_ARRAY_BUFFER)
	model.buffers[INDEX] = indexBuffer

	uvBuffer := gl.GenBuffer()
	uvBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(texUvList) * int(unsafe.Sizeof(math.Vec2{})), &texUvList[0], gl.STATIC_DRAW)
	uvBuffer.Unbind(gl.ARRAY_BUFFER)
	model.buffers[UV] = uvBuffer

	model.attribArray = [NUMATTR]gl.AttribLocation{}
	for i := range attribArray {
		model.attribArray[i] = attribArray[i]
	}

	model.uniforms[WORLD] 	   = program.GetUniformLocation("modelToCamera")
	model.uniforms[PROJ] 	   = program.GetUniformLocation("projection")
	model.uniforms[NORMALMV]   = program.GetUniformLocation("normalModelToCamera")
	model.uniforms[DIRLIGHT]   = program.GetUniformLocation("dirToLight")
	model.uniforms[LINTENSE]   = program.GetUniformLocation("lightIntensity")
	model.uniforms[LPOSITION]  = program.GetUniformLocation("lightPosition")
	model.uniforms[AMBINTENSE] = program.GetUniformLocation("ambientIntensity")

	return model
}

func (mh *Model) Delete() {
	for i := range mh.buffers {
		mh.buffers[i].Delete()
	}
	mh.program.Delete()
	mh.vao.Delete()
}