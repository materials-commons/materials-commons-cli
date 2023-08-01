package mcapi

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/materials-commons/hydra/pkg/mcdb/mcmodel"
)

type Client struct {
	APIToken string
	BaseURL  string
	rc       *resty.Client
}

var ErrAuth = errors.New("authentication failed")
var tlsConfig = tls.Config{InsecureSkipVerify: true}

func NewClient(token, baseUrl string) *Client {
	fmt.Printf("NewClient token = '%s', baseUrl = '%s'\n", token, baseUrl)
	return &Client{
		APIToken: token,
		BaseURL:  baseUrl,
		rc:       resty.New().SetTLSClientConfig(&tlsConfig),
	}
}

func (c *Client) ListDirectoryByPath(projectID int, path string) ([]mcmodel.File, error) {
	var files []mcmodel.File
	resp, err := c.rc.R().SetAuthToken(c.APIToken).
		SetQueryParam("path", path).
		SetResult(files).
		Get(fmt.Sprintf("%s/projects/%d/directories_by_path", c.BaseURL, projectID))
	if resp.IsError() {
		return files, fmt.Errorf("api call failed: %d/%s", resp.StatusCode(), resp.Status())
	}
	return files, err
}
