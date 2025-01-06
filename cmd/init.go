package cmd

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
	//"cuelang.org/go/encoding/yaml"
	"github.com/bootengine/boot/internal/helper"
	"github.com/bootengine/boot/internal/model"
	"github.com/spf13/cobra"
)

type initCmdFlags struct {
	outputFilename string
	outputType     helper.SupportedFileType
}

var initFlags initCmdFlags

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:           "init",
	Short:         "Init help you generate a default boot workflow file.",
	Long:          `Init will help you generate a default boot workflow file into the output filename given an output type [json, yaml]`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		workflow := newDefaultWorkflow()

		gen := model.GeneratingWorkflow{
			Config:       workflow.Config,
			Vars:         workflow.Vars,
			Steps:        workflow.Steps,
			FolderStruct: workflow.FolderStruct.Convert(),
		}

		var (
			marshaled []byte
			err       error
		)
		switch initFlags.outputType {
		case helper.JSON:
			marshaled, err = json.Marshal(gen)
		case helper.YAML:
			marshaled, err = yaml.Marshal(gen)
		}

		if err != nil {
			return err
		}

		return os.WriteFile(initFlags.outputFilename, marshaled, 0664)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&initFlags.outputFilename, "output", "o", "", "filename of the output.")
	initCmd.Flags().StringVarP((*string)(&initFlags.outputType), "type", "t", "yaml", "type of the output file.")
	initCmd.MarkFlagRequired("output")
}

func newDefaultWorkflow() model.Workflow {
	return model.Workflow{
		Config: model.Config{
			CreateRoot: true,
		},
		Vars: model.Vars{
			model.Var{
				Name:     "project_name",
				Required: true,
				Type:     model.String,
			},
			model.Var{
				Name:     "license",
				Required: true,
				Type:     model.License,
			},
		},
		Steps: []model.Step{
			{
				Name:   "git init",
				Module: "git",
				Action: model.InitAction,
			},
			{
				Name:   "create folder structure",
				Module: "filer",
				Action: model.CreateFolderStructAction,
			},
		},
		FolderStruct: model.FolderStruct{
			model.File{
				Name: ".gitignore",
			},
			model.File{
				Name: "README.md",
			},
		},
	}
}
