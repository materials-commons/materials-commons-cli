/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"github.com/materials-commons/materials-commons-cli/pkg/util"
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

	if len(args) != 0 {
		fa.addSpecifiedFiles(args)
		return
	}

	allFlag, _ := addCmd.Flags().GetBool("all")

	if allFlag {
		fa.addFiles(true, true)
		return
	}

	changedFlag, _ := addCmd.Flags().GetBool("changed")
	unknownFlag, _ := addCmd.Flags().GetBool("unknown")

	if changedFlag || unknownFlag {
		fa.addFiles(changedFlag, unknownFlag)
		return
	}
}

type fileAdder struct {
	db              *gorm.DB
	ignoredFileStor stor.IgnoredFileStor
	addedFileStor   stor.AddedFileStor
	fileStor        stor.FileStor
}

func newFileAdder(db *gorm.DB) *fileAdder {
	return &fileAdder{
		db:              db,
		ignoredFileStor: stor.NewGormIgnoredFileStor(db),
		addedFileStor:   stor.NewGormAddedFileStor(db),
		fileStor:        stor.NewGormFileStor(db),
	}
}

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

		projectPath := util.ToProjectPath(fullPath)

		if a.ignoredFileStor.FileIsIgnored(projectPath) {
			continue
		}

		if a.fileIsKnownAndMTimeIsUnchanged(projectPath, fullPath) {
			continue
		}

		// if we are here then this is an unknown file or a file that is known but its MTime
		// changed from what is stored in the database.
		if _, err := a.addedFileStor.AddFile(projectPath); err != nil {
			log.Printf("Unable to add file %q: %s", fullPath, err)
		}
	}
}

func (a *fileAdder) addFiles(changedFiles bool, unknownFiles bool) {
	projectWalker := mcc.NewProjectWalker(a.db, a.changedFileHandler, a.unknownFileHandler)
	if err := projectWalker.Walk(config.GetProjectRootPath()); err != nil {
		log.Fatalf("Unable to add files: %s", err)
	}
}

// fileIsKnownAndMTimeIsUnchanged returns true if the file is in the sqlite database
// and the mtime for the file and what is stored in the database are the same. Otherwise,
// it will return false. There is one special case: if the file is in the database, but
// the stat failed, we return true, meaning the file
func (a *fileAdder) fileIsKnownAndMTimeIsUnchanged(projectPath, path string) bool {
	f, err := a.fileStor.GetFileByPath(projectPath)
	if err != nil {
		// Couldn't retrieve, assume unknown
		return true
	}

	finfo, err := os.Stat(path)
	if err != nil {
		// stat failed, but file exists in database. Print a warning to the user
		// and return false meaning that the file is known, and acting like the
		// mtime is unchanged.
		log.Printf("The file %q does not appear to exist: %s", path, err)
		return false
	}

	if f.LMtime.Before(finfo.ModTime()) {
		// file has newer mtime than what is stored in database, so return false so file can be added
		return false
	}

	// If we are here, then the file exists in the database and on the file system and the mtimes match, so
	// return true signifying that the file shouldn't be added.
	return true
}

func (a *fileAdder) changedFileHandler(projectPath, _ string, _ os.FileInfo) error {
	_, _ = a.addedFileStor.AddFile(projectPath)
	return nil
}

func (a *fileAdder) unknownFileHandler(projectPath, _ string, _ os.FileInfo) error {
	_, _ = a.addedFileStor.AddFile(projectPath)
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("all", "a", false, "Add all unknown and changed files")
	addCmd.Flags().BoolP("unknown", "u", false, "Add all unknown files")
	addCmd.Flags().BoolP("changed", "c", false, "Add all changed files")
}
