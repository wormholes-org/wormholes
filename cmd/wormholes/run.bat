@echo off
set rootPath=%~dp0
set nodePath=%rootPath%.wormholes
if exist %nodePath% (
	rd /s/q %nodePath%
)

echo %1

if "%1" == "" (
	goto a
) else (
    md %nodePath%\geth
    echo %1 > %nodePath%\geth\nodekey
)

:a
wormholes.exe --devnet --datadir %nodePath% --mine --syncmode=full