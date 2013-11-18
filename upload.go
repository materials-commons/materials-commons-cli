package materials

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Project2DatadirIds struct {
	ProjectId string `json:"project_id"`
	DatadirId string `json:"datadir_id"`
}

type McId struct {
	Id string `json:"id"`
}

func (p Project) Upload() error {
	ids, err := createProject(p.Name)
	if err != nil {
		return err
	}

	dir2id := make(map[string]string)
	dir2id[p.Path] = ids.DatadirId

	filepath.Walk(p.Path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if path != p.Path {
				parentId, _ := dir2id[filepath.Dir(path)]
				id, _ := createDataDir(ids.ProjectId, p.Path, path, parentId)
				dir2id[path] = id
			}
		} else {
			// Loading a file
			ddirid := dir2id[filepath.Dir(path)]
			res, err := postFile(ddirid, ids.ProjectId, path, "http://localhost:5000/import?apikey=4a3ec8f43cc511e3ba368851fb4688d4")
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("status code = %d\n", res.StatusCode)
			}
		}

		return nil
	})
	return nil
}

func createDataDir(projectId, projectPath, dirPath, parentId string) (string, error) {
	ddirName := makeDatadirName(projectPath, dirPath)
	j := `{"name":"` + ddirName + `", "parent":"` + parentId + `", "project":"` + projectId + `"}`
	b := strings.NewReader(j)
	resp, err := http.Post("http://localhost:5000/datadirs?apikey=4a3ec8f43cc511e3ba368851fb4688d4",
		"application/json", b)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var data McId
	json.Unmarshal(body, &data)
	return data.Id, nil
}

func makeDatadirName(projectPath, dirPath string) string {
	projectPathParent := filepath.Dir(projectPath) + "/"
	return strings.Replace(dirPath, projectPathParent, "", 1)
}

func createProject(projectName string) (*Project2DatadirIds, error) {
	j := `{"name":"` + projectName + `", "description":"Newly created project"}`
	b := strings.NewReader(j)

	resp, err := http.Post("http://localhost:5000/projects?apikey=4a3ec8f43cc511e3ba368851fb4688d4",
		"application/json", b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var data Project2DatadirIds
	json.Unmarshal(body, &data)
	return &data, nil
}

func postFile(ddirId, projectId, filename, uri string) (*http.Response, error) {
	body := bytes.NewBufferString("")
	writer := multipart.NewWriter(body)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	file, _ := os.Open(filename)
	defer file.Close()

	fileContents, err := ioutil.ReadAll(file)
	part.Write(fileContents)

	boundary := writer.Boundary()
	closeStr := fmt.Sprintf("\r\n--%s--\r\n", boundary)

	writer.WriteField("datadir", ddirId)
	writer.WriteField("project", projectId)

	closeBuf := bytes.NewBufferString(closeStr)
	reader := io.MultiReader(body, file, closeBuf)

	req, err := http.NewRequest("POST", uri, reader)

	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = int64(body.Len()) + int64(closeBuf.Len())

	return http.DefaultClient.Do(req)
}
