package cmd

import (
	"fmt"
	"log"

	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"github.com/spf13/cobra"
)

// listConflictsCmd represents the conflicts command
var listConflictsCmd = &cobra.Command{
	Use:   "list-conflicts",
	Short: "Lists all files that were changed locally that pull would have overwritten.",
	Long:  `Lists all files that were changed locally that pull would have overwritten.`,
	Run:   runListConflictsCmd,
}

func runListConflictsCmd(cmd *cobra.Command, args []string) {
	conflictStor := stor.NewGormConflictStor(mcdb.MustConnectToDB())
	err := conflictStor.ListPaged(func(f *model.Conflict) error {
		fmt.Printf("%s\n", mcc.ToFullPath(f.File.Path))
		return nil
	})

	if err != nil {
		log.Fatalf("Error retrieving conflicts: %s", err)
	}
}

func init() {
	rootCmd.AddCommand(listConflictsCmd)
}
