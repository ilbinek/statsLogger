package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ilbinek/statsLogger/db"

	"github.com/rs/zerolog"

	"github.com/indig0fox/a3go/a3interface"
	"github.com/indig0fox/a3go/assemblyfinder"
)

// modulePath is the absolute path to the compiled DLL, which should be the addon folder
var modulePath string = assemblyfinder.GetModulePath()

// modulePathDir is the containing folder
var modulePathDir string = filepath.Dir(modulePath)

var EXTENSION_NAME = "STATS_LOGGER"

var mission Mission

var logFile *os.File
var logger zerolog.Logger

func init() {
	a3interface.SetVersion("1.0.0")
	a3interface.NewRegistration(":GET:TIME:").
		SetFunction(onTimeNowUTC).
		SetRunInBackground(false).
		Register()

		// when a3 sends this, set logging to INFO+
	a3interface.NewRegistration(":SET:DEBUG:OFF:").
		SetFunction(func(ctx a3interface.ArmaExtensionContext, data string) (string, error) {
			thisLogger := logger.With().
				Interface("context", ctx).
				Str("call_type", "RvExtension").
				Str("command", ":SET:DEBUG:OFF:").
				Str("data", data).
				Logger()

			thisLogger.Debug().Send()

			zerolog.SetGlobalLevel(zerolog.InfoLevel)

			return "", nil
		}).
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

	logFile, err := os.Create(filepath.Join(modulePathDir, "stats.log"))
	if err != nil {
		panic(err)
	}

	// default to debug on logger
	logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        logFile,
		TimeFormat: time.RFC3339,
		NoColor:    true,
	}).With().Timestamp().Caller().Logger().Level(zerolog.DebugLevel)
}

// onTimeNowUTC :GET:TIME: returns the current time in UTC
func onTimeNowUTC(
	ctx a3interface.ArmaExtensionContext,
	data string,
) (string, error) {
	thisLogger := logger.With().
		Interface("context", ctx).
		Str("call_type", "RvExtension").
		Str("command", ":GET:TIME:").
		Str("data", data).
		Logger()

	thisLogger.Debug().Send()

	// get time
	t := time.Now().UTC()
	// format time
	timeNow := t.Format("2006-01-02 15:04:05")
	// send data back to Arma
	thisLogger.Debug().Msgf("Returning time: %s", timeNow)
	return string(timeNow), nil
}

