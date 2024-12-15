package helper

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/bootengine/boot/internal/gateway"
	"github.com/bootengine/boot/internal/usecase"
)

type useCaseFunc func(ctx context.Context, use *usecase.ModuleUsecase) error

func WithModuleUsecase(exec useCaseFunc) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	datastore, err := gateway.NewModuleGateway()
	if err != nil {
		cancel()
		return err
	}

	dbUrl, err := getDbURL()
	if err != nil {
		cancel()
		return err
	}

	err = datastore.OpenDatabase(*dbUrl)
	if err != nil {
		cancel()
		return err
	}

	defer func() {
		err = datastore.CloseDatabase()
		signal.Stop(c)
		if err != nil {
			// TODO: panic ?
		}
	}()

	go func() {
		<-c
		cancel()
	}()

	u := usecase.NewModuleUsecase(datastore)

	return exec(ctx, u)
}

func getDbURL() (*string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	dataDir := filepath.Join(home, "bootengine", "data")
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dataDir, "db")

	return &path, nil
}
