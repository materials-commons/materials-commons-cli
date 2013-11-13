package main

import (
	//"github.com/emicklei/go-restful"
	//"github.com/emicklei/go-restful/swagger"
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/materials-commons/materials"
	"github.com/materials-commons/materials/wsmaterials"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

var mcurl = ""
var usr, _ = user.Current()

type ServerOptions struct {
	AsServer bool   `long:"server" description:"Run as webserver"`
	Port     int    `long:"port" default:"8080" description:"The port the server listens on"`
	Address  string `long:"address" default:"127.0.0.1" description:"The address to bind to"`
}

type ProjectOptions struct {
	Project   string `long:"project" description:"Specify the project"`
	Directory string `long:"directory" description:"The directory path to the project"`
	Add       bool   `long:"add" description:"Add the project to the project config file"`
	Delete    bool   `long:"delete" description:"Delete the project from the project config file"`
	List      bool   `long:"list" description:"List all known projects and their locations"`
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
		downloadWebsite(dirPath)
	}
}

type MaterialsWebsiteInfo struct {
	Version     string `json:"version"`
	Description string `json:"description"`
}

func newVersionOfWebsite() bool {
	resp, _ := http.Get(mcurl + "/materials_website.json")
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var websiteInfo MaterialsWebsiteInfo
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

func runWebServer() {
	wsContainer := wsmaterials.NewRegisteredServicesContainer()
	http.Handle("/", wsContainer)
	dir := http.Dir("/home/gtarcea/GIT/prisms/materialscommons.org/website")
	http.Handle("/materials/", http.StripPrefix("/materials/", http.FileServer(dir)))
	http.ListenAndServe(":8081", nil)
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
		runWebServer()
	}
}
