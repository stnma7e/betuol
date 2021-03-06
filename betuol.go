package main

import (
	"bufio"
	"fmt"
	//"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"flag"
	"runtime/pprof"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/event"
	"github.com/stnma7e/betuol/instance"
	"github.com/stnma7e/betuol/math"
	"github.com/stnma7e/betuol/net"
	"github.com/stnma7e/betuol/res"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			common.LogErr.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	//rand.Seed(time.Now().UnixNano())

	em := event.MakeEventManager()
	rm := res.MakeResourceManager("./data/")

	eventlink := make(chan event.Event)
	nm := net.MakeNetworkManager(em, "localhost:13560", eventlink)
	var is instance.IInstance = instance.MakeInstance(rm, em.Send)

	em.RegisterListeningFunction(func(evt event.Event) {
		rcce, ok := evt.(*event.RequestCharacterCreationEvent)
		if !ok {
			common.LogErr.Println("here2")
			return
		}

		char, err := is.CreateObject(rcce.Type, rcce.Location)
		if err != nil {
			common.LogErr.Println(err)
		}

		net.SendEvent(&event.ApproveCharacterCreationRequestEvent{
			char,
		}, rcce.RequestOrigin)
	}, "requestCharacterCreation")

	em.RegisterListeningFunction(func(evt event.Event) {
		cme, ok := evt.(*event.CharacterMoveEvent)
		if !ok {
			common.LogErr.Println("here3")
			return
		}

		is.MoveGameObject(cme.CharID, cme.NewLocation)
	}, "characterMoved")

	go func() {
		em.RegisterListeningChannel(eventlink, []string{
			"attack",
			"death",
			"characterMoved",
		}...)
		nm.Tick()
	}()

	player, err := is.CreateObject("player", math.Vec3{0, 0, 0})
	if err != nil {
		common.LogErr.Println(err)
		return
	}
	is.MoveGameObject(player, math.Vec3{1, 0, 0})

	is.AddQuest(player, "attack")
	if err := instance.StartScript(is, player); err != nil {
		common.LogErr.Println(err)
	}

	commandlink := make(chan string)
	go commandLineReader(is, player, commandlink)

	oldTime := time.Now()
	ticks := time.NewTicker(time.Second / 60)

	for {
		select {
		case <-ticks.C:
			if newTime := time.Since(oldTime).Seconds() * 1000; newTime > 16 {
				//common.LogWarn.Printf("last tick took longer than 1/60th of a second. it took %vms.", newTime)
			}
		default:
			<-ticks.C
		}
		newTime := time.Since(oldTime)
		oldTime = time.Now()

		if exit := parseSysConsole(is, player, commandlink); exit == true {
			break
		}

		dtime := newTime.Seconds()
		em.Tick(dtime)
		is.Tick(dtime)
	}
}

// ParseSysConsole is used to parse the commands input into the console.  // These commands can be used to control the internal state of the instance.
func parseSysConsole(is instance.IInstance, player component.GOiD, commandlink chan string) bool {
	command := ""
	select {
	case command = <-commandlink:
	default:
		return false
	}

	args := strings.SplitAfter(command, " ")
	for i := range args {
		if i == len(args)-1 {
			continue
		}
		args[i] = args[i][:len(args[i])-1]
	}

	switch args[0] {
	case "quit":
		return true
	case "loadmap":
		if !(len(args) >= 2) {
			common.LogErr.Print("not enough arguments to 'loadmap'")
			break
		}
		is.CreateFromMap(args[1])
	case "loadobj":
		if !(len(args) >= 3) {
			common.LogErr.Print("not enough arguments to 'loadobj'")
			break
		}
		//is.CreateObject(breed, location)
		is.CreateObject(args[1], decodeVec3String(args[2]))
	case "runai":
		if !(len(args) >= 2) {
			common.LogErr.Print("not enough arguments to 'runai'")
			break
		}
		arg, err := strconv.Atoi(args[1])
		if err != nil {
			common.LogErr.Println(err)
			break
		}
		is.RunAi(component.GOiD(arg))
	case "player":
		is.RunAi(player)
	case "render":
		if !(len(args) >= 2) {
			common.LogErr.Print("not enough arguments to 'render'")
			break
		}
		arg, err := strconv.Atoi(args[1])
		if err != nil {
			common.LogErr.Println(err)
			break
		}
		is.RenderFromPerspective(component.GOiD(arg))
	default:
		fmt.Println("\tInvalid command. Type \"help\" for choices.")
	}

	commandlink <- ""
	return false
}

func commandLineReader(is instance.IInstance, player component.GOiD, commandlink chan string) {
	for {
		r := bufio.NewReaderSize(os.Stdin, 4*1024)
		line, err := r.ReadString('\n')
		if err != nil {
			common.LogErr.Println(err)
			continue
		}
		s := string(line)
		if len(s) > 0 && s[len(s)-1] == '\n' {
			s = line[:len(s)-1]
		}
		if len(s) > 0 && s[len(s)-1] == '\r' {
			s = line[:len(s)-1]
		}
		if s != "" {
			commandlink <- s
			s = ""
			<-commandlink
		}
	}
}

// decodeVec3String is a helper function to extract a 3 dimensional vector from a string
func decodeVec3String(vec3 string) math.Vec3 {
	strLoc := strings.Split(vec3, ",")
	f1, err := strconv.ParseFloat(strLoc[0], 32)
	f2, err := strconv.ParseFloat(strLoc[1], 32)
	f3, err := strconv.ParseFloat(strLoc[2], 32)
	if err != nil {
		common.LogErr.Println(err)
	}

	return math.Vec3{float32(f1), float32(f2), float32(f3)}
}
