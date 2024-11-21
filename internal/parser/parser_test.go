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
			"project_name": model.Var{
				Type:     model.String,
				Required: true,
			},
			"author_github_name": model.Var{
				Type:     model.String,
				Required: true,
			},
			"license": model.Var{
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
