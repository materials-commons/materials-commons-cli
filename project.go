package materials

import (
	"github.com/syndtr/goleveldb/leveldb"
	"path/filepath"
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
	Name        string
	Path        string
	Status      string
	ModTime     time.Time
	MCId        string
	Changes     map[string]ProjectFileChange
	Ignore      []string
	*leveldb.DB `json:"-"`
}

func NewProject(name, path, status string) (*Project, error) {
	p := &Project{
		Name:    name,
		Path:    path,
		Status:  status,
		ModTime: time.Now(),
		Changes: map[string]ProjectFileChange{},
		Ignore:  []string{},
	}

	if err := p.OpenDB(); err != nil {
		return nil, err
	}

	return p, nil
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

func (p *Project) OpenDB() error {
	path := filepath.Join(Config.User.DotMaterialsPath(), "projectdb", p.Name+".db")
	var err error
	p.DB, err = leveldb.OpenFile(path, nil)
	if err == nil {
		p.CompactRange(leveldb.Range{nil, nil})
	}
	return err
}
