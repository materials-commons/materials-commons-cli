package materials

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/materials-commons/gohandy/handyfile"
	"io/ioutil"
	"os"
	"path/filepath"
)

// MaterialsProjects contains a list of user projects and information that
// is needed by the methods to load the projects file.
type ProjectDB struct {
	path     string
	projects []Project
}

// CurrentUserProjects opens the project database for a user contained in
// $HOME/.materials/projects
func CurrentUserProjects() (*ProjectDB, error) {
	projectsPath := filepath.Join(Config.User.DotMaterialsPath(), "projectdb")
	return OpenProjectDB(projectsPath)
}

// Load projects from the database directory at path. Project files are
// JSON files ending with a .project extension.
func OpenProjectDB(path string) (*ProjectDB, error) {
	projectDB := ProjectDB{path: path}
	err := projectDB.loadProjects()
	if err != nil {
		return nil, err
	}
	return &projectDB, err
}

// Reload re-reads and loads the projects file.
func (p *ProjectDB) Reload() error {
	return p.loadProjects()
}

// loadProjects reads the projects directory, and loads each *.project file found in it.
func (p *ProjectDB) loadProjects() error {
	if !handyfile.IsDir(p.path) {
		return fmt.Errorf("ProjectDB must be a directory: '%s'", p.path)
	}

	finfos, err := ioutil.ReadDir(p.path)
	if err != nil {
		return err
	}

	p.projects = []Project{}
	for _, finfo := range finfos {
		if isProjectFile(finfo) {
			proj, err := readProjectFile(filepath.Join(p.path, finfo.Name()))
			if err == nil {
				p.projects = append(p.projects, *proj)
			}
		}
	}

	return nil
}

// isProjectFile tests if a FileInfo project points to a project file.
func isProjectFile(finfo os.FileInfo) bool {
	if !finfo.IsDir() {
		if ext := filepath.Ext(finfo.Name()); ext == ".project" {
			return true
		}
	}

	return false
}

// readProjectFile reads a a project file, parses the JSON in a Project.
func readProjectFile(filepath string) (*Project, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var project Project
	if err := json.Unmarshal(b, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

// Projects returns the list of loaded projects.
func (p *ProjectDB) Projects() []Project {
	return p.projects
}

// Add adds a new project to and writes the corresponding project file.
func (p *ProjectDB) Add(proj Project) error {
	if p.Exists(proj.Name) {
		return errors.New(fmt.Sprintf("Project already exists: %s", proj.Name))
	}

	if err := p.writeProject(proj); err != nil {
		return err
	}

	p.projects = append(p.projects, proj)
	return nil
}

// Remove removes a project and its file.
func (p *ProjectDB) Remove(projectName string) error {
	projects, projectFound := p.projectsExceptFor(projectName)

	// We found the entry to remove, so we attempt to remove the project file.
	if projectFound {
		if err := os.Remove(p.projectFilePath(projectName)); err != nil {
			return err
		}
	}

	p.projects = projects
	return nil
}

// Update updates an existing project and its file.
func (p *ProjectDB) Update(proj Project) error {
	projects, found := p.projectsExceptFor(proj.Name)
	if found {
		if err := p.writeProject(proj); err != nil {
			return err
		}
		projects = append(projects, proj)
		p.projects = projects
		return nil
	}

	return errors.New(fmt.Sprintf("Project not found: %s", proj.Name))
}

// projectsExceptFor returns a new list of projects except for the project
// matching projectName. It returns true if it found a project matching
// projectName.
func (p *ProjectDB) projectsExceptFor(projectName string) ([]Project, bool) {
	projects := []Project{}
	found := false
	for _, project := range p.projects {
		if project.Name != projectName {
			projects = append(projects, project)
		} else {
			found = true
		}
	}
	return projects, found
}

// writeProject writes a project to a project file.
func (p *ProjectDB) writeProject(project Project) error {
	b, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return err
	}

	filename := p.projectFilePath(project.Name)
	return ioutil.WriteFile(filename, b, os.ModePerm)
}

// projectFilePath creates the path to a projects file.
func (p *ProjectDB) projectFilePath(projectName string) string {
	return filepath.Join(p.path, projectName+".project")
}

// Exists returns true if there is a project matching
// the given Name.
func (p *ProjectDB) Exists(projectName string) bool {
	_, found := p.Find(projectName)
	return found
}

// Find returns (Project, true) if the project is found otherwise
// it returns (Project{}, false).
func (p *ProjectDB) Find(projectName string) (Project, bool) {
	project, index := p.find(projectName)
	return project, index != -1
}

// find returns (Project, index) where index is -1 if
// the project wasn't found, otherwise it is the index
// in the Projects array.
func (p *ProjectDB) find(projectName string) (Project, int) {
	for index, project := range p.projects {
		if project.Name == projectName {
			return project, index
		}
	}

	return Project{}, -1
}
