package model

import (
	"time"

	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
)

type AddedFile struct {
	ID        uint      `json:"id"`
	Path      string    `json:"path" gorm:"unique"`
	FType     string    `json:"ftype" gorm:"column:ftype"`
	Reason    string    `json:"reason"`
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
