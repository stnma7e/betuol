package math

import (
	"fmt"
	"math"
)

type Vec3 [3]float32

func Mag3v(vec Vec3) float32 {
	return float32(math.Sqrt(float64(MagSqrd3v(vec))))
}

func MagSqrd3v(vec Vec3) float32 {
	return vec[0]*vec[0] + vec[1]*vec[1] + vec[2]*vec[2]
}

func Dist3v3v(vec1, vec2 Vec3) float32 {
	return float32(math.Sqrt(float64(DistSqrd3v3v(vec1, vec2))))
}

func DistSqrd3v3v(vec1, vec2 Vec3) float32 {
	split := Sub3v3v(vec1, vec2)
	return MagSqrd3v(split)
}

func (vec Vec3) ToJson() []byte {
	return []byte(fmt.Sprintf("[%v,%v,%v]", vec[0], vec[1], vec[2]))
}

func Normalize3v(vec Vec3) Vec3 {
	mag := Mag3v(vec)
	if mag != 0 {
		for i := range vec {
			vec[i] /= float32(mag)
		}
		return vec
	}
	return Vec3{}
}

func Sub3v3v(vec1, vec2 Vec3) Vec3 {
	return Vec3{vec1[0] - vec2[0], vec1[1] - vec2[1], vec1[2] - vec2[2]}
}

func Add3v3v(vec1, vec2 Vec3) Vec3 {
	return Vec3{vec1[0] + vec2[0], vec1[1] + vec2[1], vec1[2] + vec2[2]}
}

func Mult3vf(vec Vec3, scalar float32) Vec3 {
	return Vec3{
		vec[0] * scalar,
		vec[1] * scalar,
		vec[2] * scalar,
	}
}

func Dot3v3v(vec1, vec2 Vec3) float32 {
	return vec1[0]*vec2[0] + vec1[1]*vec2[1] + vec1[2]*vec2[2]
}

func Cross3v3v(vec1, vec2 Vec3) (ret Vec3) {
	ret[0] = vec1[1]*vec2[2] - vec1[2]*vec2[1]
	ret[1] = vec1[2]*vec2[0] - vec1[0]*vec2[2]
	ret[2] = vec1[0]*vec2[1] - vec1[1]*vec2[0]
	return
}

func Mult4m3v(mat Mat4x4, vec Vec3) Vec3 {
	var row1, row2, row3 float32

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

	return Vec3{row1, row2, row3}
}

func Equal3v3v(vec1, vec2 Vec3) bool {
	for i := range vec1 {
		if vec1[i] != vec2[i] {
			return false
		}
	}
	return true
}
