package graphics

import (
	"unsafe"
	"github.com/go-gl/gl"

	"smig/math"
)

const (
	// buffers / attributes
	VERT	= iota
	NORM	= iota
	UV	= iota
	COLOR	= iota
	INDEX	= iota

	NUMBUFFERS = iota
	NUMATTR = 3
)
const (
	//uniforms
	WORLD		= iota
	PROJ		= iota
	NORMALMV	= iota

	DIRLIGHT	= iota
	LINTENSE	= iota
	LPOSITION       = iota
	AMBINTENSE      = iota
	LMODELPOS	= iota

        SAMPLER = iota

	NUMUNIFORMS = iota
)

type Model struct {
	buffers		[NUMBUFFERS]gl.Buffer
	uniforms	[NUMUNIFORMS]gl.UniformLocation
	attribArray     [NUMATTR]gl.AttribLocation
	vao		gl.VertexArray
        tex             gl.Texture
}

func MakeModel(vertexList  []math.Vec3,
	        indexList   []uint32,
	        normalList  []math.Vec3,
                texUvList   []math.Vec2,
		textureList []gl.Texture) Model {

	model := Model{}

	vao := gl.GenVertexArray()
	vao.Bind()
	model.vao = vao

        //model.tex = GlLoadTexture("/home/sam/downloads/tower/tower_diffuse.png")
        model.tex = textureList[0]

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

	uvBuffer := gl.GenBuffer()
	uvBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(texUvList) * int(unsafe.Sizeof(math.Vec2{})), &texUvList[0], gl.STATIC_DRAW)
	uvBuffer.Unbind(gl.ARRAY_BUFFER)
	model.buffers[UV] = uvBuffer

	indexBuffer := gl.GenBuffer()
	indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indexList) * int(unsafe.Sizeof(uint32(1))), &indexList[0], gl.STATIC_DRAW)

	indexBuffer.Unbind(gl.ELEMENT_ARRAY_BUFFER)
	model.buffers[INDEX] = indexBuffer

	return model
}

func (mh *Model) Delete() {
	for i := range mh.buffers {
		mh.buffers[i].Delete()
	}
	mh.vao.Delete()
}
