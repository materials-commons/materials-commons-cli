package wsmaterials

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/materials-commons/gohandy/handyfile"
	"github.com/materials-commons/materials"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ProjectFileStatus struct {
	FilePath string `json:"filepath"`
	Status   string `json:"status"`
}

type ProjectResource struct {
	*materials.MaterialsProjects
	events []ProjectFileStatus
}

func newProjectResource(container *restful.Container) error {
	p, err := materials.CurrentUserProjects()
	if err != nil {
		return err
	}

	projectResource := ProjectResource{
		MaterialsProjects: p,
		events:            make([]ProjectFileStatus, 10),
	}
	projectResource.register(container)

	//go projectResource.monitorEventLoop()

	return nil
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

	ws.Route(ws.GET("/{project-name}/tree").Filter(JsonpFilter).To(p.getProjectTree).
		Doc("Retrieve the directory/file tree for the project").
		Param(ws.PathParameter("project-name", "name of the project").DataType("string")))

	ws.Route(ws.POST("").To(p.newProject).
		Doc("Create a new project").
		Reads(materials.Project{}))

	ws.Route(ws.PUT("/{project-name}").To(p.updateProject).
		Doc("Updates the project").
		Param(ws.PathParameter("project-name", "name of the project").DataType("string")).
		Reads(materials.Project{}).
		Writes(materials.Project{}))

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
		ParentId    string   `json:"parent_id"`
		Name        string   `json:"name"`
		HrefPath    string   `json:"hrefpath"`
		Displayname string   `json:"displayname"`
		Dtype       string   `json:"type"`
		Children    []*ditem `json:"children"`
	}

	dirs := make(map[string]*ditem)
	var currentDir *ditem
	var topLevelDirs []*ditem
	projectName := request.PathParameter("project-name")
	project, found := p.Find(projectName)
	nextId := 0

	createTopLevelEntry := func(path string) {
		item := &ditem{
			Id:          strconv.Itoa(nextId),
			Name:        path,
			Displayname: path,
			Dtype:       "datadir",
			Children:    []*ditem{},
		}

		nextId++
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
			panic("No parent found (there should always be a parent)")
		}

		// Create the ditem
		item := ditem{
			Id:          strconv.Itoa(nextId),
			ParentId:    d.Id,
			Name:        path,
			HrefPath:    hrefPath(projectName, path),
			Displayname: filepath.Base(path),
			Children:    []*ditem{},
		}

		nextId++

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
		if handyfile.Exists(project.Path) {

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
			response.WriteErrorString(http.StatusNotFound,
				fmt.Sprintf("Project path does not exist: %s", project.Path))
		}
	} else {
		response.WriteErrorString(http.StatusNotFound,
			fmt.Sprintf("Project not found: %s", projectName))
	}
}

func hrefPath(projectName, path string) string {
	i := strings.Index(path, projectName)
	return path[i:]
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
		return
	}

	err = p.Add(*project)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotAcceptable, err.Error())
		return
	}

	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(project)
}

func (p *ProjectResource) updateProject(request *restful.Request, response *restful.Response) {
	projectName := request.PathParameter("project-name")
	project := new(materials.Project)
	err := request.ReadEntity(&project)
	fmt.Println("project-name", projectName, project)
	if err != nil {
		response.WriteErrorString(http.StatusNotAcceptable, err.Error())
		return
	}

	originalProject, found := p.Find(projectName)
	if !found {
		response.WriteErrorString(http.StatusNotFound, fmt.Sprintf("Project not found '%s'", projectName))
		return
	}

	if project.Name != projectName {
		p.Remove(projectName)
		project.Status = originalProject.Status
		err = p.Add(*project)
	} else {
		err = p.Update(*project)
	}

	if err != nil {
		response.WriteErrorString(http.StatusNotAcceptable, err.Error())
	} else {
		response.WriteEntity(project)
	}
}

func (p *ProjectResource) uploadProject(request *restful.Request, response *restful.Response) {
	projectName := request.PathParameter("project-name")
	project, found := p.Find(projectName)
	if found {
		err := project.Upload()
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
