package wsmaterials

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/materials-commons/materials"
	"net/http"
)

type ProjectFileStatus struct {
	FilePath string `json:"filepath"`
	Status   string `json:"status"`
}

type ProjectResource struct {
	*materials.MaterialsProjects
}

func newProjectResource(container *restful.Container) error {
	p, err := materials.CurrentUserProjects()
	if err != nil {
		return err
	}

	projectResource := ProjectResource{
		MaterialsProjects: p,
	}
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

	ws.Route(ws.GET("/{project-name}/tree").Filter(JsonpFilter).To(p.getProjectTree).
		Doc("Retrieve the directory/file tree for the project").
		Param(ws.PathParameter("original-project-name", "original name of the project").
		DataType("string")))

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
	projectName := request.PathParameter("project-name")

	if project, found := p.Find(projectName); !found {
		response.WriteErrorString(http.StatusNotFound,
			fmt.Sprintf("Project not found: %s", projectName))
	} else if tree, err := project.Tree(); err != nil {
		response.WriteErrorString(http.StatusNotAcceptable, err.Error())
	} else {
		response.WriteEntity(tree)
	}
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
	originalProjectName := request.PathParameter("original-project-name")
	project := new(materials.Project)
	err := request.ReadEntity(&project)
	if err != nil {
		response.WriteErrorString(http.StatusNotAcceptable, err.Error())
		return
	}

	originalProject, found := p.Find(originalProjectName)
	if !found {
		response.WriteErrorString(http.StatusNotFound,
			fmt.Sprintf("Project not found '%s'", originalProjectName))
		return
	}

	if project.Name != originalProjectName {
		p.Remove(originalProjectName)
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
