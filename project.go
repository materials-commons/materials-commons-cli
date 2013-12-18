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
}

// Project describes the information we track about a users
// projects. Here we keep the name of the project and the
// directory path. The name of the project is the top level
// directory of the project. The path is the full path to
// the project including the name (top level directory).
type Project struct {
	Name    string
	Path    string
	Status  string
	ModTime time.Time
	MCId    string
	Changes map[string]ProjectFileChange
	Ignore  []string
}

func (p *Project) AddFileChange(fileChange ProjectFileChange) {
	entry, found := p.Changes[fileChange.Path]
	switch {
	case found:
		entry.Type = fileChange.Type
		entry.When = fileChange.When
		p.Changes[entry.Path] = entry
	case !found:
		p.Changes[fileChange.Path] = fileChange
	}
}

func (p *Project) RemoveFileChange(path string) {
	delete(p.Changes, path)
}
