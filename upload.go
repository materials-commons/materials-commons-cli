package materials

import (
	"fmt"
	"github.com/materials-commons/gohandy/ezhttp"
	"os"
	"path/filepath"
	"strings"
)

type project2DatadirIDs struct {
	ProjectID string `json:"project_id"`
	DatadirID string `json:"datadir_id"`
}

type mcID struct {
	ID string `json:"id"`
}

var client = ezhttp.NewInsecureClient()

// Upload uploads a project by doing posts. No longer supported.
func (p Project) Upload() error {
	ids, err := createProject(p.Name)
	if err != nil {
		return err
	}

	dir2id := make(map[string]string)
	dir2id[p.Path] = ids.DatadirID

	filepath.Walk(p.Path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if path != p.Path {
				parentID, _ := dir2id[filepath.Dir(path)]
				id, err := createDataDir(ids.ProjectID, p.Path, path, parentID)
				if err != nil {
					return err
				}
				dir2id[path] = id
			}
		} else {
			// Loading a file
			ddirid := dir2id[filepath.Dir(path)]
			if !fileAlreadyUploaded(ddirid, path) {
				uri := Config.APIURLPath("/import")
				var params = map[string]string{
					"datadir": ddirid,
					"project": ids.ProjectID,
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

func createProject(projectName string) (*project2DatadirIDs, error) {
	j := `{"name":"` + projectName + `", "description":"Newly created project"}`

	uri := Config.APIURLPath("/projects")
	var data project2DatadirIDs
	_, err := client.JSONStr(j).JSONPost(uri, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func createDataDir(projectID, projectPath, dirPath, parentID string) (string, error) {
	ddirName := makeDatadirName(projectPath, dirPath)
	j := `{"name":"` + ddirName + `", "parent":"` + parentID + `", "project":"` + projectID + `"}`
	var data mcID
	uri := Config.APIURLPath("/datadirs")
	_, err := client.JSONStr(j).JSONPost(uri, &data)

	if err != nil {
		return "", err
	}

	return data.ID, nil
}

func makeDatadirName(projectPath, dirPath string) string {
	projectPathParent := filepath.Dir(projectPath) + "/"
	return strings.Replace(dirPath, projectPathParent, "", 1)
}

func fileAlreadyUploaded(ddirID, filename string) bool {
	uri := Config.APIURLPath("/datafiles/" + ddirID + "/" + filepath.Base(filename))
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
