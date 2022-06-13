# Arma 3 Stats Plugin
This plugin is still in **early development** and will change going forward.<br>
This plugin is mainly intended to be used in [OFCRA](https://ofcrav2.org/forum/index.php) PvP games. But is simple to use and can be used by anyone.

# **CURENTLY TESTED AND DEVELOPED ONLY ON LINUX**

## Current features
- tracks all players, their role, side and squad
- tracks all kills, distance, weapon used (filters out bots)
- tracks number of shots fired per each user
- ability to add winning side and the points

## Planned features
- add number of hits for every player (accuracy can be calculated afterwards)
- simple windows compilation
- config file so output folder can be easily configured

# Usage
**Curently, the default output folder is `/stats-output/`, user running the server needs permission to write into this folder.**<br>
This folder can be edited in `extension/main.cpp` on line 150 and needs to be recompiled afterwards - **THIS WILL CHANGE**.<br>
This is a server sided mod, so clients are not supposed to have it.
- Download the current released version
- Add the folder into your mods folder
- Add `@statslogger` into you -mod in start script
- Play a mission
- If you want to add a winning side, execute `["WINNING SIDE", "BLUEFOR POINTS", "REDFOR POINTS"] remoteExec "statslogger_fnc_mission_end", 2];`
- Execute `call statslogger_fnc_export;` on the server before mission end (either from debug console, or add it into your framework to be called automatically)

# Currently known bugs
- If a player respawns, some things break (not high priority, mainly used in 1 life PvP games right now)

# Extension compilation
- In case you want to edit the path in the excention and compile it yourself, go into `extension/` and if on Linux, use `make` (x64 version only, you shouldn't use 32 bit anymore anyway)
- If you are on Windows, or want 32 bit version, compilation is on you (for now)

# Editing your framework
- If you want to add a call to this mod into your framework, you can use `if (isClass(configFile >> "CfgPatches" >> "STATSLOGGER")) then {// your code};` to make sure this plugin is loaded.

# Special thanks
Special thanks go to [Indigo](https://github.com/indig0fox), for his help in the creation of this plugin.
