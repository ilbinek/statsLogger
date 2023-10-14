# Arma 3 Stats Plugin

This plugin is still in **early development** and will change going forward.<br>
This plugin is mainly intended to be used in one-life PvP games. Currently in use mainly by  [TBD](https://tbdevent.eu) and  [OFCRA](https://ofcrav2.org/forum/index.php).

## Current features

- tracks all players, their role, side and squad
- tracks all kills, distance, weapon used (filters out bots)
- tracks number of shots fired per each player
- tracks number of hits for each player (if a bullet hits multiple parts of the body dealing multiple damages, it's counted as multiple hits)
- ability to add winning side and the points

## Planned features

- config file so output folder can be easily configured
- upload to a remote server
- automatic database ingest
- basic webserver to view statistics for each mission
- overall improvements

# Usage

**Currently, the default output folder is `$ArmaServerExecutable/stats/`**<br>

This is a server side mod, so clients are not supposed to have it.

- Download the current released version
- Add the folder into your mods folder
- Add `@statslogger` into you -servermod= in start script
- Play a mission
- If you want to add a winning side, execute:

  ```sqf
  [
    WINNING SIDE <string | side>, 
    BLUEFOR POINTS <number>,
    REDFOR POINTS <number>
  ] remoteExec "statslogger_fnc_mission_end", 2];
  ```

- Execute `call statslogger_fnc_export;` on the server before mission end (either from debug console, or add it into your framework to be called automatically)

# Currently known bugs

- If a player respawns, some things break (not high priority, mainly used in 1 life PvP games right now)

# Extension compilation

## Using docker

### Windows compilation

```bash
docker pull x1unix/go-mingw:1.20

# Compile x64 Windows DLL
docker run --rm -it -v ${PWD}:/go/work -w /go/work x1unix/go-mingw:1.20 go build -o stats_logger_x64.dll -buildmode=c-shared .

# Compile x86 Windows DLL
docker run --rm -it -v ${PWD}:/go/work -w /go/work -e GOARCH=386 x1unix/go-mingw:1.20 go build -o stats_logger.dll -buildmode=c-shared .

# Compile x64 Windows EXE
docker run --rm -it -v ${PWD}:/go/work -w /go/work x1unix/go-mingw:1.20 go build -o stats_logger_x64.dll .
```

### COMPILING FOR LINUX

```bash
docker build -t indifox926/build-a3go:linux-so -f ./build/Dockerfile.build ./cmd

# Compile x64 Linux .so
docker run --rm -it -v ${PWD}:/app -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=1 -e CC=gcc indifox926/build-a3go:linux-so go build -o stats_logger_x64.so -linkshared .

# Compile x86 Linux .so
docker run --rm -it -v ${PWD}:/app -e GOOS=linux -e GOARCH=386 -e CGO_ENABLED=1 -e CC=gcc indifox926/build-a3go:linux-so go build -o stats_logger.so -linkshared .
```

## Direct Windows/Linux build with proper toolchain

- if you have the proper toolchain that can cmpile CGO
  - On windows set `-buildmode=c-shared`
  - On Linux set `-linkshared`

# Editing your framework

- If you want to add a call to this mod into your framework, you can use `if (isClass(configFile >> "CfgPatches" >> "STATSLOGGER")) then {// your code};` to make sure this plugin is loaded.

# Special thanks

Special thanks go to [Indigo](https://github.com/indig0fox), for his help in the creation of this plugin. Also for the creation of [a3interface](https://github.com/indig0fox/a3go) that is what is powering the new Go extension!
