package wsmaterials

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/materials-commons/materials"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ProjectFileStatus struct {
	FilePath string `json:"filepath"`
	Status   string `json:"status"`
}

type ProjectResource struct {
	*materials.MaterialsProjects
	materialsCommons *materials.MaterialsCommons
	events           []ProjectFileStatus
}

func newProjectResource(container *restful.Container) error {
	p, err := materials.CurrentUserProjects()
	if err != nil {
		return err
	}

	u, err := materials.NewCurrentUser()
	if err != nil {
		return err
	}

	mc := materials.NewMaterialsCommons(u)

	projectResource := ProjectResource{
		MaterialsProjects: p,
		materialsCommons:  mc,
		events:            make([]ProjectFileStatus, 10),
	}
	projectResource.register(container)

	if false {
		go projectResource.monitorEventLoop()
	}
	return nil
}

func (p *ProjectResource) monitorEventLoop() {
	var projectPaths []string
	for _, project := range p.Projects() {
		projectPaths = append(projectPaths, project.Path)
	}
	watcher, err := materials.NewRecursiveWatcher("/tmp/a")
	if err != nil {
		return
	}
	watcher.Run()
	defer watcher.Close()

	//out:
	for {
		select {
		case file := <-watcher.Files:
			fmt.Printf("File changed: %s\n", file)
			p.events[0] = ProjectFileStatus{
				FilePath: file,
				Status:   "File Changed",
			}
		case folder := <-watcher.Folders:
			fmt.Printf("Folder changed: %s\n", folder)
			p.events[0] = ProjectFileStatus{
				FilePath: folder,
				Status:   "Directory Changed",
			}
		}
	}
}

func (p ProjectResource) register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/projects").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("").Filter(JsonpFilter).To(p.allProjects).
		Doc("list all projects").
		Writes([]materials.Project{}))

	ws.Route(ws.GET("/{project-name}").Filter(JsonpFilter).To(p.getProject).
		Doc("Retrieve a particular project").
		Param(ws.PathParameter("project-name", "name of the project").DataType("string")).
		Writes(materials.Project{}))

	ws.Route(ws.GET("/{project-name}/tree").To(p.getProjectTree).
		Doc("Retrieve the directory/file tree for the project").
		Param(ws.PathParameter("project-name", "name of the project").DataType("string")))

	ws.Route(ws.POST("").To(p.newProject).
		Doc("Create a new project").
		Reads(materials.Project{}))

	ws.Route(ws.GET("/changes").Filter(JsonpFilter).To(p.projectChanges).
		Doc("list all file system changes for all the projects").
		Writes(ProjectFileStatus{}))

	ws.Route(ws.GET("/{project-name}/upload").To(p.uploadProject).
		Doc("Uploads/imports a project to Materials Commons").
		Param(ws.PathParameter("project-name", "name of the project").DataType("string")))

	container.Add(ws)
}

func (p ProjectResource) allProjects(request *restful.Request, response *restful.Response) {
	if len(p.Projects()) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "User has no projects.")
	} else {
		response.WriteEntity(p.Projects())
	}
}

func (p ProjectResource) getProject(request *restful.Request, response *restful.Response) {
	projectName := request.PathParameter("project-name")
	project, found := p.Find(projectName)
	if found {
		response.WriteEntity(project)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound,
			fmt.Sprintf("Project not found: %s", projectName))
	}
}

func (p ProjectResource) getProjectTree(request *restful.Request, response *restful.Response) {
	type ditem struct {
		Id          string   `json:"id"`
		Name        string   `json:"name"`
		Displayname string   `json:"displayname"`
		Dtype       string   `json:"type"`
		Children    []*ditem `json:"children"`
	}

	dirs := make(map[string]*ditem)
	var currentDir *ditem
	var topLevelDirs []*ditem
	projectName := request.PathParameter("project-name")
	project, found := p.Find(projectName)

	createTopLevelEntry := func(path string) {
		item := &ditem{
			Id:          path2Id(path),
			Name:        path,
			Displayname: path,
			Dtype:       "datadir",
			Children:    []*ditem{},
		}

		dirs[path] = item
		currentDir = item
		topLevelDirs = append(topLevelDirs, item)
	}

	// addChild adds a child to the currentDir. If currentDir
	// is different than the childs parent it first updates
	// currentDir to the ditem for the parent path.
	addChild := func(path string, info os.FileInfo) {
		parent := filepath.Dir(path)
		d, found := dirs[parent]

		// There should always be a parent
		if !found {
			panic("d should not be null")
		}

		// Create the ditem
		item := ditem{
			Id:          path2Id(path),
			Name:        path,
			Displayname: filepath.Base(path),
			Children:    []*ditem{},
		}

		// What type of entry is this?
		if info.IsDir() {
			item.Dtype = "datadir"
			dirs[path] = &item
		} else {
			item.Dtype = "datafile"
		}

		// Update the currentDir if needed
		if currentDir.Name != parent {
			currentDir = d
		}

		// Append new entry to the currentDir list of children
		currentDir.Children = append(currentDir.Children, &item)
	}

	if found {
		filepath.Walk(project.Path, func(path string, info os.FileInfo, err error) error {
			if path == project.Path {
				createTopLevelEntry(path)
			} else {
				addChild(path, info)
			}
			return nil
		})

		response.WriteEntity(topLevelDirs)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound,
			fmt.Sprintf("Project not found: %s", projectName))
	}
}

func path2Id(path string) string {
	return strings.Replace(strings.Replace(path, "/", "_", -1), ".", "_", -1)
}

func (p *ProjectResource) newProject(request *restful.Request, response *restful.Response) {
	project := new(materials.Project)
	err := request.ReadEntity(&project)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotAcceptable, err.Error())
	}

	err = p.Add(*project)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotAcceptable, err.Error())
	}

	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(project)
}

func (p *ProjectResource) projectChanges(request *restful.Request, response *restful.Response) {
	fmt.Println(p.events[0])
	response.WriteEntity(p.events[0])
}

func (p *ProjectResource) uploadProject(request *restful.Request, response *restful.Response) {
	projectName := request.PathParameter("project-name")
	project, found := p.Find(projectName)
	if found {
		err := project.Upload(p.materialsCommons)
		if err != nil {
			response.WriteErrorString(http.StatusServiceUnavailable, "Unable to upload project")
		} else {
			project.Status = "Loaded"
			p.Update(project)
			response.WriteErrorString(http.StatusCreated, "Project uploaded")
		}
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound,
			fmt.Sprintf("Project not found: %s", projectName))
	}
}
