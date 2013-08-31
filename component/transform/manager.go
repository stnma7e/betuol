package transform

import (
	"fmt"
	"encoding/json"

	"smig/component"
	"smig/common"
	"smig/math"
	"smig/graphics"
)

const (
	WMAT = 0
	LMAT = 1
	ROOTNODE = 0
	RESIZESTEP = 1
)

type SceneManager struct {
	compList 		[2][]math.Mat4x4 // compList[WMAT] == world transform matrices
								 	 // compList[LMAT] == local transform matrices
	parentChildMap 	map[component.GOiD]common.IntQueue
	childParentMap  map[component.GOiD]component.GOiD
	boundingSpheres []math.Sphere
	meshList 		[]graphics.Mesh
	movedQueue  	common.IntQueue
	returnlink  	chan int
}

func MakeSceneManager() *SceneManager {
	tm := SceneManager{}
	for i := range tm.compList {
		tm.compList[i] = make([]math.Mat4x4,5)
	}
	tm.compList[WMAT][ROOTNODE].MakeIdentity()
	tm.returnlink = make(chan int)
	tm.parentChildMap = make(map[component.GOiD]common.IntQueue)
	tm.childParentMap = make(map[component.GOiD]component.GOiD)
	return &tm
}

func (tm *SceneManager) Tick(delta float64) {
	const BLOCKSIZE = 250
	var numberCompleted int
	for ; !tm.movedQueue.IsEmpty(); numberCompleted++ {

		var index [BLOCKSIZE]int
		var err error
		for i := range index {
			index[i],err = tm.movedQueue.Dequeue()
			if err != nil {
				break
			}
		}
		
		go func(compid [BLOCKSIZE]int) {
			for i := range compid {		
				if compid[i] == 0 {
					break
				}
				lmat 		:= tm.compList[LMAT][compid[i]]
				parentIndex := tm.childParentMap[component.GOiD(compid[i])]
				parent 		:= tm.compList[WMAT][int(parentIndex)]

				wmat := *parent.Mult(&lmat)
				tm.compList[WMAT][compid[i]] = wmat
				// fmt.Println(compid[i], "wmat", wmat.ToString())

				sp := &tm.boundingSpheres[compid[i]]
				sp.Center = math.Vec3 {
					wmat[3], wmat[7], wmat[11],
				}
			}

			tm.returnlink <- 1
		}(index)
	}
	for i := 0; i < numberCompleted; i++ {
		<-tm.returnlink
	}
}

func (tm *SceneManager) JsonCreate(index component.GOiD, compData []byte) error {
	var comp struct {
		Location math.Vec3
		Radius float32
	}
	err := json.Unmarshal(compData, &comp)
	if err != nil {
		common.Log.Error(err)
	}

	sp := math.Sphere {
		comp.Location,
		comp.Radius,
	}

	err = tm.CreateComponent(index, ROOTNODE, sp)
	if err != nil {
		common.Log.Error(err)
	}

	return nil
}
func (tm *SceneManager) CreateComponent(index, parent component.GOiD, bound math.Sphere) error {
	tm.resizeArray(index)	

	if !(tm.compList[WMAT][index].IsEmpty()) {
		return fmt.Errorf("attempt to reuse component.GOiD %d", index)
	}
	for i := range tm.compList {
		tm.compList[i][index].MakeIdentity()
	}

	tm.boundingSpheres[index] = bound

	q, ok := tm.parentChildMap[parent]
	if !ok {
		tm.parentChildMap[parent] = common.IntQueue{}
	}
	q.Queue(int(index))

	mat := math.Mat4x4{}
	mat.MakeIdentity()
	mat[3]  = bound.Center[0]
	mat[7]  = bound.Center[1]
	mat[11] = bound.Center[2]
	tm.compList[LMAT][index] = mat
	tm.compList[WMAT][index] = *tm.compList[WMAT][parent].Mult(&tm.compList[WMAT][index])

	return nil
}
func (tm *SceneManager) resizeArray(index component.GOiD) {
	if cap(tm.compList[WMAT]) - 1 < int(index) {

		for i := range tm.compList {
			newCompList := make([]math.Mat4x4, index + RESIZESTEP)
			for j := range tm.compList[i] {
				newCompList[j] = tm.compList[i][j]
			}
			tm.compList[i] = newCompList
		}

	}

	if cap(tm.boundingSpheres) - 1 < int(index) {
		tmp := tm.boundingSpheres
		tm.boundingSpheres = make([]math.Sphere, index + RESIZESTEP)
		for i := range tmp {
			tm.boundingSpheres[i] = tmp[i]
		}
	}
}

func (tm *SceneManager) DeleteComponent(index component.GOiD) {
	for i := range tm.compList {
		tm.compList[i][index] = math.Mat4x4{}
	}
}

func (tm *SceneManager) Transform(index component.GOiD, newLocalMat *math.Mat4x4) {
	tm.compList[LMAT][index] = *newLocalMat
	// go func() {
	// 	lmat 	:= tm.compList[LMAT][index]
	// 	parent 	:= tm.compList[WMAT][int(tm.parentMap[index])]

	// 	tm.compList[WMAT][index] = *parent.Mult(&lmat)
	// }()
	// fmt.Println(index, "newLocalMat ", newLocalMat.ToString())
	tm.movedQueue.Queue(int(index))
}
func (tm *SceneManager) GetTransform(index component.GOiD) (*math.Mat4x4, error) {
	if int(index) >= len(tm.compList[WMAT]) {
		common.Log.Error("invalid GOiD %v", index)
	}
	if tm.compList[WMAT][index].IsEmpty() == true {
		common.Log.Error("invalid GOiD: %v", index)
	}
		return &tm.compList[WMAT][index], nil
}

func (tm *SceneManager) GetObjectsInLocationRange(loc math.Vec3, lookRange float32) *common.IntQueue {
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