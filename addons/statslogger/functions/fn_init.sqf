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
"stats_logger" callExtension [":RESET:", []];
private _tmp = [briefingName, worldName, getMissionConfigValue ["author", ""], "public", _tmpTime];
"stats_logger" callExtension [":MISSION:", _tmp];
diag_log(text ('[STATS] Called Stats with ' + str(_tmp)));
// Start fps loop
diag_log(text ('[STATS] Starting FPS LOOP'));
call statslogger_fnc_fpsLoop;
diag_log(text ('[STATS] Started FPS LOOP'));
