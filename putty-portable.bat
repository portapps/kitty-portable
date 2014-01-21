@ECHO OFF

::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
::                                                                                ::
::  Putty Portable                                                                ::
::                                                                                ::
::  A DOS Batch script to make Putty portable on Windows systems.                 ::
::                                                                                ::
::  Copyright (C) 2013-2014 Cr@zy <webmaster@crazyws.fr>                          ::
::                                                                                ::
::  Putty Portable is free software; you can redistribute it and/or modify        ::
::  it under the terms of the GNU Lesser General Public License as published by   ::
::  the Free Software Foundation, either version 3 of the License, or             ::
::  (at your option) any later version.                                           ::
::                                                                                ::
::  Putty Portable is distributed in the hope that it will be useful,             ::
::  but WITHOUT ANY WARRANTY; without even the implied warranty of                ::
::  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the                  ::
::  GNU Lesser General Public License for more details.                           ::
::                                                                                ::
::  You should have received a copy of the GNU Lesser General Public License      ::
::  along with this program. If not, see http://www.gnu.org/licenses/.            ::
::                                                                                ::
::  Usage: putty-portable.bat                                                     ::
::                                                                                ::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

:: Config files
SET ttyInvis=%TEMP%\puttyinvis_%RANDOM%.vbs
SET ttyRnd=%TEMP%\puttyrnd_%RANDOM%.tmp
SET ttyDel=%TEMP%\puttydel_%RANDOM%.tmp
SET ttyReg=%TEMP%\puttyreg_%RANDOM%.tmp

IF "%1"=="" GOTO LAUNCH
IF "%1"=="process" GOTO PROCESS

:: Run the batch file in the background
:LAUNCH
ECHO set args = WScript.Arguments >%ttyInvis%
ECHO num = args.Count >>%ttyInvis%
ECHO. >>%ttyInvis%
ECHO if num = 0 then >>%ttyInvis%
ECHO    WScript.Quit 1 >>%ttyInvis%
ECHO end if >>%ttyInvis%
ECHO. >>%ttyInvis%
ECHO sargs = "" >>%ttyInvis%
ECHO if num ^> 1 then >>%ttyInvis%
ECHO    sargs = " " >>%ttyInvis%
ECHO    for k = 1 to num - 1 >>%ttyInvis%
ECHO        anArg = args.Item(k) >>%ttyInvis%
ECHO        sargs = sargs ^& anArg ^& " " >>%ttyInvis%
ECHO    next >>%ttyInvis%
ECHO end if >>%ttyInvis%
ECHO. >>%ttyInvis%
ECHO Set WshShell = WScript.CreateObject("WScript.Shell") >>%ttyInvis%
ECHO. >>%ttyInvis%
ECHO WshShell.Run """" ^& WScript.Arguments(0) ^& """" ^& sargs, 0, False >>%ttyInvis%

wscript.exe %ttyInvis% %~n0.bat process
GOTO DONE

:: Write config to disk (putty.ini)
:PROCESS
ECHO REGEDIT4>%ttyRnd%
ECHO [HKEY_CURRENT_USER\Software\SimonTatham\PuTTY]>>%ttyRnd%
ECHO "RandSeedFile"="%TEMP:\=\\%\\putty.rnd">>%ttyRnd%
regedit /s %ttyRnd%
DEL %ttyRnd%
SET ttyRnd=

regedit /s putty.ini
start /w putty.exe

DEL %TEMP%\putty.rnd
regedit /ea %ttyReg% HKEY_CURRENT_USER\Software\SimonTatham\PuTTY
fc putty.ini %ttyReg% | find "FC: no dif" > NUL
IF ERRORLEVEL 1 COPY %ttyReg% putty.ini
DEL %ttyReg%
SET ttyReg=

ECHO REGEDIT4>%ttyDel%
ECHO.>>%ttyDel%
ECHO [-HKEY_CURRENT_USER\Software\SimonTatham\PuTTY]>>%ttyDel%
ECHO.>>%ttyDel%
type %ttyDel%
regedit /s %ttyDel%
DEL %ttyDel%
DEL %ttyInvis%
SET ttyDel=
GOTO DONE

:DONE