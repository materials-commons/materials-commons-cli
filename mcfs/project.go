package mcfs

import (
	"fmt"
	"github.com/materials-commons/mcfs/protocol"
)

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

func (c *Client) IndexProject(path string) error {
	return nil
}
