package cmd

import (
	"context"
	"errors"
	"regexp"
	"sync"

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
	Use:           "gen",
	Short:         "Generate a new project from a config file.",
	Long:          `Generate a new project from a config file. This file can be either on your local computer or it can be a repository.`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var once sync.Once
		onceFunc := func() {
			defer func() {
				err := cleanup()
				if err != nil {
					// panic ?
				}
			}()
		}
		return helper.WithModuleUsecase(func(ctx context.Context, use *usecase.ModuleUsecase) error {
			reg := regexp.MustCompile("^(http|https)://.*$")
			if reg.Match([]byte(genFlags.pathOrURL)) {
				// dl repo in tmp folder
				// compute path to tmp/repo/whatever.[yaml|json|toml|...]
				// defer cleanup
				once.Do(onceFunc)
			}
			work, err := parser.NewParser().Parse(genFlags.pathOrURL)
			if err != nil {
				return err
			}

			worker := runner.NewRunner(use, *work)
			err = worker.Run()
			if err != nil {
				if errors.Is(err, runner.NoKeepGoingError(false)) {
					once.Do(onceFunc)
				}
				return err
			}
			return nil
		})
	},
}

func cleanup() error {
	return nil
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringVarP(&genFlags.pathOrURL, "file", "f", "", `config file for the generation process, can be either a repo url or a local path.
If it's a repo url, the repo will be downloaded in a tmp dir and removed afterward.
		`)

	genCmd.MarkFlagRequired("file")
}
