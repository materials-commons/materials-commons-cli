package mcfs

import (
	"crypto/md5"
	"fmt"
	"github.com/materials-commons/gohandy/handyfile"
	"github.com/materials-commons/gohandy/marshaling"
	"github.com/materials-commons/materials/transfer"
	"github.com/materials-commons/materials/util"
	"io"
	"net"
	"os"
	"path/filepath"
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

func (c *Client) RestartFileUpload(dataFileID, path string) (bytesUploaded int64, err error) {
	checksum, size, err := fileInfo(path)
	if err != nil {
		return 0, err
	}

	return c.uploadFile(dataFileID, path, checksum, size)
}

func (c *Client) uploadFile(dataFileID, path, checksum string, size int64) (bytesUploaded int64, err error) {
	uploadReq := &transfer.UploadReq{
		DataFileID: dataFileID,
		Checksum:   checksum,
		Size:       size,
	}

	uploadResp, err := c.startUpload(uploadReq)
	switch {
	case err != nil:
		return 0, err
	case uploadResp.DataFileID != dataFileID:
		return 0, fmt.Errorf("DataFileIDs don't match")
	default:
		n, err := c.sendFile(dataFileID, path, uploadResp.Offset)
		c.endUpload()
		return n, err
	}

}

func (c *Client) UploadNewFile(projectID, dataDirID, path string) (bytesUploaded int64, dataFileID string, err error) {
	checksum, size, err := fileInfo(path)
	if err != nil {
		return 0, "", err
	}

	createFileReq := &transfer.CreateFileReq{
		ProjectID: projectID,
		DataDirID: dataDirID,
		Name:      filepath.Base(path),
		Checksum:  checksum,
		Size:      size,
	}

	dataFileID, err = c.createFile(createFileReq)
	if err != nil {
		return 0, "", err
	}

	n, err := c.uploadFile(dataFileID, path, checksum, size)
	return n, dataFileID, err
}

func fileInfo(path string) (checksum string, size int64, err error) {
	checksum, err = handyfile.HashStr(md5.New(), path)
	if err != nil {
		return
	}

	finfo, err := os.Stat(path)
	if err != nil {
		return
	}
	size = finfo.Size()
	return
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
