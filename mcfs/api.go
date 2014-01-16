package mcfs

import (
	"fmt"
	"github.com/materials-commons/materials/transfer"
	"github.com/materials-commons/materials/util"
	"net"
)

func NewClient(host string, port uint) (*Client, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	m := util.NewGobMarshaler(conn)
	c := &Client{
		MarshalUnmarshaler: m,
		conn:               conn,
	}
	return c, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Login(user, apikey string) error {
	req := transfer.LoginReq{
		User:   user,
		ApiKey: apikey,
	}

	_, err := c.doRequest(req)
	return err
}

func (c *Client) Logout() error {
	req := transfer.LogoutReq{}
	_, err := c.doRequest(req)
	return err
}

func (c *Client) CreateProject(projectName string) (*Project, error) {
	req := &transfer.CreateProjectReq{
		Name: projectName,
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	switch t := resp.(type) {
	case transfer.CreateProjectResp:
		return &Project{t.ProjectID, t.DataDirID}, nil
	default:
		fmt.Printf("1 %s %T\n", ErrBadResponseType, t)
		return nil, ErrBadResponseType
	}
}

func (c *Client) doRequest(arg interface{}) (interface{}, error) {
	req := &transfer.Request{
		Req: arg,
	}

	if err := c.Marshal(req); err != nil {
		return nil, err
	}

	var resp transfer.Response

	if err := c.Unmarshal(&resp); err != nil {
		return nil, err
	}

	if resp.Type != transfer.ROk {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	return resp.Resp, nil
}
