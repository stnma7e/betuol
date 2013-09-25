package physics

import (
	"smig/component"
	"smig/math"
	"smig/common"
)

type PhysicsManager struct {
	sm *component.SceneManager
	linearForces map[component.GOiD][]math.Vec3

	returnlink chan int
}

func MakePhysicsManager(sm *component.SceneManager) *PhysicsManager {
	pm := PhysicsManager {
		sm,
		make(map[component.GOiD][]math.Vec3),
		make(chan int),
	}
	return &pm
}

func (pm *PhysicsManager) JsonCreate(index component.GOiD, compData []byte) error {
	err := pm.CreateComponent(index)
	if err != nil {
		common.Log.Error(err)
	}
	return nil
}
func (pm *PhysicsManager) CreateComponent(index component.GOiD) error {
	pm.linearForces[index] = make([]math.Vec3, 1)
	return nil
}

func (pm *PhysicsManager) DeleteComponent(index component.GOiD) {
	pm.linearForces[index] = nil
}

func (pm *PhysicsManager) Tick(delta float64) {
	for k,v := range pm.linearForces {
		if v == nil {
			continue
		}
		var force math.Vec3
		for j := range v {
			newForce := v[j]
			force = math.Add3v3v(force, newForce)
		}
		force = math.Mult3vf(force, float32(delta))
		// fmt.Println("force ",force)
		transMat := pm.sm.GetTransformMatrix(k)
		transMat[3]  += force[0]
		transMat[7]  += force[1]
		transMat[11] += force[2]
		pm.sm.Transform(k, transMat)
	}
}

func (pm *PhysicsManager) AddForce(index component.GOiD, newForce math.Vec3) {
	length   := len(pm.linearForces[index])
	capacity := cap(pm.linearForces[index])
	if length >= capacity - 2 {
		newlist := make([]math.Vec3,capacity + 2)
		for i := 0; i < length; i++ {
			newlist[i] = pm.linearForces[index][i]
		}
		pm.linearForces[index] = newlist
	}
	pm.linearForces[index][length] = newForce
}
func (pm *PhysicsManager) RemoveForce(index component.GOiD, force math.Vec3) {
	for i := range pm.linearForces[index] {
		if math.Equal3v3v(pm.linearForces[index][i], force) {
			pm.linearForces[index][i] = math.Vec3{}
			return
		}
	}
}