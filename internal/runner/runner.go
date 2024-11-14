package runner

import (
	"context"
	"fmt"

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

func (r Runner) Run() error {
	err := r.handleVars()
	if err != nil {
		return err
	}

	return r.handleSteps()
}

func (r *Runner) handleVars() error {
	values := make(map[string]any)
	for k, v := range r.workflow.Vars {
		switch v.Type {
		case model.String:
			var val string
			err := huh.NewInput().Title(fmt.Sprintf("what is your %s ?", k)).Value(&val).Run()
			if err != nil {
				return err
			}
			values[k] = val
		case model.Password:
			var val string
			err := huh.NewInput().Title(fmt.Sprintf("what is your %s ?", k)).EchoMode(huh.EchoModePassword).Value(&val).Run()
			if err != nil {
				return err
			}
			values[k] = val
		case model.License:
			var val string
			err := huh.NewSelect[string]().Title("what is the prefered license ?").Options(
				huh.NewOption("MIT", "mit"),
				huh.NewOption("GNU GPL v3", "gnugpl3"),
			).Run()
			if err != nil {
				return err
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

	return nil
}