// onReset :RESET: resets the mission struct
func onReset(
	ctx a3interface.ArmaExtensionContext,
	data string,
) (string, error) {
	thisLogger := logger.With().
		Interface("context", ctx).
		Str("call_type", "RvExtension").
		Str("command", ":RESET:").
		Str("data", data).
		Logger()

	thisLogger.Debug().Send()
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

	thisLogger := logger.With().
		Interface("context", ctx).
		Str("call_type", "RvExtensionArgs").
		Str("command", ":MISSION:").
		Logger()

	thisLogger.Debug().Send()

	// Check size
	if len(args) != 5 {
		thisLogger.Error().Err(errors.New("mission array size is not 5")).Strs("args", args).Send()
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

	thisLogger.Info().Errs(
		"migrations", []error{
			db.Client().AutoMigrate(&Mission{}),
			db.Client().AutoMigrate(&Player{}),
			db.Client().AutoMigrate(&Kill{}),
			db.Client().AutoMigrate(&FPSRecord{}),
		},
	).Msg("Migrating tables")

	// Create new mission
	thisLogger.Debug().Interface("mission", &mission).Msg("Creating mission")

	err := db.Client().Create(&mission).Error
	if err != nil {
		thisLogger.Error().Err(err).Msg("Could not create mission")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "MISSION ERROR", err.Error())
		return "", err
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

	thisLogger := logger.With().
		Interface("context", ctx).
		Str("call_type", "RvExtensionArgs").
		Str("command", ":WIN:").
		Logger()

	thisLogger.Debug().Send()

	// Check size
	if len(args) != 4 {
		thisLogger.Error().Err(errors.New("win array size is not 4")).Strs("args", args).Send()
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

	thisLogger.Trace().Interface("mission", &mission).Msg("Updating mission")

	// Update mission
	err := db.Client().Save(&mission).Error
	if err != nil {
		thisLogger.Error().Err(err).Msg("Could not update mission")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "WIN ERROR", err.Error())
		return "", err
	}

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

	thisLogger := logger.With().
		Interface("context", ctx).
		Str("call_type", "RvExtensionArgs").
		Str("command", ":PLAYER:").
		Logger()

	thisLogger.Debug().Send()

	// Check size
	if len(args) != 6 {
		thisLogger.Error().Err(errors.New("player array size is not 6")).Strs("args", args).Send()
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "PLAYER ERROR", "WRONG PLAYER PARAMS COUNT")
		return "", errors.New("player array size is not 6")
	}

	for i, v := range args {
		args[i] = a3interface.RemoveEscapeQuotes(v)
	}

	// Create new player
	receivedPlayer := Player{
		PlayerUID: args[0],
		Name:      args[1],
		Role:      args[2],
		Class:     args[3],
		Side:      args[4],
		Squad:     args[5],
		Mission:   &mission,
	}

	var dbPlayer Player
	db.Client().Model(&Player{}).Where(
		"player_uid = ? AND mission_id = ?",
		receivedPlayer.PlayerUID,
		mission.ID,
	).First(&dbPlayer)
	if dbPlayer.PlayerUID == "" {
		err := db.Client().Create(&receivedPlayer).Error
		if err != nil {
			thisLogger.Error().Err(err).
				Interface("player", &receivedPlayer).
				Msg("Could not create player")
			a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "PLAYER ERROR", err.Error())
			return "", err
		}
	} else {
		err := db.Client().Model(&dbPlayer).Updates(Player{
			Name:  receivedPlayer.Name,
			Side:  receivedPlayer.Side,
			Squad: receivedPlayer.Squad,
			Role:  receivedPlayer.Role,
			Class: receivedPlayer.Class,
		}).Error
		if err != nil {
			thisLogger.Error().Err(err).
				Interface("player", &dbPlayer).
				Msg("Could not update player")
			a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "PLAYER ERROR", err.Error())
			return "", err
		}
	}

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

	thisLogger := logger.With().
		Interface("context", ctx).
		Str("call_type", "RvExtensionArgs").
		Str("command", ":KILL:").
		Logger()

	thisLogger.Debug().Send()

	// Check size
	if len(args) != 5 {
		thisLogger.Error().Err(errors.New("kill array size is not 5")).Strs("args", args).Send()
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

	// Add kill to mission associations
	err := db.Client().Model(&mission).Association("Kills").Append(&kill)
	if err != nil {
		thisLogger.Error().Err(err).
			Interface("kill", &kill).
			Msg("Could not add kill to mission")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "KILL ERROR", err.Error())
		return "", err
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

	thisLogger := logger.With().
		Interface("context", ctx).
		Str("call_type", "RvExtensionArgs").
		Str("command", ":SHOT:").
		Logger()

	thisLogger.Debug().Send()

	// Check size
	if len(args) != 1 {
		thisLogger.Error().Err(errors.New("shot array size is not 1")).Strs("args", args).Send()
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "SHOT ERROR", "WRONG SHOT PARAMS COUNT")
		return "", errors.New("shot array size is not 1")
	}

	for i, v := range args {
		args[i] = a3interface.RemoveEscapeQuotes(v)
	}

	// Find the player
	shooter := Player{}
	db.Client().Where(&Player{
		PlayerUID: args[0],
		MissionID: mission.ID,
	}).First(&shooter)

	if shooter.ID == 0 {
		thisLogger.Error().Err(db.Client().Error).
			Str("player_uid", args[0]).
			Msg("Could not find shooter")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "SHOT ERROR", "COULD NOT FIND SHOOTER")
		return "", errors.New("could not find shooter")
	}
	if db.Client().Error != nil {
		thisLogger.Error().Err(db.Client().Error).
			Str("player_uid", args[0]).
			Msg("DB Error")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "SHOT ERROR", "DB ERROR")
		return "", errors.New("db error")
	}

	// Increment shots
	shooter.Shots++

	// Update player
	if err := db.Client().Save(&shooter).Error; err != nil {
		thisLogger.Error().Err(err).
			Interface("player", &shooter).
			Msg("Could not update player")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "SHOT ERROR", err.Error())
		return "", err
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

	thisLogger := logger.With().
		Interface("context", ctx).
		Str("call_type", "RvExtensionArgs").
		Str("command", ":HIT:").
		Logger()

	thisLogger.Debug().Send()

	// Check size
	if len(args) != 1 {
		thisLogger.Error().Err(errors.New("hit array size is not 1")).Strs("args", args).Send()
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "HIT ERROR", "WRONG HIT PARAMS COUNT")
		return "", errors.New("hit array size is not 1")
	}

	for i, v := range args {
		args[i] = a3interface.RemoveEscapeQuotes(v)
	}

	// Find the shooter
	shooter := Player{}
	db.Client().Where(&Player{
		PlayerUID: args[0],
		MissionID: mission.ID,
	}).First(&shooter)

	if shooter.ID == 0 {
		thisLogger.Error().Err(db.Client().Error).
			Str("player_uid", args[0]).
			Msg("Could not find shooter")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "HIT ERROR", "COULD NOT FIND SHOOTER")
		return "", errors.New("could not find shooter")
	}
	if db.Client().Error != nil {
		thisLogger.Error().Err(db.Client().Error).
			Str("player_uid", args[0]).
			Msg("DB Error")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "HIT ERROR", "DB ERROR")
		return "", errors.New("db error")
	}

	// Increment hits
	shooter.Hits++

	// Update player
	if err := db.Client().Save(&shooter).Error; err != nil {
		thisLogger.Error().Err(err).
			Interface("player", &shooter).
			Msg("Could not update player")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "HIT ERROR", err.Error())
		return "", err
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

	thisLogger := logger.With().
		Interface("context", ctx).
		Str("call_type", "RvExtensionArgs").
		Str("command", ":FPS:").
		Logger()

	thisLogger.Debug().Send()

	// Check size
	if len(args) != 1 {
		thisLogger.Error().Err(errors.New("fps array size is not 1")).Strs("args", args).Send()
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "FPS ERROR", "WRONG FPS PARAMS COUNT")
		return "", errors.New("fps array size is not 1")
	}

	for i, v := range args {
		args[i] = a3interface.RemoveEscapeQuotes(v)
	}

	// Convert to float64
	f, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		thisLogger.Error().Err(err).
			Str("fps", args[0]).
			Msg("FPS is not a float64")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "FPS ERROR", "FPS IS NOT A FLOAT64")
		return "", errors.New("fps is not a float64")
	}

	// Add FPS to mission
	fps := &FPSRecord{
		FPS:     f,
		TimeUTC: time.Now().UTC(),
	}

	// Add FPS to mission associations
	err = db.Client().Model(&mission).Association("FPSRecords").Append(fps)
	if err != nil {
		thisLogger.Error().Err(err).
			Interface("fps", fps).
			Msg("Could not add FPS to mission")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "FPS ERROR", err.Error())
		return "", err
	}

	// Send data back to Arma
	return `["Saved FPS data!"]`, nil
}

