package model

import (
	"time"

	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
)

type File struct {
	ID        uint      `json:"id"`
	RemoteID  uint      `json:"remote_id"`
	Path      string    `json:"path" gorm:"unique"`
	LMtime    time.Time `json:"lmtime" gorm:"column:lmtime"`
	RMtime    time.Time `json:"rmtime" gorm:"column:rmtime"`
	LChecksum string    `json:"lchecksum" gorm:"column:lchecksum"`
	RChecksum string    `json:"rchecksum" gorm:"column:rchecksum"`
	FType     string    `json:"ftype" gorm:"column:ftype"`
	State     string    `json:"state" gorm:"column:state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f File) IsDir() bool {
	return f.FType == mcc.FTypeDirectory
}

func (f File) IsFile() bool {
	return f.FType == mcc.FTypeFile
}
