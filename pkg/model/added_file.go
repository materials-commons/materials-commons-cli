package model

import (
	"time"
)

type AddedFile struct {
	ID        uint      `json:"id"`
	Path      string    `json:"path" gorm:"unique"`
	Reason    string    `json:"reason"`
	StartedAt time.Time `json:"started_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
