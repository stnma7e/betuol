package scenemanager

import (
	"fmt"

	"smig/common"
	"smig/component"
	"smig/event"
	"smig/math"
)

const (
	ROOTNODE   = 0
	RESIZESTEP = 1
)

type TransformManager struct {
	em *event.EventManager

	matList []math.Mat4x4
	moving  []moveOverTime
}

type moveOverTime struct {
	movementAxis math.Vec3
	timeToMove   float64
}

func MakeTransformManager(em *event.EventManager) *TransformManager {
	tm := TransformManager{
		em,
		make([]math.Mat4x4, 5),
		make([]moveOverTime, 5),
	}
	tm.matList[ROOTNODE].MakeIdentity()
	return &tm
}

func (tm *TransformManager) Tick(delta float64) {
	for i := range tm.moving {
		if tm.moving[i].timeToMove < 0 {
			continue
		}
		split := math.Mult3vf(tm.moving[i].movementAxis, float32(delta))
		tm.matList[i][3] += split[0]
		tm.matList[i][7] += split[1]
		tm.matList[i][11] += split[2]
		tm.moving[i].timeToMove -= delta
	}
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
	if cap(tm.matList)-1 < int(index) {
		newCompList := make([]math.Mat4x4, index+RESIZESTEP)
		for i := range tm.matList {
			newCompList[i] = tm.matList[i]
		}
		tm.matList = newCompList
	}
	if cap(tm.moving)-1 < int(index) {
		newMoveList := make([]moveOverTime, index+RESIZESTEP)
		for i := range tm.moving {
			newMoveList[i] = tm.moving[i]
		}
		tm.moving = newMoveList
	}
}

func (tm *TransformManager) DeleteComponent(index component.GOiD) {
	tm.matList[index] = math.Mat4x4{}
}

func (tm *TransformManager) SetTransform(index component.GOiD, newLocalMat math.Mat4x4) {
	tm.matList[index] = newLocalMat
	newLocation := math.Vec3{newLocalMat[3], newLocalMat[7], newLocalMat[11]}
	tm.em.Send(event.CharacterMoveEvent{index, newLocation})
}

func (tm *TransformManager) SetLocation(index component.GOiD, newLocation math.Vec3) {
	tm.matList[index][3] = newLocation[0]
	tm.matList[index][7] = newLocation[1]
	tm.matList[index][11] = newLocation[2]
	tm.em.Send(event.CharacterMoveEvent{index, newLocation})
}

func (tm *TransformManager) SetLocationOverTime(id component.GOiD, newLocation math.Vec3, timeToMove float64) {
	if len(tm.moving)-1 < int(id) {
		common.LogErr.Printf("invalid id, %v", id)
	}
	originalLocation := tm.GetObjectLocation(id)
	mot := moveOverTime{math.Mult3vf(math.Sub3v3v(newLocation, originalLocation), float32(1/timeToMove)), timeToMove}
	tm.moving[id] = mot
}

func (tm *TransformManager) GetTransform4m(index component.GOiD) (math.Mat4x4, error) {
	return tm.GetTransformMatrix(index)
}

func (tm *TransformManager) GetTransformMatrix(index component.GOiD) (math.Mat4x4, error) {
	if int(index) >= len(tm.matList) {
		return math.Mat4x4{}, fmt.Errorf("invalid component.GOiD, %v: not in list")
	}
	if tm.matList[index].IsEmpty() {
		return math.Mat4x4{}, fmt.Errorf("invalid component.GOiD, %v: empty matrix", index)
	}
	return tm.matList[index], nil
}

func (tm *TransformManager) GetObjectLocation(index component.GOiD) math.Vec3 {
	locMat := tm.matList[index]
	return math.Mult4m3v(locMat, math.Vec3{})
}

func (tm *TransformManager) GetObjectsInLocationRadius(loc math.Vec3, lookRange float32) *common.IntQueue {
	sp := math.Sphere{
		loc, lookRange,
	}
	stk := common.IntQueue{}

	for i := range tm.matList {
		if tm.matList[i].IsEmpty() || i == 0 {
			continue
		}
		loc2 := math.Mult4m3v(tm.matList[i], math.Vec3{0, 0, 0})
		sp2 := math.Sphere{
			loc2, 1,
		}
		if sp.Intersects(sp2) {
			stk.Queue(i)
		}
	}

	return &stk
}


















































































































































