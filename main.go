//go:generate go install -v github.com/kevinburke/go-bindata/go-bindata
//go:generate go-bindata -pkg assets -o assets/assets.go res/kitty.ini
//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico
package main

import (
	"io/ioutil"
	"os"

	_ "github.com/kevinburke/go-bindata"
	"github.com/portapps/kitty-portable/assets"
	. "github.com/portapps/portapps"
	"github.com/portapps/portapps/pkg/proc"
	"github.com/portapps/portapps/pkg/utl"
)

var (
	app *App
)

func init() {
	var err error

	// Init app
	if app, err = New("kitty-portable", "KiTTY"); err != nil {
		Log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "kitty.exe")

	configPath := utl.CreateFolder(app.DataPath, "config")
	iniFile := utl.PathJoin(app.DataPath, "kitty.ini")

	if !utl.Exists(iniFile) {
		Log.Info().Msg("Creating default ini file...")
		kittyIni, err := assets.Asset("res/kitty.ini")
		if err != nil {
			Log.Fatal().Err(err).Msg("Cannot load asset kitty.ini")
		}
		err = ioutil.WriteFile(iniFile, kittyIni, 0644)
		if err != nil {
			Log.Fatal().Err(err).Msg("Cannot write kitty.ini")
		}
	}

	Log.Info().Msg("Updating configuration...")
	if err := utl.ReplaceByPrefix(iniFile, "savemode=", "savemode=dir"); err != nil {
		Log.Fatal().Err(err).Msg("Cannot set savemode")
	}
	if err := utl.ReplaceByPrefix(iniFile, "#savemode=", "savemode=dir"); err != nil {
		Log.Fatal().Err(err).Msg("Cannot set savemode")
	}
	if err := utl.ReplaceByPrefix(iniFile, "configdir=", "configdir="+utl.FormatWindowsPath(configPath)); err != nil {
		Log.Fatal().Err(err).Msg("Cannot set configdir")
	}
	if err := utl.ReplaceByPrefix(iniFile, "#configdir=", "configdir="+utl.FormatWindowsPath(configPath)); err != nil {
		Log.Fatal().Err(err).Msg("Cannot set configdir")
	}

	Log.Info().Msg("Setting environment...")
	utl.OverrideEnv("KITTY_INI_FILE", utl.FormatWindowsPath(iniFile))

	configPathEmpty, _ := utl.IsDirEmpty(configPath)
	if configPathEmpty {
		Log.Info().Msg("Converting registry settings to dir mode...")
		if err := proc.QuickCmd(app.Process, []string{"-convert-dir"}); err != nil {
			Log.Error().Err(err).Msg("Cannot convert registry settings to dir mode")
		}
	}

	app.Launch(os.Args[1:])
}
