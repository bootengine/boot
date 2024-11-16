package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/bootengine/boot/internal/model"
	"github.com/bootengine/boot/internal/usecase"
	"github.com/charmbracelet/huh"
	extism "github.com/extism/go-sdk"
)

type Runner struct {
	ctx      context.Context
	workflow model.Workflow
	modCase  *usecase.ModuleUsecase
}

type ValueKey struct{}

func NewRunner(use *usecase.ModuleUsecase) *Runner {
	return &Runner{
		ctx:     context.Background(),
		modCase: use,
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
		mod, err := r.modCase.RetrieveModule(r.ctx, step.Name)
		if err != nil {
			return err
		}

		// check type and action
		fmt.Println(step)

		manifest := extism.Manifest{

			Wasm: []extism.Wasm{
				extism.WasmFile{
					Name: mod.Name,
					Path: mod.Path,
				},
			},
		}

		config := extism.PluginConfig{
			EnableWasi: true,
		}

		plugin, err := extism.NewPlugin(r.ctx, manifest, config, []extism.HostFunction{})
		if err != nil {
			return err
		}
		var params []byte
		if step.Params != nil {
			params, err = json.Marshal(step.Params)
			if err != nil {
				return err
			}
		}

		// if action is folder_struct, need to get template before

		exit, out, err := plugin.CallWithContext(r.ctx, string(step.Action), params)
		if err != nil {
			return err
		}
		switch mod.Type {
		case model.FilerType, model.VCSType:
			if exit != 0 {
				errString := plugin.GetErrorWithContext(r.ctx)
				return errors.New(errString)
			}
		// log success
		case model.CmdType:
			cwd, err := func() (string, error) {
				if step.CurrentWorkingDir != "" {
					return step.CurrentWorkingDir, nil
				}
				if r.workflow.Config.CreateRoot {
					return r.ctx.Value(ValueKey{}).(map[string]any)["project_name"].(string), nil
				}
				return os.Getwd()
			}()
			if err != nil {
				return err
			}
			err = r.executeCommand(string(out), cwd)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (r Runner) executeCommand(cmd string, cwd string) error {
	if !r.workflow.Config.Unrestricted && !r.checkCommandContent(cmd) {
		return fmt.Errorf("plugin is trying to execute a suspicious command: %s", cmd)
	}

	command := exec.CommandContext(r.ctx, cmd)
	command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr

	command.Dir = cwd
	return command.Run()
}

func (r Runner) checkCommandContent(cmd string) bool {
	if strings.Contains(cmd, "sudo") || strings.Contains(cmd, "rm") {
		return false
	}
	return true
}
