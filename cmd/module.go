package cmd

import (
	"github.com/spf13/cobra"
)

// moduleCmd represents the module command
var moduleCmd = &cobra.Command{
	Use:     "module",
	Aliases: []string{"mod", "m"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(moduleCmd)
}
