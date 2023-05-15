package model

import (
	"time"

	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
)

// AddedFile represents a file that is not in the files table that user wants to push, or a file that
// is in the files table but has changed.
type AddedFile struct {
	ID   uint   `json:"id"`
	Path string `json:"path" gorm:"unique"`

	// FType can be FTypeFile or FTypeDirectory
	FType string `json:"ftype" gorm:"column:ftype"`

	// Reason is the reason the file was added. A file can be in the files table and have
	// changed (MTime different, Checksum different), or the file can be unknown, which
	// means it is not in the files table.
	Reason string `json:"reason"`

	// StartedAt will be set when a file is being pushed (uploaded). When the file has been
	// successfully pushed it will be removed from the added_files table. If a user goes
	// to push files and a file has StartedAt set, that means it's push did not complete.
	// This field exists to enable restart on push/upload.
	StartedAt time.Time `json:"started_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f AddedFile) IsChanged() bool {
	switch {
	case f.Reason == mcc.FileMTimeChanged:
		return true
	case f.Reason == mcc.FileChanged:
		return true
	default:
		return false
	}
}

func (f AddedFile) IsUnknown() bool {
	return f.Reason == mcc.FileUnknown
}

func (f AddedFile) IsDir() bool {
	return f.FType == mcc.FTypeDirectory
}

func (f AddedFile) IsFile() bool {
	return f.FType == mcc.FTypeFile
}
