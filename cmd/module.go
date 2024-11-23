package cmd

import (
	"github.com/spf13/cobra"
)

// moduleCmd represents the module command
var moduleCmd = &cobra.Command{
	Use:     "module",
	Aliases: []string{"mod", "m"},
	Short:   "Command to manipulate installed module",
}

func init() {
	rootCmd.AddCommand(moduleCmd)
}
