package model

import (
	"time"
)

type IgnoredFile struct {
	ID        uint      `json:"id"`
	Path      string    `json:"path"`
	IsDir     bool      `json:"is_dir"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
