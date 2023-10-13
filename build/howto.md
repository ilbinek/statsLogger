# Building

*From <https://github.com/indig0fox/a3go>*

## COMPILING FOR WINDOWS

These compile commands should be run from the project root.

Set version in shell:

```powershell
# version is semantic + build date 
# e.g. 1.0.0-20210530
$version = "0.3.0-20231012"
```

```bash
export version="0.3.0-20231012"
```

```ps1
docker pull x1unix/go-mingw:1.20

# Compile x64 Windows DLL
docker run --rm -it -v ${PWD}\extension\:/go/work -w /go/work -e GOARCH=amd64 -e CGO_ENABLED=1 x1unix/go-mingw:1.20 go build -o ./dist/stats_logger_x64.dll -buildmode=c-shared -ldflags "-w -s -X main.EXTENSION_VERSION=$version" .

rm ./stats_logger_x64.dll
mv ./extension/dist/stats_logger_x64.dll ./stats_logger_x64.dll

# Compile x86 Windows DLL
docker run --rm -it -v ${PWD}\extension:/go/work -w /go/work -e GOARCH=386 -e CGO_ENABLED=1 x1unix/go-mingw:1.20 go build -o ./dist/stats_logger.dll -buildmode=c-shared -ldflags "-w -s -X main.EXTENSION_VERSION=$version" .

rm ./stats_logger.dll
mv ./extension/dist/stats_logger.dll ./stats_logger.dll

# Compile x64 Windows EXE
docker run --rm -it -v ${PWD}:/go/work -w /go/work -e GOARCH=amd64 -e CGO_ENABLED=1 x1unix/go-mingw:1.20 go build -o ./dist/stats_logger_x64.exe -ldflags "-w -s -X main.EXTENSION_VERSION=$version" .
```

## COMPILING FOR LINUX

```ps1
docker build -t indifox926/build-a3go:linux-so -f ./build/Dockerfile.build .

# Compile x64 Linux .so
docker run --rm -it -v ${PWD}\extension:/app -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=1 indifox926/build-a3go:linux-so go build -o ./dist/stats_logger_x64.so -linkshared -ldflags "-w -s -X main.EXTENSION_VERSION=$version" .

rm ./stats_logger_x64.so
mv ./extension/dist/stats_logger_x64.so ./stats_logger_x64.so

# Compile x86 Linux .so
docker run --rm -it -v ${PWD}\extension:/app -e GOOS=linux -e GOARCH=386 -e CGO_ENABLED=1 indifox926/build-a3go:linux-so go build -o ./dist/stats_logger.so -linkshared -ldflags "-w -s -X main.EXTENSION_VERSION=$version" .

rm ./stats_logger.so
mv ./extension/dist/stats_logger.so ./stats_logger.so
```

## Compile Addon

First, move the compiled dlls from `extension/dist` to the project root. Or use the provided commands.

To prepare the addon, you'll need to download the [HEMTT](https://brettmayson.github.io/HEMTT/commands/build.html) binary, place it in the project root, and run the following command:

```bash
./HEMTT.exe release
```

The PBOs and relevant files will be placed in the ./.hemmttout/build directory.

---
