package materials

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (p Project) Upload() error {
	err := createProject(p.Name)
	if err != nil {
		return err
	}

	filepath.Walk(p.Path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			createDataDir(p.Name, p.Path, path)
		}

		return nil
	})
	return nil
}

func createDataDir(projectName string, projectPath string, dirPath string) {
	fmt.Printf("Create datadir: %s for project %s with project path %s\n", dirPath, projectName, projectPath)
}

func createProject(projectName string) error {
	//	user := NewCurrentUser()
	json := `{"name":"` + projectName + `", "description":"Newly created project"}`
	b := strings.NewReader(json)
	_, err := http.Post("http://localhost:5000/projects?apikey=4a3ec8f43cc511e3ba368851fb4688d4",
		"application/json", b)

	return err
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
