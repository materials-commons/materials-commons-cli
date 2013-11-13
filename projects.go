package materials

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Project describes the information we track about a users
// projects. Here we keep the name of the project and the
// directory path. The name of the project is the top level
// directory of the project. The path is the full path to
// the project excluding the name.
type Project struct {
	Name string `json:"name" xml:"name"`
	Path string `json:"path" xml:"path"`
}

// Projects contains a list of user projects and information that
// is needed by the methods to load the projects file.
type MaterialsProjects struct {
	dir      string
	projects []Project
}

// CurrentUserProjects retrieves the projects for the currently logged in user.
func CurrentUserProjects() (*MaterialsProjects, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	return loadProjectsFrom(u.HomeDir)
}

// ProjectsForUser retrieves the projects for the named user.
func ProjectsForUser(username string) (*MaterialsProjects, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return nil, err
	}
	return loadProjectsFrom(u.HomeDir)
}

// ProjectsFromHome retrieves the projects located in the directory
// pointed at by the HOME environment variable.
func ProjectsFromHome() (*MaterialsProjects, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return nil, errors.New("HOME Environment variable is undefined")
	}
	return loadProjectsFrom(home)
}

// ProjectsFrom retrieves the projects from the specified directory. The
// path given cannot contain the .materials subdirectory. A .materials
// subdirectory must exist in the path.
func ProjectsFrom(dir string) (*MaterialsProjects, error) {
	return loadProjectsFrom(dir)
}

// Load projects using the specified directory as a base.
// The projects file is located in {dir}/.materials/projects.
func loadProjectsFrom(dir string) (*MaterialsProjects, error) {
	userProjects := MaterialsProjects{dir: dir}
	err := userProjects.loadProjects()
	if err != nil {
		return nil, err
	}
	return &userProjects, err
}

// Reload re-reads and loads the projects file.
func (p *MaterialsProjects) Reload() error {
	return p.loadProjects()
}

// Reads the projects file, parses its contents and loads it into
// the Projects struct.
func (p *MaterialsProjects) loadProjects() error {
	projectsFile, err := os.Open(p.projectsFile())
	if err != nil {
		return err
	}
	defer projectsFile.Close()

	projects := []Project{}
	scanner := bufio.NewScanner(projectsFile)
	for scanner.Scan() {
		splitLine := strings.Split(scanner.Text(), "|")
		if len(splitLine) == 2 {
			projects = append(projects, Project{splitLine[0], splitLine[1]})
		}
	}
	p.projects = projects
	return nil
}

// Creates the path to the projects file:
// {dir}/.materias/projects
func (p *MaterialsProjects) projectsFile() string {
	return filepath.Join(p.dir, ".materials", "projects")
}

// Projects returns the list of loaded projects.
func (p *MaterialsProjects) Projects() []Project {
	return p.projects
}

// Add adds a new project and updates the projects file.
func (p *MaterialsProjects) Add(proj Project) error {
	if p.Exists(proj.Name) {
		return errors.New(fmt.Sprintf("Project already exists: %s", proj.Name))
	}
	p.projects = append(p.projects, proj)
	return p.writeToProjectsFile(p.projects)
}

// Remove removes a project and updates the projects file.
func (p *MaterialsProjects) Remove(projectName string) error {
	projects, projectsUpdated := p.projectsExceptFor(projectName)

	// We found the entry to remove. Thus we need to
	// update the projects file.
	if projectsUpdated {
		err := p.writeToProjectsFile(projects)
		if err == nil {
			p.projects = projects
		}
		return err
	} else {
		// project wasn't found so we don't update anything and we return no error
		return nil
	}
}

// projectsExceptFor returns a new list of projects except for the project
// matching projectName. It returns true if it found a project matching
// projectName.
func (p *MaterialsProjects) projectsExceptFor(projectName string) ([]Project, bool) {
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

// writeToProjectsFile overwrites the projects file with the list
// of projects.
func (p *MaterialsProjects) writeToProjectsFile(projects []Project) error {
	file, err := os.Create(p.projectsFile())
	if err != nil {
		return err
	}
	defer file.Close()
	for _, project := range projects {
		projectLine := fmt.Sprintf("%s|%s\n", project.Name, project.Path)
		file.WriteString(projectLine)
	}
	return nil
}

// Exists returns true if there is a project matching
// the given Name.
func (p *MaterialsProjects) Exists(projectName string) bool {
	for _, project := range p.projects {
		if project.Name == projectName {
			return true
		}
	}
	return false
}
