package graphics

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl"

	//"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/math"
	"github.com/stnma7e/betuol/res"
)

// Renderer represents an interface for GlGraphicsManager to utilize to render a Model based on the program being used.
type Renderer interface {
	Render(mod Model, transMat, camMat, projectMat math.Mat4x4)
}

type renderer struct {
	program     gl.Program
	attribArray [NUMATTR]gl.AttribLocation
}

type lightingRenderer struct {
	lightIntensity   math.Vec4
	ambientIntensity math.Vec4
}

// FragmentPointLightingRenderer uses a program which lights a scene based on a Fragment shader based Point Lighting scheme.
type FragmentPointLightingRenderer struct {
	renderer
	lightingRenderer

	uniforms      [NUMUNIFORMS]gl.UniformLocation
	lightPosition math.Vec3
}

// MakeFragmentPointLightingRenderer returns a pointer to a FragmentPointLightingRenderer.
func MakeFragmentPointLightingRenderer(rm *res.ResourceManager, glg *GlGraphicsManager) (Renderer, error) {
	fplr := FragmentPointLightingRenderer{}
	vert, err := rm.GetFileContents("graphics/shader/FragmentLighting.vert")
	if err != nil {
		return &FragmentPointLightingRenderer{}, fmt.Errorf("failed to open vertex shader, error: %s", err.Error())
	}
	frag, err := rm.GetFileContents("graphics/shader/FragmentLighting.frag")
	if err != nil {
		return &FragmentPointLightingRenderer{}, fmt.Errorf("failed to open fragment shader, error: %s", err.Error())
	}
	vertStr, fragStr := string(vert), string(frag)
	shaders := [2]gl.Shader{}
	shaders[0], err = LoadShader(gl.VERTEX_SHADER, vertStr)
	if err != nil {
		return &FragmentPointLightingRenderer{}, fmt.Errorf("failed to compile vertex shader, error: %s", err.Error())
	}
	shaders[1], err = LoadShader(gl.FRAGMENT_SHADER, fragStr)
	if err != nil {
		return &FragmentPointLightingRenderer{}, fmt.Errorf("failed to compile fragment shader, error: %s", err.Error())
	}
	fplr.program = gl.CreateProgram()
	fplr.attribArray, err = LinkProgram(fplr.program, shaders[:])
	if err != nil {
		return &FragmentPointLightingRenderer{}, fmt.Errorf("failed to link program, error: %s", err.Error())
	}

	fplr.lightIntensity = math.Vec4{0.8, 0.8, 0.8, 0.8}
	fplr.ambientIntensity = math.Vec4{0.2, 0.2, 0.2, 0.2}

	fplr.uniforms[WORLD] = fplr.program.GetUniformLocation("modelToCamera")
	fplr.uniforms[PROJ] = fplr.program.GetUniformLocation("projection")
	fplr.uniforms[LMODELPOS] = fplr.program.GetUniformLocation("modelSpaceLightPosition")
	fplr.uniforms[LINTENSE] = fplr.program.GetUniformLocation("lightIntensity")
	fplr.uniforms[AMBINTENSE] = fplr.program.GetUniformLocation("ambientIntensity")

	return &fplr, nil
}

// Render implements the Renderer interface and can be used by GlGraphicsManager to render a Momdel without any extra specification of data.
func (fplr *FragmentPointLightingRenderer) Render(mod Model, transMat, camMat, projectMat math.Mat4x4) {
	mod.vao.Bind()
	defer NOVAO.Bind()
	fplr.program.Use()
	defer NOPROGRAM.Use()

	//common.LogInfo.Println(transMat, camMat, projectMat)

	lightPosition := math.Vec3{40, 20, 20}
	//lightPosCameraSpace := math.Mult4m3v(camMat, lightPosition)
	lightPosModelSpace := math.Mult4m3v(transMat.Inverse(), lightPosition)

	transArray := [16]float32(math.Mult4m4m(camMat, transMat))
	projection := [16]float32(projectMat)
	var normArray [9]float32
	for i := range normArray {
		normArray[i] = transArray[i]
	}

	fplr.uniforms[WORLD].UniformMatrix4f(false, &transArray)
	fplr.uniforms[PROJ].UniformMatrix4f(false, &projection)
	fplr.uniforms[LMODELPOS].Uniform4fv(1, lightPosModelSpace[:])
	fplr.uniforms[LINTENSE].Uniform4fv(1, fplr.lightIntensity[:])
	fplr.uniforms[AMBINTENSE].Uniform4fv(1, fplr.ambientIntensity[:])
	fplr.uniforms[LPOSITION].Uniform4fv(1, lightPosition[:])
	//fplr.uniforms[LCAMPOS].Uniform4fv(1, lightPosCameraSpace[:])

	mod.buffers[VERT].Bind(gl.ARRAY_BUFFER)
	fplr.attribArray[VERT].EnableArray()
	fplr.attribArray[VERT].AttribPointer(3, gl.FLOAT, false, 0, nil)
	mod.buffers[VERT].Unbind(gl.ARRAY_BUFFER)
	defer fplr.attribArray[VERT].DisableArray()

	mod.buffers[NORM].Bind(gl.ARRAY_BUFFER)
	fplr.attribArray[NORM].EnableArray()
	fplr.attribArray[NORM].AttribPointer(3, gl.FLOAT, false, 0, nil)
	mod.buffers[NORM].Unbind(gl.ARRAY_BUFFER)
	defer fplr.attribArray[NORM].DisableArray()

	fplr.renderModel(mod)
}

func (rd *renderer) renderModel(mod Model) {
	mod.buffers[INDEX].Bind(gl.ELEMENT_ARRAY_BUFFER)
	defer mod.buffers[INDEX].Unbind(gl.ELEMENT_ARRAY_BUFFER)

	// gl.DrawArrays(gl.TRIANGLES,0,3)
	sizeptr := gl.GetBufferParameteriv(gl.ELEMENT_ARRAY_BUFFER, gl.BUFFER_SIZE)
	size := sizeptr / int32(unsafe.Sizeof(float32(1)))
	//fmt.Println(size)
	gl.DrawElements(gl.TRIANGLES, int(size), gl.UNSIGNED_INT, nil)
}
