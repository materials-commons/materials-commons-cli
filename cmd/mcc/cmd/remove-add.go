/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"github.com/spf13/cobra"
)

// removeAddCmd represents the removeAdd command
var removeAddCmd = &cobra.Command{
	Use:   "remove-add",
	Short: "Remove added files(s).",
	Long:  `Remove added files(s).`,
	Args: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if all && len(args) != 0 {
			return fmt.Errorf("args not allowed if --all flag is specified")
		}

		if !all && len(args) == 0 {
			return fmt.Errorf("you must specify files to remove when the --all flag is not used")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		if all {
			removeAllAddedFiles()
			return
		}

		removeSpecifiedFiles(args)
	},
}

func removeAllAddedFiles() {
	addedFileStor := stor.NewGormAddedFileStor(mcdb.MustConnectToDB())
	if err := addedFileStor.RemoveAll(); err != nil {
		log.Fatalf("Error removing all added files: %s", err)
	}
}

func removeSpecifiedFiles(files []string) {
	addedFileStor := stor.NewGormAddedFileStor(mcdb.MustConnectToDB())
	for _, filePath := range files {
		if err := addedFileStor.RemoveByPath(filePath); err != nil {
			log.Printf("Error deleting %q: %s\n", filePath, err)
		}
	}
}

func init() {
	rootCmd.AddCommand(removeAddCmd)
	removeAddCmd.Flags().BoolP("all", "a", false, "Remove all added files")
}
