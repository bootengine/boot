package cmd

import (
	"context"
	"path/filepath"
	"regexp"

	"github.com/bootengine/boot/internal/helper"
	"github.com/bootengine/boot/internal/model"
	"github.com/bootengine/boot/internal/usecase"
	"github.com/spf13/cobra"
)

type installCmdFlags struct {
	name       string
	pathOrURL  string
	moduleType string
}

var installFlags installCmdFlags

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:           "install",
	Aliases:       []string{"i"},
	SilenceErrors: true,
	Short:         "Install a module",
	Long: `Install a module, given a path (or url) and a type (cmd, filer, vcs, template_engine).
	----
	cmd: module that will return a command to be executed.
	filer: module that will create files/folder or bootstrap the folder_struct definition.
	vcs: module that will run vcs based command like commit and push code.
	template_engine: module that will handle templating in the folder_struct definition.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return helper.WithModuleUsecase(func(ctx context.Context, use *usecase.ModuleUsecase) error {
			var (
				moduleType model.ModuleType
				err        error
			)
			moduleType.FromString(installFlags.moduleType)
			reg := regexp.MustCompile("^(http|https)://.*$")
			if reg.Match([]byte(installFlags.pathOrURL)) {
				err = use.InstallModuleFromURL(ctx, installFlags.name, moduleType, installFlags.pathOrURL)
				if err != nil {
					return err
				}
				return nil
			} else {
				installFlags.pathOrURL, err = filepath.Abs(installFlags.pathOrURL)
				if err != nil {
					return err
				}
			}
			return use.InstallModuleFromFS(ctx, installFlags.name, moduleType, installFlags.pathOrURL)
		})
	},
}

func init() {
	moduleCmd.AddCommand(installCmd)

	installCmd.Flags().StringVarP(&installFlags.name, "name", "n", "", "module's name - don't forget that the name is UNIQUE.")
	installCmd.Flags().StringVarP(&installFlags.pathOrURL, "location", "l", "", "module's location - it can be either a path or a URL to a .wasm file.")
	installCmd.Flags().StringVarP(&installFlags.moduleType, "type", "t", "", "module's type - one of [filer,cmd,vcs,template_engine].")

	installCmd.MarkFlagRequired("name")
	installCmd.MarkFlagRequired("location")
	installCmd.MarkFlagRequired("type")
}
