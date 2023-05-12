package model

import (
	"time"
)

type Conflict struct {
	ID         uint      `json:"id"`
	RemoteID   uint      `json:"remote_id"`
	FileID     uint      `json:"file_id"`
	File       *File     `json:"file" gorm:"foreignKey:FileID;references:ID"`
	ResolvedAt time.Time `json:"resolved_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
