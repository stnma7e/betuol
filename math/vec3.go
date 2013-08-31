package math

import (
	"math"
	"fmt"
)

type Vec3 [3]float32

func (vec Vec3) Magnitude() float32 {
	return float32((math.Sqrt(float64(vec[A]*vec[A] + vec[B]*vec[B] + vec[C]*vec[C]))))
}
func (vec Vec3) MagSqrd() float32 {
	return vec[A]*vec[A] + vec[B]*vec[B] + vec[C]*vec[C]
}
func (vec *Vec3) Normalize() *Vec3 {
	mag := float32((math.Sqrt(float64(vec[A]*vec[A] + vec[B]*vec[B] + vec[C]*vec[C]))))
	if mag != 0 {
		for i := range vec {
			vec[i] /= mag
		}
	}
	return vec
}
func (vec Vec3) Distance(vec2 *Vec3) float32 {
	return float32(math.Sqrt(float64(vec.DistanceSqrd(vec2))))
}
func (vec Vec3) DistanceSqrd(vec2 *Vec3) float32 {
	split := vec.Subtract(vec2)
	return split.MagSqrd()
}
func (vec Vec3) Subtract(vec2 *Vec3) Vec3 {
	return Vec3{ vec[0]-vec2[0], vec[1]-vec2[1], vec[2]-vec2[2] }
}
func (vec Vec3) Add(vec2 *Vec3) Vec3 {
	return Vec3{ vec[0]+vec2[0], vec[1]+vec2[1], vec[2]+vec2[2]}
}

func Mult(vec Vec3, mat *Mat4x4) Vec3 {
	var row1,row2,row3 float32

	for i := range vec {
		row1 += vec[i]*mat[i]
	}
	row1 += 1 * mat[3]
	for i := range vec {
		row2 += vec[i]*mat[i+4]
	}
	row2 += 1 * mat[7]
	for i := range vec {
		row3 += vec[i]*mat[i+8]
	}
	row3 += 1 * mat[11]

	return Vec3{row1,row2,row3}
}
func (vec Vec3) MultScalar(scalar float32) Vec3 {
	return Vec3{vec[0] * scalar, vec[1] * scalar, vec[2] * scalar}
}
func (vec Vec3) Dot(vec2 *Vec3) (dot float32) {
	for i := range vec {
		dot += vec[i] * vec2[i]
	}
	return
}
func (vec Vec3) Cross(vec2 *Vec3) (ret Vec3) {
	ret[0] = vec[1]*vec2[2] - vec[2]*vec2[1]
	ret[1] = vec[2]*vec2[0] - vec[0]*vec2[2]
	ret[2] = vec[0]*vec2[1] - vec[1]*vec2[0]
	return
}

func (vec Vec3) Equals(vec2 *Vec3) bool {
	for i := range vec {
		if vec[i] != vec2[i] {
			return false
		}
	}
	return true
}

func (vec Vec3) ToJson() []byte {
	return []byte(fmt.Sprintf("[%v,%v,%v]", vec[0], vec[1], vec[2]))
}