package config

import (
	"encoding/json"
	"fmt"
	"log"
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

	configPath string

	remote ConfigRemote

	remoteRead bool
)

// GetTxRetry returns the number of times to retry a failed transaction. The minimum is 3. It uses
// the value for MCTXRETRY if set. If MCTXRETRY < 3, then 3 is used instead.
func GetTxRetry() int {
	if txRetry != 0 {
		return txRetry
	}

	txRetryCount64, err := strconv.ParseInt(os.Getenv("MCTXRETRY"), 10, 32)
	if err != nil || txRetryCount64 < 3 {
		txRetryCount64 = 3
	}

	txRetry = int(txRetryCount64)

	return txRetry
}

func GetProjectMCConfig() string {
	if configPath != "" {
		return configPath
	}

	config := filepath.Join(GetProjectMCDirPath(), "config.json")

	if _, err := os.Stat(config); err == nil {
		fmt.Printf("found %s\n", config)
		configPath = config
		return configPath
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	config = filepath.Join(homeDir, ".materialscommons", "config.json")

	if _, err := os.Stat(config); err == nil {
		configPath = config
		return configPath
	}

	return ""
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

func MustReadAPIKey() string {
	apikey := os.Getenv("MCAPIKEY")
	if apikey != "" {
		return apikey
	}

	if _, err := GetRemote(); err != nil {
		log.Fatalf("Failed to read config: %s", err)
	}

	if remote.DefaultRemote.MCAPIKey != "" {
		return remote.DefaultRemote.MCAPIKey
	}

	log.Fatalf("No default mcapikey set")

	// Not reachable
	return ""
}

func MustReadMCUrl() string {
	mcurl := os.Getenv("MCURL")
	if mcurl != "" {
		return mcurl
	}

	if _, err := GetRemote(); err != nil {
		log.Fatalf("Failed to read config: %s", err)
	}

	if remote.DefaultRemote.MCUrl != "" {
		return remote.DefaultRemote.MCUrl
	}

	log.Fatalf("No default mcurl set")

	// Not reachable
	return ""
}

func GetRemote() (ConfigRemote, error) {
	if remoteRead {
		return remote, nil
	}

	config := GetProjectMCConfig()

	if config == "" {
		return remote, fmt.Errorf("cannot find config.json")
	}

	contents, err := os.ReadFile(config)
	if err != nil {
		return remote, fmt.Errorf("unable to read %q: %s", config, err)
	}

	if err := json.Unmarshal(contents, &remote); err != nil {
		return remote, fmt.Errorf("unable to unmarshal %q: %s", config, err)
	}

	fmt.Printf("remote = %#v\n", remote)

	remoteRead = true

	return remote, nil
}

func GetWSScheme() string {
	wsScheme := os.Getenv("MC_WS_SCHEME")
	if wsScheme == "" {
		return "wss"
	}

	return wsScheme
}
