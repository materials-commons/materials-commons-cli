package project

import (
	"os"

	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"gorm.io/gorm"
)

type FileStatusDeterminer struct {
	ignoredFileStor stor.IgnoredFileStor
	fileStor        stor.FileStor
	addedFileStor   stor.AddedFileStor
}

func NewFileStatusDeterminer(db *gorm.DB) *FileStatusDeterminer {
	return &FileStatusDeterminer{
		ignoredFileStor: stor.NewGormIgnoredFileStor(db),
		fileStor:        stor.NewGormFileStor(db),
		addedFileStor:   stor.NewGormAddedFileStor(db),
	}
}

// DetermineFileStatus returns the status of the file. The states are FileUnknown,
// FileIgnored, FileMissing, FileMTimeChanged, FileKnownAndUnchanged.
func (d *FileStatusDeterminer) DetermineFileStatus(projectPath, path string) (string, string) {
	var (
		err   error
		f     *model.File
		finfo os.FileInfo
	)

	if af, err := d.addedFileStor.GetFileByPath(projectPath); err == nil {
		// We found a file matching path in the added_files table
		return mcc.FileAlreadyAdded, af.FType
	}

	if d.ignoredFileStor.FileIsIgnored(projectPath) {
		return mcc.FileIgnored, mcc.FTypeUnknown
	}

	if finfo, err = os.Stat(path); err != nil {
		// stat failed, but file exists in database.
		return mcc.FileMissing, mcc.FTypeUnknown
	}

	ftype := mcc.FTypeFile
	if finfo.IsDir() {
		ftype = mcc.FTypeDirectory
	}

	if f, err = d.fileStor.GetFileByPath(projectPath); err != nil {
		// Couldn't retrieve, assume unknown
		return mcc.FileUnknown, ftype
	}

	if f.LMTime.Before(finfo.ModTime()) {
		// file has newer mtime than what is stored in database, so return FileChanged
		return mcc.FileMTimeChanged, f.FType
	}

	// If we are here, then the file exists in the database and on the file system and the mtimes match, so
	// the file is both known and hasn't changed.
	return mcc.FileKnownAndUnchanged, f.FType
}
