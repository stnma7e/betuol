package component

import (
	"fmt"

	"smig/common"
	"smig/math"
)

const SceneType = "transform"

type CreationFunction func(GOiD, []byte) error

type CreationManager struct {
	mang ComponentManager
	create CreationFunction
}

type GameObjectFactory struct {
	topIndex GOiD
	EventManagers map[string]CreationManager
	vacantIndices common.IntQueue
}

func MakeGameObjectFactory() *GameObjectFactory {
	gof := GameObjectFactory{ 1, make(map[string]CreationManager), common.IntQueue{} }
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
	var count uint
	smap := *sceneMap
	for i := range smap {
		for j := range smap[i].Entities {
			for k := 0; k < smap[i].Entities[j].Quantity; k++ {

				bytes := smap[i].Location.ToJson()
				data := makeSceneTypeData(bytes)
				smap[i].Entities[j].CompList[SceneType] = data

				id, err := gof.Create(smap[i].Entities[j].CompList)
				if err != nil {
					common.Log.Error(err)
				}
				idQueue.Queue(int(id))

				count++
			}
		}
	}

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
func (gof *GameObjectFactory) Create(compList GameObject) (GOiD, error) {
	id := gof.getNewGOiD()
	if id < 1 {
		common.Log.Error("invalid id: %v", id)
	}

	transformed := false
	for k,v := range compList {
		if k == SceneType {
			transformed = true
		}

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

	if !transformed {
		mang, ok := gof.EventManagers[SceneType]
		if !ok {
			common.Log.Error("unregistered componentManager (%s) in compList", SceneType)
		}

		loc := math.Vec3{}
		err := mang.create(id, makeSceneTypeData(loc.ToJson()))
		if err != nil {
			gof.Delete(id)
			common.Log.Error(err)
			return NULLINDEX, err
		}
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

func makeSceneTypeData(loc []byte) []byte {
	data := make([]byte,len(loc)+13)
	data[0] = '{'
	data[1] = '"'
	data[2] = 'l'
	data[3] = 'o'
	data[4] = 'c'
	data[5] = 'a'
	data[6] = 't'
	data[7] = 'i'
	data[8] = 'o'
	data[9] = 'n'
	data[10] = '"'
	data[11] = ':'
	for i := range loc {
		data[i+12] = loc[i]
	}
	data[len(loc)+12] = '}'

	return data
}