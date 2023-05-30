package cmd

import (
	"log"

	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls (downloads) files from the server.",
	Long:  `Pulls (downloads) files from the server.`,
	Run:   runPullCmd,
}

func runPullCmd(cmd *cobra.Command, args []string) {
	remoteStor := stor.MustLoadJsonRemoteStor()
	defaultRemote, err := remoteStor.GetDefaultRemote()
	if err != nil {
		log.Fatalf("No default remote set: %s", err)
	}

	_ = defaultRemote

	db := mcdb.MustConnectToDB()
	projectStor := stor.NewGormProjectStor(db)
	p, err := projectStor.GetProject()
	if err != nil {
		log.Fatalf("Unable to retrieve project: %s", err)
	}

	_ = p

	if len(args) != 0 {
		pullSpecificFiles()
		return
	}

	pullDownloadedDirs()
}

func pullSpecificFiles() {

}

func pullDownloadedDirs() {

}

func init() {
	rootCmd.AddCommand(pullCmd)
}
