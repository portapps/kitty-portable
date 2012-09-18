putty-launcher.bat
==================

A DOS Batch script to launch [PuTTY](http://www.chiark.greenend.org.uk/~sgtatham/putty/).

Your PuTTY config is saved to disk (putty.ini) instead of registry.

Tested on Windows XP, Windows Vista and Windows 7.

Requirements
------------

* Latest version of [PuTTY](http://www.chiark.greenend.org.uk/~sgtatham/putty/).
* [WSH (Windows Script Host)](http://support.microsoft.com/kb/232211) : Open a command prompt and type ``wscript`` to check.
* Access to the [Windows registry](http://support.microsoft.com/kb/256986) : Open a command prompt and type ``regedit`` to check.

Installation
------------

* Download the [latest version of PuTTY](http://the.earth.li/~sgtatham/putty/latest/x86/putty.exe).
* Put the ``putty-launcher.bat`` in the same directory as ``putty.exe``.
* Run ``putty-launcher.bat``.

Note
----

If you have already sessions saved in the registry, they will be copied automatically to the portable version.

More infos
----------

http://www.crazyws.fr/dev/applis-et-scripts/putty-launcher-garder-la-configuration-sur-disque-dans-un-fichier-UBVQA.html