package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/component/gofactory"
	"github.com/stnma7e/betuol/component/scenemanager"
	"github.com/stnma7e/betuol/event"
	"github.com/stnma7e/betuol/math"
	"github.com/stnma7e/betuol/res"
)

func TestPlayerCreation(t *testing.T) {
	rm := res.MakeResourceManager("./data/")
	em := event.MakeEventManager()
	tm := scenemanager.MakeTransformManager(em)
	gof := gofactory.MakeGameObjectFactory(tm)

	player, err := CreateObject(gof, rm, "player", math.Vec3{0, 0, 0})
	if err != nil {
		common.LogErr.Println(err)
		t.Fail()
	}
	tm.SetLocation(player, math.Vec3{1, 0, 0})

	dtime := (16 * time.Millisecond).Seconds()
	em.Tick(dtime)
	tm.Tick(dtime)
}

// CreateObject is a helper function to create a GameObject with a starting location.
func CreateObject(gof *gofactory.GameObjectFactory, rm *res.ResourceManager, objName string, loc math.Vec3) (component.GOiD, error) {
	components, err := rm.LoadGameObject(objName)
	if err != nil {
		return 0, err
	}
	id, err := gof.Create(components, loc)
	if err != nil {
		return 0, fmt.Errorf("gameobject %s creation failed, error: %s", objName, err.Error())
	}

	// is.pm.AddForce(id, math.Vec3{0,0.5,0})
	common.LogInfo.Println("entity created, id:", id)

	return id, nil
}
