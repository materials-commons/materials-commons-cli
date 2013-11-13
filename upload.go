package materials

import (
	"bytes"
	"io"
	"mime/multipart"
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
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile(paramName, filepath.Base(path))

	_, err = io.Copy(part, file)
	writer.WriteField("fullpath", path)

	for key, val := range params {
		writer.WriteField(key, val)
	}

	writer.Close()
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", uri, body)
}
