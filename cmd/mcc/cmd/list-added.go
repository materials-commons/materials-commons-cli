/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listAddedCmd represents the listAdded command
var listAddedCmd = &cobra.Command{
	Use:   "list-added",
	Short: "List added files.",
	Long:  `List added files.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listAdded called")
	},
}

func init() {
	rootCmd.AddCommand(listAddedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listAddedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listAddedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
