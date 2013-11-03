package gofactory

import (
	"fmt"

	"smig/component"
	"smig/common"
	"smig/math"
	"smig/event"
)

const SceneType = "scene"

type CreationFunction func(component.GOiD, []byte) error

type CreationManager struct {
	mang component.ComponentManager
	create CreationFunction
}

type GameObjectFactory struct {
	topIndex component.GOiD
	EventManagers map[string]CreationManager
	vacantIndices common.IntQueue
	tm *component.TransformManager
}

func MakeGameObjectFactory(tm *component.TransformManager) *GameObjectFactory {
	gof := GameObjectFactory{ 1, make(map[string]CreationManager), common.IntQueue{}, tm }
	return &gof
}

func (gof *GameObjectFactory) Register(compType string, mang component.ComponentManager, creationFunc CreationFunction) error {
	_, isthere := gof.EventManagers[compType]
	if isthere {
		return fmt.Errorf("attempt to register manager (%s) failed; space already filled", compType)
	} else { // manager does not already exist
		gof.EventManagers[compType] = CreationManager {
			mang,
			creationFunc,
		}
	}
	return nil
}

func (gof *GameObjectFactory) CreateFromMap(sceneMap *component.Map) []component.GOiD {
	var idQueue common.IntQueue
	smap := *sceneMap
	for i := range smap {
		for j := range smap[i].Entities {
			for k := 0; k < smap[i].Entities[j].Quantity; k++ {
				id, err := gof.Create(smap[i].Entities[j].CompList, smap[i].Location)
				if err != nil {
					common.LogErr.Print(err)
				} else {
					idQueue.Queue(int(id))
				}
			}
		}
	}

	count := idQueue.Size
	idList := make([]component.GOiD,count)
	for i := 0; !idQueue.IsEmpty(); i++ {
		num, err := idQueue.Dequeue()
		if err != nil {
			common.LogErr.Print(err)
		}
		idList[i] = component.GOiD(num)
	}

	return idList
}
func (gof *GameObjectFactory) Create(compList component.GameObject, location math.Vec3) (component.GOiD, error) {
	id := gof.getNewGOiD()
	if id < 1 {
		common.LogErr.Print("invalid id: %v", id)
	}

	err := gof.tm.CreateComponent(id)
	if err != nil {
		gof.Delete(id)
		common.LogErr.Print(err)
		return component.NULLINDEX, err
	}

	for k,v := range compList {
		mang, ok := gof.EventManagers[k]
		if !ok {
			common.LogWarn.Print("unregistered component.ComponentManager (%s) in compList", k)
		}

		err := mang.create(id, v)
		if err != nil {
			gof.Delete(id)
			fmt.Println(err)
			return component.NULLINDEX, err
		}
	}

	return id, nil
}
func (gof *GameObjectFactory) Delete(index component.GOiD) {
	if index == 0 {
		return
	}
	for _,v := range gof.EventManagers {
		v.mang.DeleteComponent(index)
	}
	gof.tm.DeleteComponent(index)
	gof.vacantIndices.Queue(int(index))
}

func (gof *GameObjectFactory) getNewGOiD() component.GOiD {
	var idToUse component.GOiD
	if !(gof.vacantIndices.IsEmpty()) { // there are availiable pre-owned IDs
		id, err := gof.vacantIndices.Dequeue()
		if err != nil {
			fmt.Println(err)
		}
		idToUse = component.GOiD(id)
		// if an error is received, then the paired ID will be 0
		// 0 will be rejected by all of the component managers
	} else {
		idToUse = gof.topIndex
		gof.topIndex++
	}

	return idToUse
}

func (gof *GameObjectFactory) HandleDeath(evt event.Event) {
	devt := evt.(event.DeathEvent)
	fmt.Printf("%v died.\n", devt.Id)
	gof.Delete(devt.Id)
}
