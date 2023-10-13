// Will change, first test if it even freaking works
{
	player createDiarySubject ["StatsLoggerInfo", "StatsPlugin", "\A3\ui_f\data\igui\cfg\simpleTasks\types\whiteboard_ca.paa"];

	player createDiaryRecord [
		"StatsLoggerInfo",
		[
			"About",
			(
				"Simple plugin to capture all important stats about a PvP mission.<br/>" +
				"This plugin captures all players, kills, how many times they fired and server FPS.<br/>" +
				"<br/>" +
				"This plugin is open-source and available for everyone to use.<br/>" +
				"All the necassary information can be found on https://github.com/ilbinek/statsLogger"
			)
		]
	];

	player createDiaryRecord [
		"StatsLoggerInfo",
		[
			"Status",
			(
				"Stats Addon initialised<br/>" +
				"Version 0.3<br/>" +
				"Capture is running."
			)
		]
	];
	
} remoteExecCall ["call", 0, true];

// Register mission handlers
call statslogger_fnc_addEventMission;

// Call the basic ocnfiguration for the extension with the mission info
private _time = dayTime;
private _tmpTime = [_time, "HH:MM"] call BIS_fnc_timeToString;
// private _utcTime = "stats_logger" callExtension ":GET:TIME:"; // 2015-01-01 00:00:00
"stats_logger" callExtension ":RESET:";
private _tmp = [briefingName, worldName, getMissionConfigValue ["author", ""], "public", _tmpTime];

private _response = "stats_logger" callExtension [":MISSION:", _tmp];
private _realResponse = parseSimpleArray (_response#0);
if (count _realResponse isEqualTo 0) exitWith {
	diag_log ("[STATS] Error: Bad response from stats logger");
};
// if error, first element is command name, second is error message
if (_realResponse#0 isEqualTo ":MISSION:") then {
	diag_log formatText[
		"[STATS] %1 ",
		_realResponse#1
	];
} else {
	// logs "Saved mission data!"
	diag_log formatText[
		"[STATS] %1",
		_realResponse#0
	];
};


diag_log(text ('[STATS] Called Stats with ' + str(_tmp)));
// Start fps loop
diag_log(text ('[STATS] Starting FPS LOOP'));
call statslogger_fnc_fpsLoop;
diag_log(text ('[STATS] Started FPS LOOP'));
