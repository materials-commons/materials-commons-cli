package cmd

import (
	"fmt"

	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// listAddedCmd represents the listAdded command
var listAddedCmd = &cobra.Command{
	Use:   "list-added",
	Short: "List added files.",
	Long:  `List added files.`,
	Run:   runListAddedCmd,
}

var listAddedFlags *pflag.FlagSet

func runListAddedCmd(cmd *cobra.Command, args []string) {
	unknownFlag, _ := listAddedFlags.GetBool("unknown")
	changedFlag, _ := listAddedFlags.GetBool("changed")

	if !unknownFlag && !changedFlag && len(args) == 0 {
		// If neither flag is specified, then the command was run without any flags. We treat
		// this as if the user wants to see all added file types.
		showStatusAllAddedFiles()
		return
	}

	if len(args) > 0 {
		// The user has asked to see the status on specific files
		showStatusSpecificFiles(args)
		return
	}

	// If we are here then the user has specified at least one flag
	showStatusForAddedFileByReason(unknownFlag, changedFlag)
}

func showStatusAllAddedFiles() {
	fmt.Println("showStatusAllAddedFiles")
	db := mcdb.MustConnectToDB()

	var addedFiles []model.AddedFile
	offset := 0
	pageSize := 100
	for {
		if err := db.Offset(offset).Limit(pageSize).Find(&addedFiles).Error; err != nil {
			break
		}

		if len(addedFiles) == 0 {
			break
		}

		for _, f := range addedFiles {
			fmt.Printf("%s %s (%s)", f.Reason, f.Path, f.FType)
		}
		offset = offset + pageSize
	}
}

func showStatusSpecificFiles(paths []string) {
	db := mcdb.MustConnectToDB()
	addedFileStor := stor.NewGormAddedFileStor(db)

	var (
		err error
		f   *model.AddedFile
	)

	for _, p := range paths {
		projectPath := mcc.ToProjectPath(p)

		if f, err = addedFileStor.GetFileByPath(projectPath); err != nil {
			fmt.Printf("%s not in added files\n", p)
			continue
		}

		fmt.Printf("%s %s (%s)", f.Reason, f.Path, f.FType)
	}
}

func showStatusForAddedFileByReason(unknown bool, changed bool) {
	db := mcdb.MustConnectToDB()

	var addedFiles []model.AddedFile
	offset := 0
	pageSize := 100
	for {
		if err := db.Offset(offset).Limit(pageSize).Find(&addedFiles).Error; err != nil {
			break
		}

		if len(addedFiles) == 0 {
			break
		}

		for _, f := range addedFiles {
			switch {
			case unknown && f.IsUnknown():
				fmt.Printf("%s %s (%s)", f.Reason, f.Path, f.FType)

			case changed && f.IsChanged():
				fmt.Printf("%s %s (%s)", f.Reason, f.Path, f.FType)
			}
		}
		offset = offset + pageSize
	}
}

func init() {
	rootCmd.AddCommand(listAddedCmd)
	listAddedCmd.Flags().BoolP("unknown", "u", false, "Add all unknown files")
	listAddedCmd.Flags().BoolP("changed", "c", false, "Add all changed files")

	listAddedFlags = listAddedCmd.Flags()
}
