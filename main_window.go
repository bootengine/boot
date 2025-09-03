//go:build windows

/*
Copyright © 2024 BootEngine <mathob.jehanno@hotmail.fr>
*/
package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/bootengine/boot/cmd"

	"github.com/bootengine/boot/internal/assets"
	"github.com/bootengine/boot/internal/gateway"
	"github.com/bootengine/boot/internal/model"
	"github.com/charmbracelet/log"
)

func create_config_folder() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dataDir := filepath.Join(configDir, "bootengine", "data")
	err = os.MkdirAll(dataDir, 0775)
	if err != nil {
		log.Error("boot-create:" + err.Error())
		return "", err
	}
	return dataDir, nil
}

func init_db(mg *gateway.ModuleGateway) error {
	return mg.InitDatabase()
}

func init_plugins(mg *gateway.ModuleGateway, dataDir string) error {
	err := os.CopyFS(dataDir, assets.PluginFS)
	if err != nil {
		return errors.New("failed to create default plugin directory")
	}

	err = mg.AddModule(context.Background(), model.Module{
		Name: "git",
		Type: model.VCSType,
		Path: filepath.Join(dataDir, "plugins", "boot-git.wasm"),
	})
	if err != nil {
		return nil
	}

	err = mg.AddModule(context.Background(), model.Module{
		Name: "filer",
		Type: model.FilerType,
		Path: filepath.Join(dataDir, "plugins", "boot-filer.wasm"),
	})
	if err != nil {
		return nil
	}

	err = mg.AddModule(context.Background(), model.Module{
		Name: "jinja",
		Type: model.TempEngineType,
		Path: filepath.Join(dataDir, "plugins", "boot-template.wasm"),
	})
	return err
}

func bootstrap() error {
	dataDir, err := create_config_folder()
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(dataDir, "INSTALL")); err == nil {
		return nil
	}
	mg, err := gateway.NewModuleGateway()
	if err != nil {
		return err
	}
	err = mg.OpenDatabase(filepath.Join(dataDir, "db"))
	if err != nil {
		log.Error("boot-open:" + err.Error())
		return err
	}
	defer func() {
		err = mg.CloseDatabase()
		if err != nil {
			log.Fatalf("failed to close database: %s", err.Error())
		}
	}()

	err = init_db(mg)
	if err != nil {
		return err
	}

	err = init_plugins(mg, dataDir)
	if err != nil {
		return err
	}

	_, err = os.Create(filepath.Join(dataDir, "INSTALL"))
	return err
}

func init() {
	err := bootstrap()
	if err != nil {
		log.Fatalf("failed to run init phase of boot application: %s", err.Error())
	}
}

func main() {
	cmd.Execute()
}
