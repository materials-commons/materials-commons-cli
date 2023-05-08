package config

import (
	"os"
	"path/filepath"
	"strconv"
)

var (
	txRetry         int
	projectRootPath string
)

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

func GetProjectDBPath() string {
	return filepath.Join(GetProjectMCDirPath(), "project.db")
}

func GetProjectMCDirPath() string {
	return filepath.Join(GetProjectRootPath(), ".mc")
}

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
