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
)

type SceneManager struct {
	compList 		[2][]math.Mat4x4 // compList[WMAT] == world transform matrices
								 	 // compList[LMAT] == local transform matrices
	parentChildMap 	map[component.GOiD]common.IntQueue
	childParentMap  map[component.GOiD]component.GOiD
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

func (tm *SceneManager) Render(glg graphics.GraphicsManager) {
	go func() {
		for i := range tm.meshList {
			glg.Render(&tm.meshList[i], &tm.compList[WMAT][i])
		}
	}()
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
		Location [3]float32
	}
	err := json.Unmarshal(compData, &comp)
	if err != nil {
		panic(err)
	}

	err = tm.CreateComponent(index, ROOTNODE)
	if err != nil {
		panic(err)
	}

	mat := math.Mat4x4{}
	mat.MakeIdentity()
	mat[3]  = comp.Location[0]
	mat[7]  = comp.Location[1]
	mat[11] = comp.Location[2]
	tm.Transform(index, &mat)

	return nil
}
func (tm *SceneManager) CreateComponent(index, parent component.GOiD) error {
	if component.GOiD(cap(tm.compList[WMAT]))-1 < index {
		for i := range tm.compList {
			newCompList := make([]math.Mat4x4, index + 25)
			for j := range tm.compList[i] {
				newCompList[j] = tm.compList[i][j]
			}
			tm.compList[i] = newCompList
		}
	}

	if !(tm.compList[WMAT][index].IsEmpty()) {
		return fmt.Errorf("attempt to reuse component.GOiD %d", index)
	}

	for i := range tm.compList {
		tm.compList[i][index].MakeIdentity()
	}

	q, ok := tm.parentChildMap[parent]
	if !ok {
		tm.parentChildMap[parent] = common.IntQueue{}
	}

	q.Queue(int(index))

	return nil
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
	if tm.compList[WMAT][index].IsEmpty() != true {
		return &tm.compList[WMAT][index], nil
	}
	panic(fmt.Sprintf("invalid GOiD: %v", index))
}

func (tm *SceneManager) FindObjectsInRange(rootId component.GOiD, radius float32) []component.GOiD {
	return []component.GOiD{}
}