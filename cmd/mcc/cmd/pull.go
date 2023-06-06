package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/project"
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

var mu sync.Mutex

func pullDownloadedDirs() {
	changedFileHandler := func(projectPath, path string, finfo os.FileInfo) error {
		mu.Lock()
		defer mu.Unlock()
		fmt.Printf("changed file: %s\n", path)
		return nil
	}

	unknownFileHandler := func(projectPath, path string, finfo os.FileInfo) error {
		mu.Lock()
		defer mu.Unlock()
		fmt.Printf("unknown file: %s\n", path)
		return nil
	}

	unchangedFileHandler := func(projectPath, path string, finfo os.FileInfo) error {
		mu.Lock()
		defer mu.Unlock()
		fmt.Printf("unchanged file: %s\n", path)
		return nil
	}
	db := mcdb.MustConnectToDB()

	//projectStor := stor.NewGormProjectStor(db)
	//p, err := projectStor.GetProject()
	//c := mcapi.NewClient("", "")
	//files, err := c.ListDirectoryByPath(int(p.ID), "/")
	//_ = files
	//
	//if err != nil {
	//	log.Fatalf("Unable")
	//}

	projectWalker := project.NewWalker(db).
		WithChangedFileHandler(changedFileHandler).
		WithUnknownFileHandler(unknownFileHandler).
		WithUnchangedFileHandler(unchangedFileHandler).
		WithSkipUnknownDirs(false)
	if err := projectWalker.Walk(config.GetProjectRootPath()); err != nil {
		log.Fatalf("Unable to add files: %s", err)
	}
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
