package wsmaterials

import (
	"github.com/emicklei/go-restful"
)

func NewRegisteredServicesContainer() *restful.Container {
	wsContainer := restful.NewContainer()

	if err := newProjectResource(wsContainer); err != nil {
		panic("Could not register ProjectResource")
	}

	return wsContainer
}
