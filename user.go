package materials

import (
	"bufio"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Project struct {
	Name      string `json:"name" xml:"name"`
	Directory string `json:"directory" xml:"directory"`
}

type User struct {
	user     *user.User
	Projects []Project
}

func CurrentUser() (*User, error) {
	u, err := user.Current()

	if err != nil {
		return nil, err
	}

	usr := User{
		user: u,
	}

	usr.loadProjects()

	return &usr, nil
}

func (u *User) loadProjects() {
	projectsFile, err := os.Open(u.projectsFile())

	if err != nil {
		return
	}

	defer projectsFile.Close()

	scanner := bufio.NewScanner(projectsFile)
	for scanner.Scan() {
		splitLine := strings.Split(scanner.Text(), "|")
		if len(splitLine) == 2 {
			u.Projects = append(u.Projects, Project{splitLine[0], splitLine[1]})
		}
	}
}

func (u *User) projectsFile() string {
	return filepath.Join(u.user.HomeDir, ".materials", "projects")
}
