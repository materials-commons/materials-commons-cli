package mcfs

import (
	"fmt"
	"github.com/materials-commons/gohandy/marshaling"
	"github.com/materials-commons/materials/transfer"
	"io"
	"net"
	"os"
)

const readBufSize = 1024 * 1024 * 20

type Client struct {
	marshaling.MarshalUnmarshaler
	conn net.Conn
}

type Project struct {
	ProjectID string
	DataDirID string
}

var ErrBadResponseType = fmt.Errorf("Unexpected Response Type")

func NewClient(host string, port uint) (*Client, error) {
	return nil, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Login(user, apikey string) error {
	req := &transfer.LoginReq{
		User:   user,
		ApiKey: apikey,
	}

	_, err := c.doRequest(req)
	return err
}

func (c *Client) Logout() error {
	req := &transfer.LogoutReq{}
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
		return nil, ErrBadResponseType
	}
}

func (c *Client) CreateDir(projectID, path string) (dataDirID string, err error) {
	req := &transfer.CreateDirReq{
		ProjectID: projectID,
		Path:      path,
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	switch t := resp.(type) {
	case transfer.CreateResp:
		return t.ID, nil
	default:
		return "", ErrBadResponseType
	}
}

func (c *Client) UploadNewFile(projectID, dataDirID, path string) (bytesUploaded int64, err error) {
	return 0, nil
}

func (c *Client) createFile(req *transfer.CreateFileReq) (dataFileID string, err error) {
	resp, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	switch t := resp.(type) {
	case transfer.CreateResp:
		return t.ID, nil
	default:
		return "", ErrBadResponseType
	}
}

func (c *Client) startUpload(req *transfer.UploadReq) (*transfer.UploadResp, error) {
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	switch t := resp.(type) {
	case transfer.UploadResp:
		return &t, nil
	default:
		return nil, ErrBadResponseType
	}
}

func (c *Client) endUpload() {
	c.doRequest(&transfer.DoneReq{})
}

func (c *Client) sendFile(dataFileID, path string, offset int64) (bytesSent int64, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	_, err = f.Seek(offset, 0)
	if err != nil {
		return 0, err
	}

	return c.sendFileBytes(f, dataFileID)
}

func (c *Client) sendFileBytes(f *os.File, dataFileID string) (totalSent int64, err error) {
	sendReq := transfer.SendReq{
		DataFileID: dataFileID,
	}

	buf := make([]byte, readBufSize)
	for {
		var bytesSent int
		n, err := f.Read(buf)
		if n != 0 {
			sendReq.Bytes = buf[:n]
			bytesSent, err = c.sendBytes(&sendReq)
			if err != nil {
				break
			}
			totalSent = totalSent + int64(bytesSent)
		}
		if err != nil {
			break
		}
	}

	if err != nil && err != io.EOF {
		return totalSent, err
	}

	return totalSent, nil
}

func (c *Client) sendBytes(sendReq *transfer.SendReq) (bytesSent int, err error) {
	resp, err := c.doRequest(sendReq)
	if err != nil {
		return 0, err
	}

	switch t := resp.(type) {
	case transfer.SendResp:
		return t.BytesWritten, nil
		// there are other cases we need to check for
	default:
		return 0, ErrBadResponseType
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
