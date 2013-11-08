package materials

import (
	"os/user"
)

type Project struct {
	Name string
	dir string
}

type User struct {
	user     *user.User
	projects map[string]string
}

func New() (User, error) {
	u, err := user.Current()

	if err != nil {
		return User{}, err
	}

	usr := User{
		user:     u,
		projects: make(map[string]string),
	}

	usr.loadProjects()
}

func (u *User) loadProjects() {

}
