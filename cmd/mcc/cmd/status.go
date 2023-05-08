/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"github.com/materials-commons/materials-commons-cli/pkg/util"
	"github.com/saracen/walker"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Gets the status of local files and identifies what should be uploaded.",
	Long:  `Gets the status of local files and identifies what should be uploaded.`,
	Run: func(cmd *cobra.Command, args []string) {
		db := mcdb.MustConnectToDB()
		sw := newStatusWalker(db)
		if err := sw.Walk(config.GetProjectRootPath()); err != nil {
			log.Fatalf("Unable to get status: %s", err)
		}

		sw.ShowStatus()
	},
}

type statusWalker struct {
	ignoredFileStor stor.IgnoredFileStor
	fileStor        stor.FileStor
	unknownFiles    sync.Map
	changedFiles    sync.Map
}

func newStatusWalker(db *gorm.DB) *statusWalker {
	return &statusWalker{
		ignoredFileStor: stor.NewGormIgnoredFileStor(db),
		fileStor:        stor.NewGormFileStor(db),
	}
}

func (w *statusWalker) Walk(path string) error {
	return walker.Walk(path, w.walkCallback, walker.WithErrorCallback(w.walkerErrorCallback))
}

func (w *statusWalker) walkCallback(path string, finfo os.FileInfo) error {
	if finfo.IsDir() && path == config.GetProjectMCDirPath() {
		return filepath.SkipDir
	}

	filename := filepath.Base(path)
	if filename == ".mcignore" {
		return nil
	}

	projectPath := util.ToProjectPath(path)

	if w.ignoredFileStor.FileIsIgnored(projectPath) {
		return nil
	}

	f, err := w.fileStor.GetFileByPath(projectPath)
	if err != nil {
		// Error looking up file, assume it's an unknown file
		w.unknownFiles.Store(projectPath, path)
		return nil
	}

	if finfo.IsDir() {
		// If we are here then this is a known directory. There is nothing
		// we need to do for the directory.
		return nil
	}

	if !fileIsChanged(f, finfo) {
		return nil
	}

	// If we are here then the file is known and the mtime indicates its probably
	// changed so add it to the changed file list
	w.changedFiles.Store(projectPath, path)

	return nil
}

func (w *statusWalker) walkerErrorCallback(path string, err error) error {
	if os.IsPermission(err) {
		return nil
	}

	// Halt on any other error
	return err
}

func (w *statusWalker) ShowStatus() {
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

func fileIsChanged(f *model.File, finfo os.FileInfo) bool {
	if f.LMtime.Before(finfo.ModTime()) {
		// If the Local MTime we have for this file is before the MTime in the file system
		// then the file has potentially changed. We will only know for sure by computing
		// the checksum, but for now we can just set this file as changed and a candidate
		// to upload.
		return true
	}

	return false
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
