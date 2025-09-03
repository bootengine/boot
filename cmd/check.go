package cmd

import (
	"cuelang.org/go/cue/cuecontext"
	"github.com/bootengine/boot/internal/helper"
	"github.com/bootengine/boot/internal/parser"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type checkCmdsFlags struct {
	filename string
}

var checkFlags checkCmdsFlags

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:           "check",
	Short:         "Check that the given file is valid.",
	Long:          `Check that the given file is a valid boot workflow. It will not check that selected module are installed, it will just check that it`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cuecontext.New()
		p := parser.NewParser()
		cueValue, err := helper.CueUnmarshalFile(ctx, checkFlags.filename)
		if err != nil {
			return err
		}

		if err = p.Check(ctx, *cueValue); err != nil {
			return err
		}

		// TODO: log success
		log.Info("everything is fine !")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringVarP(&checkFlags.filename, "filename", "f", "", `the path to the config file you want to check.`)
	checkCmd.MarkFlagFilename("filename", []string{string(helper.JSON), string(helper.YAML), string(helper.YML)}...)
	checkCmd.MarkFlagRequired("filename")
}
