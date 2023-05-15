package model

import (
	"time"
)

// Conflict represents a known file (file in the files tables) that has changed on the local host and on the server,
// where the server download will overwrite the local file. This entry is removed once a conflict has been fully
// resolved by either uploading
type Conflict struct {
	ID uint `json:"id"`

	// The id of the file on the server
	RemoteID uint `json:"remote_id"`

	// The local file id for the file in conflict
	FileID uint `json:"file_id"`

	// The file referenced by FileID
	File *File `json:"file" gorm:"foreignKey:FileID;references:ID"`

	// If a conflict has been marked as resolved, then ResolvedAt will be set. This happens
	// when the user decides that they will allow the file to be overwritten. Conflicts are
	// removed either at push time (if the ResolvedAt is field is set, once the file has been
	// uploaded the conflict is removed from the database), or a pull time (if the ResolvedAt
	// field is set, then once the file has been replaced during the pull the conflict will
	// be removed from the database).
	ResolvedAt time.Time `json:"resolved_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
