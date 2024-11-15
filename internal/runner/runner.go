package runner

import (
	"context"
	"fmt"
	"slices"

	"github.com/bootengine/boot/internal/model"
	"github.com/charmbracelet/huh"
)

type Runner struct {
	ctx      context.Context
	workflow model.Workflow
}

type ValueKey struct{}

func NewRunner() *Runner {
	return &Runner{
		ctx: context.Background(),
	}
}

type HuhError struct {
	Err error
}

func (h HuhError) Error() string {
	return fmt.Sprintf("'huh' error: %s", h.Err.Error())
}

func (r Runner) checkFolderStructCreation() error {
	if r.workflow.FolderStruct != nil && !slices.ContainsFunc(r.workflow.Steps, func(elem model.Step) bool {
		return elem.Module == "filer" && elem.Action == ""
	}) {
		err := huh.NewConfirm().Description("a folder_struct is set without explicit step to create it").
			Title("/!\\ Are you sure ?").
			Affirmative("Yes !").
			Negative("Oups ... No!").
			Run()
		if err != nil {
			return HuhError{Err: err}
		}
	}
	return nil
}

func (r Runner) Run() error {
	err := r.checkFolderStructCreation()
	if err != nil {
		return err
	}

	err = r.handleVars()
	if err != nil {
		return err
	}

	return r.handleSteps()
}

func (r *Runner) handleVars() error {
	values := make(map[string]any)
	if _, ok := r.workflow.Vars["project_name"]; !ok && r.workflow.Config.CreateRoot {
		// TODO: handle cancellation
		err := huh.NewConfirm().Description("the project needs a name, convention is forcing on a 'project_name' var").
			Title("/!\\ Caution !").
			Affirmative("Ok").
			Negative("Cancel").
			Run()
		if err != nil {
			return HuhError{Err: err}
		}
	}

	for k, v := range r.workflow.Vars {
		switch v.Type {
		case model.String:
			var val string
			err := huh.NewInput().Title(fmt.Sprintf("what is your %s ?", k)).Value(&val).Run()
			if err != nil {
				return HuhError{Err: err}
			}
			values[k] = val
		case model.Password:
			var val string
			err := huh.NewInput().Title(fmt.Sprintf("what is your %s ?", k)).EchoMode(huh.EchoModePassword).Value(&val).Run()
			if err != nil {
				return HuhError{Err: err}
			}
			values[k] = val
		case model.License:
			var val string
			err := huh.NewSelect[string]().Title("what is the prefered license ?").Options(
				huh.NewOption("MIT", "mit"),
				huh.NewOption("GNU GPL v3", "gnugpl3"),
			).Run()
			if err != nil {
				return HuhError{Err: err}
			}
			values[k] = val
		default:
			return fmt.Errorf("%s type of var is not managed by boot", v.Type)
		}
	}
	r.ctx = context.WithValue(r.ctx, ValueKey{}, values)

	return nil
}

func (r Runner) handleSteps() error {
	for _, step := range r.workflow.Steps {
		// find module

		fmt.Println(step)
	}

	return nil
}
