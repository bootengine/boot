package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bootengine/boot/internal/license"
	"github.com/bootengine/boot/internal/model"
	"github.com/bootengine/boot/internal/usecase"
	"github.com/charmbracelet/huh"
	extism "github.com/extism/go-sdk"
	"github.com/pterm/pterm"
)

type Runner struct {
	ctx      context.Context
	workflow model.Workflow
	modCase  *usecase.ModuleUsecase
}

type StepError struct {
	err                error
	moduleName, action string
}

func (s StepError) Error() string {
	return fmt.Sprintf("failed to execute runner for module %s, action %s: %s", s.moduleName, s.action, s.err.Error())
}

func (s StepError) GetType() string {
	return "steps"
}

type VarError struct {
	err  error
	vars string
}

func (v VarError) Error() string {
	return fmt.Sprintf("failed to handle var %s: %s", v.vars, v.err)
}
func (v VarError) GetType() string {
	return "vars"
}

type RunnerError interface {
	GetType() string
	Error() string
}

type ValueKey struct{}

var keepGoing bool = true

type NoKeepGoingError bool

func (n NoKeepGoingError) Error() string {
	return "no keep going"
}

func NewRunner(use *usecase.ModuleUsecase, workflow model.Workflow) *Runner {
	return &Runner{
		ctx:      context.Background(),
		modCase:  use,
		workflow: workflow,
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
		return elem.Module == "filer" && elem.Action == model.CreateFolderStructAction
	}) {
		err := huh.NewConfirm().Description("a folder_struct is set without explicit step to create it").
			Title("/!\\ Are you sure ?").
			Affirmative("Yes !").
			Negative("Oups ... No!").
			Value(&keepGoing).
			Run()
		if err != nil {
			return HuhError{Err: err}
		}
		if !keepGoing {
			return NoKeepGoingError(keepGoing)
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

	if r.workflow.Config.CreateRoot {
		projectName := r.ctx.Value(ValueKey{}).(map[string]any)["project_name"].(string)
		err = os.MkdirAll(projectName, 0775)
		if err != nil {
			return fmt.Errorf("failed to create root project directory: %w", err)
		}
	}

	return r.handleSteps()
}

var NoLicenseSelected = fmt.Errorf("no license selected")

func (r *Runner) createLicense() error {
	contextValue := r.ctx.Value(ValueKey{}).(map[string]any)
	selectedLicense, ok := contextValue["license"]
	if !ok {
		return NoLicenseSelected
	}

	licensePath := filepath.Join(contextValue["project_name"].(string), "LICENSE")

	r.ctx = context.WithValue(r.ctx, ValueKey{}, contextValue)
	content, err := license.GetLicenseContent(r.ctx, selectedLicense.(string))
	if err != nil {
		return err
	}

	return os.WriteFile(licensePath, []byte(*content), 0664)
}

func (r *Runner) handleVars() error {
	values := make(map[string]any)
	if _, ok := r.workflow.Vars["project_name"]; !ok && r.workflow.Config.CreateRoot {
		err := huh.NewConfirm().Description("the project needs a name, convention is forcing on a 'project_name' var").
			Title("/!\\ Caution !").
			Affirmative("Ok").
			Negative("Cancel").
			Value(&keepGoing).
			Run()
		if err != nil {
			return HuhError{Err: err}
		}
		if !keepGoing {
			return NoKeepGoingError(keepGoing)
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
			err := huh.NewSelect[string]().Title("what is the prefered license ?").
				Description(`
				more info here : https://choosealicense.com/licenses,
				if the license you want is not in the list, please create an issue () or contribute ()
				`).
				Options(
					huh.NewOption("MIT", "mit"),
					huh.NewOption("GNU GPL v3", "gnugpl3"),
					huh.NewOption("GNU AGPL v3", "gnuagpl3"),
					huh.NewOption("GNU LGPL v3", "gnulgpl3"),
					huh.NewOption("Mozilla Public License", "mozillapublic"),
					huh.NewOption("Apache 2.0", "apache2"),
					huh.NewOption("Boost Software License", "boostsoftware"),
					huh.NewOption("Unlicense", "unlicense"),
				).Run()
			if err != nil {
				return HuhError{Err: err}
			}
			values[k] = val
		default:
			return VarError{
				err:  fmt.Errorf("%s type of var is not managed by boot", v.Type),
				vars: k,
			}
		}
	}
	r.ctx = context.WithValue(r.ctx, ValueKey{}, values)

	return nil
}

func (r Runner) handleSteps() error {
	for _, step := range r.workflow.Steps {
		spin, _ := pterm.DefaultSpinner.WithShowTimer(true).Start(step.Name)
		if step.Module == "license" {
			err := r.createLicense()
			if err != nil && errors.Is(err, NoLicenseSelected) {
				// warn user
				spin.Warning("failed to create license: no selected license")
			}
			if err != nil {
				spin.Fail()
				return StepError{
					moduleName: step.Module,
					action:     string(step.Action),
					err:        err,
				}
			}
			continue
		}

		var config map[string]string
		mod, err := r.modCase.RetrieveModule(r.ctx, step.Module)
		if err != nil {
			spin.Fail()
			return StepError{
				moduleName: step.Module,
				action:     string(step.Action),
				err:        err,
			}
		}

		if !slices.Contains(model.Capabilities[mod.Type], step.Action) {
			spin.Fail()
			return StepError{
				moduleName: step.Module,
				action:     string(step.Action),
				err:        fmt.Errorf("this type of plugin (%s) can't run this action (%s)", mod.Type, step.Action),
			}
		}

		if mod.Type == model.FilerType {
			data, err := json.Marshal(r.workflow.FolderStruct)
			if err != nil {
			}
			config["folder_struct"] = string(data)
		}

		plugin, err := r.createPlugin(step, *mod, config)
		if err != nil {
			spin.Fail()
			return StepError{
				moduleName: step.Module,
				action:     string(step.Action),
				err:        err,
			}
		}

		var params []byte
		if step.Params != nil {
			params, err = json.Marshal(step.Params)
			if err != nil {
				spin.Fail()
				return StepError{
					moduleName: step.Module,
					action:     string(step.Action),
					err:        err,
				}
			}
		}

		exit, out, err := plugin.CallWithContext(r.ctx, string(step.Action), params)
		if err != nil {
			spin.Fail()
			return StepError{
				moduleName: step.Module,
				action:     string(step.Action),
				err:        err,
			}
		}

		if err = r.handleOutput(step, mod.Type, plugin, exit, out); err != nil {
			spin.Fail()
			return StepError{
				moduleName: step.Module,
				action:     string(step.Action),
				err:        err,
			}
		}
		spin.Success()
		spin.Stop()
	}

	return nil
}

func (r Runner) createPlugin(step model.Step, mod model.Module, config map[string]string) (*extism.Plugin, error) {
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Name: mod.Name,
				Path: mod.Path,
			},
		},
		Config: config,
	}

	hostFunctions := []extism.HostFunction{}
	pluginConfig := extism.PluginConfig{}
	if mod.Type == model.FilerType {
		callTemplate := extism.NewHostFunctionWithStack(
			"callTemplate",
			r.callTemplateFunc, []extism.ValueType{extism.ValueTypePTR, extism.ValueTypePTR}, []extism.ValueType{extism.ValueTypePTR, extism.ValueTypePTR})

		hostFunctions = append(hostFunctions, callTemplate)
		if mod.Type == model.FilerType || mod.Type == model.VCSType {
			pluginConfig.EnableWasi = true
		}
	}

	plugin, err := extism.NewPlugin(r.ctx, manifest, pluginConfig, hostFunctions)
	if err != nil {
		return nil, StepError{
			moduleName: step.Module,
			action:     string(step.Action),
			err:        err,
		}
	}
	return plugin, err
}

