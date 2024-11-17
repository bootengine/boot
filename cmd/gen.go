package cmd

import (
	"context"
	"regexp"

	"github.com/bootengine/boot/internal/helper"
	"github.com/bootengine/boot/internal/parser"
	"github.com/bootengine/boot/internal/runner"
	"github.com/bootengine/boot/internal/usecase"
	"github.com/spf13/cobra"
)

type genCmdFlags struct {
	pathOrURL string
}

var genFlags genCmdFlags

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate a new project from a config file.",
	Long:  `generate a new project from a config file. This file can be either on your local computer or it can be a repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		helper.WithModuleUsecase(func(ctx context.Context, use *usecase.ModuleUsecase) {
			reg := regexp.MustCompile("^(http|https)://.*$")
			if reg.Match([]byte(genFlags.pathOrURL)) {
				// dl repo in tmp folder
				// compute path to tmp/repo/whatever.[yaml|json|toml|...]

			}
			work, err := parser.NewParser().Parse(genFlags.pathOrURL)
			if err != nil {
			}

			worker := runner.NewRunner(use, *work)
			err = worker.Run()
			if err != nil {
			}
		})
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringVarP(&genFlags.pathOrURL, "file", "f", "", `config file for the generation process, can be either a repo url or a local path.
If it's a repo url, the repo will be downloaded in a tmp dir and removed afterward.
		`)

	genCmd.MarkFlagRequired("file")
}
