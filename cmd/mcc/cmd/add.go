/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/project"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add unknown and/or changed files to be uploaded.",
	Long:  `Add unknown and/or changed files to be uploaded.`,
	Run:   runAddCmd,
}

func runAddCmd(cmd *cobra.Command, args []string) {

	db := mcdb.MustConnectToDB()

	fa := newFileAdder(db)

	// If there are any args then the user is adding specific files
	if len(args) != 0 {
		fa.addSpecifiedFiles(args)
		return
	}

	// If we are here then the user specified types of files to add - unknown, changed or both

	allFlag, _ := cmd.Flags().GetBool("all")

	if allFlag {
		fa.addFiles(true, true)
		return
	}

	changedFlag, _ := cmd.Flags().GetBool("changed")
	unknownFlag, _ := cmd.Flags().GetBool("unknown")

	if changedFlag || unknownFlag {
		fa.addFiles(changedFlag, unknownFlag)
		return
	}
}

// fileAdder handles adding files. These are either specific files the user
// has specified, or files that the project walker has identified.
type fileAdder struct {
	db                   *gorm.DB
	ignoredFileStor      stor.IgnoredFileStor
	addedFileStor        stor.AddedFileStor
	fileStor             stor.FileStor
	fileStatusDeterminer *project.FileStatusDeterminer
}

func newFileAdder(db *gorm.DB) *fileAdder {
	return &fileAdder{
		db:                   db,
		ignoredFileStor:      stor.NewGormIgnoredFileStor(db),
		addedFileStor:        stor.NewGormAddedFileStor(db),
		fileStor:             stor.NewGormFileStor(db),
		fileStatusDeterminer: project.NewFileStatusDeterminer(db),
	}
}

// addSpecifiedFiles adds files the user has passed in as args. It checks each of the files to
// make sure they are either changed or unknown. If a file is ignored then it is not added.
func (a *fileAdder) addSpecifiedFiles(args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to determine current directory: %s", err)
	}

	for _, filePath := range args {
		fullPath := filePath
		if fullPath[0] != '/' {
			// if fullPath doesn't start with '/' then it's a relative path. Turn it
			// into a full path.
			fullPath = filepath.Join(cwd, fullPath)
		}

		projectPath := mcc.ToProjectPath(fullPath)

		fileState, ftype := a.fileStatusDeterminer.DetermineFileStatus(projectPath, fullPath)
		switch fileState {
		case mcc.FileAlreadyAdded:
			// Nothing to do

		case mcc.FileIgnored:
			// Nothing to do

		case mcc.FileKnownAndUnchanged:
			// Nothing to do

		case mcc.FileMTimeChanged:
			fmt.Printf("Adding changed file: %q\n", fullPath)

			if _, err := a.addedFileStor.AddFile(projectPath, mcc.FileChanged, ftype); err != nil {
				log.Printf("Unable to add file %q: %s", fullPath, err)
			}

		case mcc.FileUnknown:
			fmt.Printf("Adding unknown file: %q\n", fullPath)
			if _, err := a.addedFileStor.AddFile(projectPath, mcc.FileUnknown, ftype); err != nil {
				log.Printf("Unable to add file %q: %s", fullPath, err)
			}

		case mcc.FileMissing:
			log.Printf("File %q is in the project database, but appears to be deleted\n", fullPath)

		default:
			// Shouldn't happen - nothing to do
		}
	}
}

func (a *fileAdder) addFiles(changedFiles bool, unknownFiles bool) {
	var (
		changedFileHandler project.WalkerHandlerFn = nil
		unknownFileHandler project.WalkerHandlerFn = nil
	)

	if changedFiles {
		changedFileHandler = a.changedFileHandler
	}

	if unknownFiles {
		unknownFileHandler = a.unknownFileHandler
	}

	projectWalker := project.NewWalker(a.db).
		WithChangedFileHandler(changedFileHandler).
		WithUnknownFileHandler(unknownFileHandler)
	if err := projectWalker.Walk(config.GetProjectRootPath()); err != nil {
		log.Fatalf("Unable to add files: %s", err)
	}
}

func (a *fileAdder) changedFileHandler(projectPath, path string, finfo os.FileInfo) error {
	fmt.Printf("Adding changed file %q\n", path)

	ftype := mcc.FTypeFile
	if finfo.IsDir() {
		ftype = mcc.FTypeDirectory
	}
	_, _ = a.addedFileStor.AddFile(projectPath, mcc.FileChanged, ftype)

	return nil
}

func (a *fileAdder) unknownFileHandler(projectPath, path string, finfo os.FileInfo) error {
	fmt.Printf("Adding unknown file %q\n", path)

	ftype := mcc.FTypeFile
	if finfo.IsDir() {
		ftype = mcc.FTypeDirectory
	}
	_, _ = a.addedFileStor.AddFile(projectPath, mcc.FileUnknown, ftype)

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("all", "a", false, "Add all unknown and changed files")
	addCmd.Flags().BoolP("unknown", "u", false, "Add all unknown files")
	addCmd.Flags().BoolP("changed", "c", false, "Add all changed files")
}
