package scenemanager

import (
	"fmt"

	"smig/common"
	"smig/math"
        "smig/event"
        "smig/component"
)

const (
	ROOTNODE   = 0
	RESIZESTEP = 1
)

type TransformManager struct {
        em *event.EventManager

	matList		[]math.Mat4x4
	movedQueue	common.IntQueue
	returnlink	chan int
}

func MakeTransformManager(em *event.EventManager) *TransformManager {
	tm := TransformManager{}
        tm.em = em
	tm.matList = make([]math.Mat4x4,5)
	tm.matList[ROOTNODE].MakeIdentity()
	tm.returnlink = make(chan int)
	return &tm
}

func (tm *TransformManager) Tick(delta float64) {
}

func (tm *TransformManager) CreateComponent(index component.GOiD) error {
	tm.resizeArray(index)

	if !(tm.matList[index].IsEmpty()) {
		return fmt.Errorf("attempt to reuse component.GOiD %d", index)
	}
	tm.matList[index].MakeIdentity()

	return nil
}
func (tm *TransformManager) resizeArray(index component.GOiD) {
	if cap(tm.matList) - 1 < int(index) {
		newCompList := make([]math.Mat4x4, index + RESIZESTEP)
		for i := range tm.matList {
			newCompList[i] = tm.matList[i]
		}
		tm.matList = newCompList
	}
}

func (tm *TransformManager) DeleteComponent(index component.GOiD) {
	tm.matList[index] = math.Mat4x4{}
}

func (tm *TransformManager) SetTransform(index component.GOiD, newLocalMat math.Mat4x4) {
	tm.matList[index] = newLocalMat
        newLocation := math.Vec3{ newLocalMat[3], newLocalMat[7], newLocalMat[11] }
        tm.em.Send(event.CharacterMoveEvent{ index, newLocation })
}
func (tm *TransformManager) SetLocation(index component.GOiD, newLocation math.Vec3) {
    tm.matList[index][3] = newLocation[0]
    tm.matList[index][7] = newLocation[1]
    tm.matList[index][11] = newLocation[2]
    tm.em.Send(event.CharacterMoveEvent{ index, newLocation })
}
func (tm *TransformManager) GetTransform4m(index component.GOiD) math.Mat4x4 {
    return tm.GetTransformMatrix(index)
}
func (tm *TransformManager) GetTransformMatrix(index component.GOiD) math.Mat4x4 {
	if int(index) >= len(tm.matList) {
            common.LogErr.Printf("invalid component.GOiD %v: not in list", index)
	}
	if tm.matList[index].IsEmpty() {
            common.LogErr.Printf("invalid component.GOiD: %v: empty matrix", index)
	}
	return tm.matList[index]
}
func (tm *TransformManager) GetObjectLocation(index component.GOiD) math.Vec3 {
	locMat := tm.matList[index]
	return math.Mult4m3v(locMat, math.Vec3{})
}

func (tm *TransformManager) GetObjectsInLocationRadius(loc math.Vec3, lookRange float32) *common.IntQueue {
	sp := math.Sphere {
		loc, lookRange,
	}
	stk := common.IntQueue{}

	for i := range tm.matList {
                loc2 := math.Mult4m3v(tm.matList[i], math.Vec3{0,0,0})
                sp2 := math.Sphere {
                    loc2, 1,
                }
		if sp.Intersects(sp2) {
			stk.Queue(i)
		}
	}

	return &stk
}
