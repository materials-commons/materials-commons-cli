package util

import (
	"os"
	"path/filepath"
)

func ToProjectPath(path string) string {
	return filepath.Clean(filepath.Join(string(os.PathSeparator), path))
}
