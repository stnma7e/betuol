package res

import (
	"os"
	"io"
	"encoding/json"

	"smig/component"
)

type ResourceManager struct {
	fileDepot string
	resMap map[string][]byte
}

func MakeResourceManager(fileDepot string) *ResourceManager {
	return &ResourceManager {
		fileDepot,
		make(map[string][]byte),
	}
}

func (rm *ResourceManager) GetFileContents(fileName string) []byte {
	file, err := os.Open(rm.fileDepot + fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	stats, err := file.Stat()
	buf := make([]byte, stats.Size())
	_, err = file.Read(buf)
	if err != io.EOF && err != nil {
		panic(err)
	}
	return buf
}

func (rm *ResourceManager) LoadJsonMap(mapName string) component.Map {
	mapJson := rm.GetFileContents("map/" + mapName + ".json")
	var obj []component.MapLocation
	err := json.Unmarshal(mapJson, &obj)
	if err != nil {
		panic(err)
	}

	for i := range obj {
		for j := range obj[i].Entities {
			for k := 0; k < obj[i].Entities[j].Quantity; k++ {
				breedStr := rm.GetFileContents("map/gameobject/" + obj[i].Entities[j].Breed + "/obj.json")
				obj[i].Entities[j].CompList = rm.LoadGameObject(breedStr, obj[i].Entities[j].Breed)
			}
		}
	}

	return obj
}

func (rm *ResourceManager) LoadGameObject(objJson []byte, objType string) component.GameObject {
	var obj struct {
		Component []string
	}
	err := json.Unmarshal(objJson, &obj)
	if err != nil {
		panic(err)
	}

	gameobj := make(map[string][]byte)
	for i := range obj.Component {
		gameobj[obj.Component[i]] = nil
	}

	for k,_ := range gameobj {
		gameobj[k] = rm.GetFileContents("map/gameobject/" + objType + "/" + k + ".json")
	}

	return gameobj
}