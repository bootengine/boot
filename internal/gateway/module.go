package gateway

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bootengine/boot/internal/model"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "modernc.org/sqlite"
)

type ModuleGateway struct {
	db *sql.DB
	DB *goqu.Database
}

const tablename = "module"

// TODO: it's not its place !
func meh() error {
	home, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	dataDir := filepath.Join(home, "boot", "data")
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		return err
	}
	filepath.Join(dataDir, "db")
	return nil
}

func (m *ModuleGateway) OpenDatabase(databaseUrl string) error {
	d, err := sql.Open("sqlite", databaseUrl)
	if err == nil {
		m.db = d
	}

	m.DB = goqu.New("sqlite3", m.db)

	return err
}

func (m *ModuleGateway) CloseDatabase() error {
	return m.db.Close()
}

func NewModuleGateway() (*ModuleGateway, error) {
	return &ModuleGateway{}, nil
}

func (m ModuleGateway) AddModule(ctx context.Context, module model.Module) error {
	ex := m.DB.Insert(tablename).Prepared(true).Rows(module).Executor()
	if _, err := ex.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (m ModuleGateway) GetModule(ctx context.Context, moduleName string) (*model.Module, error) {
	var res model.Module
	found, err := m.DB.From(tablename).Select(goqu.Star()).Where(
		goqu.C("module_name").Eq(moduleName),
	).ScanStructContext(ctx, &res)

	if !found {
		return nil, fmt.Errorf("module not found")
	}

	return &res, err
}

func (m ModuleGateway) ListModules(ctx context.Context) ([]model.Module, error) {
	var res []model.Module
	err := m.DB.From(tablename).Select(goqu.Star()).ScanStructsContext(ctx, &res)
	return res, err
}

func (m ModuleGateway) UpdateModulePath(ctx context.Context, moduleName, modulePath string) error {
	ex := m.DB.Update(tablename).Prepared(true).Where(goqu.C("module_name").Eq(moduleName)).Set(goqu.Record{
		"module_path": modulePath,
	}).Executor()
	if r, err := ex.ExecContext(ctx); err != nil {
		return err
	} else if affected, err := r.RowsAffected(); affected == 0 || err != nil {
		if affected == 0 {
			return fmt.Errorf("no module found with this name")
		}
		return err
	}
	return nil
}

func (m ModuleGateway) RemoveModule(ctx context.Context, moduleName string) error {
	ex := m.DB.From(tablename).Delete().Prepared(true).Where(goqu.C("module_name").Eq(moduleName)).Executor()
	if r, err := ex.ExecContext(ctx); err != nil {
		return err
	} else if affected, err := r.RowsAffected(); affected == 0 || err != nil {
		if affected == 0 {
			return fmt.Errorf("no module found with this name")
		}
		return err
	}
	return nil
}
