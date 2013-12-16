package materials

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MaterialsProjects contains a list of user projects and information that
// is needed by the methods to load the projects file.
type ProjectDB struct {
	path     string
	projects []Project
}

func CurrentUserProjects() (*ProjectDB, error) {
	projectsPath := filepath.Join(Config.User.DotMaterialsPath(), "projects")
	return OpenProjectDB(projectsPath)
}

// Load projects from the database file.
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

// Reads the projects file, parses its contents and loads it into
// the Projects struct.
func (p *ProjectDB) loadProjects() error {
	projectsFile, err := os.Open(p.path)
	if err != nil {
		return err
	}
	defer projectsFile.Close()

	projects := []Project{}
	scanner := bufio.NewScanner(projectsFile)
	for scanner.Scan() {
		splitLine := strings.Split(scanner.Text(), "|")
		if len(splitLine) == 3 {
			projects = append(projects, Project{
				Name:   strings.TrimSpace(splitLine[0]),
				Path:   strings.TrimSpace(splitLine[1]),
				Status: strings.TrimSpace(splitLine[2]),
			})
		}
	}
	p.projects = projects
	return nil
}

// Attempts to create an empty projects file.
func (p *ProjectDB) createEmptyProjectsFile() error {
	file, err := os.Create(p.path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// Projects returns the list of loaded projects.
func (p *ProjectDB) Projects() []Project {
	return p.projects
}

// Add adds a new project and updates the projects file.
func (p *ProjectDB) Add(proj Project) error {
	if p.Exists(proj.Name) {
		return errors.New(fmt.Sprintf("Project already exists: %s", proj.Name))
	}
	p.projects = append(p.projects, proj)
	return p.writeToProjectsFile(p.projects)
}

// Remove removes a project and updates the projects file.
func (p *ProjectDB) Remove(projectName string) error {
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

func (p *ProjectDB) Update(proj Project) error {
	_, index := p.find(proj.Name)
	if index != -1 {
		projects := p.Projects()
		projects[index].Path = proj.Path
		projects[index].Status = proj.Status
		return p.writeToProjectsFile(p.Projects())
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

// writeToProjectsFile overwrites the projects file with the list
// of projects.
func (p *ProjectDB) writeToProjectsFile(projects []Project) error {
	file, err := os.Create(p.path)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, project := range projects {
		projectLine := fmt.Sprintf("%s|%s|%s\n",
			strings.TrimSpace(project.Name),
			strings.TrimSpace(project.Path),
			strings.TrimSpace(project.Status))
		file.WriteString(projectLine)
	}
	return nil
}

// Exists returns true if there is a project matching
// the given Name.
func (p *ProjectDB) Exists(projectName string) bool {
	_, found := p.Find(projectName)
	return found
}

// Find returns (Project, true) if the project is found otherwise
// it returns (Project{}, false)
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
