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

func Download() error {
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

func unpackWebsite(path string) {
	file, _ := os.Open(path)
	defer file.Close()

	zhandle, _ := gzip.NewReader(file)
	defer zhandle.Close()

	thandle := tar.NewReader(zhandle)
	for {
		hdr, err := thandle.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		if hdr.Typeflag == tar.TypeDir {
			dirpath := filepath.Join(mcuser.DotMaterialsPath(), hdr.Name)
			os.MkdirAll(dirpath, 0777)
		} else if hdr.Typeflag == tar.TypeReg || hdr.Typeflag == tar.TypeRegA {
			filepath := filepath.Join(mcuser.DotMaterialsPath(), hdr.Name)
			out, _ := os.Create(filepath)
			if _, err := io.Copy(out, thandle); err != nil {
				fmt.Println(err)
			}
			out.Close()
		}
	}
	8
}
