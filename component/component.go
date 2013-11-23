package component

import "betuol/math"

type GOiD uint32

const NULLINDEX = 0

type ComponentManager interface {
	DeleteComponent(GOiD)
}

type SceneManager interface {
	GetTransform4m(index GOiD) (math.Mat4x4, error)
	GetObjectLocation(index GOiD) math.Vec3
}

type GameObject map[string][]byte

type Map []MapLocation
type MapLocation struct {
	Location math.Vec3
	Entities []MapEntity
}
type MapEntity struct {
	Breed    string
	CompList GameObject
	Quantity int
}
