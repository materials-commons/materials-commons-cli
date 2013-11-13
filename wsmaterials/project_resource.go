package wsmaterials

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/materials-commons/materials"
	"net/http"
	"os"
	"path/filepath"
)

type ProjectResource struct {
	*materials.MaterialsProjects
}

func newProjectResource(container *restful.Container) error {
	p, err := materials.CurrentUserProjects()
	if err != nil {
		return err
	}
	projectResource := ProjectResource{p}
	projectResource.register(container)
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

	ws.Route(ws.GET("/{project-name}/tree").To(p.getProjectTree).
		Doc("Retrieve the directory/file tree for the project").
		Param(ws.PathParameter("project-name", "name of the project").DataType("string")))

	ws.Route(ws.POST("/projects").To(p.newProject).
		Doc("Create a new project").
		Reads(materials.Project{}))

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
	for _, p := range p.Projects() {
		if p.Name == projectName {
			response.WriteEntity(p)
			return
		}
	}

	response.AddHeader("Content-Type", "text/plain")
	response.WriteErrorString(http.StatusNotFound, fmt.Sprintf("Project not found: %s", projectName))
}

func (p ProjectResource) getProjectTree(request *restful.Request, response *restful.Response) {
	type ditem struct {
		Id          string  `json:"id"`
		Name        string  `json:"name"`
		Displayname string  `json:"displayname"`
		Dtype       string  `json:"type"`
		Children    []ditem `json:"children"`
	}

	sayHello := func() {
		fmt.Println("hello")
	}

	dirs := make(map[string]*ditem)
	var currentDir *ditem
	var topLevelDirs []*ditem

	filepath.Walk("/tmp/a", func(path string, info os.FileInfo, err error) error {
		sayHello()
		if path == "/tmp/a" {
			item := &ditem{
				Id:          path,
				Name:        path,
				Displayname: path,
				Dtype:       "datadir",
				Children:    []ditem{},
			}
			dirs[path] = item
			currentDir = item
			topLevelDirs = append(topLevelDirs, item)
		} else {

			parent := filepath.Dir(path)
			d, found := dirs[parent]
			if !found {
				panic("d should not be null")
			}
			item := ditem{
				Id:          path,
				Name:        path,
				Displayname: filepath.Base(path),
				Children:    []ditem{},
			}
			if info.IsDir() {
				item.Dtype = "datadir"
				dirs[path] = &item
			} else {
				item.Dtype = "datafile"
			}
			// ******************************
			// Put in logic around currentDir and adding
			// children to the currentDir
			// ******************************
			fmt.Printf("Appending to children of %s\n", parent)
			fmt.Println(item)
			d.Children = append(d.Children, item)
			//dirs[parent] = d
			fmt.Println(d.Children)
		}
		fmt.Printf("path = %s\n", path)
		return nil
	})

	fmt.Println(dirs)
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
