package cmd

import (
	"encoding/json"
	"os"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/yaml"
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
		var (
			marshaled []byte
			err       error
		)
		switch initFlags.outputType {
		case helper.JSON:
			marshaled, err = json.Marshal(workflow)
		case helper.YAML:
			ctx := cuecontext.New()
			val := ctx.Encode(workflow)
			marshaled, err = yaml.Encode(val)
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
			"project_name": model.Var{
				Required: true,
				Type:     model.String,
			},
			"license": model.Var{
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
