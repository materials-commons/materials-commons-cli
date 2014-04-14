package mcfs

import (
	"github.com/materials-commons/base/dir"
	"github.com/materials-commons/mcfs/protocol"
)

type ProjectStat struct {
	ID      string
	Name    string
	Entries []dir.FileInfo
}

func (c *Client) StatProject(projectName string) (*ProjectStat, error) {
	req := protocol.StatProjectReq{
		Name: projectName,
		Base: "/home/gtarcea",
	}

	resp, err := c.doRequest(req)
	if resp == nil {
		return nil, err
	}

	switch t := resp.(type) {
	case protocol.StatProjectResp:
		return &ProjectStat{
			Name:    projectName,
			ID:      t.ProjectID,
			Entries: t.Entries,
		}, nil
	default:
		return nil, ErrBadResponseType
	}
}
