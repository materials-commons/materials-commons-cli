package materials

import (
	"fmt"
	"github.com/materials-commons/gohandy/ezhttp"
	"os"
	"path/filepath"
	"strings"
)

type Project2DatadirIds struct {
	ProjectId string `json:"project_id"`
	DatadirId string `json:"datadir_id"`
}

type MCId struct {
	Id string `json:"id"`
}

var client = ezhttp.NewInsecureClient()

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
				id, err := createDataDir(ids.ProjectId, p.Path, path, parentId)
				if err != nil {
					return err
				}
				dir2id[path] = id
			}
		} else {
			// Loading a file
			ddirid := dir2id[filepath.Dir(path)]
			if !fileAlreadyUploaded(ddirid, path) {
				uri := Config.ApiUrlPath("/import")
				var params = map[string]string{
					"datadir": ddirid,
					"project": ids.ProjectId,
				}
				_, err := client.PostFile(uri, path, "file", params)
				//resp, err := postFile(ddirid, ids.ProjectId, path, uri, client)
				if err != nil {
					fmt.Printf("Unable to import file %s, error: %s\n", path, err.Error())
				} else {
					fmt.Printf("Imported file %s\n", path)
				}
			} else {
				fmt.Printf("File already uploaded: %s\n", path)
			}
		}

		return nil
	})
	return nil
}

func createProject(projectName string) (*Project2DatadirIds, error) {
	j := `{"name":"` + projectName + `", "description":"Newly created project"}`

	uri := Config.ApiUrlPath("/projects")
	var data Project2DatadirIds
	_, err := client.JSONStr(j).JSONPost(uri, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func createDataDir(projectId, projectPath, dirPath, parentId string) (string, error) {
	ddirName := makeDatadirName(projectPath, dirPath)
	j := `{"name":"` + ddirName + `", "parent":"` + parentId + `", "project":"` + projectId + `"}`
	var data MCId
	uri := Config.ApiUrlPath("/datadirs")
	_, err := client.JSONStr(j).JSONPost(uri, &data)

	if err != nil {
		return "", err
	}

	return data.Id, nil
}

func makeDatadirName(projectPath, dirPath string) string {
	projectPathParent := filepath.Dir(projectPath) + "/"
	return strings.Replace(dirPath, projectPathParent, "", 1)
}

func fileAlreadyUploaded(ddirId, filename string) bool {
	uri := Config.ApiUrlPath("/datafiles/" + ddirId + "/" + filepath.Base(filename))
	var rv map[string]interface{}
	status, err := client.JSONGet(uri, &rv)

	if err != nil {
		return false
	}

	if status > 499 {
		// Server error, assume it is uploaded for now
		return true
	} else if status > 299 {
		return false
	}

	return true
}

/*
func postFile(ddirId, projectId, filename, uri string, client *http.Client) (*http.Response, error) {
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
	req.Close = true

	return client.Do(req)
}
*/
