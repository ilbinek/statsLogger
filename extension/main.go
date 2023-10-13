package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/indig0fox/a3go/a3interface"
	"github.com/indig0fox/a3go/assemblyfinder"
)

// modulePath is the absolute path to the compiled DLL, which should be the addon folder
var modulePath string = assemblyfinder.GetModulePath()

// modulePathDir is the containing folder
var modulePathDir string = filepath.Dir(modulePath)

var EXTENSION_NAME = "STATS_LOGGER"

var mission Mission

func init() {
	a3interface.SetVersion("1.0.0")
	a3interface.NewRegistration(":GET:TIME:").
		SetFunction(onTimeNowUTC).
		SetRunInBackground(false).
		Register()

	a3interface.NewRegistration(":RESET:").
		SetFunction(onReset).
		SetRunInBackground(false).
		Register()

	a3interface.NewRegistration(":MISSION:").
		SetArgsFunction(onSetupMissionArgs).
		SetRunInBackground(false).
		Register()

	a3interface.NewRegistration(":WIN:").
		SetArgsFunction(onWinArgs).
		SetRunInBackground(false).
		Register()

	a3interface.NewRegistration(":PLAYER:").
		SetArgsFunction(onAddPlayerArgs).
		SetRunInBackground(false).
		Register()

	a3interface.NewRegistration(":SHOT:").
		SetArgsFunction(onAddShotArgs).
		// may have multiple at once - run async
		SetRunInBackground(true).
		SetDefaultResponse(`["Saving shot data..."]`).
		Register()

	a3interface.NewRegistration(":HIT:").
		SetArgsFunction(onAddHitArgs).
		// may have multiple at once - run async
		SetRunInBackground(true).
		SetDefaultResponse(`["Saving hit data..."]`).
		Register()

	a3interface.NewRegistration(":KILL:").
		SetArgsFunction(onAddKillArgs).
		SetRunInBackground(false).
		Register()

	a3interface.NewRegistration(":FPS:").
		SetArgsFunction(onAddFPSArgs).
		SetRunInBackground(false).
		Register()

	a3interface.NewRegistration(":EXPORT:").
		SetFunction(onExport).
		SetRunInBackground(false).
		Register()
}

// onTimeNowUTC :GET:TIME: returns the current time in UTC
func onTimeNowUTC(
	ctx a3interface.ArmaExtensionContext,
	data string,
) (string, error) {
	t := time.Now().UTC()
	// format time
	timeNow := t.Format("2006-01-02 15:04:05")
	// send data back to Arma
	return string(timeNow), nil
}

// onReset :RESET: resets the mission struct
func onReset(
	ctx a3interface.ArmaExtensionContext,
	data string,
) (string, error) {
	// reset mission struct
	mission = Mission{}
	// send data back to Arma
	return "", nil
}

// onSetupMissionArgs :MISSION: defines info about the mission during init
func onSetupMissionArgs(
	ctx a3interface.ArmaExtensionContext,
	command string,
	args []string,
) (string, error) {
	// args are:
	// 0: mission name (briefingName)
	// 1: world name
	// 2: mission author
	// 3: mission type (ex. public)
	// 4: time of day at mission start (daytime)

	// Check size
	if len(args) != 5 {
		log.Println("Error: Mission array size is not 5")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "MISSION ERROR", "WRONG MISSION PARAMS COUNT - ["+strings.Join(args, ", ")+"]")
		return "", errors.New("mission array size is not 5")
	}

	for i, v := range args {
		args[i] = a3interface.RemoveEscapeQuotes(v)
	}

	// Set up new mission
	mission = Mission{
		MissionName:   args[0],
		Worldname:     args[1],
		MissionAuthor: args[2],
		MissionType:   args[3],
		// Victory:       args[4], // not in sqf
		MissionStart: args[4],
	}
	// send data back to Arma
	return `["Saved mission data!"]`, nil
}

