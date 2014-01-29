package main

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/materials-commons/contrib/mc"
	"github.com/materials-commons/materials"
	"github.com/materials-commons/materials/autoupdate"
	_ "github.com/materials-commons/materials/db"
	"github.com/materials-commons/materials/mcfs"
	"github.com/materials-commons/materials/site"
	"github.com/materials-commons/materials/wsmaterials"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

var mcuser, _ = materials.NewCurrentUser()

type ServerOptions struct {
	AsServer bool   `long:"server" description:"Run as webserver"`
	Port     uint   `long:"port" description:"The port the server listens on"`
	Address  string `long:"address" description:"The address to bind to"`
	Retry    int    `long:"retry" description:"Number of times to retry connecting to address/port"`
}

type ProjectOptions struct {
	Project   string   `long:"project" description:"Specify the project"`
	Directory string   `long:"directory" description:"The directory path to the project"`
	Add       bool     `long:"add" description:"Add the project to the project config file"`
	Delete    bool     `long:"delete" description:"Delete the project from the project config file"`
	List      bool     `long:"list" description:"List all known projects and their locations"`
	Upload    bool     `long:"upload" description:"Uploads a new project. Cannot be used on existing projects"`
	Convert   bool     `long:"convert" description:"Converts projects to new layout"`
	Files     []string `long:"file" description:"comma separated list of files to operate on"`
	Tracking  bool     `long:"tracking" description:"Display tracking information for specified files"`
	Options   []string `long:"option" description:"Options for tracking"`
}

type Options struct {
	Server     ServerOptions  `group:"Server Options"`
	Project    ProjectOptions `group:"Project Options"`
	Initialize bool           `long:"init" description:"Create configuration"`
	Config     bool           `long:"config" description:"Show server configuration"`
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
	projects := materials.CurrentUserProjectDB()
	for _, p := range projects.Projects() {
		fmt.Printf("%s, %s\n", p.Name, p.Path)
	}
}

func addProject(projectName, projectPath string) error {
	if materials.CurrentUserProjectDB().Exists(projectName) {
		fmt.Printf("Project %s already exists\n", projectName)
		return mc.ErrExists
	}

	p, err := materials.NewProject(projectName, projectPath, "Unloaded")
	if err != nil {
		fmt.Printf("Unable to create project %s: %s\n", projectName, err)
		return err
	}

	if err = materials.CurrentUserProjectDB().Add(*p); err != nil {
		fmt.Printf("Unable to add project %s: %s\n", projectName, err)
		return err
	}

	fmt.Println("Created new project:", projectName)
	return nil
}

func convertProjects() {
	setupProjectsDir()
	convertProjectsFile()
}

func setupProjectsDir() {
	projectDB := filepath.Join(materials.Config.User.DotMaterialsPath(), "projectdb")
	err := os.MkdirAll(projectDB, os.ModePerm)
	checkError(err)
}

func convertProjectsFile() {
	projectsPath := filepath.Join(materials.Config.User.DotMaterialsPath(), "projects")
	projectsFile, err := os.Open(projectsPath)
	projectdbPath := filepath.Join(materials.Config.User.DotMaterialsPath(), "projectdb")
	checkError(err)
	defer projectsFile.Close()

	scanner := bufio.NewScanner(projectsFile)
	for scanner.Scan() {
		splitLine := strings.Split(scanner.Text(), "|")
		if len(splitLine) == 3 {
			project := materials.Project{
				Name:    strings.TrimSpace(splitLine[0]),
				Path:    strings.TrimSpace(splitLine[1]),
				Status:  strings.TrimSpace(splitLine[2]),
				ModTime: time.Now(),
				Changes: map[string]materials.ProjectFileChange{},
				Ignore:  []string{},
			}
			b, err := json.MarshalIndent(&project, "", "  ")
			if err != nil {
				fmt.Printf("Could not convert '%s' to new project format\n", scanner.Text())
				continue
			}
			path := filepath.Join(projectdbPath, project.Name+".project")
			if err := ioutil.WriteFile(path, b, os.ModePerm); err != nil {
				fmt.Printf("Unable to write project file %s\n", path)
			}
		}
	}
}

func uploadProject(projectName string) {
	project, found := materials.CurrentUserProjectDB().Find(projectName)
	if !found {
		fmt.Println("No such project:", projectName)
		return
	}

	c, err := mcfs.NewClient("localhost", 35862)
	if err != nil {
		fmt.Println("Unable create client", err)
		return
	}

	err = c.Login(materials.Config.User.Username, materials.Config.User.Apikey)

	if err != nil {
		fmt.Println("Unable to login", err)
		return
	} else {
		err = c.UploadNewProject(project.Path)
		if err != nil {
			fmt.Println("Error on upload", err)
		}
		projects := materials.CurrentUserProjectDB()
		projects.Update(func() *materials.Project {
			project.Status = "Loaded"
			return project
		})
	}
}

func startServer(serverOpts ServerOptions) {
	autoupdate.StartUpdateMonitor()

	if serverOpts.Address != "" {
		materials.Config.Server.Address = serverOpts.Address
	}

	if serverOpts.Port != 0 {
		materials.Config.Server.Port = serverOpts.Port
	}

	if serverOpts.Retry != 0 {
		wsmaterials.StartRetry(serverOpts.Retry)
	} else {
		wsmaterials.Start()
	}
}

func showConfig() {
	fmt.Println("Configuration:")
	if b, err := json.MarshalIndent(&materials.Config, "", "  "); err != nil {
		fmt.Printf("Unable to display configuration: %s\n", err)
	} else {
		fmt.Println(string(b))
	}

	fmt.Printf("\nEnvironment Variables:\n")
	for _, envVarName := range materials.EnvironmentVariables {
		value := os.Getenv(envVarName)
		if value == "" {
			value = "Not Set"
		}
		fmt.Println("  ", envVarName+":", value)
	}

	fmt.Printf("\nPath to .materials: %s\n\n", materials.Config.User.DotMaterialsPath())
}

func showTracking(projectName string, files []string) {
	projects := materials.CurrentUserProjectDB()
	project, found := projects.Find(projectName)
	if !found {
		fmt.Printf("Unknown project %s\n", projectName)
		return
	}

	for _, filename := range files {
		fullpath := filepath.Join(project.Path, filename)
		key := md5.New().Sum([]byte(fullpath))
		data, err := project.Get(key, nil)
		if err != nil {
			fmt.Println("")
			fmt.Println("Unknown file:", fullpath)
			continue
		}

		var p materials.ProjectFileInfo
		json.Unmarshal(data, &p)
		b, _ := json.MarshalIndent(&p, "", "  ")
		fmt.Println("")
		fmt.Println("Tracking information for", fullpath, ":")
		fmt.Println(string(b))
	}
}

func main() {
	materials.ConfigInitialize(mcuser)
	var opts Options
	flags.Parse(&opts)

	switch {
	case opts.Initialize:
		initialize()
	case opts.Project.List:
		listProjects()
	case opts.Project.Add:
		addProject(opts.Project.Project, opts.Project.Directory)
	case opts.Project.Convert:
		convertProjects()
	case opts.Project.Upload:
		uploadProject(opts.Project.Project)
	case opts.Server.AsServer:
		startServer(opts.Server)
	case opts.Config:
		showConfig()
	case opts.Project.Tracking:
		showTracking(opts.Project.Project, opts.Project.Files)
	}
}
