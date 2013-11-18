package materials

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"io"
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
			//importDatafile(ids.ProjectId, ddirid, path)
			
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
	fmt.Printf("Create datadir: %s for project %s with project path %s, parent: %s\n", dirPath, projectId, projectPath, parentId)
	ddirName := makeDatadirName(projectPath, dirPath)
	j := `{"name":"` + ddirName + `", "parent":"` + parentId + `", "project":"` + projectId + `"}`
	b := strings.NewReader(j)
	resp, err := http.Post("http://localhost:5000/datadirs?apikey=4a3ec8f43cc511e3ba368851fb4688d4",
		"application/json", b)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	fmt.Printf("createDataDir call = %d\n", resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)
	var data McId
	json.Unmarshal(body, &data)
	fmt.Println(data)
	return data.Id, nil
}

func makeDatadirName(projectPath, dirPath string) string {
	projectPathParent := filepath.Dir(projectPath) + "/"
	return strings.Replace(dirPath, projectPathParent, "", 1)
}

func createProject(projectName string) (*Project2DatadirIds, error) {
	//	user := NewCurrentUser()
	j := `{"name":"` + projectName + `", "description":"Newly created project"}`
	b := strings.NewReader(j)
	resp, err := http.Post("http://localhost:5000/projects?apikey=4a3ec8f43cc511e3ba368851fb4688d4",
		"application/json", b)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("status = %d\n", resp.StatusCode)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("body ", string(body))
	var data Project2DatadirIds
	json.Unmarshal(body, &data)
	fmt.Println(data)

	return &data, nil
}

func importDatafile(projectId, ddirId, path string) (string, error) {
	uri := "http://localhost:5000/import?apikey=4a3ec8f43cc511e3ba368851fb4688d4"
	otherParams := map[string]string{
		"project": projectId,
		"datadir": ddirId,
		"name":    filepath.Base(path),
	}
	request, err := newFileUploadRequest(uri, otherParams, "file", path)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	fmt.Printf("importDatafile call = %d\n", resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)	
	fmt.Println("body ", string(body))
	var data McId
	json.Unmarshal(body, &data)
	fmt.Println(data)
	return data.Id, nil
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
	//finfo, err := file.Stat()
	req, err := http.NewRequest("POST", uri, reader)
	req.Header.Add("Content-Type", "multipart/form-data; boundary=" + boundary)
	req.ContentLength = int64(body.Len()) + int64(closeBuf.Len())
	fmt.Println(body)
	return http.DefaultClient.Do(req)
}

func newFileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(paramName, filepath.Base(path))

	if err != nil {
		return nil, err
	}

	part.Write(fileContents)

	//for key, val := range params {
	//	writer.WriteField(key, val)
	//}

	writer.Close()
	if err != nil {
		return nil, err
	}

	fmt.Println("===============")
	fmt.Println(body)
	fmt.Println("===============")
	return http.NewRequest("POST", uri, body)
}
