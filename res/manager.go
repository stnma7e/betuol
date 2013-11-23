// Package res implements functions to handle file system and resource management.
package res

import (
	"encoding/json"
	"fmt"
	"io"
	gomath "math"
	"os"
	"strconv"
	"strings"

	"betuol/common"
	"betuol/component"
	"betuol/math"
)

// ResourceManager is a struct to handle resource management and to prevent multiple identical resources being loaded in memory simultaneously.
type ResourceManager struct {
	fileDepot string
	resMap    map[string][]byte
}

// MakeResourceManager returns a pointer to a ResourceManager.
func MakeResourceManager(fileDepot string) *ResourceManager {
	return &ResourceManager{
		fileDepot,
		make(map[string][]byte),
	}
}

// GetFileContents is a static wrapper for ResourceManager.GetFileContents.
func GetFileContents(fileName string) []byte {
	rm := MakeResourceManager("/home/sam/go/src/betuol/data/")
	return rm.GetFileContents(fileName)
}

// GetFileContents is a function to retreive the contols of a file returned as a byte array.
func (rm *ResourceManager) GetFileContents(fileName string) []byte {
	file, err := os.Open(rm.fileDepot + fileName)
	if err != nil {
		common.LogErr.Fatal(err)
	}
	defer file.Close()
	stats, err := file.Stat()
	buf := make([]byte, stats.Size())
	_, err = file.Read(buf)
	if err != io.EOF && err != nil {
		common.LogErr.Print(err)
	}
	return buf
}

// LoadJsonMap will return a component.Map used by the GameObjectFactory to create a list of GameObjects organized by location.
func (rm *ResourceManager) LoadJsonMap(mapName string) component.Map {
	mapJson := rm.GetFileContents("map/" + mapName + ".json")
	var obj []component.MapLocation
	err := json.Unmarshal(mapJson, &obj)
	if err != nil {
		common.LogErr.Print(err)
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

// LoadGameObject is used to retrieve and compile the creation data for making a GameObject with a GameObjectFactory.
func (rm *ResourceManager) LoadGameObject(objType string) component.GameObject {
	objJson := rm.GetFileContents("map/gameobject/" + objType + "/" + "obj.json")
	var obj struct {
		Component []string
	}
	err := json.Unmarshal(objJson, &obj)
	if err != nil {
		common.LogErr.Print(err)
	}

	gameobj := make(map[string][]byte)
	for i := range obj.Component {
		gameobj[obj.Component[i]] = nil
	}

	for k, _ := range gameobj {
		gameobj[k] = rm.GetFileContents("map/gameobject/" + objType + "/" + k + ".json")
	}

	return gameobj
}

// LoadModelWavefront is a function to load a Wavefront (.obj) style 3D mesh.
// This function is used to load models created using blender or another 3d mdoeling program outside the game engine.
// The function parses a text file and extracts a list of vertices, indices, normal vectors, and texture coordinates.
func (rm *ResourceManager) LoadModelWavefront(modelName string) (*common.Vector, *common.Vector, *common.Vector, *common.Vector, float32) {
	modelStr := GetFileContents("graphics/mesh/" + modelName + ".obj")
	verts := common.MakeVector()
	norms := common.MakeVector()
	texes := common.MakeVector()
	index := common.MakeVector()
	lines := strings.SplitAfter(string(modelStr), "\n")
	var maxDistanceSqrd float32
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
			vec := math.Vec3{float32(x), float32(y), float32(z)}
			verts.Push_back(vec, 1, 1)
			for i := range verts.Array() {
				distSqrd := math.DistSqrd3v3v(vec, verts.Array()[i].(math.Vec3))
				if distSqrd > maxDistanceSqrd {
					maxDistanceSqrd = distSqrd
				}
			}
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
			numAttributes := 1
			if strings.Contains(strings.Join(words[1:], " "), "/") {
				numAttributes = 3
				if strings.Contains(strings.Join(words[1:], " "), "//") {
					numAttributes = 2
				}
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
			one, err := strconv.ParseInt(ints[0*numAttributes], 10, 32)
			two, err := strconv.ParseInt(ints[1*numAttributes], 10, 32)
			three, err := strconv.ParseInt(ints[2*numAttributes], 10, 32)
			if err != nil {
				fmt.Println(err)
			}
			index.Push_back(uint32(one-1), 1, 1)
			index.Push_back(uint32(two-1), 1, 1)
			index.Push_back(uint32(three-1), 1, 1)
		}

	}

	if verts.IsEmpty() {
		verts.Push_back(math.Vec3{}, 1, 1)
	}
	if index.IsEmpty() {
		index.Push_back(0, 1, 1)
	}
	if norms.IsEmpty() {
		norms.Push_back(math.Vec3{}, 1, 1)
	}
	if texes.IsEmpty() {
		texes.Push_back(math.Vec2{}, 1, 1)
	}

	boundingRadius := float32(gomath.Sqrt(float64(maxDistanceSqrd))) / 2
	return verts, index, norms, texes, boundingRadius
}
