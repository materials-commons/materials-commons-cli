package site

import (
	"fmt"
	"github.com/materials-commons/gohandy/ezhttp"
	"github.com/materials-commons/materials"
	"github.com/materials-commons/materials/wsmaterials"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

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

func Download() (to string, err error) {
	const MaterialsFile = "materials.tar.gz"
	client := ezhttp.NewClient()
	url := fmt.Sprintf("%s/%s", materials.Config.MCDownload(), MaterialsFile)
	to = filepath.Join(materials.Config.DotMaterials(), MaterialsFile)
	status, err := client.FileGet(url, to)
	switch {
	case err != nil:
		return
	case status != 200:
		return to, fmt.Errorf("Download failed with HTTP status code %d", status)
	default:
		return to, nil
	}
}
