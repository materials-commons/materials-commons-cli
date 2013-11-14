package materials

import (
	"os/user"
	"path/filepath"
)

func dotmaterialsForCurrentUser() string {
	u, _ := user.Current()
	return dotmaterials(u.HomeDir)
}

func dotmaterialsForUser(username string) string {
	u, _ := user.Lookup(username)
	return dotmaterials(u.HomeDir)
}

func dotmaterialsFrom(path string) string {
	return dotmaterials(path)
}

func dotmaterials(path string) string {
	return filepath.Join(path, ".materials")
}
