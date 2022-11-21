package terraClient

import (
	"os"
	"terraClient/pkg/stateImporter"

	"github.com/spf13/cobra"
)

var sourceFolder string
var destinationFolder string
var updateConfigFiles bool
var importCmd = &cobra.Command{
	Use:   "stateImport",
	Short: "Does import of state",
	Long: `Does import of item from state list in src state into destination state. 
	Fetches id for import from state output. 
	Delete original state at end.
	Error: States can be dependant on eachother. Import does not remove in original state on fail.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(sourceFolder); os.IsNotExist(err) {
			return err
		}

		if _, err := os.Stat(destinationFolder); os.IsNotExist(err) {
			return err
		}

		stateImporter.Execute(sourceFolder, destinationFolder)
		return nil
	},
}

func init() {
	importCmd.PersistentFlags().StringVarP(&sourceFolder, "src", "s", "", "folder with terragruntstate to import from")
	importCmd.PersistentFlags().StringVarP(&destinationFolder, "dest", "d", "", "folder with terragruntstate to import to")
	importCmd.PersistentFlags().BoolVarP(&updateConfigFiles, "update-config", "u", false, "also update config file (default: false)")
	importCmd.MarkPersistentFlagRequired("src")
	importCmd.MarkPersistentFlagRequired("dest")
	rootCmd.AddCommand(importCmd)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
