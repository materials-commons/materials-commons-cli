package schema

import (
	"time"
)

// A Project is an instance of a users project.
type Project struct {
	ID   int    // Primary key
	Name string // Name of project
	Path string // Path to project
	MCID string `db:"mcid"` // Materials Commons id for project
}

// A ProjectEvent is a file change event in the project.
type ProjectEvent struct {
	ID        int       // Primary key
	ProjectID int       `db:"project_id"` // Foreign key to project
	Path      string    // Path of file/directory this event pertains to
	Event     string    // Type of event
	EventTime time.Time `db:"event_time"` // Time event occurred
}

// A ProjectFile is a file or directory in the project
type ProjectFile struct {
	ID        int       // Primary key
	ProjectID int       `db:"project_id"` // Foreign key to project
	Path      string    // Full path to file/directory
	Size      int64     // Size of file (valid only for files)
	Checksum  string    // MD5 Hash of file (valid only for files)
	MTime     time.Time // Last known Modification time
	IsDir     bool      // True if this entry is a directory
}
