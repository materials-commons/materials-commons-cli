package mcc

import (
	"os"
	"path/filepath"

	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"github.com/materials-commons/materials-commons-cli/pkg/util"
	"github.com/saracen/walker"
	"gorm.io/gorm"
)

type ProjectWalkerHandlerFn func(projectPath, realPath string, finfo os.FileInfo) error

type ProjectWalker struct {
	ignoredFileStor      stor.IgnoredFileStor
	fileStor             stor.FileStor
	ChangedFileHandlerFn ProjectWalkerHandlerFn
	UnknownFileHandlerFn ProjectWalkerHandlerFn
	SkipUnknownDirs      bool
}

func NewProjectWalker(db *gorm.DB, changedFileHandlerFn, unknownFileHandlerFn ProjectWalkerHandlerFn) *ProjectWalker {
	return &ProjectWalker{
		ignoredFileStor:      stor.NewGormIgnoredFileStor(db),
		fileStor:             stor.NewGormFileStor(db),
		ChangedFileHandlerFn: changedFileHandlerFn,
		UnknownFileHandlerFn: unknownFileHandlerFn,
		SkipUnknownDirs:      true,
	}
}

func (w *ProjectWalker) Walk(path string) error {
	return walker.Walk(path, w.walkCallback, walker.WithErrorCallback(w.walkerErrorCallback))
}

func (w *ProjectWalker) walkCallback(path string, finfo os.FileInfo) error {
	if finfo.IsDir() && path == config.GetProjectMCDirPath() {
		return filepath.SkipDir
	}

	filename := filepath.Base(path)
	if filename == ".mcignore" {
		// Skip .mcignore files. These are project local and are used to specify what files
		// to skip processing.
		return nil
	}

	projectPath := util.ToProjectPath(path)

	if w.ignoredFileStor.FileIsIgnored(projectPath) {
		return nil
	}

	f, err := w.fileStor.GetFileByPath(projectPath)
	if err != nil {
		// Error looking up file, assume it's an unknown file
		if err := w.UnknownFileHandlerFn(projectPath, path, finfo); err != nil {
			// do something
		}

		if finfo.IsDir() && w.SkipUnknownDirs {
			// If it's an unknown directory, and is not the project root then we can skip it.
			if path != config.GetProjectRootPath() {
				return filepath.SkipDir
			}
		}

		return nil
	}

	if finfo.IsDir() {
		// If we are here then this is a known directory. There is nothing
		// we need to do for the directory.
		return nil
	}

	if w.fileIsChanged(f, finfo) {
		if err := w.ChangedFileHandlerFn(projectPath, path, finfo); err != nil {
			// do something
		}
		return nil
	}

	return nil
}

func (w *ProjectWalker) walkerErrorCallback(path string, err error) error {
	if os.IsPermission(err) {
		return nil
	}

	// Halt on any other error
	return err
}

func (w *ProjectWalker) fileIsChanged(f *model.File, finfo os.FileInfo) bool {
	if f.LMtime.Before(finfo.ModTime()) {
		// If the Local MTime we have for this file is before the MTime in the file system
		// then the file has potentially changed. We will only know for sure by computing
		// the checksum, but for now we can just set this file as changed and a candidate
		// to upload.
		return true
	}

	return false
}
