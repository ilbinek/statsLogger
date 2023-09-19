package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/indig0fox/a3go/a3interface"
	"github.com/indig0fox/a3go/assemblyfinder"
)

// modulePath is the absolute path to the compiled DLL, which should be the addon folder
var modulePath string = assemblyfinder.GetModulePath()
// modulePathDir is the containing folder
var modulePathDir string = path.Dir(modulePath)

var EXTENSION_NAME = "STATS_LOGGER"

var mission Mission

var RVExtensionChannels = map[string]chan string {
	":timeNow:" : make(chan string),
}

var RVExtensionArgsChannels = map[string]chan []string{
	":RESET:":    	make(chan []string),
	":MISSION:":   	make(chan []string),
	":WIN:":   		make(chan []string),
	":PLAYER:": 	make(chan []string),
	":KILL:":   	make(chan []string),
	":SHOT:":   	make(chan []string),
	":HIT:":   		make(chan []string),
	":EXPORT:":   	make(chan []string),
	":FPS:":   		make(chan []string),
}

var a3ErrorChan = make(chan error)

func init() {
	a3interface.SetVersion("1.0.0")
	a3interface.RegisterRvExtensionArgsChannels(RVExtensionArgsChannels)

	go func() {
		for {
			select {
			case v := <-RVExtensionChannels[":timeNow:"]:
				go writeTimeNow(v)
			case arg := <-RVExtensionArgsChannels[":RESET:"]:
				go reset(arg)
			case arg := <-RVExtensionArgsChannels[":MISSION:"]:
				go setUpMission(arg)
			case arg := <-RVExtensionArgsChannels[":WIN:"]:
				go win(arg)
			case arg := <-RVExtensionArgsChannels[":PLAYER:"]:
				go addPlayer(arg)
			case arg := <-RVExtensionArgsChannels[":KILL:"]:
				go addKill(arg)
			case arg := <-RVExtensionArgsChannels[":SHOT:"]:
				go addShot(arg)
			case arg := <-RVExtensionArgsChannels[":HIT:"]:
				go addHit(arg)
			case arg := <-RVExtensionArgsChannels[":FPS:"]:
				go addFPS(arg)
			case arg := <-RVExtensionArgsChannels[":EXPORT:"]:
				go export(arg)
			}
		}
	}()
}

func main() {
	fmt.Scanln()
}

func writeTimeNow(id string) {
	t := time.Now()
	// format time
	timeNow := t.Format("2006-01-02 15:04:05")
	// send data back to Arma
	a3interface.WriteArmaCallback(EXTENSION_NAME, "timeNow", id, string(timeNow))
}

func reset(arg []string) {
	// reset mission struct
	mission = Mission{}
}

func setUpMission(arg []string) {
	// Check size
	if len(arg) != 5 {
		log.Println("Error: Mission array size is not 5")
		return
	}	
	// Set up new mission
	mission = Mission{
		MissionName: arg[0],
		Worldname: arg[1],
		MissionAuthor: arg[2],
		MissionType: arg[3],
		Victory: arg[4],
		MissionStart: arg[5],
	}
}

func win(arg []string) {
	// Check size
	if len(arg) != 4 {
		log.Println("Error: Win array size is not 4")
		return
	}
	// Set win
	mission.Victory = arg[0]
	mission.MissionEnd = arg[1]
	mission.ScoreBlue = arg[2]
	mission.ScoreRed = arg[3]
}

func addPlayer(arg []string) {
	// Check size
	if len(arg) != 6 {
		log.Println("Error: Player array size is not 6")
		return
	}
	// Create new player
	player := Player{
		UID: arg[0],
		Name: arg[1],
		Role: arg[2],
		Class: arg[3],
		Side: arg[4],
		Squad: arg[5],
	}

	// Check if this player is already in the mission
	for _, p := range mission.Players {
		if p.UID == player.UID {
			// Player already in mission so just update
			p.Name = player.Name
			p.Side = player.Side
			p.Squad = player.Squad
			p.Role = player.Role
			p.Class = player.Class
			return
		}
	}

	// Add player to mission
	mission.Players = append(mission.Players, player)
}

func addKill(arg []string) {
	// Check size
	if len(arg) != 5 {
		log.Println("Error: Kill array size is not 5")
		return
	}
	// Create new kill
	kill := Kill{
		Time: arg[0],
		Victim: arg[1],
		Killer: arg[2],
		Weapon: arg[3],
		Distance: arg[4],
	}
	// Add kill to mission
	mission.Kills = append(mission.Kills, kill)
}

func addShot(arg []string) {
	// Check size
	if len(arg) != 1 {
		log.Println("Error: Shot array size is not 1")
		return
	}

	// Find the player
	for i, p := range mission.Players {
		if p.UID == arg[0] {
			// Increment shots
			mission.Players[i].Shots++
			return
		}
	}
}

func addHit(arg []string) {
	// Check size
	if len(arg) != 1 {
		log.Println("Error: Hit array size is not 1")
		return
	}

	// Find the player
	for i, p := range mission.Players {
		if p.UID == arg[0] {
			// Increment hits
			mission.Players[i].Hits++
			return
		}
	}
}

func addFPS(arg []string) {
	// For every FPS value
	for _, fps := range arg {
		// Convert to float64
		f, err := strconv.ParseFloat(fps, 64)
		if err != nil {
			log.Println("Error: FPS is not a float64")
			return
		}

		// Add FPS to mission
		mission.FPS = append(mission.FPS, f)
	}
}

func export(arg []string) {
	// Export mission
	// Get executable path
	p := modulePathDir + "\\stats\\"
	os.MkdirAll(p, os.ModePerm)

	// Generate filename
	t := time.Now().Format(time.RFC3339)
	filename := p + t + "_" + mission.MissionName + ".json"

	// Write mission to file
	data, err := json.MarshalIndent(mission, "", "    ")
	if err != nil {
		log.Println("Error: Could not marshal mission")
		return
	}

	// Write to file
	err = os.WriteFile(filename, data, os.ModePerm)

	// Check for errors
	if err != nil {
		log.Println("Error: Could not write mission to file")
		return
	}
}