func (r Runner) handleOutput(step model.Step, modType model.ModuleType, plugin *extism.Plugin, exit uint32, out []byte) error {
	switch modType {
	case model.FilerType, model.VCSType:
		if exit != 0 {
			errString := plugin.GetErrorWithContext(r.ctx)
			return StepError{
				moduleName: step.Module,
				action:     string(step.Action),
				err:        errors.New(errString),
			}
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
			return StepError{
				moduleName: step.Module,
				action:     string(step.Action),
				err:        err,
			}
		}
		err = r.executeCommand(string(out), cwd)
		if err != nil {
			return StepError{
				moduleName: step.Module,
				action:     string(step.Action),
				err:        err,
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

func (r Runner) callTemplateFunc(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
	handleErr := func(msg []byte) {
		var err error
		stack[1], err = p.WriteBytes(msg)
		if err != nil {
			panic(fmt.Errorf("an error as occured while parsing template and the caller failed to handle it: %w", err))
		}
	}

	tempPath, err := p.ReadString(stack[1])
	if err != nil {
		handleErr([]byte(err.Error()))
	}
	tempEngine, err := p.ReadString(stack[0])
	if err != nil {
		handleErr([]byte(err.Error()))
	}

	mod, err := r.modCase.RetrieveModule(ctx, tempEngine)
	if err != nil {
		handleErr([]byte(err.Error()))
	}

	if mod.Type != model.TempEngineType {
		err = fmt.Errorf("")
		handleErr([]byte(err.Error()))
	}
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Name: mod.Name,
				Path: mod.Path,
			},
		},
	}
	config := extism.PluginConfig{}
	plugin, err := extism.NewPlugin(r.ctx, manifest, config, []extism.HostFunction{})
	if err != nil {
		handleErr([]byte(err.Error()))
	}

	fileContent, err := os.ReadFile(tempPath)
	if err != nil {
		handleErr([]byte(err.Error()))
	}

	ex, out, err := plugin.CallWithContext(ctx, string(model.FormatTemplAction), fileContent)
	if err != nil {
		handleErr([]byte(err.Error()))
	}

	if ex == 0 {
		stack[0], err = p.WriteBytes(out)
		if err != nil {
			panic(fmt.Sprintf("the filer plugin failed to retrieve data from the template engine %s: %s", mod.Name, err.Error()))
		}
	} else {
		errString := plugin.GetErrorWithContext(r.ctx)
		handleErr([]byte(fmt.Sprintf("template plugin exited with error status: %s", errString)))
	}
}
