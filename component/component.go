package component

import "smig/math"

type GOiD uint32
const NULLINDEX = 0

type ComponentManager interface {
	DeleteComponent(index GOiD)
}

type GameObject map[string][]byte

type Map []MapLocation
type MapLocation struct {
	Location math.Vec3
	Entities []MapEntity
}
type MapEntity struct {
	Breed string
	CompList GameObject
	Quantity int
}