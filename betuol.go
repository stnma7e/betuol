package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/instance"
	"github.com/stnma7e/betuol/math"
	"github.com/stnma7e/betuol/res"
)

func main() {
	rm := res.MakeResourceManager("./data/")
	var is instance.IInstance = instance.MakeInstance(rm)

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

	exitlink := make(chan bool)
	go func() {
		for {
			commandLineReader(is, player, exitlink)
		}
	}()

	oldTime := time.Now()
	ticks := time.NewTicker(time.Second / 60)

	for {
		select {
		case <-ticks.C:
			common.LogWarn.Printf("last tick took longer than 1/60th of a second. it took %v", time.Since(oldTime))
		default:
			<-ticks.C
		}
		newTime := time.Since(oldTime)
		oldTime = time.Now()

		select {
		case <-exitlink:
			return
		default:
		}

		is.Tick(newTime.Seconds())
	}
}

// ParseSysConsole is used to parse the commands input into the console.  // These commands can be used to control the internal state of the instance.
func parseSysConsole(is instance.IInstance, player component.GOiD, exitlink chan bool, command string) {
	args := strings.SplitAfter(command, " ")
	for i := range args {
		if i == len(args)-1 {
			continue
		}
		args[i] = args[i][:len(args[i])-1]
	}
	fmt.Println(args)

	switch args[0] {
	case "exit":
		exitlink <- true
		return
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
		}
		is.RunAi(component.GOiD(arg))
	case "player":
		is.RunAi(player)
	default:
		fmt.Println("\tInvalid command. Type \"help\" for choices.")
	}
}

func commandLineReader(is instance.IInstance, player component.GOiD, exitlink chan bool) {
	r := bufio.NewReaderSize(os.Stdin, 4*1024)
	line, err := r.ReadString('\n')
	if err != nil {
		common.LogErr.Println(err)
	}
	s := string(line)
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = line[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '\r' {
		s = line[:len(s)-1]
	}
	if s != "" {
		parseSysConsole(is, player, exitlink, s)
		s = ""
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
