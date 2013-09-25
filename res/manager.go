package res

import (
	"os"
	"io"
	"fmt"
	"encoding/json"
	"strings"
	"strconv"

	"smig/component"
	"smig/common"
	"smig/math"
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
	objJson := rm.GetFileContents("map/gameobject/" + objType + "/" + "obj.json")
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

func (rm *ResourceManager) LoadModelWavefront(modelName string) (*common.Vector, *common.Vector, *common.Vector, *common.Vector) {
	modelStr := rm.GetFileContents("graphics/mesh/" + modelName + ".obj")
	verts := common.MakeVector()
	norms := common.MakeVector()
	texes := common.MakeVector()
	index := common.MakeVector()
	lines := strings.SplitAfter(string(modelStr), "\n")
	for i := range lines {
		words := strings.Fields(lines[i])
		if len(words) < 1 {
			break
		}

		switch words[0] {
		case "v":
			x, err := strconv.ParseFloat(words[1], 32)
			y, err := strconv.ParseFloat(words[2], 32)
			z, err := strconv.ParseFloat(words[3], 32)
			if err != nil {
				fmt.Println(err)
			}
			verts.Push_back(math.Vec3{float32(x), float32(y), float32(z)}, 1, 1)
		case "vn":
			x, err := strconv.ParseFloat(words[1], 32)
			y, err := strconv.ParseFloat(words[2], 32)
			z, err := strconv.ParseFloat(words[3], 32)
			if err != nil {
				fmt.Println(err)
			}
			norms.Push_back(math.Vec3{float32(x), float32(y), float32(z)}, 1, 1)
		case "vt":
			u, err := strconv.ParseFloat(words[1], 32)
			v, err := strconv.ParseFloat(words[2], 32)
			if err != nil {
				fmt.Println(err)
			}
			texes.Push_back(math.Vec2{float32(u), float32(v)}, 1, 1)
		case "f":
			numAttributes := 2
			if !strings.Contains(strings.Join(words[1:], " "), "//") {
				numAttributes = 3
			}
			ints := strings.FieldsFunc(strings.Join(words[1:], " "), func(c rune) bool {
				if c == '/' {
					return true
				}
				if c == ' ' {
					return true
				}
				return false
			})
			one,   err 	:= strconv.ParseInt(ints[0*numAttributes], 10, 32)
			two,   err 	:= strconv.ParseInt(ints[1*numAttributes], 10, 32)
			three, err 	:= strconv.ParseInt(ints[2*numAttributes], 10, 32)
			if err != nil {
				fmt.Println(err)
			}
			index.Push_back(uint32(one-1), 1, 1)
			index.Push_back(uint32(two-1), 1, 1)
			index.Push_back(uint32(three-1), 1, 1)
		}
		
	}

	return verts, index, norms, texes
}