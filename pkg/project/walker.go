package project

import (
	"os"
	"path/filepath"

	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"github.com/saracen/walker"
	"gorm.io/gorm"
)

// WalkerHandlerFn is the type definition for the callback functions used by the Walker.
// This methods must be thread safe as they can be called in parallel.
type WalkerHandlerFn func(projectPath, realPath string, finfo os.FileInfo) error

// Walker is a file walker for the project space. It will walk the entire local project
// space calling the ChangedFileHandlerFn and UnknownFileHandlerFn whenever it encounters a
// file or directory that meets the criteria for unknown or changed. Additionally, it will ignore
// files and directories that have been marked as ignored. This can happen because the files
// are marked as ignored in an .mcignore file, or because they are in the database ignored_files
// table. The project walker has an SkipUnknownDirs flag that changes how it handles directories
// that are unknown. By default, this flag is set to true. This means when it encounters an
// unknown directory it will return filepath.SkipDir and not process any of the files or directories
// under that directory. If SkipUnknownDirs is false then it will descend into the directory.
type Walker struct {
	ignoredFileStor        stor.IgnoredFileStor
	fileStor               stor.FileStor
	ChangedFileHandlerFn   WalkerHandlerFn
	UnknownFileHandlerFn   WalkerHandlerFn
	UnchangedFileHandlerFn WalkerHandlerFn
	SkipUnknownDirs        bool
}

func NewWalker(db *gorm.DB) *Walker {
	return &Walker{
		ignoredFileStor:        stor.NewGormIgnoredFileStor(db),
		fileStor:               stor.NewGormFileStor(db),
		ChangedFileHandlerFn:   nil,
		UnchangedFileHandlerFn: nil,
		UnknownFileHandlerFn:   nil,
		SkipUnknownDirs:        true,
	}
}

func (w *Walker) WithChangedFileHandler(changedFileHandlerFn WalkerHandlerFn) *Walker {
	w.ChangedFileHandlerFn = changedFileHandlerFn
	return w
}

func (w *Walker) WithUnchangedFileHandler(unchangedFileHandlerFn WalkerHandlerFn) *Walker {
	w.UnchangedFileHandlerFn = unchangedFileHandlerFn
	return w
}

func (w *Walker) WithUnknownFileHandler(unknownFileHandlerFn WalkerHandlerFn) *Walker {
	w.UnknownFileHandlerFn = unknownFileHandlerFn
	return w
}

func (w *Walker) WithSkipUnknownDirs(skip bool) *Walker {
	w.SkipUnknownDirs = skip
	return w
}

// Walk will walk the directory path. It uses parallel file walker underneath.
func (w *Walker) Walk(path string) error {
	return walker.Walk(path, w.walkCallback, walker.WithErrorCallback(w.walkerErrorCallback))
}

// walkCallback is called by walker.Walk for each file/directory it encounters. This
// method must be thread safe.
func (w *Walker) walkCallback(path string, finfo os.FileInfo) error {
	// If the directory is the <project>/.mc directory then skip it. This directory is
	// where the mcc command stores its metadata and database.
	if finfo.IsDir() && path == config.GetProjectMCDirPath() {
		return filepath.SkipDir
	}

	filename := filepath.Base(path)
	if filename == ".mcignore" {
		// Skip .mcignore files. These are project local and are used to specify what files
		// to skip processing.
		return nil
	}

	// We want two representations of the file. It's full path and its project path. The
	// project path starts with a / (slash), whereas the full path is the local file system
	// full path to the file/directory.
	projectPath := mcc.ToProjectPath(path)

	if w.ignoredFileStor.FileIsIgnored(projectPath) {
		return nil
	}

	f, err := w.fileStor.GetFileByPath(projectPath)
	if err != nil {
		// Error looking up file, assume it's an unknown file
		if w.UnknownFileHandlerFn != nil {
			if err := w.UnknownFileHandlerFn(projectPath, path, finfo); err != nil {
				return err
			}
		}

		if finfo.IsDir() && w.SkipUnknownDirs {
			// If it's an unknown directory, and *is not* the project root then we can skip it.
			if path != config.GetProjectRootPath() {
				return filepath.SkipDir
			}
		}

		return nil
	}

	if finfo.IsDir() {
		// If we are here then this is a **KNOWN** directory. There is nothing
		// we need to do for the directory.
		return nil
	}

	if w.fileMTimeIsChanged(f, finfo) {
		if w.ChangedFileHandlerFn != nil {
			return w.ChangedFileHandlerFn(projectPath, path, finfo)
		}
		return nil
	} else {
		if w.UnchangedFileHandlerFn != nil {
			return w.UnchangedFileHandlerFn(projectPath, path, finfo)
		}
		return nil
	}
}

// walkerErrorCallback is called whenever the parallel walker encounters an error.
// We skip permission errors, and only return an error that will cause walking to
// stop if it's not a permission error.
func (w *Walker) walkerErrorCallback(_ string, err error) error {
	if os.IsPermission(err) {
		return nil
	}

	// Halt on any other error
	return err
}

// fileMTimeIsChanged compares the mtime from the database with the current file system mtime.
// If these are different then the file has potentially changed. Potentially means that
// a determination of if it has changed can only be made by seeing if the sizes are
// different, or if they are the same, if the checksums have changed. We leave this
// determination to the callback for changed files to give them flexibility in how
// to handle this.
func (w *Walker) fileMTimeIsChanged(f *model.File, finfo os.FileInfo) bool {
	if f.LMTime.Before(finfo.ModTime()) {
		// If the Local MTime we have for this file is before the MTime in the file system
		// then the file has potentially changed. We will only know for sure by computing
		// the checksum, but for now we can just set this file as changed and a candidate
		// to upload.
		return true
	}

	return false
}
