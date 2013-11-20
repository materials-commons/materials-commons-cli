package materials

import (
	"os"
)

type MaterialsCommons struct {
	user    *User
	baseUri string
}

func NewMaterialsCommons(user *User) *MaterialsCommons {
	mcurl := os.Getenv("MCAPIURL")
	if mcurl == "" {
		mcurl = "https://api.materialscommons.org"
	}

	return &MaterialsCommons{
		user:    user,
		baseUri: mcurl,
	}
}

func (mc *MaterialsCommons) ApiUrlPath(service string) string {
	uri := mc.baseUri + service + "?apikey=" + mc.user.Apikey
	return uri
}
