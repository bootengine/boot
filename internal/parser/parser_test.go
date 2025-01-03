package parser_test

import (
	"testing"

	"github.com/bootengine/boot/internal/model"
	"github.com/bootengine/boot/internal/parser"
	"github.com/maxatome/go-testdeep/td"
)

func TestParser_Parse(t *testing.T) {
	p := parser.NewParser()

	expected := model.Workflow{
		Config: model.Config{
			CreateRoot:   true,
			Unrestricted: false,
		},
		Vars: model.Vars{
			model.Var{
				Name:     "project_name",
				Type:     model.String,
				Required: true,
			},
			model.Var{
				Name:     "author_github_name",
				Type:     model.String,
				Required: true,
			},
			model.Var{
				Name:     "license",
				Type:     model.License,
				Required: false,
			},
		},
		Steps: []model.Step{
			{
				Name:   "git init",
				Module: "git",
				Action: model.InitAction,
			},
			{
				Name:   "go mod init",
				Module: "go",
				Action: model.InitAction,
			},
			{
				Name:              "go get deps",
				Module:            "go",
				Action:            model.InstallLocalDepsAction,
				CurrentWorkingDir: "frontend",
				Params: []string{
					"github.com/charmbracelet/bubbletea",
					"github.com/charmbracelet/log",
				},
			},
		},
		FolderStruct: model.FolderStruct{
			model.Folder{
				Name: "cmd",
				Filers: model.FolderStruct{
					model.Folder{
						Name: "install",
					},
					model.Folder{
						Name: "remove",
					},
					model.File{
						Name: "root.go",
					},
				},
			},
			model.Folder{
				Name: "internal",
			},
			model.File{
				Name: "main.go",
				TempWrapper: &model.TempWrapper{
					TemplateDef: model.TemplateDef{
						Filepath: "./temp.go",
						Engine:   "jinja2",
					},
				},
			},
		},
	}

	got, err := p.Parse("../mocks/workflow.yaml")
	td.CmpNoError(t, err)
	if err == nil {
		td.Cmp(t, got, &expected)
	}

}
