package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/materials-commons/materials"
	"github.com/materials-commons/materials/site"
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

	if downloadedTo, err := site.Download(); err == nil {
		if site.IsNew(downloadedTo) {
			site.Deploy(downloadedTo)
		}
	}
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
