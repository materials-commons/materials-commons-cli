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

// resolveCmd represents the resolve command
var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve a conflict by marking a conflict as able to be overwritten.",
	Long:  `Resolve a conflict by marking a conflict as able to be overwritten.`,
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
	Run: runResolveCmd,
}

func runResolveCmd(cmd *cobra.Command, args []string) {
	allFlag, _ := cmd.Flags().GetBool("all")
	if allFlag {
		resolveAllConflicts()
		return
	}

	if len(args) != 0 {
		resolveSpecifiedConflicts(args)
	}
}

func resolveAllConflicts() {
	conflictStor := stor.NewGormConflictStor(mcdb.MustConnectToDB())
	if err := conflictStor.ResolveAllConflicts(); err != nil {
		log.Printf("Error resolving all conflicts: %s\n", err)
	}
}

func resolveSpecifiedConflicts(paths []string) {
	conflictStor := stor.NewGormConflictStor(mcdb.MustConnectToDB())

	for _, p := range paths {
		if err := conflictStor.ResolveConflictByPath(p); err != nil {
			log.Printf("Error resolving conflict for %q: %s\n", p, err)
		}
	}
}

func init() {
	rootCmd.AddCommand(resolveCmd)
	resolveCmd.Flags().BoolP("all", "a", false, "Mark all conflicts as resolved")
}
