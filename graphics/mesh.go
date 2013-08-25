package graphics

import (
	"unsafe"

	"github.com/go-gl/gl"

	"smig/math"
)

const (
	// buffers
	VERT	= 0
	INDEX 	= 1
	NORMAL  = 2
	COLOR 	= 3
	UV 		= 4

	//uniforms
	MVP 	= 0
)

type Mesh struct {
	buffers		[5]gl.Buffer
	uniforms 	[1]gl.UniformLocation
	attribArray []gl.AttribLocation
	program 	gl.Program
	vao 		gl.VertexArray
}

func MakeMesh(program gl.Program,
			  attribArray []gl.AttribLocation,
			  vertexList  []math.Vec3,
			  indexList   []uint32,
			  texUvList   []math.Vec2,
			  normalList  []math.Vec3,
			  textureList []gl.Texture) *Mesh {

	vao := gl.GenVertexArray()
	vao.Bind()

	vertBuffer := gl.GenBuffer()
	vertBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexList) * int(unsafe.Sizeof(math.Vec3{})), &vertexList[0], gl.STATIC_DRAW)
	defer vertBuffer.Unbind(gl.ARRAY_BUFFER)

	var a uint32
	indexBuffer := gl.GenBuffer()
	indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indexList) * int(unsafe.Sizeof(a)), &indexList[0], gl.STATIC_DRAW)
	indexBuffer.Unbind(gl.ELEMENT_ARRAY_BUFFER)

	mesh := Mesh{}
	mesh.program = program
	mesh.vao = vao
	mesh.buffers[VERT]  = vertBuffer
	mesh.buffers[INDEX] = indexBuffer
	mesh.attribArray = make([]gl.AttribLocation, len(attribArray))
	for i := range attribArray {
		mesh.attribArray[i] = attribArray[i]
	}

	mesh.uniforms[MVP] = program.GetUniformLocation("mvp")

	return &mesh
}

func (mh *Mesh) Delete() {
	for i := range mh.buffers {
		mh.buffers[i].Delete()
	}
	mh.program.Delete()
	mh.vao.Delete()
}