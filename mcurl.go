package materials

import (
	"os"
)

type MaterialsCommons struct {
	user    *User
	baseUri string
}

func NewMaterialsCommons(user *User) *MaterialsCommons {
	mcurl := os.Getenv("MCURL")
	if mcurl == "" {
		mcurl = "https://api.materialscommons.org"
	}

	return &MaterialsCommons{
		user:    user,
		baseUri: mcurl,
	}
}

func (mc *MaterialsCommons) UrlPath(service string) string {
	return mc.baseUri + service + "?apikey=" + mc.user.Apikey
}
