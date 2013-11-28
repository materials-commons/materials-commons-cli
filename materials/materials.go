package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/materials-commons/materials"
	"github.com/materials-commons/materials/site"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

var mcuser, _ = materials.NewCurrentUser()

type ServerOptions struct {
	AsServer bool   `long:"server" description:"Run as webserver"`
	Port     uint   `long:"port" description:"The port the server listens on"`
	Address  string `long:"address" description:"The address to bind to"`
	Retry    int    `long:"retry" description:"Number of times to retry connecting to address/port"`
}

type ProjectOptions struct {
	Project   string `long:"project" description:"Specify the project"`
	Directory string `long:"directory" description:"The directory path to the project"`
	Add       bool   `long:"add" description:"Add the project to the project config file"`
	Delete    bool   `long:"delete" description:"Delete the project from the project config file"`
	List      bool   `long:"list" description:"List all known projects and their locations"`
	Upload    bool   `long:"upload" description:"Uploads a new project. Cannot be used on existing projects"`
}

type Options struct {
	Server     ServerOptions  `group:"Server Options"`
	Project    ProjectOptions `group:"Project Options"`
	Initialize bool           `long:"init" description:"Create configuration"`
}

func initialize() {
	usr, err := user.Current()
	checkError(err)

	dirPath := filepath.Join(usr.HomeDir, ".materials")
	err = os.MkdirAll(dirPath, 0777)
	checkError(err)

	if newVersionOfWebsite() {
		websiteFilepath := filepath.Join(dirPath, "website")
		os.RemoveAll(websiteFilepath)
		downloadWebsite(dirPath)
	}
}

type MaterialsWebsiteInfo struct {
	Version     string `json:"version"`
	Description string `json:"description"`
}

func newVersionOfWebsite() bool {

	/*
		resp, _ := http.Get(mcurl + "/materials_website.json")
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var websiteInfo MaterialsWebsiteInfo
		json.Unmarshal(body, &websiteInfo)
	*/
	return true
}

func downloadWebsite(dirPath string) {
	getDownloadedVersionOfWebsite()
	websiteTarPath := filepath.Join(dirPath, "materials.tar.gz")
	out, _ := os.Create(websiteTarPath)
	defer out.Close()

	client := makeClient()

	resp, _ := client.Get(materials.Config.MCUrl() + "/materials.tar.gz")
	defer resp.Body.Close()
	io.Copy(out, resp.Body)
	unpackWebsite(websiteTarPath)
}

func makeClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
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
}

func getDownloadedVersionOfWebsite() int {
	//content := ioutil.ReadFile()
	return 0
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func listProjects() {
	projects, err := materials.CurrentUserProjects()
	if err != nil {
		return
	}
	for _, p := range projects.Projects() {
		fmt.Printf("%s, %s\n", p.Name, p.Path)
	}
}

func uploadProject(projectName string) {
	projects, _ := materials.CurrentUserProjects()
	project, _ := projects.Find(projectName)
	err := project.Upload()
	if err != nil {
		fmt.Println(err)
	} else {
		project.Status = "Loaded"
		projects.Update(project)
	}
}

func main() {
	materials.ConfigInitialize(mcuser)
	var opts Options
	_, err := flags.Parse(&opts)

	if err != nil {
		panic(err)
		os.Exit(1)
	}

	if opts.Initialize {
		initialize()
	}

	if opts.Project.List {
		listProjects()
	}

	if opts.Server.AsServer {
		if opts.Server.Address != "" {
			materials.Config.SetServerAddress(opts.Server.Address)
		}

		if opts.Server.Port != 0 {
			materials.Config.SetServerPort(opts.Server.Port)
		}

		if opts.Server.Retry != 0 {
			site.StartRetry(opts.Server.Retry)
		} else {
			site.Start()
		}
	}

	if opts.Project.Upload {
		uploadProject(opts.Project.Project)
	}
}
