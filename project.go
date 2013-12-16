package materials

import (
	"time"
)

// ProjectFileChange contains information about the change that
// occurred to a file in a project.
type ProjectFileChange struct {
	Path string
	Type string
	When time.Time
	Hash []byte
}

// Project describes the information we track about a users
// projects. Here we keep the name of the project and the
// directory path. The name of the project is the top level
// directory of the project. The path is the full path to
// the project excluding the name.
type Project struct {
	Name    string
	Path    string
	Status  string
	ModTime time.Time
	MCId    string
	Changes map[string]ProjectFileChange
	Ignore  []string
}
