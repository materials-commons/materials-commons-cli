package wsmaterials

import (
	"github.com/emicklei/go-restful"
)

func NewRegisteredServicesContainer() *restful.Container {
	wsContainer := restful.NewContainer()

	if err := newProjectResource(wsContainer); err != nil {
		panic("Could not register ProjectResource")
	}

	if err := newAdminResource(wsContainer); err != nil {
		panic("Could not register AdminResource")
	}

	return wsContainer
}
