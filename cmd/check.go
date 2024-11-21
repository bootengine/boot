package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type checkCmdsFlags struct {
	filename string
}

var checkFlags checkCmdsFlags

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check that the given file is valid.",
	Long:  `Check that the given file is a valid boot workflow. It will not check that selected module are installed, it will just check that it`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("check called")
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringVarP(&checkFlags.filename, "filename", "f", "", "the path to the config file you want to check.")
	checkCmd.MarkFlagRequired("filename")
}
