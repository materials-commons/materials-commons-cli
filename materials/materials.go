package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/materials-commons/materials"
	"github.com/materials-commons/materials/wsmaterials"
	"io"
	"archive/tar"
	"compress/gzip"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

var mcurl = ""
var usr, _ = user.Current()
var mcuser, _ = materials.NewCurrentUser()
var commons = materials.NewMaterialsCommons(mcuser)

//var user = NewCurrentUser()

type ServerOptions struct {
	AsServer bool   `long:"server" description:"Run as webserver"`
	Port     int    `long:"port" default:"8081" description:"The port the server listens on"`
	Address  string `long:"address" default:"127.0.0.1" description:"The address to bind to"`
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

	resp, _ := http.Get(mcurl + "/materials.tar.gz")
	defer resp.Body.Close()
	io.Copy(out, resp.Body)
	unpackWebsite(websiteTarPath)
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

func setup() {
	envMCURL := os.Getenv("MCURL")
	if envMCURL == "" {
		mcurl = "https://materialscommons.org"
	} else {
		mcurl = envMCURL
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

func runWebServer(address string, port int) {
	wsContainer := wsmaterials.NewRegisteredServicesContainer()
	http.Handle("/", wsContainer)
	mcwebdir := os.Getenv("MCWEBDIR")
	if mcwebdir == "" {
		mcwebdir = mcuser.DotMaterialsPath()
	}
	websiteDir := filepath.Join(mcwebdir, "website")
	dir := http.Dir(websiteDir)
	http.Handle("/materials/", http.StripPrefix("/materials/", http.FileServer(dir)))
	addr := fmt.Sprintf("%s:%d", address, port)
	http.ListenAndServe(addr, nil)
}

func uploadProject(projectName string) {
	projects, _ := materials.CurrentUserProjects()
	project, _ := projects.Find(projectName)
	project.Upload(commons)
	project.Status = "Loaded"
	projects.Update(project)
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)

	if err != nil {
		panic(err)
		os.Exit(1)
	}

	setup()

	if opts.Initialize {
		initialize()
	}

	if opts.Project.List {
		listProjects()
	}

	if opts.Server.AsServer {
		runWebServer(opts.Server.Address, opts.Server.Port)
	}

	if opts.Project.Upload {
		uploadProject(opts.Project.Project)
	}
}
