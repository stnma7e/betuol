// Package scenemanager implements a basic location manager that satisfies the component.SceneManager interface.
package scenemanager

import (
	"fmt"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/event"
	"github.com/stnma7e/betuol/math"
)

// TransformManager implements a basic location manager that satisfies the component.SceneManager interface
type TransformManager struct {
	em *event.EventManager

	matList []math.Mat4x4
	moving  []moveOverTime
}

type moveOverTime struct {
	movementAxis math.Vec3
	timeToMove   float64
}

// MakeTransformManager returns a pointer to a TransformManager
func MakeTransformManager(em *event.EventManager) *TransformManager {
	tm := TransformManager{
		em,
		make([]math.Mat4x4, 5),
		make([]moveOverTime, 5),
	}
	tm.matList[0].MakeIdentity()
	return &tm
}

// Tick is used to update the locations of components moving using the SetLocationOverTime function and execute a basic collision detection and resolution algorithm.
func (tm *TransformManager) Tick(delta float64) {
	for i := range tm.moving {
		if tm.moving[i].timeToMove < 0 {
			continue
		}
		split := math.Mult3vf(tm.moving[i].movementAxis, float32(delta))
		trans := tm.matList[i]
		trans[3] += split[0]
		trans[7] += split[1]
		trans[11] += split[2]
		tm.SetTransform(component.GOiD(i), trans)
		tm.moving[i].timeToMove -= delta
	}
}

// CreateComponent is used to initialize a matrix to store the location of a transform component.
func (tm *TransformManager) CreateComponent(index component.GOiD) error {
	tm.resizeArray(index)

	if !(tm.matList[index].IsEmpty()) {
		return fmt.Errorf("attempt to reuse component.GOiD %d before deleting it", index)
	}
	tm.matList[index].MakeIdentity()

	return nil
}

// resizeArray is a helper function to resize the array of components to accomodate a new component.
// If the GOiD of the new component is larger than the size of the array, then resizeArrays will grow the array and copy data over in order to fit the new component.
func (tm *TransformManager) resizeArray(index component.GOiD) {
	const RESIZESTEP = 1
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

// DeleteComponent implements the component.ComponentManager interface and deletes character component data from the manager.
func (tm *TransformManager) DeleteComponent(index component.GOiD) {
	tm.matList[index] = math.Mat4x4{}
}

// SetTransform is used to set the location of a transform component using a new matrix.
func (tm *TransformManager) SetTransform(index component.GOiD, newLocalMat math.Mat4x4) {
	tm.matList[index] = newLocalMat
	newLocation := math.Vec3{newLocalMat[3], newLocalMat[7], newLocalMat[11]}
	tm.em.Send(event.CharacterMoveEvent{index, newLocation})
}

// SetLocation is used to set the location of a transform component using a 3 dimensional vector.
func (tm *TransformManager) SetLocation(index component.GOiD, newLocation math.Vec3) {
	tm.matList[index][3] = newLocation[0]
	tm.matList[index][7] = newLocation[1]
	tm.matList[index][11] = newLocation[2]
	tm.em.Send(event.CharacterMoveEvent{index, newLocation})
}

// SetLocationOverTime is used to set the location of a transform component using a 3 dimensional vector.
// This function will update the location interpolated across the timespan specified.
func (tm *TransformManager) SetLocationOverTime(id component.GOiD, newLocation math.Vec3, timeToMove float64) error {
	if len(tm.moving)-1 < int(id) {
		common.LogErr.Printf("invalid id, %v", id)
	}
	originalLocation, err := tm.GetObjectLocation(id)
	if err != nil {
		return err
	}
	mot := moveOverTime{math.Mult3vf(math.Sub3v3v(newLocation, originalLocation), float32(1/timeToMove)), timeToMove}
	tm.moving[id] = mot
	return nil
}

// GetTransform4m implements the component.SceneManager interface and returns a matrix of the location of an object.
func (tm *TransformManager) GetTransform4m(index component.GOiD) (math.Mat4x4, error) {
	return tm.GetTransformMatrix(index)
}

// GetTransformMatrix returns a matrix of the location of an object.
func (tm *TransformManager) GetTransformMatrix(index component.GOiD) (math.Mat4x4, error) {
	if int(index) >= len(tm.matList) {
		return math.Mat4x4{}, fmt.Errorf("invalid component.GOiD, %v: not in list")
	}
	if tm.matList[index].IsEmpty() {
		return math.Mat4x4{}, fmt.Errorf("invalid component.GOiD, %v: empty matrix", index)
	}
	return tm.matList[index], nil
}

// GetObjectLocation returns the location of an object in a 3 dimensional vector.
func (tm *TransformManager) GetObjectLocation(index component.GOiD) (math.Vec3, error) {
	if int(index) >= len(tm.matList) {
		return math.Vec3{}, fmt.Errorf("invalid component.GOiD, %v: not in list")
	}
	if tm.matList[index].IsEmpty() {
		return math.Vec3{}, fmt.Errorf("invalid component.GOiD, %v: empty matrix", index)
	}
	locMat := tm.matList[index]
	return math.Mult4m3v(locMat, math.Vec3{}), nil
}

// GetObjectsInLocationRadius returns a list of GOiD's within a radius around a location.
func (tm *TransformManager) GetObjectsInLocationRadius(loc math.Vec3, lookRange float32) *common.Queue {
	sp := math.Sphere{
		loc, lookRange,
	}
	stk := common.Queue{}

	for i := range tm.matList {
		if tm.matList[i].IsEmpty() || i == 0 {
			continue
		}
		loc2 := math.Mult4m3v(tm.matList[i], math.Vec3{0, 0, 0})
		sp2 := math.Sphere{
			loc2, 1,
		}
		if sp.Intersects(sp2) {
			stk.Queue(component.GOiD(i))
		}
	}

	return &stk
}

// GetMatrixList returns the internal list of location matrices used by the TransformManager.
func (tm *TransformManager) GetMatrixList() []math.Mat4x4 {
	return tm.matList
}
