package physics

import (
	"fmt"

	"smig/component"
	"smig/component/transform"
	"smig/math"
	"smig/common"
)

type PhysicsManager struct {
	tm *transform.SceneManager
	linearForces map[component.GOiD][]math.Vec3

	returnlink chan int
}

func MakePhysicsManager(tm *transform.SceneManager) *PhysicsManager {
	pm := PhysicsManager {
		tm,
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
		index := k
			var force math.Vec3
			for j := range v {
				newForce := v[j]
				force = force.Add(&newForce)
			}
			force = force.MultScalar(float32(delta))
			// fmt.Println("force ",force)
			transMat,err := pm.tm.GetTransform(index)
			if err != nil {
				fmt.Println(err)
			}
			var newMat math.Mat4x4
			for i := range transMat {
				newMat[i] = transMat[i]
			}
			newMat[3]  += force[0]
			newMat[7]  += force[1]
			newMat[11] += force[2]
			pm.tm.Transform(k, &newMat)
	}
}

func (pm *PhysicsManager) AddForce(index component.GOiD, newForce *math.Vec3) {
	length   := len(pm.linearForces[index])
	capacity := cap(pm.linearForces[index])
	if length >= capacity - 2 {
		newlist := make([]math.Vec3,capacity + 2)
		for i := 0; i < length; i++ {
			newlist[i] = pm.linearForces[index][i]
		}
		pm.linearForces[index] = newlist
	}
	pm.linearForces[index][length] = *newForce
}
func (pm *PhysicsManager) RemoveForce(index component.GOiD, force *math.Vec3) {
	for i := range pm.linearForces[index] {
		if pm.linearForces[index][i].Equals(force) {
			pm.linearForces[index][i] = math.Vec3{}
			return
		}
	}
}