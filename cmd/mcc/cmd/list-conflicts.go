package cmd

import (
	"fmt"

	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/model"
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
	db := mcdb.MustConnectToDB()
	var conflictFiles []model.Conflict

	offset := 0
	pageSize := 100
	for {
		if err := db.Offset(offset).Limit(pageSize).Preload("Files").Find(&conflictFiles).Error; err != nil {
			break
		}

		if len(conflictFiles) == 0 {
			break
		}

		for _, f := range conflictFiles {
			fmt.Printf("%s\n", mcc.ToFullPath(f.File.Path))
		}
		offset = offset + pageSize
	}
}

func init() {
	rootCmd.AddCommand(listConflictsCmd)
}
