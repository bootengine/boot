package gateway

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

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

type DBError struct {
	moduleName, action string
	err                error
}

//go:embed init.sql
var initSQL string

var errNoModuleFound = fmt.Errorf("no module found with this name")

func (d DBError) Error() string {
	if d.action == "open" || d.action == "close" {
		return fmt.Sprintf("failed to %s module database: %s", d.action, d.err.Error())
	}
	if d.action == "list" {
		return fmt.Sprintf("failed to %s modules: %s", d.action, d.err.Error())
	}
	return fmt.Sprintf("failed to %s module %s: %s", d.action, d.moduleName, d.err.Error())
}

func (m *ModuleGateway) OpenDatabase(databaseUrl string) error {
	d, err := sql.Open("sqlite", databaseUrl)
	if err != nil {
		return DBError{
			action: "open",
			err:    err,
		}
	}
	m.db = d

	m.DB = goqu.New("sqlite3", m.db)
	return nil

}

func (m ModuleGateway) InitDatabase() error {
	if m.DB == nil {
		fmt.Println("wuuuut")
	}
	_, err := m.DB.Exec(initSQL)
	if err != nil {
		fmt.Println("error in init")
	}

	return err
}

func (m *ModuleGateway) CloseDatabase() error {
	err := m.db.Close()
	if err != nil {
		return DBError{
			action: "close",
			err:    err,
		}
	}
	return nil
}

func NewModuleGateway() (*ModuleGateway, error) {
	return &ModuleGateway{}, nil
}

func (m ModuleGateway) AddModule(ctx context.Context, module model.Module) error {
	ex := m.DB.Insert(tablename).Prepared(true).Rows(module).Executor()
	if _, err := ex.ExecContext(ctx); err != nil {
		return DBError{
			action:     "add",
			err:        err,
			moduleName: module.Name,
		}
	}

	return nil
}

func (m ModuleGateway) GetModule(ctx context.Context, moduleName string) (*model.Module, error) {
	var res model.Module
	found, err := m.DB.From(tablename).Select(goqu.Star()).Where(
		goqu.C("module_name").Eq(moduleName),
	).ScanStructContext(ctx, &res)

	if err != nil {
		return nil, DBError{
			action:     "find",
			err:        err,
			moduleName: moduleName,
		}
	}

	if !found {
		return nil,
			DBError{
				action:     "find",
				err:        errNoModuleFound,
				moduleName: moduleName,
			}

	}

	return &res, err
}

func (m ModuleGateway) ListModules(ctx context.Context) ([]model.Module, error) {
	var res []model.Module
	err := m.DB.From(tablename).Select(goqu.Star()).ScanStructsContext(ctx, &res)
	if err != nil {
		return nil, DBError{
			action: "list",
			err:    err,
		}
	}
	return res, err
}

func (m ModuleGateway) UpdateModulePath(ctx context.Context, moduleName, modulePath string) error {
	ex := m.DB.Update(tablename).Prepared(true).Where(goqu.C("module_name").Eq(moduleName)).Set(goqu.Record{
		"module_path": modulePath,
	}).Executor()
	if r, err := ex.ExecContext(ctx); err != nil {
		return DBError{
			action:     "update",
			moduleName: moduleName,
			err:        err,
		}

	} else if affected, err := r.RowsAffected(); affected == 0 || err != nil {
		if affected == 0 {
			return DBError{
				action:     "update",
				moduleName: moduleName,
				err:        errNoModuleFound,
			}
		}
		return DBError{
			action:     "update",
			moduleName: moduleName,
			err:        err,
		}
	}
	return nil
}

func (m ModuleGateway) RemoveModule(ctx context.Context, moduleName string) error {
	ex := m.DB.From(tablename).Delete().Prepared(true).Where(goqu.C("module_name").Eq(moduleName)).Executor()
	if r, err := ex.ExecContext(ctx); err != nil {
		return DBError{
			action:     "remove",
			moduleName: moduleName,
			err:        err,
		}
	} else if affected, err := r.RowsAffected(); affected == 0 || err != nil {
		if affected == 0 {
			return DBError{
				action:     "remove",
				moduleName: moduleName,
				err:        errNoModuleFound,
			}
		}
		return DBError{
			action:     "remove",
			moduleName: moduleName,
			err:        err,
		}
	}
	return nil
}
