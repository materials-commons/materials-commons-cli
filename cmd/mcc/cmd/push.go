/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/project"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushes (uploads) added files to the server.",
	Long:  `Pushes (uploads) added files to the server.`,
	Run:   runPushCmd,
}

func runPushCmd(cmd *cobra.Command, args []string) {
	db := mcdb.MustConnectToDB()

	projectWalker := project.NewWalker(db, nil, nil)
	if err := projectWalker.Walk(config.GetProjectRootPath()); err != nil {

	}
}

type uploader struct {
}

func newUploader() *uploader {
	return &uploader{}
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