// onWinArgs :WIN: fills in victory data for the mission
func onWinArgs(
	ctx a3interface.ArmaExtensionContext,
	command string,
	args []string,
) (string, error) {
	// args are:
	// 0: winner
	// 1: end time
	// 2: blue score
	// 3: red score

	// Check size
	if len(args) != 4 {
		log.Println("Error: Win array size is not 4")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "WIN ERROR", "WRONG WIN PARAMS COUNT")
		return "", errors.New("win array size is not 4")
	}

	for i, v := range args {
		args[i] = a3interface.RemoveEscapeQuotes(v)
	}

	// Set win
	mission.Victory = args[0]
	mission.MissionEnd = args[1]
	mission.ScoreBlue = args[2]
	mission.ScoreRed = args[3]

	// send data back to Arma
	return `["Saved win data!"]`, nil
}

// onAddPlayerArgs :PLAYER: adds a player to the mission
func onAddPlayerArgs(
	ctx a3interface.ArmaExtensionContext,
	command string,
	args []string,
) (string, error) {
	// args are:
	// 0: playerUID
	// 1: playerName
	// 2: roleDescription (prefix only)
	// 3: class (typeOf)
	// 4: side (WEST, EAST, etc.)
	// 5: groupID (str group _x)

	// Check size
	if len(args) != 6 {
		log.Println("Error: Player array size is not 6")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "PLAYER ERROR", "WRONG PLAYER PARAMS COUNT")
		return "", errors.New("player array size is not 6")
	}

	for i, v := range args {
		args[i] = a3interface.RemoveEscapeQuotes(v)
	}

	// Create new player
	player := Player{
		UID:   args[0],
		Name:  args[1],
		Role:  args[2],
		Class: args[3],
		Side:  args[4],
		Squad: args[5],
	}

	// Check if this player is already in the mission
	for _, p := range mission.Players {
		if p.UID == player.UID {
			// Player already in mission so just update
			p.Update(PlayerUpdateOptions{
				Name:  player.Name,
				Side:  player.Side,
				Squad: player.Squad,
				Role:  player.Role,
				Class: player.Class,
			})
			return `["Updated player data!"]`, nil
		}
	}

	// Add player to mission
	mission.AddPlayer(&player)

	// Send data back to Arma
	return `["Saved player data!"]`, nil
}

// onAddKillArgs :KILL: adds a kill to the mission
func onAddKillArgs(
	ctx a3interface.ArmaExtensionContext,
	command string,
	args []string,
) (string, error) {
	// args are:
	// 0: killerUID
	// 1: victimUID
	// 2: weapon
	// 3: distance
	// 4: time

	// Check size
	if len(args) != 5 {
		log.Println("Error: Kill array size is not 5")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "KILL ERROR", "WRONG KILL PARAMS COUNT")
		return "", errors.New("kill array size is not 5")
	}

	for i, v := range args {
		args[i] = a3interface.RemoveEscapeQuotes(v)
	}

	// Create new kill
	kill := Kill{
		Killer:   args[0],
		Victim:   args[1],
		Weapon:   args[2],
		Distance: args[3],
		Time:     args[4],
	}
	// Add kill to mission
	mission.AddKill(&kill)

	// Find the killer
	for i, p := range mission.Players {
		if p.UID == args[0] {
			// Increment kills
			mission.Players[i].AddKill()
			return `["Saved kill data!"]`, nil
		}
	}

	// Send data back to Arma
	return `["Saved kill data!"]`, nil
}

// addShot :SHOT: adds a shot to the mission
func onAddShotArgs(
	ctx a3interface.ArmaExtensionContext,
	command string,
	args []string,
) (string, error) {
	// args are:
	// 0: playerUID

	// Check size
	if len(args) != 1 {
		log.Println("Error: Shot array size is not 1")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "SHOT ERROR", "WRONG SHOT PARAMS COUNT")
		return "", errors.New("shot array size is not 1")
	}

	for i, v := range args {
		args[i] = a3interface.RemoveEscapeQuotes(v)
	}

	// Find the player
	for i, p := range mission.Players {
		if p.UID == args[0] {
			// Increment shots
			// do so using a method which uses the mutex
			mission.Players[i].AddShot()
			return `["Saved shot data!"]`, nil
		}
	}

	// Send data back to Arma
	return `["Saved shot data!"]`, nil
}

