# Introduction
This QuakeLog Parser get kill events from `quake 3` log file and agregate the infos into json separate by match, eg.:
```javascript
{
  "game_1": {
        "total_kills": 0,
        "players": [],
        "kills": {},
        "kills_by_means": {}
    },
    "game_10": {
        "total_kills": 60,
        "players": [
            "Chessus",
            "Dono da Bola",
            "Oootsimo",
            "Zeh",
            "Assasinu Credi",
            "Mal",
            "Isgalamido"
        ],
        "kills": {
            "Assasinu Credi": 3,
            "Chessus": 5,
            "Dono da Bola": 3,
            "Isgalamido": 6,
            "Mal": 1,
            "Oootsimo": -1,
            "Zeh": 7
        },
        "kills_by_means": {
            "MOD_BFG": 2,
            "MOD_BFG_SPLASH": 2,
            "MOD_CRUSH": 1,
            "MOD_MACHINEGUN": 1,
            "MOD_RAILGUN": 7,
            "MOD_ROCKET": 4,
            "MOD_ROCKET_SPLASH": 1,
            "MOD_TELEFRAG": 25,
            "MOD_TRIGGER_HURT": 17
        }
    },
}
```
# Implemantation

The way that I've decided to work on this parser was using AVL Tree for everything, to games, players kills and for kills by means. The main decision for that was because AVLs are optimal for find elements, and in concurrency contex is more easier to work on(in my experience).

# Requirements

The requirement to compile project

- golang 1.17+
- make(optional)

# Compiling

## Makefile

The makefile for this project had this commands:
- all: same as make build, compile project and put binaries into bin folder
- build: compile project and put binaries into bin folder  
- clean_builds: remove all binaries from bin folder 
- test: run all test files without caching
- test_coverage: same as test but output a coverage per function

After build will have 4 files inside bin folder:
- parse_quake_log_darwin(for macOS intel)
- parse_quake_log_m1(for macOS M1)
- parse_quake_log_linux(for Linux distributions)
- parse_quake_log_windows(for Windows)

## Without Makefile

If you want to build this without make file, you only had to know the architecture of your cpu and the SO you are running:

* MacOS M1

> `GOARCH=arm64 GOOS=darwin go build -o bin/parse_quake_log cmd/parse_quake_log.go`

* MacOS Intel

> `GOARCH=amd64 GOOS=darwin go build -o bin/parse_quake_log cmd/parse_quake_log.go`

* Linux
> `GOARCH=amd64 GOOS=linux go build -o bin/parse_quake_log cmd/parse_quake_log.go`]

* Windows
> `GOARCH=amd64 GOOS=windows go build -o bin/parse_quake_log cmd/parse_quake_log.go`

# Executing

After compile project, execute based in your system the correspondent binary:

```bin/parse_quake_log [flags] [log_filepath]```

If you type `bin/parse_quake_log -help` will show flag paramters. The current flags are:
- -out: filepath to put the result of parser(stdout is default)
> `bin/parse_quake_log -out=file.json qgames.log`
- -duration: mesure the parse duration
> `bin/parse_quake_log -duration qgames.log`