// onExport :EXPORT: exports the mission to a file. run synchronously to catch errors
func onExport(
	ctx a3interface.ArmaExtensionContext,
	data string,
) (string, error) {

	thisLogger := logger.With().
		Interface("context", ctx).
		Str("call_type", "RvExtension").
		Str("command", ":EXPORT:").
		Str("data", data).
		Logger()

	thisLogger.Debug().Send()

	// Export mission
	// Get executablepath/stats
	p := filepath.Join(modulePathDir, "stats-output")
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		thisLogger.Error().Err(err).
			Str("path", p).
			Msg("Could not create stats-output folder")
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
	// catch error
	if err != nil {
		thisLogger.Error().Err(err).
			Interface("mission", &mission).
			Msg("Could not marshal mission to JSON")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "EXPORT ERROR", err.Error())
		return "", err
	}

	// Write to file
	err = os.WriteFile(fileDestinationAbsolute, missionBytes, os.ModePerm)

	// Check for errors
	if err != nil {
		thisLogger.Error().Err(err).
			Str("path", fileDestinationAbsolute).
			Msg("Could not write mission to file")
		a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "EXPORT ERROR", err.Error())
		return "", err
	}

	thisLogger.Info().Str("path", fileDestinationAbsolute).Msg("Exported mission to file")
	a3interface.WriteArmaCallback(EXTENSION_NAME, "DEBUG", "EXPORT DONE", "EXPORT FINISHED")
	return fmt.Sprintf(
			`["Successfully exported data to %s"]`,
			fileDestinationAbsolute,
		),
		nil
}

func main() {
	db.Client().AutoMigrate(&Mission{})
	db.Client().AutoMigrate(&Player{})
	db.Client().AutoMigrate(&Kill{})
	db.Client().AutoMigrate(&FPSRecord{})
	fmt.Scanln()
}
