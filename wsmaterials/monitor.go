package wsmaterials

/*
 * monitor.go contains the routines for monitoring project files and directory changes. It
 * communicates with the frontend using socket.io. Each project is monitored separately.
 */

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/materials-commons/gohandy/fs"
	"github.com/materials-commons/materials"
	"net/http"
	"os"
	"time"
)

// projectFileStatus communicates the types of changes that
// have occured.
type projectFileStatus struct {
	Project  string `json:"project"`
	FilePath string `json:"filepath"`
	Event    string `json:"event"`
}

// startMonitor starts the monitor service and the HTTP and SocketIO connections.
func startMonitor() {
	sio := socketio.NewSocketIOServer(&socketio.Config{})
	sio.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// Nothing to do
	})
	go monitorProjectChanges(sio)
	go startHttp(10, sio)
}

// startHttp starts up a HTTP server. It will attempt to start the server
// retryCount times. The retry on server startup handles the case where
// the old materials service is stopping, and the new one is starting.
func startHttp(retryCount int, sio *socketio.SocketIOServer) {
	for i := 0; i < retryCount; i++ {
		fmt.Println(http.ListenAndServe(":8082", sio))
		time.Sleep(1000 * time.Millisecond)
	}
	os.Exit(1)
}

// monitorProjectChanges starts a separate go thread for each project to monitor.
func monitorProjectChanges(sio *socketio.SocketIOServer) {
	p, _ := materials.CurrentUserProjects()

	for _, project := range p.Projects() {
		go projectWatcher(project, sio)
	}
}

// projectWatcher starts the file system monitor. It watches for file system
// events and then communicates them along the SocketIOServer. It sends events
// to the front end as projectFileStatus messages encoded in JSON.
func projectWatcher(project materials.Project, sio *socketio.SocketIOServer) {
	watcher, err := fs.NewRecursiveWatcher(project.Path)
	if err != nil {
		fmt.Println(err)
		return
	}
	watcher.Start()
	defer watcher.Close()

	for {
		event := <-watcher.Events
		pfs := &projectFileStatus{
			Project:  project.Name,
			FilePath: event.Name,
			Event:    eventType(event),
		}
		sio.Broadcast("file", pfs)
	}
}

// eventType takes an event determines what type of event occurred and
// returns the corresponding string.
func eventType(event fs.Event) string {
	switch {
	case event.IsCreate():
		return "Created"
	case event.IsDelete():
		return "Deleted"
	case event.IsModify():
		return "Modified"
	case event.IsRename():
		return "Renamed"
	default:
		return "Unknown"
	}
}
