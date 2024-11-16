package helper

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/bootengine/boot/internal/gateway"
	"github.com/bootengine/boot/internal/usecase"
)

func WithModuleUsecase(exec func(ctx context.Context, use *usecase.ModuleUsecase)) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	datastore, err := gateway.NewModuleGateway()
	if err != nil {
	}

	dbUrl, err := getDbURL()
	if err != nil {
	}

	err = datastore.OpenDatabase(*dbUrl)
	if err != nil {
	}

	defer func() {
		err = datastore.CloseDatabase()
		signal.Stop(c)
		if err != nil {
		}
	}()

	go func() {
		<-c
		cancel()
	}()

	u := usecase.NewModuleUsecase(datastore)

	exec(ctx, u)

}

func getDbURL() (*string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	dataDir := filepath.Join(home, "boot", "data")
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dataDir, "db")

	return &path, nil
}
