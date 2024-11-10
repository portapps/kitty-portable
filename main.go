//go:generate go install -v github.com/kevinburke/go-bindata/v4/go-bindata
//go:generate go-bindata -pkg assets -o assets/assets.go res/kitty.ini
//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"os"

	"github.com/portapps/kitty-portable/assets"
	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/proc"
	"github.com/portapps/portapps/v3/pkg/utl"
)

var (
	app *portapps.App
)

func init() {
	var err error

	// Init app
	if app, err = portapps.New("kitty-portable", "KiTTY"); err != nil {
		log.Debug().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "kitty.exe")

	configPath := utl.CreateFolder(app.DataPath, "config")
	iniFile := utl.PathJoin(app.DataPath, "kitty.ini")

	if !utl.Exists(iniFile) {
		log.Info().Msg("Creating default ini file...")
		kittyIni, err := assets.Asset("res/kitty.ini")
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot load asset kitty.ini")
		}
		err = os.WriteFile(iniFile, kittyIni, 0644)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot write kitty.ini")
		}
	}

	log.Info().Msg("Updating configuration...")
	if err := utl.ReplaceByPrefix(iniFile, "savemode=", "savemode=dir"); err != nil {
		log.Fatal().Err(err).Msg("Cannot set savemode")
	}
	if err := utl.ReplaceByPrefix(iniFile, ";savemode=", "savemode=dir"); err != nil {
		log.Fatal().Err(err).Msg("Cannot set savemode")
	}
	if err := utl.ReplaceByPrefix(iniFile, "configdir=", "configdir="+utl.FormatWindowsPath(configPath)); err != nil {
		log.Fatal().Err(err).Msg("Cannot set configdir")
	}
	if err := utl.ReplaceByPrefix(iniFile, ";configdir=", "configdir="+utl.FormatWindowsPath(configPath)); err != nil {
		log.Fatal().Err(err).Msg("Cannot set configdir")
	}

	log.Info().Msg("Setting environment...")
	os.Setenv("KITTY_INI_FILE", utl.FormatWindowsPath(iniFile))

	configPathEmpty, _ := utl.IsDirEmpty(configPath)
	if configPathEmpty {
		log.Info().Msg("Converting registry settings to dir mode...")
		if err := proc.QuickCmd(app.Process, []string{"-convert-dir"}); err != nil {
			log.Error().Err(err).Msg("Cannot convert registry settings to dir mode")
		}
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
