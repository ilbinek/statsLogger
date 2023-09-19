/*
** indigo, thank you so much for your help, this was... Crazy...
** What you do for the community is awesome! Please keep it up.
*/

// exclude server connected event
if ((_this#4) != 2) then {
    _this spawn {
        params ["_id", "_uid", "_name", "_jip", "_owner", "_idstr"];
        
        private "_unit";
        
        // wait until the unit associated with the player's netId is not null
        waitUntil {
            uiSleep 1;
            _unit = getUserinfo _idstr select 10;
            !isNull _unit;
        };
        
        private _strClass = typeOf _unit;
        
        private _strRole = gettext(configFile >> "Cfgvehicles" >> typeOf(_unit) >> "displayname");
        if ((roleDescription _unit) != "") then {
            _nbr = (roleDescription _unit) find "@";
            if (_nbr < 0) then {
                _strRole = (roleDescription _unit);
            } else {
                _strRole = ((roleDescription _unit) select [0, _nbr]);
            };
        };
        
        private _side = str (side _unit);
        private _group = str (group _unit);
        private _tmp = [_uid, _name, _strRole, _strClass, _side, _group];
        diag_log(text ('[STATS] ' + str(_tmp)));
        "stats_logger" callExtension [":PLAYER:", _tmp];

        _unit addEventHandler ["firedMan", {
            _this call statslogger_fnc_eh_fired;
        }];
    };
};
