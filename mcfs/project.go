package mcfs

import (
	"errors"
	"fmt"
	"github.com/materials-commons/contrib/mc"
	"github.com/materials-commons/materials/db"
	"github.com/materials-commons/materials/db/schema"
	"github.com/materials-commons/mcfs/protocol"
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
		fmt.Println("not pointer")
		return &t, nil
	case *protocol.ProjectEntriesResp:
		fmt.Println("pointer")
		return t, nil
	default:
		return nil, ErrBadResponseType
	}
}

func (c *Client) CreateProject(projectName string) (*Project, error) {
	req := &protocol.CreateProjectReq{
		Name: projectName,
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	switch t := resp.(type) {
	case protocol.CreateProjectResp:
		return &Project{t.ProjectID, t.DataDirID}, nil
	default:
		fmt.Printf("1 %s %T\n", ErrBadResponseType, t)
		return nil, ErrBadResponseType
	}
}

func (c *Client) UploadProject(projectName string) {

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
