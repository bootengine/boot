/*
Copyright © 2024 BootEngine <mathob.jehanno@hotmail.fr>
*/
package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "boot",
	Short: "boot is a higthly customizable project bootstrapper.",
	Long: `boot is a higthly customizable project bootstrapper.
It comes with a plugin system that accepts .wasm files meaning that you can create plugins in anything that compile to wasm.
	`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		log.Errorf("something bad happened: %s", err.Error())
		os.Exit(1)
	}
}

func init() {
}
