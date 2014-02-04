package mcfs

import (
	"errors"
	"fmt"
	"github.com/materials-commons/base/mc"
	"github.com/materials-commons/gohandy/file"
	"github.com/materials-commons/materials/db"
	"github.com/materials-commons/materials/db/schema"
	"github.com/materials-commons/mcfs/protocol"
	"os"
	"path/filepath"
)

var (
	ErrPathsDiffer = errors.New("Paths differ")
)

func (c *Client) projectEntries(projectName string) (*protocol.ProjectEntriesResp, error) {
	req := &protocol.ProjectEntriesReq{
		Name: projectName,
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	switch t := resp.(type) {
	case protocol.ProjectEntriesResp:
		return &t, nil
	default:
		return nil, ErrBadResponseType
	}
}

func (c *Client) CreateProject(projectName string) (*Project, error) {
	req := protocol.CreateProjectReq{
		Name: projectName,
	}

	resp, err := c.doRequest(req)
	if resp == nil {
		return nil, err
	}

	switch t := resp.(type) {
	case protocol.CreateProjectResp:
		return &Project{t.ProjectID, t.DataDirID}, err
	default:
		return nil, ErrBadResponseType
	}
}

func (c *Client) UploadNewProject(path string) error {
	var dataDirs = map[string]string{}
	projectName := filepath.Base(path)
	project, err := c.CreateProject(projectName)
	if err != nil && err != mc.ErrExists {
		return err
	}

	dataDirs[path] = project.DataDirID

	filepath.Walk(path, func(fpath string, info os.FileInfo, err error) error {
		switch info.IsDir() {
		case true:
			if fpath == path {
				// Top level project dir already created
				return nil
			}
			// Create Directory
			dataDirID, err := c.CreateDir(project.ProjectID, projectName, fpath)

			if err != nil {
				fmt.Println("CreateDir failure", err)
			} else {
				fmt.Printf("Created New Directory %s with ID %s\n", fpath, dataDirID)
				dataDirs[file.NormalizePath(fpath)] = dataDirID
			}
		case false:
			// Upload File
			dir := file.NormalizePath(filepath.Dir(fpath))
			fmt.Println("Upload file looking up directory", dir)
			dataDirID, ok := dataDirs[dir]
			if !ok {
				fmt.Println("  Couldn't find directory id for", dir)
				return nil
			}
			fmt.Printf("  Uploading file %s for dataDir %s and project %s\n", fpath, dataDirID, project.ProjectID)
			bytes, dataFileID, err := c.UploadNewFile(project.ProjectID, dataDirID, fpath)
			if err != nil {
				fmt.Printf("Upload file %s failed %s\n", fpath, err)
			}
			fmt.Printf("  Done with upload of %s datafileid %s bytes %d\n", fpath, dataFileID, bytes)
		}
		return nil
	})

	return nil
}

func (c *Client) LoadFromRemote(path string) error {
	return nil
}

func (c *Client) IndexProject(path string) error {
	var project *schema.Project
	var err error
	project, err = projectByPath(path)
	switch {
	case err == mc.ErrNotFound:
		return c.loadNewProject(path)
	case err != nil:
		return err
	}

	entries, err := c.projectEntries(project.Name)

	var _ = entries

	return nil
}

func (c *Client) loadNewProject(path string) error {
	project, err := createNewProject(path) // TODO: Need MC ProjectID
	if err != nil {
		return err
	}

	entryResp, err := c.projectEntries(project.Name)
	if err != nil {
		return nil
	}

	for _, entry := range entryResp.Entries {
		switch {
		case entry.DataFileName == "":
			// This is just a datadir
			dataDir := schema.DataDir{
				ProjectID:  project.Id,
				MCId:       entry.DataDirID,
				Name:       entry.DataDirName,
				Path:       "", // TODO: Create the path
				ParentMCId: "", //TODO: We aren't sending this yet
				Parent:     0,  // This needs to be computed...

			}
			err := db.DataDirs.Insert(dataDir)
			if err != nil {
				fmt.Println("err on insert into database %s", err)
			}
		default:
			// This is a datafile
		}
	}

	return nil
}

func projectByPath(path string) (*schema.Project, error) {
	var project schema.Project
	projectName := filepath.Base(path)
	err := db.Projects.Get(&project, "select * from projects where name=$1", projectName)
	switch {
	case err != nil:
		return nil, mc.ErrNotFound
	default:
		return &project, nil
	}
}

func createNewProject(path string) (*schema.Project, error) {
	project := schema.Project{
		Name: filepath.Base(path),
		Path: path,
	}
	err := db.Projects.Insert(project)
	if err != nil {
		return nil, err
	}

	return projectByPath(path)
}
