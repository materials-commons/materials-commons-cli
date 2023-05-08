/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Gets the status of local files and identifies what should be uploaded.",
	Long:  `Gets the status of local files and identifies what should be uploaded.`,
	Run: func(cmd *cobra.Command, args []string) {
		db := mcdb.MustConnectToDB()
		sw := &statusWalkerState{}
		projectWalker := mcc.NewProjectWalker(db, sw.changedFileHandler, sw.unknownFileHandler)
		if err := projectWalker.Walk(config.GetProjectRootPath()); err != nil {
			log.Fatalf("Unable to get status: %s", err)
		}

		sw.ShowStatus()
	},
}

type statusWalkerState struct {
	unknownFiles sync.Map
	changedFiles sync.Map
}

func (w *statusWalkerState) changedFileHandler(projectPath, path string, _ os.FileInfo) error {
	w.changedFiles.Store(projectPath, path)
	return nil
}

func (w *statusWalkerState) unknownFileHandler(projectPath, path string, _ os.FileInfo) error {
	w.unknownFiles.Store(projectPath, path)
	return nil
}

func (w *statusWalkerState) ShowStatus() {
	fmt.Println("Unknown:")
	w.unknownFiles.Range(func(key any, value any) bool {
		path := value.(string)
		fmt.Println(path)
		return true
	})

	fmt.Println("Changed:")
	w.changedFiles.Range(func(key any, value any) bool {
		path := value.(string)
		fmt.Println(path)
		return true
	})
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
