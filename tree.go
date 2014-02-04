package materials

import (
	"fmt"
	"github.com/materials-commons/gohandy/file"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ProjectFileEntry struct {
	Id          string              `json:"id"`
	ParentId    string              `json:"parent_id"`
	Level       int                 `json:"level"`
	Path        string              `json:"path"`
	HrefPath    string              `json:"hrefpath"`
	DisplayName string              `json:"displayname"`
	Type        string              `json:"type"`
	Children    []*ProjectFileEntry `json:"children"`
}

type treeState struct {
	nextId           int
	projectName      string
	knownDirectories map[string]*ProjectFileEntry
	currentDir       *ProjectFileEntry
	topLevelDirs     []*ProjectFileEntry
}

func (p Project) Tree() ([]*ProjectFileEntry, error) {
	ts := &treeState{
		projectName:      p.Name,
		knownDirectories: make(map[string]*ProjectFileEntry),
	}

	if !file.Exists(p.Path) {
		return nil, fmt.Errorf("Directory path '%s' doesn't exist for project %s", p.Path, p.Name)
	}

	filepath.Walk(p.Path, func(path string, info os.FileInfo, err error) error {
		normalizedPath := file.NormalizePath(path)
		if normalizedPath == p.Path {
			ts.createTopLevelEntry(normalizedPath)
		} else {
			ts.addChild(normalizedPath, info)
		}
		return nil
	})

	return ts.topLevelDirs, nil
}

func (ts *treeState) createTopLevelEntry(path string) {
	item := &ProjectFileEntry{
		Id:          strconv.Itoa(ts.nextId),
		Path:        path,
		DisplayName: path,
		Type:        "datadir",
		Children:    []*ProjectFileEntry{},
	}

	ts.nextId++
	ts.knownDirectories[path] = item
	ts.currentDir = item
	ts.topLevelDirs = append(ts.topLevelDirs, item)
}

func (ts *treeState) addChild(path string, info os.FileInfo) {
	parent := file.NormalizePath(filepath.Dir(path))
	d, found := ts.knownDirectories[parent]

	// There should always be a parent
	if !found {
		panic("No parent found (there should always be a parent)")
	}

	// Create the entry to add
	item := &ProjectFileEntry{
		Id:          strconv.Itoa(ts.nextId),
		Level:       ts.currentDir.Level + 1,
		ParentId:    d.Id,
		Path:        path,
		DisplayName: filepath.Base(path),
		Children:    []*ProjectFileEntry{},
	}

	ts.nextId++

	// What type of entry is this?
	if info.IsDir() {
		item.Type = "datadir"
		ts.knownDirectories[path] = item
	} else {
		item.HrefPath = ts.hrefPath(path)
		item.Type = "datafile"
		item.Level = 0
	}

	// Update the currentDir if needed
	if ts.currentDir.Path != parent {
		ts.currentDir = d
	}

	// Append new entry to the currentDir list of children
	ts.currentDir.Children = append(ts.currentDir.Children, item)
}

func (ts *treeState) hrefPath(path string) string {
	i := strings.Index(path, ts.projectName)
	return path[i:]
}
