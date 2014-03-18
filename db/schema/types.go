package schema

import (
	"time"
)

// A Project is an instance of a users project.
type Project struct {
	ID   int
	Name string
	Path string
	MCID string `db:"mcid"`
}

// A ProjectEvent is a file change event in the project.
type ProjectEvent struct {
	ID        int
	Path      string
	Event     string
	EventTime time.Time `db:"event_time"`
	ProjectID int       `db:"project_id"`
}

// A DataDir is a directory in a project.
type DataDir struct {
	ID         int
	ProjectID  int    `db:"project_id"`
	MCID       string `db:"mcid"`
	Name       string
	Path       string
	ParentMCID string `db:"parent_mcid"`
	Parent     int
}

// A DataFile is a file in a directory in a project.
type DataFile struct {
	ID         int
	MCID       string    `db:"mcid"`
	DataDirID  int       `db:"datadir_id"`
	ProjectID  int       `db:"project_id"`
	LastUpload time.Time `db:"last_upload"`
	MTime      time.Time `db:"mtime"`
	ParentMCID string    `db:"parent_mcid"`
	Parent     int
	Name       string
	Path       string
	Size       int
	Checksum   string
	Version    int
}
