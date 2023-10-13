diag_log(text ('[STATS] EXPORT CALLED!'));

private _results = parseSimpleArray ("stats_logger" callExtension ":EXPORT:");
if (count _results isEqualTo 0) exitWith {
	diag_log formatText ["%1", '[STATS] EXPORT FAILED!'];
};
// check for error
if (_result#0 isEqualTo ":EXPORT:") exitWith {
	diag_log formatText[
		"[STATS] EXPORT FAILED! %1",
		_results#1
	];
};

diag_log formatText [
	"[STATS] EXPORT SUCCESS! %1",
	_results#0
];




