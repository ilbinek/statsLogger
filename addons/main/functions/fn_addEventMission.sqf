addMissionEventHandler ["PlayerConnected", {
	_this call statslogger_fnc_eh_connected;
}];

addMissionEventHandler ["EntityKilled", {
	_this call statslogger_fnc_eh_killed;
}];
