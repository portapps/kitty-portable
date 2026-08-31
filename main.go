//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
package main

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/files"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/proc"
)

//go:embed res/kitty.ini
var defaultKittyIni []byte

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
	if err := os.MkdirAll(app.DataPath, 0o755); err != nil {
		log.Fatal().Err(err).Msg("Cannot create data path")
	}
	app.Process = filepath.Join(app.AppPath, "kitty.exe")

	configPath := filepath.Join(app.DataPath, "config")
	if err := os.MkdirAll(configPath, 0o755); err != nil {
		log.Fatal().Err(err).Msg("Cannot create config path")
	}
	iniFile := filepath.Join(app.DataPath, "kitty.ini")

	if !files.Exists(iniFile) {
		log.Info().Msg("Creating default ini file...")
		err := os.WriteFile(iniFile, defaultKittyIni, 0644)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot write kitty.ini")
		}
	}

	log.Info().Msg("Updating configuration...")
	if err := files.ReplaceByPrefix(iniFile, "savemode=", "savemode=dir"); err != nil {
		log.Fatal().Err(err).Msg("Cannot set savemode")
	}
	if err := files.ReplaceByPrefix(iniFile, ";savemode=", "savemode=dir"); err != nil {
		log.Fatal().Err(err).Msg("Cannot set savemode")
	}
	if err := files.ReplaceByPrefix(iniFile, "configdir=", "configdir="+filepath.FromSlash(configPath)); err != nil {
		log.Fatal().Err(err).Msg("Cannot set configdir")
	}
	if err := files.ReplaceByPrefix(iniFile, ";configdir=", "configdir="+filepath.FromSlash(configPath)); err != nil {
		log.Fatal().Err(err).Msg("Cannot set configdir")
	}

	log.Info().Msg("Setting environment...")
	os.Setenv("KITTY_INI_FILE", filepath.FromSlash(iniFile))

	configPathEmpty, _ := files.IsDirEmpty(configPath)
	if configPathEmpty {
		log.Info().Msg("Converting registry settings to dir mode...")
		if err := proc.QuickCmd(app.Process, []string{"-convert-dir"}); err != nil {
			log.Error().Err(err).Msg("Cannot convert registry settings to dir mode")
		}
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
