package repository

import (
	"context"

	"github.com/bootengine/boot/internal/model"
)

type ModuleRepository interface {
	AddModule(ctx context.Context, module model.Module) error
	GetModule(ctx context.Context, moduleName string) (*model.Module, error)
	ListModules(ctx context.Context) ([]model.Module, error)
	UpdateModulePath(ctx context.Context, moduleName, modulePath string) error
	RemoveModule(ctx context.Context, moduleName string) error
}
