package config

import (
	"os"
	"path/filepath"
	"strconv"
)

var (
	// txRetry is the number of times to retry a failed transaction. It is populated by
	// GetTxRetry the first time it is executed.
	txRetry int

	// projectRootPath is the path full path to the project root. This is defined as the first directory
	// containing the .mc directory. It is populated by GetProjectRootPath the first time it is executed.
	projectRootPath string
)

// GetTxRetry returns the number of times to retry a failed transaction. The minimum is 3. It uses
// the value for MCCLI_TX_RETRY if set. If MCCLI_TX_RETRY < 3, then 3 is used instead.
func GetTxRetry() int {
	if txRetry != 0 {
		return txRetry
	}

	txRetryCount64, err := strconv.ParseInt(os.Getenv("MCCLI_TX_RETRY"), 10, 32)
	if err != nil || txRetryCount64 < 3 {
		txRetryCount64 = 3
	}

	txRetry = int(txRetryCount64)

	return txRetry
}

// GetProjectDBPath returns the path to the project.db file. <PROJECTROOT>/.mc/project.db
func GetProjectDBPath() string {
	return filepath.Join(GetProjectMCDirPath(), "project.db")
}

// GetProjectMCDirPath returns the path to the .mc directory. <PROJECTROOT>/.mc
func GetProjectMCDirPath() string {
	return filepath.Join(GetProjectRootPath(), ".mc")
}

// GetProjectRootPath returns the path to the project root. The first directory with a .mc directory in it is
// considered the project root. It returns the full path. If a .mc directory is not found it returns "".
func GetProjectRootPath() string {
	if projectRootPath != "" {
		return projectRootPath
	}

	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Look for .mc dir
	for {
		dirToCheck := filepath.Join(dir, ".mc")
		finfo, err := os.Stat(dirToCheck)

		// If we found a .mc, and it is a directory return that path
		if err == nil && finfo.IsDir() {
			projectRootPath = dir
			return dir
		}

		// If we are here, then go up one level
		childDir := filepath.Dir(dir)

		// If we are at the root, then the filepath.Dir of root == root, so we stop searching
		if childDir == dir {
			return ""
		}

		dir = childDir
	}
}
