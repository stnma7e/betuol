package math

import "math"
import "fmt"

type Mat4x4 [16]float32

func MakePerspectiveMatrix(near, far, fov, aspect float32) (mat Mat4x4) {
	scale := float32(float64(aspect) * math.Tan(float64((fov * 0.5) * (math.Pi/180))))
	mat[0], mat[5] = scale, scale
	mat[10] = -(far / (far - near))
	mat[11] = -1
	mat[14] = -(far * near / (far - near))
	mat[15] = 0

	return mat
}

/*
	Returns matrix M = M*B
*/
func (m *Mat4x4) Mult(b *Mat4x4) *Mat4x4 { 
	for i := 0; i < 16; i += 4 {
		var val [4]float32
		for k := 0; k < 4; k++ {
			for j := 0; j < 4; j++ {
				val[k] += (b[j*4+k] * m[i+j])
			}	
		}
		for k := 0; k < 4; k++ {
			m[i+k] = val[k]
		}
	}
	return m
}
func (mat Mat4x4) Inverse() *Mat4x4 {
	var segs [4]Mat2x2 // four corner 2x2 matrices
	segs[0] = Mat2x2{mat[0],mat[1],mat[4],mat[5]}
	segs[1] = Mat2x2{mat[2],mat[3],mat[6],mat[7]}
	segs[2] = Mat2x2{mat[8],mat[9],mat[12],mat[13]}
	segs[3] = Mat2x2{mat[10],mat[11],mat[14],mat[15]}
	for i := range segs {
		segs[i].Invert()
	}

	mat[0],mat[1],mat[4],mat[5]  	= segs[0].Split()
	mat[2],mat[3],mat[6],mat[7]   	= segs[1].Split()
	mat[8],mat[9],mat[12],mat[13] 	= segs[2].Split()
	mat[10],mat[11],mat[14],mat[15]	= segs[3].Split()

	for i := range mat {
		if math.IsNaN(float64(mat[i])) {
			mat[i] = 0
		}
	}

	return &mat
}

func (m *Mat4x4) Equals(mat2 *Mat4x4) bool {
	for i := range m {
		if m[i] != mat2[i] {
			return false
		}
	}
	return true
}
func (m *Mat4x4) MakeIdentity() {
	for i := range m {
		m[i] = 0
	}
	m[0], m[5], m[10], m[15] = 1,1,1,1
}
func (m *Mat4x4) IsEmpty() bool {
	for i := range m {
		if m[i] != 0 {
			return false
		}
	}
	return true
}

func (m *Mat4x4) ToString() string {
	s := fmt.Sprint("\n",m[0],m[1],m[2],m[3],"\n",
					m[4],m[5],m[6],m[7],"\n",
					m[8],m[9],m[10],m[11],"\n",
					m[12],m[13],m[14],m[15],"\n")
	return s
}