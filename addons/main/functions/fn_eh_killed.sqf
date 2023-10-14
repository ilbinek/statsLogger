params ["_victim", "_killer", "_instigator"];

[_victim, _killer] spawn {
	params ["_victim", "_killer"];
	if (_killer == _victim) then {
		private _time = diag_ticktime;
	    [_victim, {
            _this setVariable ["ace_medical_lastdamageSource", (_this getVariable "ace_medical_lastdamageSource"), 2];
        }] remoteExec ["call", _victim];
        waitUntil {
            diag_ticktime - _time > 10 || !(isnil {
                _victim getVariable "ace_medical_lastdamageSource"
            });
        };
        _killer = _victim getVariable ["ace_medical_lastdamageSource", _killer];
    } else {
         _killer
    };
    private _time = dayTime;
    private _tmpTime = [_time, "HH:MM"] call BIS_fnc_timeToString;
    private _weapon = [_killer] call statslogger_fnc_getEventWeaponText;
    private _uid = "0";
    if (isPlayer _victim) then {
        _uid = getPlayerUID _victim;
    };
    private _uidk = "0";
    if (isPlayer _killer) then {
        _uidk = getPlayerUID _killer;
    };
    if (_uidk != "0" && _uid != "0") then {
        private _tmp = [_uidk, _uid, _weapon, round(_killer distance _victim), _tmpTime];
        diag_log(text ('[STATS] ' + str(_tmp)));
        // don't care about reply
        "stats_logger" callExtension [":KILL:", _tmp];
    }
};
