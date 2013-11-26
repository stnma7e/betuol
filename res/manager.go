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

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/math"
)

// WavefrontData is a struct to contain the data extracted from a Wavefront .obj file. This is used by the ResourceManager to return to the graphics manager.
type WavefrontData struct {
	Vertices, Indices, Normals, Uvs *common.Vector
	BoundingRadius                  float32
}

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
func GetFileContents(fileName string) ([]byte, error) {
	rm := MakeResourceManager("./data/")
	return rm.GetFileContents(fileName)
}

// GetFileContents is a function to retreive the contols of a file returned as a byte array.
func (rm *ResourceManager) GetFileContents(fileName string) ([]byte, error) {
	file, err := os.Open(rm.fileDepot + fileName)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()
	stats, err := file.Stat()
	buf := make([]byte, stats.Size())
	_, err = file.Read(buf)
	if err != io.EOF && err != nil {
		common.LogErr.Print(err)
	}
	return buf, nil
}

// LoadJsonMap will return a component.Map used by the GameObjectFactory to create a list of GameObjects organized by location.
func (rm *ResourceManager) LoadJsonMap(mapName string) (component.Map, error) {
	mapJson, err := rm.GetFileContents("game/" + mapName + ".json")
	if err != nil {
		return component.Map{}, fmt.Errorf("failed to open map file, error: %s", err.Error())
	}
	var obj []component.MapLocation
	err = json.Unmarshal(mapJson, &obj)
	if err != nil {
		return component.Map{}, err
	}

	for i := range obj {
		for j := range obj[i].Entities {
			for k := 0; k < obj[i].Entities[j].Quantity; k++ {
				obj[i].Entities[j].CompList, err = rm.LoadGameObject(obj[i].Entities[j].Breed)
				if err != nil {
					return component.Map{}, err
				}
			}
		}
	}

	return obj, nil
}

// LoadGameObject is used to retrieve and compile the creation data for making a GameObject with a GameObjectFactory.
func (rm *ResourceManager) LoadGameObject(objType string) (component.GameObject, error) {
	objJson, err := rm.GetFileContents("game/gameobject/" + objType + "/" + "obj.json")
	if err != nil {
		return component.GameObject{}, fmt.Errorf("failed to open gameobject file, error: %s", err.Error())
	}
	var obj struct {
		Component []string
	}
	err = json.Unmarshal(objJson, &obj)
	if err != nil {
		return component.GameObject{}, fmt.Errorf("json failure when loading gameobject, error: %s", err.Error())
	}

	gameobj := make(map[string][]byte)
	for i := range obj.Component {
		gameobj[obj.Component[i]] = nil
	}

	for k, _ := range gameobj {
		gameobj[k], err = rm.GetFileContents("game/gameobject/" + objType + "/" + k + ".json")
		if err != nil {
			return component.GameObject{}, fmt.Errorf("failed to load gameobject file, error: %s", err.Error())
		}
	}

	return gameobj, nil
}

// LoadModelWavefront is a function to load a Wavefront (.obj) style 3D mesh.
// This function is used to load models created using blender or another 3d mdoeling program outside the game engine.
// The function parses a text file and extracts a list of vertices, indices, normal vectors, and texture coordinates.
func (rm *ResourceManager) LoadModelWavefront(modelName string) (WavefrontData, error) {
	modelStr, err := GetFileContents("graphics/mesh/" + modelName + ".obj")
	if err != nil {
		return WavefrontData{}, fmt.Errorf("failed to open wavefront .obj file, error: %s", err.Error())
	}
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
			if err != nil {
				return WavefrontData{}, fmt.Errorf("failed to parse wavefront vertex x coordinate, error: %s", err.Error())
			}
			y, err := strconv.ParseFloat(words[2], 32)
			if err != nil {
				return WavefrontData{}, fmt.Errorf("failed to parse wavefront vertex z coordinate, error: %s", err.Error())
			}
			z, err := strconv.ParseFloat(words[3], 32)
			if err != nil {
				return WavefrontData{}, fmt.Errorf("failed to parse wavefront vertex z coordinate, error: %s", err.Error())
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
			if err != nil {
				return WavefrontData{}, fmt.Errorf("failed to parse wavefront vertex normal x coordinate, error: %s", err.Error())
			}
			y, err := strconv.ParseFloat(words[2], 32)
			if err != nil {
				return WavefrontData{}, fmt.Errorf("failed to parse wavefront vertex normal y coordinate, error: %s", err.Error())
			}
			z, err := strconv.ParseFloat(words[3], 32)
			if err != nil {
				return WavefrontData{}, fmt.Errorf("failed to parse wavefront vertex normal z coordinate, error: %s", err.Error())
			}
			norms.Push_back(math.Vec3{float32(x), float32(y), float32(z)}, 1, 1)
		case "vt":
			u, err := strconv.ParseFloat(words[1], 32)
			if err != nil {
				return WavefrontData{}, fmt.Errorf("failed to parse wavefront texture u coordinate, error: %s", err.Error())
			}
			v, err := strconv.ParseFloat(words[2], 32)
			if err != nil {
				return WavefrontData{}, fmt.Errorf("failed to parse wavefront texture v coordinate, error: %s", err.Error())
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
			if err != nil {
				return WavefrontData{}, fmt.Errorf("failed to parse wavefront vertex index, error: %s", err.Error())
			}
			two, err := strconv.ParseInt(ints[1*numAttributes], 10, 32)
			if err != nil {
				return WavefrontData{}, fmt.Errorf("failed to parse wavefront texture coordinate index, error: %s", err.Error())
			}
			three, err := strconv.ParseInt(ints[2*numAttributes], 10, 32)
			if err != nil {
				return WavefrontData{}, fmt.Errorf("failed to parse wavefront vertex normal index, error: %s", err.Error())
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
	return WavefrontData{verts, index, norms, texes, boundingRadius}, nil
}
