package remote

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Remote struct {
	MCUrl    string `json:"mcurl"`
	EMail    string `json:"email"`
	MCAPIKey string `json:"mcapikey"`
}

type Config struct {
	DefaultRemote *Remote  `json:"default_remote"`
	Remotes       []Remote `json:"remotes"`
}

type DB struct {
	config Config
	Path   string
	Loaded bool
}

func NewDB(path string) *DB {
	return &DB{Loaded: false}
}

func (db *DB) Load() error {
	if db.Loaded {
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	db.Path = filepath.Join(homeDir, ".materialscommons", "config.json")
	if _, err := os.Stat(db.Path); err != nil {
		return err
	}

	contents, err := os.ReadFile(db.Path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(contents, &db.Path); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetAPIKey(mcuser, mcurl string) (string, error) {
	mcapikeyEnv := os.Getenv("MCAPIKEY")
	if mcapikeyEnv != "" {
		// mcapikey overridden
		return mcapikeyEnv, nil
	}

	if !db.Loaded {
		return "", fmt.Errorf("remote db not loaded")
	}

	for _, r := range db.config.Remotes {
		if r.EMail == mcuser && r.MCUrl == mcurl {
			return r.MCAPIKey, nil
		}
	}

	return "", fmt.Errorf("no such entry match %q/%q", mcuser, mcurl)
}

// GetDefaultRemote returns the default remote. This is determined by first checking if both
// MCURL and MCAPIKEY are set in the environment. If they are set then it returns a Remote
// with those values. If neither env is set or only one is set then it checks if
// db.config.DefaultRemote is valid. If it is valid, then it returns that remote, overridding
// the MCUrl or MCAPIKEY from the environment if either set is set.
func (db *DB) GetDefaultRemote() (*Remote, error) {
	mcapikeyEnv := os.Getenv("MCAPIKEY")
	mcurlEnv := os.Getenv("MCURL")

	// The $HOME/.materialscommons/config.json hasn't been loaded, so return an error.
	if !db.Loaded {
		return nil, fmt.Errorf("remote db not loaded")
	}

	// The config.json was loaded, but no default_remote was set. When this happens
	// we can still return a remote without a user email if both mcapikeyEnv and
	// mcurlEnv are set.
	if db.config.DefaultRemote == nil {
		if mcapikeyEnv != "" && mcurlEnv != "" {
			return &Remote{MCAPIKey: mcapikeyEnv, MCUrl: mcurlEnv}, nil
		}

		return nil, fmt.Errorf("no default remote set")
	}

	r := &Remote{
		EMail:    db.config.DefaultRemote.EMail,
		MCUrl:    db.config.DefaultRemote.MCUrl,
		MCAPIKey: db.config.DefaultRemote.MCAPIKey,
	}

	// Override MCAPIKey or MCUrl from the environment if appropriate.

	if mcapikeyEnv != "" {
		r.MCAPIKey = mcapikeyEnv
	}

	if mcurlEnv != "" {
		r.MCUrl = mcurlEnv
	}

	// We've created a remote, and overridden either MCAPIKey and/or MCUrl from env (if set in env). If both
	// of these are still unset (For example the user only had email set in default_remote, but not mcapikey
	// and/or mcurl), then the remote is invalid and treated as unset.

	if r.MCAPIKey == "" || r.MCUrl == "" {
		return nil, fmt.Errorf("no default remote set")
	}

	// We've passed all the checks!
	return r, nil
}
