package math

import (
	"math"
)

type Mat4x4 [16]float32

func MakePerspectiveMatrix(near, far, fov, aspect float32) (mat Mat4x4) {
	fovRadians := float64(fov * math.Pi / 180)
	scale := float32(math.Tan(fovRadians / 2))
	mat[0] = 1 / scale * aspect
	mat[5] = 1 / scale
	mat[10] = -(far + near) / (far - near)
	mat[14] = -1
	mat[11] = -(2 * far * near) / (far - near)

	return
}

// m00, m01, m02, m03
// m04, m05, m06, m07
// m08, m09, m10, m11
// m12, m13, m14, m15

//func (mat Mat4x4) Inverse() Mat4x4 {
//mat2 := mat
//mat2[1] = mat[4]
//mat2[2] = mat[8]
//mat2[6] = mat[9]
//mat2[4] = mat[1]
//mat2[8] = mat[2]
//mat2[9] = mat[6]

//mat2[3] = -mat[3]
//mat2[7] = -mat[7]
//mat2[11] = -mat[11]

//return mat2
//}

func Equals4m4m(mat, mat2 Mat4x4) bool {
	for i := range mat {
		if mat[i] != mat2[i] {
			return false
		}
	}
	return true
}

func (mat Mat4x4) Inverse() Mat4x4 {
	var segs [4]Mat2x2 // four corner 2x2 matrices
	segs[0] = Mat2x2{mat[0], mat[1], mat[4], mat[5]}
	segs[1] = Mat2x2{mat[2], mat[3], mat[6], mat[7]}
	segs[2] = Mat2x2{mat[8], mat[9], mat[12], mat[13]}
	segs[3] = Mat2x2{mat[10], mat[11], mat[14], mat[15]}
	for i := range segs {
		segs[i].Invert()
	}

	mat[0], mat[1], mat[4], mat[5] = segs[0].Split()
	mat[2], mat[3], mat[6], mat[7] = segs[1].Split()
	mat[8], mat[9], mat[12], mat[13] = segs[2].Split()
	mat[10], mat[11], mat[14], mat[15] = segs[3].Split()

	for i := range mat {
		if math.IsNaN(float64(mat[i])) {
			mat[i] = 0
		}
	}

	return mat
}

func (m *Mat4x4) MakeIdentity() {
	for i := range m {
		m[i] = 0
	}
	m[0], m[5], m[10], m[15] = 1, 1, 1, 1
}

func (m *Mat4x4) IsEmpty() bool {
	for i := range m {
		if m[i] != 0 {
			return false
		}
	}
	return true
}

// Return M1*M2
func Mult4m4m(mat1, mat2 Mat4x4) Mat4x4 {
	for i := 0; i < 16; i += 4 {
		val := [4]float32{}
		for k := 0; k < 4; k++ {
			for j := 0; j < 4; j++ {
				val[k] += (mat2[j*4+k] * mat1[i+j])
			}
		}
		for k := 0; k < 4; k++ {
			mat1[i+k] = val[k]
		}
	}
	return mat1
}

// func Mult4m4mj(mat1, mat2 Mat4x4) (ret Mat4x4) {
//      var rows [4]Vec4
//      var cols [4]Vec4
//      for i := 0; i < 4; i++ {
//              rows[i] = Vec4{mat1[0+(4*i)],mat1[1+(4*i)],mat1[2+(4*i)],mat1[3+(4*i)]}
//              cols[i] = Vec4{mat2[0+i],mat2[1+i],mat2[2+i],mat2[3+i]}
//      }
//      for i := 0; i < 4; i++ {
//              for j := 0; j < 4; j++ {
//                      ret[j+(4*i)] = rows[i][0]*cols[j][0] + rows[i][1]*cols[j][1] + rows[i][2]*cols[j][2] + rows[i][3]*cols[j][3]
//              }
//      }

//      return
// }
