/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/spf13/cobra"
)

// resolveCmd represents the resolve command
var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve a conflict by marking a conflict as able to be overwritten.",
	Long:  `Resolve a conflict by marking a conflict as able to be overwritten.`,
	Run:   runResolveCmd,
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
	db := mcdb.MustConnectToDB()
	_ = db
}

func resolveSpecifiedConflicts(paths []string) {

}

func init() {
	rootCmd.AddCommand(resolveCmd)
	resolveCmd.Flags().BoolP("all", "a", false, "Mark all conflicts as resolved")
}
