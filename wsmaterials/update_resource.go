package wsmaterials

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/materials-commons/materials"
)

type UpdateResource struct {
	// empty for now
}

func newUpdateResource(container *restful.Container) error {
	updateResource := UpdateResource{}
	updateResource.register(container)
	return nil
}

func (u *UpdateResource) register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/updates")

	ws.Route(ws.GET("/restart").To(u.restart).
		Doc("Restarts the materials service"))

	container.Add(ws)
}

func (u *UpdateResource) restart(request *restful.Request, response *restful.Response) {
	fmt.Println("Restarting materials service...")
	materials.Restart()
}
