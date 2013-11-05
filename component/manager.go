package component

import (
	"fmt"

	"smig/common"
	"smig/math"
)

const (
	ROOTNODE   = 0
	RESIZESTEP = 1
)

type TransformManager struct {
	matList		[]math.Mat4x4
	boundingSpheres []math.Sphere
	movedQueue	common.IntQueue
	returnlink	chan int
}

func MakeTransformManager() *TransformManager {
	tm := TransformManager{}
	tm.matList = make([]math.Mat4x4,5)
	tm.matList[ROOTNODE].MakeIdentity()
	tm.returnlink = make(chan int)
	return &tm
}

func (tm *TransformManager) Tick(delta float64) {
	// const BLOCKSIZE = 250
	// var numberCompleted int
	// for ; !tm.movedQueue.IsEmpty(); numberCompleted++ {

	// 	var index [BLOCKSIZE]int
	// 	var err error
	// 	for i := range index {
	// 		index[i],err = tm.movedQueue.Dequeue()
	// 		if err != nil {
	// 			break
	// 		}
	// 	}
		
	// 	go func(compid [BLOCKSIZE]int) {
	// 		for i := range compid {		
	// 			if compid[i] == 0 {
	// 				break
	// 			}
	// 			lmat 		:= tm.matList[LMAT][compid[i]]
	// 			parentIndex := tm.childParentMap[GOiD(compid[i])]
	// 			parent 		:= tm.matList[WMAT][int(parentIndex)]

	// 			wmat := math.Mult4m4m(parent, lmat)
	// 			tm.matList[WMAT][compid[i]] = wmat
	// 			// fmt.Println(compid[i], "wmat", wmat.ToString())

	// 			sp := &tm.boundingSpheres[compid[i]]
	// 			sp.Center = math.Vec3 {
	// 				wmat[3], wmat[7], wmat[11],
	// 			}
	// 		}

	// 		tm.returnlink <- 1
	// 	}(index)
	// }
	// for i := 0; i < numberCompleted; i++ {
	// 	<-tm.returnlink
	// }

	for i := range tm.boundingSpheres {
		if i == 0 {
			continue
		}
		for j := range tm.boundingSpheres {
			if j == 0 || j == i {
				continue
			}
			if tm.boundingSpheres[i].Intersects(tm.boundingSpheres[j]) {
                                //bound1, bound2 := tm.boundingSpheres[i], tm.boundingSpheres[j]
                                //split := math.Sub3v3v(bound1.Center, bound2.Center)
                                //collisionOnBound2 := math.Sub3v3v(split, math.Mult3vf(math.Normalize3v(split), bound1.Radius))
                                //fmt.Printf("collision: %v,%v %v\n", i, j, split)
                                 //tm.matList[j] = math.Mult4m4m(math.Mult4m4m(tm.matList[j], math.Mat4x4 {
                                         //1,0,0,0,
                                         //0,1,0,0,
                                         //0,0,1,0,
                                         //0,0,0,1,
                                 //}), math.Mat4x4 {
                                         //1,0,0,collisionOnBound2[0],
                                         //0,1,0,collisionOnBound2[1],
                                         //0,0,1,collisionOnBound2[2],
                                         //0,0,0,1,
                                //})
                        }
		}
	}
}

func (tm *TransformManager) CreateComponent(index GOiD) error {
	tm.resizeArray(index)	

	if !(tm.matList[index].IsEmpty()) {
		return fmt.Errorf("attempt to reuse GOiD %d", index)
	}
	tm.matList[index].MakeIdentity()
        fmt.Println(tm.matList[index])

	return nil
}
func (tm *TransformManager) resizeArray(index GOiD) {
	if cap(tm.matList) - 1 < int(index) {
		newCompList := make([]math.Mat4x4, index + RESIZESTEP)
		for i := range tm.matList {
			newCompList[i] = tm.matList[i]
		}
		tm.matList = newCompList
	}

	if cap(tm.boundingSpheres) - 1 < int(index) {
		tmp := tm.boundingSpheres
		tm.boundingSpheres = make([]math.Sphere, index + RESIZESTEP)
		for i := range tmp {
			tm.boundingSpheres[i] = tmp[i]
		}
	}
}

func (tm *TransformManager) DeleteComponent(index GOiD) {
	tm.matList[index] = math.Mat4x4{}
}

func (tm *TransformManager) SetTransform(index GOiD, newLocalMat math.Mat4x4) {
	tm.matList[index] = newLocalMat
}
func (tm *TransformManager) SetLocation(index GOiD, newLocation math.Vec3) {
    tm.matList[index][3] = newLocation[0]
    tm.matList[index][7] = newLocation[1]
    tm.matList[index][11] = newLocation[2]
}
func (tm *TransformManager) GetTransform4m(index GOiD) math.Mat4x4 {
    return tm.GetTransformMatrix(index)
}
func (tm *TransformManager) GetTransformMatrix(index GOiD) math.Mat4x4 {
	if int(index) >= len(tm.matList) {
            common.LogErr.Printf("invalid GOiD %v: not in list", index)
	}
	if tm.matList[index].IsEmpty() {
            common.LogErr.Printf("invalid GOiD: %v: empty matrix", index)
	}
	return tm.matList[index]
}
func (tm *TransformManager) GetObjectLocation(index GOiD) math.Vec3 {
	locMat := tm.matList[index]
	return math.Mult4m3v(locMat, math.Vec3{})
}
func (tm *TransformManager) GetBoundingSphere(index GOiD) math.Sphere {
	return tm.boundingSpheres[index]
}

func (tm *TransformManager) SetBoundingSphere(index GOiD, bound math.Sphere) {
	mat := math.Mat4x4{}
	mat.MakeIdentity()
	mat[3]  = bound.Center[0]
	mat[7]  = bound.Center[1]
	mat[11] = bound.Center[2]
	tm.matList[index] = mat
}

func (tm *TransformManager) GetObjectsInLocationRadius(loc math.Vec3, lookRange float32) *common.IntQueue {
	sp := math.Sphere {
		loc, lookRange,
	}
	stk := common.IntQueue{}

	for i := range tm.boundingSpheres {
		bsp := tm.boundingSpheres[i]
		if sp.Intersects(bsp) {
			stk.Queue(i)
		}
	}

	return &stk
}
