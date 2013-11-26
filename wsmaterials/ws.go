package wsmaterials

import (
	"github.com/emicklei/go-restful"
)

func NewRegisteredServicesContainer() *restful.Container {
	wsContainer := restful.NewContainer()

	if err := newProjectResource(wsContainer); err != nil {
		panic("Could not register ProjectResource")
	}

	if err := newUpdateResource(wsContainer); err != nil {
		panic("Could not register UpdateResource")
	}

	return wsContainer
}