// addHit :HIT: adds a hit to the mission
func onAddHitArgs(
	ctx a3interface.ArmaExtensionContext,
	command string,
	args []string,
) (string, error) {
	// args are:
	// 0: playerUID

	// Check size
	if len(args) != 1 {
		log.Println("Error: Hit array size is not 1")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "HIT ERROR", "WRONG HIT PARAMS COUNT")
		return "", errors.New("hit array size is not 1")
	}

	for i, v := range args {
		args[i] = a3interface.RemoveEscapeQuotes(v)
	}

	// Find the player
	for i, p := range mission.Players {
		if p.UID == args[0] {
			// Increment hits
			mission.Players[i].Hits++
			return `["Saved hit data!"]`, nil
		}
	}

	// Send data back to Arma
	return `["Saved hit data!"]`, nil
}

// onAddFPSArgs :FPS: adds FPS data to the mission
func onAddFPSArgs(
	ctx a3interface.ArmaExtensionContext,
	command string,
	args []string,
) (string, error) {
	// args are:
	// 0: fps

	// Check size
	if len(args) != 1 {
		log.Println("Error: FPS array size is not 1")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "FPS ERROR", "WRONG FPS PARAMS COUNT")
		return "", errors.New("fps array size is not 1")
	}

	for i, v := range args {
		args[i] = a3interface.RemoveEscapeQuotes(v)
	}

	// Convert to float64
	f, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		log.Println("Error: FPS is not a float64")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "FPS ERROR", "FPS IS NOT A FLOAT64")
		return "", errors.New("fps is not a float64")
	}

	// Add FPS to mission
	mission.AddFPS(&FPSRecord{
		FPS:     f,
		TimeUTC: time.Now().UTC(),
	})

	// Send data back to Arma
	return `["Saved FPS data!"]`, nil
}

// onExport :EXPORT: exports the mission to a file. run synchronously to catch errors
func onExport(
	ctx a3interface.ArmaExtensionContext,
	data string,
) (string, error) {

	// Export mission
	// Get executablepath/stats
	p := filepath.Join(modulePathDir, "stats-output")
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		log.Println("Error: Could not create stats folder")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "EXPORT ERROR", err.Error())
		return "", err
	}

	// Generate filename
	// year, month, day := time.Now().Date()
	// hours, minutes, _ := time.Now().Clock()
	// filename := fmt.Sprintf("%s%d-%d-%d-%d-%d_%s.json", p, year, int(month), day, hours, minutes, mission.MissionName)
	filename := fmt.Sprintf(
		"%s_%s.json",
		time.Now().Format("2006-01-02_15-04"),
		mission.MissionName,
	)
	fileDestinationAbsolute := filepath.Join(p, filename)

	// Write mission to file

	// First, marshal to indented JSON for easy viewing
	missionBytes, err := json.MarshalIndent(&mission, "", "    ")
	// If not viewer friendly, don't use indentation
	missionBytes, err = json.Marshal(&mission)
	// catch error
	if err != nil {
		log.Println("Error: Could not marshal mission")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "EXPORT ERROR", err.Error())
		return "", err
	}

	// Write to file
	err = os.WriteFile(fileDestinationAbsolute, missionBytes, os.ModePerm)

	// Check for errors
	if err != nil {
		log.Println("Error: Could not write mission to file")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "EXPORT ERROR", err.Error())
		return "", err
	}

	a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "EXPORT DONE", "EXPORT FINISHED")
	return fmt.Sprintf(
			`["Successfully exported data to %s"]`,
			fileDestinationAbsolute,
		),
		nil
}

func main() {
	fmt.Scanln()
}
