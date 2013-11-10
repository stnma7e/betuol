package math

import (
	"math"
)

type Vec4 [4]float32

func Mult4m4v(mat Mat4x4, vec Vec4) Vec4 {
	var row1, row2, row3, row4 float32

	for i := range vec {
		row1 += vec[i] * mat[i]
	}
	row1 += 1 * mat[3]
	for i := range vec {
		row2 += vec[i] * mat[i+4]
	}
	row2 += 1 * mat[7]
	for i := range vec {
		row3 += vec[i] * mat[i+8]
	}
	row3 += 1 * mat[11]
	for i := range vec {
		row4 += vec[i] * mat[i+12]
	}
	row4 += 1 * mat[15]

	return Vec4{row1, row2, row3, row4}
}

func Mult4v4v(vec1, vec2 Vec4) (ret Vec4) {
	for i := range vec1 {
		ret[i] = vec1[i] * vec2[i]
	}
	return
}

func Mag4v(vec Vec4) float32 {
	return float32(math.Sqrt(float64(MagSqrd4v(vec))))
}

func MagSqrd4v(vec Vec4) float32 {
	return vec[0]*vec[0] + vec[1]*vec[1] + vec[2]*vec[2] + vec[3]*vec[3]
}
