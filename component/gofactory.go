package component

import (
	"fmt"

	"smig/common"
	"smig/math"
)

const SceneType = "scene"

type CreationFunction func(GOiD, []byte) error

type CreationManager struct {
	mang ComponentManager
	create CreationFunction
}

type GameObjectFactory struct {
	topIndex GOiD
	EventManagers map[string]CreationManager
	vacantIndices common.IntQueue
	sm *SceneManager
}

func MakeGameObjectFactory(sm *SceneManager) *GameObjectFactory {
	gof := GameObjectFactory{ 1, make(map[string]CreationManager), common.IntQueue{}, sm }
	return &gof
}

func (gof *GameObjectFactory) Register(compType string, mang ComponentManager, creationFunc CreationFunction) error {
	_, isthere := gof.EventManagers[compType]
	if isthere {
		return fmt.Errorf("attemp to register manager (%s) failed; space already filled", compType)
	} else { // manager does not already exist
		gof.EventManagers[compType] = CreationManager {
			mang,
			creationFunc,
		}
	}
	return nil
}

func (gof *GameObjectFactory) CreateFromMap(sceneMap *Map) []GOiD {
	var idQueue common.IntQueue
	smap := *sceneMap
	for i := range smap {
		for j := range smap[i].Entities {
			for k := 0; k < smap[i].Entities[j].Quantity; k++ {
				id, err := gof.Create(smap[i].Entities[j].CompList, smap[i].Location, smap[i].Entities[j].Radius)
				if err != nil {
					common.Log.Error(err)
				} else {
					idQueue.Queue(int(id))
				}
			}
		}
	}

	count := idQueue.Size
	idList := make([]GOiD,count)
	for i := 0; !idQueue.IsEmpty(); i++ {
		num, err := idQueue.Dequeue()
		if err != nil {
			fmt.Println(err)
		}
		idList[i] = GOiD(num)
	}

	return idList
}
func (gof *GameObjectFactory) Create(compList GameObject, location math.Vec3, radius float32) (GOiD, error) {
	id := gof.getNewGOiD()
	if id < 1 {
		common.Log.Error("invalid id: %v", id)
	}

	for k,v := range compList {
		mang, ok := gof.EventManagers[k]
		if !ok {
			common.Log.Error("unregistered componentManager (%s) in compList", k)
		}

		err := mang.create(id, v)
		if err != nil {
			gof.Delete(id)
			common.Log.Error(err)
			return NULLINDEX, err
		}
	}

	sp := math.Sphere {
		location, radius,
	}
	err := gof.sm.CreateComponent(id, ROOTNODE, sp)
	if err != nil {
		gof.Delete(id)
		common.Log.Error(err)
		return NULLINDEX, err
	}

	return id, nil
}
func (gof *GameObjectFactory) Delete(index GOiD) {
	if index == 0 {
		return
	}
	for _,v := range gof.EventManagers {
		v.mang.DeleteComponent(index)
	}
	gof.sm.DeleteComponent(index)
	gof.vacantIndices.Queue(int(index))
}

func (gof *GameObjectFactory) getNewGOiD() GOiD {
	var idToUse GOiD
	if !(gof.vacantIndices.IsEmpty()) { // there are availiable pre-owned IDs
		id, err := gof.vacantIndices.Dequeue()
		if err != nil {
			fmt.Println(err)
		}
		idToUse = GOiD(id)
		// if an error is received, then the paired ID will be 0
		// 0 will be rejected by all of the component managers
	} else {
		idToUse = gof.topIndex
		gof.topIndex++
	}

	return idToUse
}