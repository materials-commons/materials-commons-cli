package materials

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
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
	client := makeClient()
	ids, err := createProject(p.Name, client)
	if err != nil {
		return err
	}

	dir2id := make(map[string]string)
	dir2id[p.Path] = ids.DatadirId

	filepath.Walk(p.Path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if path != p.Path {
				parentId, _ := dir2id[filepath.Dir(path)]
				id, err := createDataDir(ids.ProjectId, p.Path, path, parentId, client)
				if err != nil {
					return err
				}
				dir2id[path] = id
			}
		} else {
			// Loading a file
			ddirid := dir2id[filepath.Dir(path)]
			if !fileAlreadyUploaded(ddirid, path, client) {
				uri := Config.ApiUrlPath("/import")
				resp, err := postFile(ddirid, ids.ProjectId, path, uri, client)
				if err != nil {
					fmt.Println(err)
				} else {
					if resp.StatusCode > 299 {
						body, _ := ioutil.ReadAll(resp.Body)
						resp.Body.Close()
						fmt.Printf("Unable to import file %s, error: %s\n", path, string(body))
					} else {
						fmt.Printf("Imported file %s\n", path)
					}
				}
			} else {
				fmt.Printf("File already uploaded: %s\n", path)
			}
		}

		return nil
	})
	return nil
}

func makeClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

func createDataDir(projectId, projectPath, dirPath, parentId string, client *http.Client) (string, error) {
	ddirName := makeDatadirName(projectPath, dirPath)
	j := `{"name":"` + ddirName + `", "parent":"` + parentId + `", "project":"` + projectId + `"}`
	b := strings.NewReader(j)
	uri := Config.ApiUrlPath("/datadirs")
	resp, err := client.Post(uri, "application/json", b)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 299 {
		return "", errors.New(
			fmt.Sprintf("Unable to create datadir %s, error: %s",
				dirPath, string(body)))
	}
	var data McId
	json.Unmarshal(body, &data)
	return data.Id, nil
}

func makeDatadirName(projectPath, dirPath string) string {
	projectPathParent := filepath.Dir(projectPath) + "/"
	return strings.Replace(dirPath, projectPathParent, "", 1)
}

func createProject(projectName string, client *http.Client) (*Project2DatadirIds, error) {
	j := `{"name":"` + projectName + `", "description":"Newly created project"}`
	b := strings.NewReader(j)

	uri := Config.ApiUrlPath("/projects")
	resp, err := client.Post(uri, "application/json", b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 299 {
		return nil, errors.New(
			fmt.Sprintf("Unable to create create project %s, error: %s",
				projectName, string(body)))
	}

	var data Project2DatadirIds
	json.Unmarshal(body, &data)
	return &data, nil
}

func fileAlreadyUploaded(ddirId, filename string, client *http.Client) bool {
	uri := Config.ApiUrlPath("/datafiles/" + ddirId + "/" + filepath.Base(filename))
	resp, err := client.Get(uri)
	defer resp.Body.Close()

	if err != nil {
		return false
	}

	if resp.StatusCode > 499 {
		// Server error, assume it is uploaded for now
		return true
	} else if resp.StatusCode > 299 {
		return false
	}

	return true
}

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
