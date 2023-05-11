package mcc

import (
	"os"

	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"gorm.io/gorm"
)

type FileStatusDeterminer struct {
	ignoredFileStor stor.IgnoredFileStor
	fileStor        stor.FileStor
}

func NewFileStatusDeterminer(db *gorm.DB) *FileStatusDeterminer {
	return &FileStatusDeterminer{
		ignoredFileStor: stor.NewGormIgnoredFileStor(db),
		fileStor:        stor.NewGormFileStor(db),
	}
}

// DetermineFileStatus returns the status of the file. The states are FileUnknown,
// FileIgnored, FileMissing, FileMTimeChanged, FileKnownAndUnchanged.
func (d *FileStatusDeterminer) DetermineFileStatus(projectPath, path string) string {
	var (
		err   error
		f     *model.File
		finfo os.FileInfo
	)

	if d.ignoredFileStor.FileIsIgnored(projectPath) {
		return FileIgnored
	}

	if f, err = d.fileStor.GetFileByPath(projectPath); err != nil {
		// Couldn't retrieve, assume unknown
		return FileUnknown
	}

	if finfo, err = os.Stat(path); err != nil {
		// stat failed, but file exists in database.
		return FileMissing
	}

	if f.LMtime.Before(finfo.ModTime()) {
		// file has newer mtime than what is stored in database, so return FileChanged
		return FileMTimeChanged
	}

	// If we are here, then the file exists in the database and on the file system and the mtimes match, so
	// the file is both known and hasn't changed.
	return FileKnownAndUnchanged
}
