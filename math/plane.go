package math

import "math"

// Represents the sides of a plane according to mathematical styles.
const (
	A = iota
	B = iota
	C = iota
	D = iota
)

// Plane is a wrapper for a 4 dimensional vector that is used to represent a plane in 3 dimensional space.
type Plane Vec4

// MakePlane returns a plane based on three inputs of 3D vectors.
func MakePlane(p1, p2, p3 Vec3) Plane {
	planeVec1 := Sub3v3v(p1, p2)
	planeVec2 := Sub3v3v(p1, p3)
	var normal [4]float32
	tmp := Cross3v3v(planeVec1, planeVec2)
	normal[0], normal[1], normal[2] = tmp[0], tmp[1], tmp[2]
	normal[3] = -(normal[0] + normal[1] + normal[2])
	plane := Plane{normal[0], normal[1], normal[2], normal[3]}

	return plane
}

// Normalize sets all dimensions of the plane to a value between 0 and 1.
func (pl *Plane) Normalize() {
	mag := float32((math.Sqrt(float64(pl[A]*pl[A] + pl[B]*pl[B] + pl[C]*pl[C]))))
	if mag != 0 {
		for i := range pl {
			pl[i] /= mag
		}
	}
}

// IsInside returns true if a point is on the side of the plane with a normal facing it.
func (pl *Plane) IsInside(vec Vec3) bool {
	return pl.Distance(vec) > 0.0
}

// IsOnPlane returns true if a point lies within a plane.
func (pl *Plane) IsOnPlane(vec Vec3) bool {
	var normDot float32
	for i := range vec {
		normDot += vec[i] * pl[i] // Dot product
	}
	if normDot == pl[D] {
		return true
	}
	return false
}

// Distance returns the shortest distance from the plane to a point in 3D space.
func (pl *Plane) Distance(vec Vec3) float32 {
	normal := Vec3{pl[A], pl[B], pl[C]}
	d := pl[D] / float32(math.Sqrt(float64((pl[A]*pl[A] + pl[B]*pl[B] + pl[C]*pl[C]))))
	return Dot3v3v(normal, vec) + d
}
