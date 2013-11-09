package wsmaterials

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/materials-commons/materials"
	"net/http"
)

type ProjectResource struct {
	user *materials.User
}

func newProjectResource(container *restful.Container) error {
	u, err := materials.CurrentUser()
	if err != nil {
		return err
	}
	projectResource := ProjectResource{user: u}
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

	container.Add(ws)
}

func (p ProjectResource) allProjects(request *restful.Request, response *restful.Response) {
	if len(p.user.Projects) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "User has no projects.")
	} else {
		response.WriteEntity(p.user.Projects)
	}
}

func (p ProjectResource) getProject(request *restful.Request, response *restful.Response) {
	projectName := request.PathParameter("project-name")
	for _, p := range p.user.Projects {
		if p.Name == projectName {
			response.WriteEntity(p)
			return
		}
	}

	response.AddHeader("Content-Type", "text/plain")
	response.WriteErrorString(http.StatusNotFound, fmt.Sprintf("Project not found: %s", projectName))
}

func (p ProjectResource) getProjectTree(request *restful.Request, response *restful.Response) {
}
