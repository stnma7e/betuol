// Package physics implements physics simulation and collision detection.
package physics

import (
	"encoding/json"
	"fmt"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/math"
)

// PhysicsManager implements a basic physics manager that handles collision detection and resolution.
// The structure also satisifies the component.SceneManager interface.
type PhysicsManager struct {
	sm    component.SceneManager
	radii []float32
}

// MakePhysicsManager returns a pointer to a PhysicsManager.
func MakePhysicsManager(sm component.SceneManager) *PhysicsManager {
	pm := PhysicsManager{
		sm,
		make([]float32, 0),
	}
	return &pm
}

// Tick updates the physics forces on each component, and checks for collisions between components.
func (pm *PhysicsManager) Tick(delta float64) {
	matList := pm.sm.GetMatrixList()
	for i := range pm.radii {
		if i == 0 {
			continue
		}
		if pm.radii[i] == 0 {
			continue
		}
		if matList[i].IsEmpty() {
			continue
		}
		loc1, err := pm.sm.GetObjectLocation(component.GOiD(i))
		if err != nil {
			common.LogErr.Println(err)
		}
		sp1 := math.Sphere{loc1, pm.radii[i]}
		for j := range pm.radii {
			if i == j || j == 0 {
				continue
			}
			if pm.radii[j] == 0 {
				continue
			}
			if matList[j].IsEmpty() {
				continue
			}
			loc2, err := pm.sm.GetObjectLocation(component.GOiD(j))
			if err != nil {
				common.LogErr.Println(err)
			}
			sp2 := math.Sphere{loc2, pm.radii[j]}
			if sp1.Intersects(sp2) {
				//common.LogWarn.Printf("collision between %d and %d\n", i, j)
				penetration := math.Sub3v3v(sp1.Center, sp2.Center)
				pSqrd := math.MagSqrd3v(penetration)
				if pSqrd-(sp1.Radius+sp2.Radius)*(sp1.Radius+sp2.Radius) > 0 {
					common.LogErr.Println("math fucked up")
				}
				split := math.Normalize3v(penetration)
				smallestDistanceToRemoveIntersection := math.Mult3vf(split, math.Mag3v(penetration))
				trans := matList[i]
				trans[3] += smallestDistanceToRemoveIntersection[0]
				trans[7] += smallestDistanceToRemoveIntersection[1]
				trans[11] += smallestDistanceToRemoveIntersection[2]
				pm.sm.SetTransform(component.GOiD(i), trans)
			}
		}
	}
}

// JsonCreate extracts creation data from a byte array of json text to pass to CreateComponent.
func (pm *PhysicsManager) JsonCreate(id component.GOiD, data []byte) error {
	var obj struct {
		BoundingRadius float32
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("failed to unmarshal physics component, error: %s", err.Error())
	}

	return pm.CreateComponent(id, obj.BoundingRadius)
}

// Uses extracted data from higher level component creation functions and initializes a character component based on the id passed through.
func (pm *PhysicsManager) CreateComponent(id component.GOiD, radius float32) error {
	pm.resizeArrays(id)

	pm.radii[id] = radius
	return nil
}

// resizeArray is a helper function to resize the array of components to accomodate a new component.
// If the GOiD of the new component is larger than the size of the array, then resizeArrays will grow the array and copy data over in order to fit the new component.
func (pm *PhysicsManager) resizeArrays(index component.GOiD) {
	const RESIZESTEP = 1
	if cap(pm.radii)-1 < int(index) {
		tmp := pm.radii
		pm.radii = make([]float32, index+RESIZESTEP)
		for i := range tmp {
			pm.radii[i] = tmp[i]
		}
	}
}

// DeleteComponent implements the component.ComponentManager interface and deletes character component data from the manager.
func (pm *PhysicsManager) DeleteComponent(id component.GOiD) {
	if len(pm.radii) <= int(id) {
		return
	}
	pm.radii[id] = 0
}
