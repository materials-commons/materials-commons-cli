package site

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var mcurl = ""

type WebsiteInfo struct {
	Version     string `json:"version"`
	Description string `json:"description"`
}

func Download2() error {
	return nil
}

func newVersionOfWebsite() bool {
	resp, _ := http.Get(mcurl + "/materials_website.json")
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var websiteInfo WebsiteInfo
	json.Unmarshal(body, &websiteInfo)
	return false
}

func downloadWebsite(dirPath string) {
	getDownloadedVersionOfWebsite()
	out, _ := os.Create(filepath.Join(dirPath, "materials_website.tar"))
	defer out.Close()

	resp, _ := http.Get(mcurl + "/materials_website.tar")
	defer resp.Body.Close()
	io.Copy(out, resp.Body)
}

func getDownloadedVersionOfWebsite() int {
	//content := ioutil.ReadFile()
	return 0
}
