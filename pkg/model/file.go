package model

import (
	"time"

	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
)

// File represents the known files. Known files are files that are on the server that have
// been downloaded to the local project. These files are tracked in the database. This
// allows the system to learn about new and changed files. This is used to identify files
// that can be uploaded, or changes that would be overwritten on a download.
type File struct {
	ID       uint   `json:"id"`
	RemoteID uint   `json:"remote_id"`
	Path     string `json:"path" gorm:"unique"`

	// LMTime is the local file MTime
	LMTime time.Time `json:"lmtime" gorm:"column:lmtime"`

	// RMTime is the MTime for the last known version on the server
	RMTime time.Time `json:"rmtime" gorm:"column:rmtime"`

	// LChecksum is the MD5 checksum for the local file
	LChecksum string `json:"lchecksum" gorm:"column:lchecksum"`

	// RChecksum is the MD5 checksum for the last known version on the server
	RChecksum string `json:"rchecksum" gorm:"column:rchecksum"`

	// FType FTypeFile or FTypeDirectory
	FType string `json:"ftype" gorm:"column:ftype"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f File) IsDir() bool {
	return f.FType == mcc.FTypeDirectory
}

func (f File) IsFile() bool {
	return f.FType == mcc.FTypeFile
}
