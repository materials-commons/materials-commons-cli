package materials

import (
	"net/http"
	"os"
	"path/filepath"
)

func (p Project) UploadProject() error {
	filepath.Walk(p.Path, func(path string, info os.FileInfo, err error) error {
		return nil
	})
	return nil
}

func newFileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	return nil, nil
}
