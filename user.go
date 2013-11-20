package materials

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// User struct defines the information we need to know about
// a user in order to use the materials commons web services.
type User struct {
	Username string
	Apikey   string
	path     string
}

// NewCurrentUser creates a new User, and looks up
// the materials commons information by using the path
// to the .user file based on the current users home directory.
func NewCurrentUser() (*User, error) {
	u, _ := user.Current()
	user := &User{path: u.HomeDir + "/.materials"}
	user.readUser()
	return user, nil
}

// NewUserFrom creates a new user User and reads the materials commons information
// from the .user file in the given path.
func NewUserFrom(path string) (*User, error) {
	user := &User{path: path}
	user.readUser()
	return user, nil
}

// readUser reads the .user file and fills in the materials commons
// username and apikey
func (u *User) readUser() error {
	content, err := ioutil.ReadFile(u.dotuserPath())
	if err != nil {
		return err
	}

	pieces := strings.Split(string(content), "|")
	if len(pieces) != 2 {
		return errors.New("The .user file is corrupted")
	}

	u.Username = strings.TrimSpace(pieces[0])
	u.Apikey = strings.TrimSpace(pieces[1])
	return nil
}

// dotuser constructs the path to the .user file
func dotuser(dotmaterialsPath string) string {
	return filepath.Join(dotmaterialsPath, ".user")
}

// douserPath is a method to construct the path to the .user file
func (u *User) dotuserPath() string {
	return dotuser(u.path)
}

// Save writes the materials commons user information to the .user file
func (u *User) Save() error {
	file, err := os.Create(u.dotuserPath())
	if err != nil {
		return err
	}
	defer file.Close()

	userLine := fmt.Sprintf("%s|%s", u.Username, u.Apikey)
	file.WriteString(userLine)
	return nil
}
