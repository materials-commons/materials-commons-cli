package mcapi

import (
	"crypto/tls"
	"errors"

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
	return &Client{
		APIToken: token,
		BaseURL:  baseUrl,
		rc:       resty.New().SetTLSClientConfig(&tlsConfig),
	}
}

func (c *Client) ListDirectoryByPath(projectID int, path string) ([]mcmodel.File, error) {
	_, err := c.rc.R().SetAuthToken(c.APIToken).Post("")
	return nil, err
}
