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
    _time = dayTime;
    _tmpTime = [_time, "HH:MM"] call BIS_fnc_timeToString;
    _weapon = [_killer] call statslogger_fnc_getEventWeaponText;
    _uid = "0";
    if (isPlayer _victim) then {
        _uid = getPlayerUID _victim;
    };
    _uidk = "0";
    if (isPlayer _killer) then {
        _uidk = getPlayerUID _killer;
    };
    if (_uidk != "0" && _uid != "0") then {
        _tmp = ["KILL", _uidk, _uid, _weapon, round(_killer distance _victim), _tmpTime] joinString "::";
        diag_log(text ('[STATS] ' + _tmp));
        "Stats" callExtension _tmp;
    }
};
