package cmd

import (
	"context"

	"github.com/bootengine/boot/internal/helper"
	"github.com/bootengine/boot/internal/usecase"
	"github.com/spf13/cobra"
)

type removeCmdFlags struct {
	name string
}

var removeFlags removeCmdFlags

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:           "remove",
	Aliases:       []string{"rm", "r"},
	Short:         "Remove a module from it's name.",
	Long:          `Remove a module from it's name.`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return helper.WithModuleUsecase(func(ctx context.Context, use *usecase.ModuleUsecase) error {
			return use.RemoveModule(ctx, removeFlags.name)
		})
	},
}

func init() {
	moduleCmd.AddCommand(removeCmd)
	removeCmd.Flags().StringVarP(&removeFlags.name, "name", "n", "", "name of the module you want to remove.")
	installCmd.MarkFlagRequired("name")

}
