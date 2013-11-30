package site

import (
	"fmt"
	"github.com/materials-commons/materials"
	"github.com/materials-commons/materials/wsmaterials"
	"net/http"
	"os"
	"time"
)

// Name of archive file for the materials website.
const materialsArchive = "materials.tar.gz"

// Start starts up all the webservices and the webserver.
func Start() {
	addr := setupSite()
	fmt.Println(http.ListenAndServe(addr, nil))
}

// StartRetry attempts a number of times to try connecting to the port address.
// This is useful when the server is restarting and the old server hasn't exited yet.
func StartRetry(retryCount int) {
	addr := setupSite()
	for i := 0; i < retryCount; i++ {
		fmt.Println(http.ListenAndServe(addr, nil))
		time.Sleep(1000 * time.Millisecond)
	}
	os.Exit(1)
}

// setupSite creates all the different web services for the http server.
// It returns the address and port the http server should use.
func setupSite() string {
	container := wsmaterials.NewRegisteredServicesContainer()
	http.Handle("/", container)
	dir := http.Dir(materials.Config.WebDir())
	http.Handle("/materials/", http.StripPrefix("/materials/", http.FileServer(dir)))
	addr := fmt.Sprintf("%s:%d", materials.Config.ServerAddress(), materials.Config.ServerPort())
	return addr
}
