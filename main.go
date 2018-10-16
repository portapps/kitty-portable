//go:generate go install -v github.com/kevinburke/go-bindata/go-bindata
//go:generate go-bindata -pkg assets -o assets/assets.go res/kitty.ini
//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico
package main

import (
	"io/ioutil"
	"os"

	"github.com/portapps/kitty-portable/assets"
	. "github.com/portapps/portapps"
)

func init() {
	Papp.ID = "kitty-portable"
	Papp.Name = "KiTTY"
	Init()
}

func main() {
	Papp.AppPath = AppPathJoin("app")
	Papp.DataPath = CreateFolder(AppPathJoin("data"))
	Papp.Process = PathJoin(Papp.AppPath, "kitty.exe")
	Papp.Args = nil
	Papp.WorkingDir = Papp.AppPath

	CreateFolder(PathJoin(Papp.DataPath, "config"))
	iniFile := PathJoin(Papp.DataPath, "kitty.ini")

	if !Exists(iniFile) {
		Log.Info("Creating default ini file...")
		kittyIni, err := assets.Asset("res/kitty.ini")
		if err != nil {
			Log.Error("Cannot load asset kitty.ini:", err)
		}
		err = ioutil.WriteFile(iniFile, kittyIni, 0644)
		if err != nil {
			Log.Error("Cannot write kitty.ini:", err)
		}
	}

	Log.Info("Updating configuration...")
	if err := ReplaceByPrefix(iniFile, "savemode=", "savemode=dir"); err != nil {
		Log.Error("Cannot set savemode:", err)
	}
	if err := ReplaceByPrefix(iniFile, "#savemode=", "savemode=dir"); err != nil {
		Log.Error("Cannot set savemode:", err)
	}
	if err := ReplaceByPrefix(iniFile, "configdir=", `configdir=..\data\config`); err != nil {
		Log.Error("Cannot set configdir:", err)
	}
	if err := ReplaceByPrefix(iniFile, "#configdir=", `configdir=..\data\config`); err != nil {
		Log.Error("Cannot set configdir:", err)
	}

	Log.Info("Setting environment...")
	os.Setenv("KITTY_INI_FILE", FormatWindowsPath(iniFile))

	Launch(os.Args[1:])
}
