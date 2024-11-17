package cmd

import (
	"context"

	"github.com/bootengine/boot/internal/helper"
	"github.com/bootengine/boot/internal/usecase"
	"github.com/spf13/cobra"
)

type updateCmdFlags struct {
	name string
	path string
}

var updateFlags updateCmdFlags

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update the path to a given module.",
	Long:  `update the path to a given module, assuming you're targeting a locally installed .wasm file.`,
	Run: func(cmd *cobra.Command, args []string) {
		helper.WithModuleUsecase(func(ctx context.Context, use *usecase.ModuleUsecase) {
			err := use.UpdateModule(ctx, updateFlags.name, updateFlags.path)
			if err != nil {
			}
		})
	},
}

func init() {
	moduleCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updateFlags.name, "name", "n", "", "the name of the module you want to modify.")
	updateCmd.Flags().StringVarP(&updateFlags.path, "path", "p", "", "the local path to the new module file (.wasm).")
	installCmd.MarkFlagRequired("name")
	installCmd.MarkFlagRequired("path")
}
