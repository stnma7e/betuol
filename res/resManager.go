package res

import (
	"os"
	"io"
	"fmt"
	"encoding/json"

	"smig/component"
	"smig/common"
)

const (
	OBJFILENAME = "obj.json"
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
		common.Log.Error(fmt.Sprint(err))
	}
	defer file.Close()
	stats, err := file.Stat()
	buf := make([]byte, stats.Size())
	_, err = file.Read(buf)
	if err != io.EOF && err != nil {
		common.Log.Error(err)
	}
	return buf
}

func (rm *ResourceManager) LoadJsonMap(mapName string) component.Map {
	mapJson := rm.GetFileContents("map/" + mapName + ".json")
	var obj []component.MapLocation
	err := json.Unmarshal(mapJson, &obj)
	if err != nil {
		common.Log.Error(fmt.Sprint(err))
	}

	for i := range obj {
		for j := range obj[i].Entities {
			for k := 0; k < obj[i].Entities[j].Quantity; k++ {
				obj[i].Entities[j].CompList = rm.LoadGameObject(obj[i].Entities[j].Breed)
			}
		}
	}

	return obj
}

func (rm *ResourceManager) LoadGameObject(objType string) component.GameObject {
	objJson := rm.GetFileContents("map/gameobject/" + objType + "/" + OBJFILENAME)
	var obj struct {
		Component []string
	}
	err := json.Unmarshal(objJson, &obj)
	if err != nil {
		common.Log.Error(fmt.Sprint(err))
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