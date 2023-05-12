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

// listAddedCmd represents the listAdded command
var listAddedCmd = &cobra.Command{
	Use:   "list-added",
	Short: "List added files.",
	Long:  `List added files.`,
	Run:   runListAddedCmd,
}

func runListAddedCmd(cmd *cobra.Command, args []string) {
	unknownFlag, _ := cmd.Flags().GetBool("unknown")
	changedFlag, _ := cmd.Flags().GetBool("changed")

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
	addedFileStor := stor.NewGormAddedFileStor(mcdb.MustConnectToDB())
	err := addedFileStor.ListPaged(func(f *model.AddedFile) error {
		fmt.Printf("%s %s (%s)\n", f.Reason, mcc.ToFullPath(f.Path), f.FType)
		return nil
	})

	if err != nil {
		log.Fatalf("Error retrieving added files: %s", err)
	}
}

func showStatusSpecificFiles(paths []string) {
	var (
		err error
		f   *model.AddedFile
	)

	addedFileStor := stor.NewGormAddedFileStor(mcdb.MustConnectToDB())

	for _, p := range paths {
		projectPath := mcc.ToProjectPath(p)

		if f, err = addedFileStor.GetFileByPath(projectPath); err != nil {
			fmt.Printf("%s not in added files\n", p)
			continue
		}

		fmt.Printf("%s %s (%s)\n", f.Reason, mcc.ToFullPath(f.Path), f.FType)
	}
}

func showStatusForAddedFileByReason(unknown bool, changed bool) {
	addedFileStor := stor.NewGormAddedFileStor(mcdb.MustConnectToDB())
	err := addedFileStor.ListPaged(func(f *model.AddedFile) error {
		switch {
		case unknown && f.IsUnknown():
			fmt.Printf("%s %s (%s)\n", f.Reason, mcc.ToFullPath(f.Path), f.FType)

		case changed && f.IsChanged():
			fmt.Printf("%s %s (%s)\n", f.Reason, mcc.ToFullPath(f.Path), f.FType)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error retrieving added files: %s", err)
	}
}

func init() {
	rootCmd.AddCommand(listAddedCmd)
	listAddedCmd.Flags().BoolP("unknown", "u", false, "Add all unknown files")
	listAddedCmd.Flags().BoolP("changed", "c", false, "Add all changed files")
}
