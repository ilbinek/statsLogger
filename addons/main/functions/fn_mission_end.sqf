params ["_winner", ["_scoreBlue", 0], ["_scoreRed", 0]];

// get string from side or number type in winner param
if (_winner isEqualType west || _winner isEqualType 5) then {
	_winner = _winner call BIS_fnc_sideNameUnlocalized;
};

private _time = dayTime;
private _tmpTime = [_time, "HH:MM"] call BIS_fnc_timeToString;
private _tmp = [_winner, _tmpTime, _scoreBlue, _scoreRed];
diag_log(text ('[STATS] ' + str(_tmp)));

// args format returns array - the extension data is the first element
private _response = "stats_logger" callExtension [":WIN:", _tmp];
private _realResponse = parseSimpleArray (_response#0);
if (count _realResponse isEqualTo 0) exitWith {
	diag_log ("[STATS] Error: Bad response from stats logger");
};
// if error, first element is command name, second is error message
if (_realResponse#0 isEqualTo ":WIN:") then {
	// Error: Bad response from stats logger
	diag_log ("[STATS] " + _realResponse#1);
} else {
	// logs "Saved win data!"
	diag_log ("[STATS] " + _realResponse#0);
};

