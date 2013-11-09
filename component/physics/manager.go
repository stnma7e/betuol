package physics

import (
        "fmt"
        "encoding/json"
        gomath "math"

	"smig/component"
        "smig/component/scenemanager"
	"smig/math"
)

const (
    MASS = iota
    RESTITUTION = iota
)

type PhysicsManager struct {
	tm *scenemanager.TransformManager
	//linearForces map[component.GOiD][]math.Vec3
        attrib map[component.GOiD][2]float32
        velocity map[component.GOiD]math.Vec3
        boundings map[component.GOiD]math.Sphere

	returnlink chan int
}

func MakePhysicsManager(tm *scenemanager.TransformManager) *PhysicsManager {
	pm := PhysicsManager {
		tm,
		//make(map[component.GOiD][]math.Vec3),
                make(map[component.GOiD][2]float32),
                make(map[component.GOiD]math.Vec3),
                make(map[component.GOiD]math.Sphere),
		make(chan int),
	}
	return &pm
}

func (pm *PhysicsManager) JsonCreate(index component.GOiD, compData []byte) error {
        var obj struct {
            Mass, Restitution float32
        }
        err := json.Unmarshal(compData, &obj)
        if err != nil {
            return err
        }

	return pm.CreateComponent(index, obj.Mass, obj.Restitution)
}
func (pm *PhysicsManager) CreateComponent(index component.GOiD, mass, restitution float32) error {
        if index == 0 {
            return fmt.Errorf("help the index is 0")
        }
	//pm.linearForces[index] = make([]math.Vec3, 1)
        pm.AddForce(index, math.Vec3{0,-2.95,0})
        pm.attrib[index] = [2]float32{ mass, restitution }
	return nil
}

func (pm *PhysicsManager) DeleteComponent(index component.GOiD) {
	//pm.linearForces[index] = nil
        pm.attrib[index] = [2]float32{}
        pm.velocity[index] = math.Vec3{}
}

func (pm *PhysicsManager) Tick(delta float64) {
        //for k,v := range pm.linearForces {
		//if v == nil {
			//continue
		//}
		//var force math.Vec3
		//for j := range v {
			//newForce := v[j]
			//force = math.Add3v3v(force, newForce)
		//}
		//force = math.Mult3vf(force, float32(delta))
                 //fmt.Println("force ",force)
		//transMat := pm.tm.GetTransformMatrix(k)
		//transMat[3]  += force[0]
		//transMat[7]  += force[1]
		//transMat[11] += force[2]
		//pm.tm.SetTransform(k, transMat)
	//}
        for k,v := range pm.boundings {
            for l, w := range pm.boundings {
                split := math.Sub3v3v(w.Center, v.Center)
                if math.Dot3v3v(split, pm.velocity[k]) <= 0 {
                    continue
                }
                vel := math.Mult3vf(pm.velocity[k], float32(delta))
                if math.MagSqrd3v(vel) < math.MagSqrd3v(split) {
                    continue
                }
                dist := math.Dot3v3v(math.Normalize3v(vel), split)
                f := math.MagSqrd3v(split) - (dist*dist)
                sumRadiiSqrd := (v.Radius + w.Radius)*(v.Radius + w.Radius)
                if f >= sumRadiiSqrd {
                    continue
                }
                t := sumRadiiSqrd - f
                if t < 0 {
                    continue
                }
                continueDist := dist - float32(gomath.Sqrt(float64(t)))
                mag := math.Mag3v(vel)
                if mag < continueDist {
                    continue
                }
                vel = math.Normalize3v(pm.velocity[k])
                pm.velocity[k] = math.Mult3vf(vel, continueDist)
                fmt.Println("collide", k, l)
            }
        }
        for k,v := range pm.velocity {
            v = math.Mult3vf(v, float32(delta))
            transMat := pm.tm.GetTransformMatrix(k)
            transMat[3]  += v[0]
            transMat[7]  += v[1]
            transMat[11] += v[2]
            pm.tm.SetTransform(k, transMat)
        }
}

func (pm *PhysicsManager) AddForce(index component.GOiD, newForce math.Vec3) {
	//length   := len(pm.linearForces[index])
	//capacity := cap(pm.linearForces[index])
	//if length >= capacity - 2 {
		//newlist := make([]math.Vec3,capacity + 2)
		//for i := 0; i < length; i++ {
			//newlist[i] = pm.linearForces[index][i]
		//}
		//pm.linearForces[index] = newlist
	//}
	//pm.linearForces[index][length] = newForce
        pm.velocity[index] = math.Add3v3v(pm.velocity[index], newForce)
}
func (pm *PhysicsManager) RemoveForce(index component.GOiD, force math.Vec3) {
	//for i := range pm.linearForces[index] {
		//if math.Equal3v3v(pm.linearForces[index][i], force) {
			//pm.linearForces[index][i] = math.Vec3{}
			//return
		//}
	//}
        pm.velocity[index] = math.Sub3v3v(pm.velocity[index], force)
}
