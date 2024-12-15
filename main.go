/*
Copyright © 2024 BootEngine <mathob.jehanno@hotmail.fr>
*/
package main

import (
	"os"
	"path/filepath"

	"github.com/bootengine/boot/cmd"
	"github.com/bootengine/boot/internal/gateway"
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

func init_db(dataDir string) error {
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
		}
	}()
	err = mg.InitDatabase()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func bootstrap() error {
	dataDir, err := create_config_folder()
	if err != nil {
	}

	return init_db(dataDir)
}

func init() {
	err := bootstrap()
	if err != nil {
	}
}

func main() {
	cmd.Execute()
}
