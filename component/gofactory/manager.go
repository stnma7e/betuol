// Pacakge gofactory is used for the craetion and deletion of GameObjects as a whole entity.
package gofactory

import (
	"fmt"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/component/scenemanager"
	"github.com/stnma7e/betuol/event"
	"github.com/stnma7e/betuol/math"
)

const SceneType = "scene"

// creationFunction is a function template that is used to create components based on data retrieved from the file system.
type creationFunction func(component.GOiD, []byte) error

// creationManager is a helper struct to pair a component.ComponentManager with a creationFunction.
type creationManager struct {
	mang   component.ComponentManager
	create creationFunction
}

// GameObjectFactory is used to keep track of component managers and GOiD's in use.
// The GameObjectFactory is used to create and delete GameObjects and will facilitate the reuse of GOiD's without overlap.
type GameObjectFactory struct {
	topIndex      component.GOiD
	EventManagers map[string]creationManager
	vacantIndices common.Queue
	tm            *scenemanager.TransformManager
}

// MakeGameObjectFactory returns a pointer to a GameObjectFactory.
func MakeGameObjectFactory(tm *scenemanager.TransformManager) *GameObjectFactory {
	gof := GameObjectFactory{
		1,
		make(map[string]creationManager),
		common.Queue{},
		tm,
	}
	return &gof
}

// Register registers a component.ComponentManager to the GameObjectFactory to be used for component creation.
// The compType string passed as an argument is used to assosiate a component type with the component manager.
func (gof *GameObjectFactory) Register(compType string, mang component.ComponentManager, creationFunc creationFunction) error {
	_, isthere := gof.EventManagers[compType]
	if isthere {
		return fmt.Errorf("attempt to register manager (%s) failed; space already filled", compType)
	} else { // manager does not already exist
		gof.EventManagers[compType] = creationManager{
			mang,
			creationFunc,
		}
	}
	return nil
}

// CreateFromMap is used to create a list of GameObjects based on a map format defined in component.
func (gof *GameObjectFactory) CreateFromMap(sceneMap *component.Map) ([]component.GOiD, error) {
	var idQueue common.Queue
	smap := *sceneMap
	for i := range smap {
		for j := range smap[i].Entities {
			for k := 0; k < smap[i].Entities[j].Quantity; k++ {
				id, err := gof.Create(smap[i].Entities[j].CompList, smap[i].Location)
				if err != nil {
					return []component.GOiD{}, err
				} else {
					idQueue.Queue(int(id))
				}
			}
		}
	}

	count := idQueue.Size
	idList := make([]component.GOiD, count)
	for i := 0; !idQueue.IsEmpty(); i++ {
		num, err := idQueue.Dequeue()
		if err != nil {
			common.LogErr.Print(err)
			return []component.GOiD{}, fmt.Errorf("map creation failed when dequeueing gameobject %d from map list, error: %s", i, err.Error())
		}
		idList[i] = component.GOiD(num.(int))
	}

	return idList, nil
}

// Create will return the GOiD of a GameObject after creating it. It will use the map of creation data to create various components.
func (gof *GameObjectFactory) Create(compList component.GameObject, location math.Vec3) (component.GOiD, error) {
	id := gof.getNewGOiD()
	if id < 1 {
		return 0, fmt.Errorf("invalid id: %v", id)
	}

	err := gof.tm.CreateComponent(id)
	if err != nil {
		gof.Delete(id)
		return component.NULLINDEX, err
	}
	gof.tm.SetLocation(id, location)

	for k, v := range compList {
		mang, ok := gof.EventManagers[k]
		if !ok {
			common.LogWarn.Println("unregistered component type (%s) in compList", k)
			continue
		}

		err := mang.create(id, v)
		if err != nil {
			gof.Delete(id)
			return component.NULLINDEX, err
		}
	}

	return id, nil
}

// Delete will call the DeleteComponent function of each registered component manager to remove the GOiD from use in the system.
func (gof *GameObjectFactory) Delete(index component.GOiD) {
	if index == 0 {
		return
	}
	for _, v := range gof.EventManagers {
		v.mang.DeleteComponent(index)
	}
	gof.tm.DeleteComponent(index)
	gof.vacantIndices.Queue(int(index))
}

// getNewGOiD is a helper function to facilitate the reuse of GOiD's without overlap.
func (gof *GameObjectFactory) getNewGOiD() component.GOiD {
	var idToUse component.GOiD
	if !(gof.vacantIndices.IsEmpty()) { // there are availiable pre-owned IDs
		id, err := gof.vacantIndices.Dequeue()
		if err != nil {
			fmt.Println(err)
		}
		idToUse = component.GOiD(id.(int))
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
	gof.Delete(devt.Id)
	fmt.Printf("%v died.\n", devt.Id)
}
