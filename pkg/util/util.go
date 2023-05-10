package util

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/materials-commons/materials-commons-cli/pkg/config"
)

// ToProjectPath takes a path, strips out the path to the project, replaces it with a / and then returns
// the cleaned and normalized version. For example, if your project is in /home/me/project, and the
// path passed in is /home/me/project/images/file.jpg, then ToProjectPath will return /images/file.jpg.
// This is the path relative to your project, with your project path treated as '/'
func ToProjectPath(path string) string {
	// To step through this, lets say that config.GetProjectRootPath() returns '/home/me/project',
	// and that path is /home/me/project/images/file.jpg

	// First remove the /home/me/project portion.
	// pathWithProjectDirReplaced = /images/file.jpg
	pathWithProjectDirReplaced := strings.Replace(path, config.GetProjectRootPath(), "", 1)

	// Now, join this with a '/' because we can't guarantee that pathWithProjectDirReplaced starts
	// with a '/'.
	// addSlashToPath = /images/file.jpg
	addSlashToPath := filepath.Join(string(os.PathSeparator), pathWithProjectDirReplaced)

	// Finally - Return the cleaned version. This is only important if, for some reason you ended up with a
	// constructed addSlashToPath that looked like /./images.file.jpg or similar. In that case  filepath.Clean
	// would return /images/file.jpg. That is it cleans up the path.
	return filepath.Clean(addSlashToPath)
}
