@echo off
set rootPath=%~dp0
set nodePath=%rootPath%.wormholes
if exist %nodePath% (
	rd /s/q %nodePath%
)

if "%1" == "" (
     echo "Please pass in the private key of the account to be pledged."
	 exit -1
 ) else (
     md %nodePath%\geth
     echo %1 > %nodePath%\geth\nodekey
     .\wormholes.exe --devnet --mine --syncmode=full
 )