params ["_winner", ["_scoreBlue", 0], ["_scoreRed", 0]];

private _time = dayTime;
private _tmpTime = [_time, "HH:MM"] call BIS_fnc_timeToString;
private _tmp = [_winner, _tmpTime, _scoreBlue, _scoreRed];
diag_log(text ('[STATS] ' + str(_tmp)));
"stats_logger" callExtension [":WIN:", _tmp];
