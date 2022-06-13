params ["_winner", ["_scoreBlue", 0], ["_scoreRed", 0]];

_time = dayTime;
_tmpTime = [_time, "HH:MM"] call BIS_fnc_timeToString;
_tmp = ["WIN", _winner, _tmpTime, _scoreBlue, _scoreRed] joinString "::";
diag_log(text ('[STATS] ' + _tmp));
"Stats" callExtension _tmp;
