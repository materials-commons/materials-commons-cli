package wsmaterials

import (
	"github.com/emicklei/go-restful"
	"fmt"
)

func NewRegisteredServicesContainer() *restful.Container {
	wsContainer := restful.NewContainer()

	if err := newProjectResource(wsContainer); err != nil {
		panic("Could not register ProjectResource")
	}

	return wsContainer
}

func JsonpFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	//callback := req.Request.FormValue("callback")
	fmt.Println("before ProcessFilter")
	chain.ProcessFilter(req, resp)
	fmt.Println("after ProcessFilter")
	fmt.Println(resp)
}
