/*
** indigo, thank you so much for your help, this was... Crazy...
** What you do for the community is awesome! Please keep it up.
** <3
*/

// exclude server connected event if not running local listen server
private _isLocalListenServer = (isServer && hasInterface);
if (
    // so if not local listen server, eval the owner id and skip 2 (server)
    // otherwise process all (for debugging)
    [
        (_this#4) != 2,
        true
    ] select _isLocalListenServer
) then {
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

        // args format returns array - the extension data is the first element
        private _response = "stats_logger" callExtension [":PLAYER:", _tmp];
        private _realResponse = parseSimpleArray (_response#0);
        if (count _realResponse isEqualTo 0) exitWith {
            diag_log ("[STATS] Error: Bad response from stats logger");
        };
        // if error, first element is command name, second is error message
        if (_realResponse#0 isEqualTo ":PLAYER:") then {
            // Error: Player array size is not 6
            diag_log ("[STATS] " + _realResponse#1);
        } else {
            // logs "Saved player data!"
            diag_log ("[STATS] " + _realResponse#0);
        };

        _unit addEventHandler ["firedMan", {
            _this call statslogger_fnc_eh_fired;
        }];

         _unit addMPEventHandler ["MPHit", { 
            _this call statslogger_fnc_eh_hit; 
        }]; 
    };
};
