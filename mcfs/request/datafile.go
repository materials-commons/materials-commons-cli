package request

import (
	"path/filepath"
	"strings"
)

func datafileDir(dataFileID string) string {
	pieces := strings.Split(dataFileID, "-")
	return filepath.Join("/mcfs/data/materialscommons", pieces[1][0:2], pieces[1][2:4])
}

func datafilePath(dataFileID string) string {
	return filepath.Join(datafileDir(dataFileID), dataFileID)
}
