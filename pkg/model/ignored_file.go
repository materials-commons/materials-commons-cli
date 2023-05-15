package model

import (
	"time"

	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
)

// IgnoredFile is a file that is being marked as ignored, but that the user might want
// to upload later. This is a convenience as this could also be accomplished by updating
// an .mcignore file. However .mcignore are meant as "permanent" entries, while entries
// in this table are provisionally temporary and specific to the local project.
type IgnoredFile struct {
	ID   uint   `json:"id"`
	Path string `json:"path"`

	// FType can be FTypeFile or FTypeDirectory
	FType string `json:"ftype" gorm:"column:ftype"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f IgnoredFile) IsDir() bool {
	return f.FType == mcc.FTypeDirectory
}

func (f IgnoredFile) IsFile() bool {
	return f.FType == mcc.FTypeFile
}
