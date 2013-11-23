// Package component is imported for use of entity-component system.
// This package only includes types relevant to all of the component managers (GOiD).
package component

import "betuol/math"

// GOiD stands for GameObject ID. It is the type used by all component managers to represent a specific GameObject.
type GOiD uint32

// NULLINDEX is a reserved GOiD that can be used internally by the component managers.
const NULLINDEX = 0

// ComponentManager is an interface to treat all component managers the same for the GameObject Factory.
type ComponentManager interface {
	DeleteComponent(GOiD)
}

// SceneManager is an interface used for location management. This can be a physics manager or simply a location manager.
type SceneManager interface {
	GetTransform4m(index GOiD) (math.Mat4x4, error)
	GetObjectLocation(index GOiD) math.Vec3
}

// Type used by the GameObject Factory for handling creation data for various component types.
type GameObject map[string][]byte

// Type used by the GameObject Factory for loading a list of GameObjects according to parameters specified in a map file.
type Map []MapLocation

// Type used to specify a single location or grid space in the map.
type MapLocation struct {
	Location math.Vec3
	Entities []MapEntity
}

// Type used to specify the type information and quantity of an entity on a map.
type MapEntity struct {
	Breed    string
	CompList GameObject
	Quantity int
}
