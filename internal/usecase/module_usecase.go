package usecase

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bootengine/boot/internal/model"
	"github.com/bootengine/boot/internal/repository"
)

type ModuleUsecase struct {
	Datastore repository.ModuleRepository
}

func NewModuleUsecase(datastore repository.ModuleRepository) *ModuleUsecase {
	return &ModuleUsecase{
		Datastore: datastore,
	}
}

func (m ModuleUsecase) getInstallFolder() (*string, error) {
	config, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	installPath := filepath.Join(config, "boot", "plugins")

	return &installPath, nil
}

func (m ModuleUsecase) RetrieveModule(ctx context.Context, modName string) (*model.Module, error) {
	return m.Datastore.GetModule(ctx, modName)
}

func (m ModuleUsecase) ListModules(ctx context.Context) ([]model.Module, error) {
	return m.Datastore.ListModules(ctx)
}

func (m ModuleUsecase) InstallModule(ctx context.Context, modName string, modType model.PluginType, modUrl string) error {
	// install in filesystem
	res, err := http.Get(modUrl)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	installPath, err := m.getInstallFolder()
	if err != nil {
		return err
	}
	pluginPath := filepath.Join(*installPath, modName)

	err = os.WriteFile(pluginPath, data, 0755)
	if err != nil {
		return err
	}

	mod := model.Module{
		Name: modName,
		Path: pluginPath,
		Type: modType,
	}

	// install in Db
	return m.Datastore.AddModule(ctx, mod)
}

func (m ModuleUsecase) RemoveModule(ctx context.Context, modName string) error {
	// delete in fs
	mod, err := m.Datastore.GetModule(ctx, modName)
	if err != nil {
		return err
	}
	err = os.Remove(mod.Path)
	if err != nil {
		return err
	}

	return m.Datastore.RemoveModule(ctx, modName)
}
