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
        
        _strClass = typeOf _unit;
        
        _strRole = gettext(configFile >> "Cfgvehicles" >> typeOf(_unit) >> "displayname");
        if ((roleDescription _unit) != "") then {
            _nbr = (roleDescription _unit) find "@";
            if (_nbr < 0) then {
                _strRole = (roleDescription _unit);
            } else {
                _strRole = ((roleDescription _unit) select [0, _nbr]);
            };
        };
        
        _side = str (side _unit);
        _group = str (group _unit);
        _tmp = ["PLAYER", _uid, _name, _strRole, _strClass, _side, _group] joinstring "::";
        diag_log(text ('[STATS] ' + _tmp));
        "Stats" callExtension _tmp;

        _unit addEventHandler ["firedMan", {
            _this call statslogger_fnc_eh_fired;
        }];

        // TODO Will add hits, just not working for now
        /*_unit addMPEventHandler ["MPhit", {
            _this call statslogger_fnc_eh_fired;
        }];*/

        // If unit is respawned, something weird happens, so this should get around it
        _unit addEventHandler ["Respawn", {
            _unit addEventHandler ["firedMan", {
                _this call statslogger_fnc_eh_fired;
            }];
        }];
    };
};