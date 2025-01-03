package gateway_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/bootengine/boot/internal/gateway"
	"github.com/bootengine/boot/internal/model"
	"github.com/maxatome/go-testdeep/td"
)

//go:embed init.sql
var db_init string

type key string

const (
	gtw  key = "gateway"
	test key = "testing"
)

func Suite(t *testing.T, run func(ctx context.Context)) {
	gt, err := gateway.NewModuleGateway()
	td.Require(t).CmpNoError(err)
	err = gt.OpenDatabase(":memory:")
	td.Require(t).CmpNoError(err)

	_, err = gt.DB.Exec(db_init)
	td.Require(t).CmpNoError(err)

	ctx := context.WithValue(context.Background(), gtw, gt)
	ctx = context.WithValue(ctx, test, t)

	run(ctx)

	defer func() {
		err = gt.CloseDatabase()
		td.Require(t).CmpNoError(err)
	}()
}

func Test_AddModule(t *testing.T) {
	Suite(t, func(ctx context.Context) {
		gt := ctx.Value(gtw).(*gateway.ModuleGateway)
		t := ctx.Value(test).(*testing.T)

		tests := []struct {
			name        string
			input       model.Module
			expectedErr string
		}{
			{
				name: "valid",
				input: model.Module{
					Name: "go",
					Path: "~/Documents/test.html",
					Type: "vcs",
				},
			},
			{
				name: "name already taken",
				input: model.Module{
					Name: "go",
					Path: "~/Documents/test.html",
					Type: "filer",
				},
				expectedErr: "UNIQUE constraint failed: module.module_name",
			},
			{
				name: "type not valid",
				input: model.Module{
					Name: "another",
					Path: "~/Documents/test.html",
					Type: "fziuherfuh",
				},
				expectedErr: `module_type IN ("cmd", "filer", "vcs","template_engine")`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := gt.AddModule(ctx, tt.input)
				if tt.expectedErr != "" {
					td.CmpContains(t, err, tt.expectedErr)
				} else {
					td.CmpNoError(t, err)
					res, err := gt.ListModules(ctx)
					td.CmpNoError(t, err)

					td.Cmp(t, len(res), 1)
				}
			})
		}
	})
}

func Test_ListModule(t *testing.T) {
	Suite(t, func(ctx context.Context) {
		gt := ctx.Value(gtw).(*gateway.ModuleGateway)
		t := ctx.Value(test).(*testing.T)

		m1 := model.Module{
			Name: "go",
			Path: "./tmp",
			Type: "cmd",
		}
		m2 := model.Module{
			Name: "node",
			Path: "./tmp",
			Type: "cmd",
		}

		err := gt.AddModule(ctx, m1)
		td.CmpNoError(t, err)
		err = gt.AddModule(ctx, m2)
		td.CmpNoError(t, err)

		got, err := gt.ListModules(ctx)
		td.CmpNoError(t, err)

		td.Cmp(t, len(got), 2)

	})
}

func Test_GetModule(t *testing.T) {
	Suite(t, func(ctx context.Context) {
		gt := ctx.Value(gtw).(*gateway.ModuleGateway)
		t := ctx.Value(test).(*testing.T)

		m1 := model.Module{
			Name: "go",
			Path: "./tmp",
			Type: "cmd",
		}
		m2 := model.Module{
			Name: "node",
			Path: "./tmp",
			Type: "cmd",
		}
		m3 := model.Module{
			Name: "python",
			Path: "./tmp",
			Type: "cmd",
		}

		err := gt.AddModule(ctx, m1)
		td.CmpNoError(t, err)
		err = gt.AddModule(ctx, m2)
		td.CmpNoError(t, err)
		err = gt.AddModule(ctx, m3)
		td.CmpNoError(t, err)

		got, err := gt.GetModule(ctx, "node")
		td.CmpNoError(t, err)

		td.Cmp(t, got.Name, "node")

		got, err = gt.GetModule(ctx, "vuejs")
		td.Cmp(t, got, td.Nil())
		td.CmpContains(t, err, "no module found with this name")

	})
}

func Test_UpdateModule(t *testing.T) {
	Suite(t, func(ctx context.Context) {
		gt := ctx.Value(gtw).(*gateway.ModuleGateway)
		t := ctx.Value(test).(*testing.T)
		m1 := model.Module{
			Name: "go",
			Path: "./tmp",
			Type: "cmd",
		}
		m2 := model.Module{
			Name: "node",
			Path: "./tmp",
			Type: "cmd",
		}
		m3 := model.Module{
			Name: "python",
			Path: "./tmp",
			Type: "cmd",
		}

		err := gt.AddModule(ctx, m1)
		td.CmpNoError(t, err)
		err = gt.AddModule(ctx, m2)
		td.CmpNoError(t, err)
		err = gt.AddModule(ctx, m3)
		td.CmpNoError(t, err)

		err = gt.UpdateModulePath(ctx, "node", "/bin")
		td.CmpNoError(t, err)
		got, err := gt.GetModule(ctx, "node")
		td.CmpNoError(t, err)
		td.Cmp(t, got.Path, "/bin")

		err = gt.UpdateModulePath(ctx, "cargo", "/bin")
		td.CmpContains(t, err, "no module found with this name")
	})
}

func Test_DeleteModule(t *testing.T) {
	Suite(t, func(ctx context.Context) {
		gt := ctx.Value(gtw).(*gateway.ModuleGateway)
		t := ctx.Value(test).(*testing.T)
		m1 := model.Module{
			Name: "go",
			Path: "./tmp",
			Type: "cmd",
		}
		m2 := model.Module{
			Name: "node",
			Path: "./tmp",
			Type: "cmd",
		}
		m3 := model.Module{
			Name: "python",
			Path: "./tmp",
			Type: "cmd",
		}

		err := gt.AddModule(ctx, m1)
		td.CmpNoError(t, err)
		err = gt.AddModule(ctx, m2)
		td.CmpNoError(t, err)
		err = gt.AddModule(ctx, m3)
		td.CmpNoError(t, err)

		err = gt.RemoveModule(ctx, "node")
		td.CmpNoError(t, err)
		err = gt.RemoveModule(ctx, "cargo")
		td.CmpContains(t, err, "no module found with this name")
	})
}
