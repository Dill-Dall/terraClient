package stateImporter

import (
	"terraClient/pkg/stateImporter"

	"github.com/spf13/cobra"
)

var reverseCmd = &cobra.Command{
	Use:   "importfromto",
	Short: "Does import of state",
	//Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stateImporter.Execute()
	},
}

func init() {
	rootCmd.AddCommand(reverseCmd)
}
