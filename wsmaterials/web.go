package wsmaterials

import (
	"fmt"
	"github.com/materials-commons/materials"
	"net/http"
	"os"
	"time"
)

// Start starts up all the webservices and the webserver.
func Start() {
	startMonitor()
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
	container := NewRegisteredServicesContainer()
	http.Handle("/", container)
	dir := http.Dir(materials.Config.Server.Webdir)
	http.Handle("/materials/", http.StripPrefix("/materials/", http.FileServer(dir)))
	setupProjects()
	addr := fmt.Sprintf("%s:%d", materials.Config.Server.Address, materials.Config.Server.Port)
	return addr
}

func setupProjects() {
	projects, _ := materials.CurrentUserProjects()
	for _, project := range projects.Projects() {
		projectUrlPath := fmt.Sprintf("/%s/", project.Name)
		fmt.Printf("Setting up '%s' path '%s'\n", projectUrlPath, project.Path)
		dir := http.Dir(project.Path)
		http.Handle(projectUrlPath, http.StripPrefix(projectUrlPath, http.FileServer(dir)))
	}
}
