package stor

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/materials-commons/materials-commons-cli/pkg/model"
)

type RemotesDB struct {
	DefaultRemote *model.Remote  `json:"default_remote"`
	Remotes       []model.Remote `json:"remotes"`
}

type JsonRemoteStor struct {
	db     RemotesDB
	Path   string
	Loaded bool
}

// MustLoadJsonRemoteStor attempts to read the config.json file
// at $HOME/.materialscommons/config.json. If this fails it
// immediately exits with a log message about the failure.
func MustLoadJsonRemoteStor() *JsonRemoteStor {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to get user home dir: %s", err)
	}

	storPath := filepath.Join(homeDir, ".materialscommons", "config.json")
	s := &JsonRemoteStor{
		Path:   storPath,
		Loaded: false,
	}

	if err := s.load(); err != nil {
		log.Fatalf("Unable to load %q: %s", storPath, err)
	}

	return s
}

func NewJsonRemoteStor(path string) *JsonRemoteStor {
	return &JsonRemoteStor{Path: path, Loaded: false}
}

// GetDefaultRemote returns the default remote. This is determined by first checking if both
// MCURL and MCAPIKEY are set in the environment. If they are set then it returns a Remote
// with those values. If neither env is set or only one is set then it checks if
// db.config.DefaultRemote is valid. If it is valid, then it returns that remote, overridding
// the MCUrl or MCAPIKEY from the environment if either set is set.
func (s *JsonRemoteStor) GetDefaultRemote() (*model.Remote, error) {
	mcapikeyEnv := os.Getenv("MCAPIKEY")
	mcurlEnv := os.Getenv("MCURL")

	// Load the database if it hasn't been loaded.
	if err := s.load(); err != nil {
		return nil, err
	}

	// The config.json was loaded, but no default_remote was set. When this happens
	// we can still return a remote without a user email if both mcapikeyEnv and
	// mcurlEnv are set.
	if s.db.DefaultRemote == nil {
		if mcapikeyEnv != "" && mcurlEnv != "" {
			return &model.Remote{MCAPIKey: mcapikeyEnv, MCUrl: mcurlEnv}, nil
		}

		return nil, fmt.Errorf("no default remote set")
	}

	r := &model.Remote{
		EMail:    s.db.DefaultRemote.EMail,
		MCUrl:    s.db.DefaultRemote.MCUrl,
		MCAPIKey: s.db.DefaultRemote.MCAPIKey,
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

// GetRemoteByUserServerUrl iterates through the list of remotes looking for match.
func (s *JsonRemoteStor) GetRemoteByUserServerUrl(email, serverUrl string) (*model.Remote, error) {
	// Load the database if it hasn't been loaded.
	if err := s.load(); err != nil {
		return nil, err
	}

	for _, r := range s.db.Remotes {
		if r.EMail == email && r.MCUrl == serverUrl {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("no such email/serverurl (%s/%s)", email, serverUrl)
}

// ListPaged calls fn for each remote. It stops iterating if fn returns an error.
func (s *JsonRemoteStor) ListPaged(fn func(remote *model.Remote) error) error {
	// Load the database if it hasn't been loaded.
	if err := s.load(); err != nil {
		return err
	}

	for _, r := range s.db.Remotes {
		if err := fn(&r); err != nil {
			return err
		}
	}

	return nil
}

// load will attempt to load the given json file in s.Path. If successful it
// will mark the s.Loaded as true.
func (s *JsonRemoteStor) load() error {
	if s.Loaded {
		return nil
	}

	contents, err := os.ReadFile(s.Path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(contents, &s.db); err != nil {
		return err
	}

	s.Loaded = true

	return nil
}
