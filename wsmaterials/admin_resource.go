package wsmaterials

import (
	"github.com/emicklei/go-restful"
	"github.com/materials-commons/materials"
	"net/http"
	"os"
	"time"
)

type adminResource struct {
	// empty for now
}

func newAdminResource(container *restful.Container) error {
	adminResource := adminResource{}
	adminResource.register(container)
	return nil
}

func (ar *adminResource) register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/admin")

	ws.Route(ws.GET("/restart").To(ar.restart).
		Doc("Restarts the materials service"))

	ws.Route(ws.GET("/update").To(ar.update).
		Doc("If updates are available downloads, installs and restarts the server."))

	ws.Route(ws.GET("/stop").To(ar.stop).
		Doc("Stops the server."))

	container.Add(ws)
}

func (ar *adminResource) restart(request *restful.Request, response *restful.Response) {
	response.WriteErrorString(http.StatusOK, "Restarting materials service\n")
	go func() {
		sleep(1)
		materials.Restart()
	}()
}

func (ar *adminResource) update(request *restful.Request, response *restful.Response) {

}

func (ar *adminResource) stop(request *restful.Request, response *restful.Response) {
	response.WriteErrorString(http.StatusOK, "Stopping materials service\n")
	go func() {
		sleep(1)
		os.Exit(0)
	}()
}

func sleep(seconds time.Duration) {
	time.Sleep(seconds * time.Millisecond)
}
