// Will change, first test if it even freaking works
{
	player createDiarySubject ["StatsLoggerInfo", "StatsPlugin", "\A3\ui_f\data\igui\cfg\simpleTasks\types\whiteboard_ca.paa"];

	player createDiaryRecord [
		"StatsLoggerInfo",
		[
			"About",
			(
				"Because our lord and savior was asking for this for so long,<br/>" +
				"this thing finally exists and maybe, just maybe, works.<br/>" +
				"this is still in a 'testing' phase and if it breaks, it's not my fault.<br/>" +
				"In case you see something weird this game, please report it direcly to<br/>" +
				"me (your friend ilbinek), or Manchot and he'll relay it to me."
			)
		]
	];

	player createDiaryRecord [
		"StatsLoggerInfo",
		[
			"Status",
			(
				"Probably working<br/>" +
				"Version 0.1<br/>" +
				"Really prone to breaking."
			)
		]
	];
	
} remoteExecCall ["call", 0, true];

// Register mission handlers
call statslogger_fnc_addEventMission;

// Call the basic ocnfiguration for the extension with the mission info
_time = dayTime;
_tmpTime = [_time, "HH:MM"] call BIS_fnc_timeToString;
"Stats" callExtension "RESET";
_tmp = ["MISSION", briefingName, worldName, getMissionConfigValue ["author", ""], "public", _tmpTime] joinString "::";
"Stats" callExtension _tmp;
diag_log(text ('[STATS] Called Stats with ' + _tmp));
