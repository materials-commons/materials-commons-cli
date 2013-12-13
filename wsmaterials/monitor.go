package wsmaterials

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/materials-commons/materials"
	"net/http"
	"os"
	"time"
)

type ProjectFileStatus struct {
	Project  string `json:"project"`
	FilePath string `json:"filepath"`
	Status   string `json:"status"`
}

func startMonitor() {
	sio := socketio.NewSocketIOServer(&socketio.Config{})
	sio.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// Nothing to do right now
	})
	go monitorProjectChanges(sio)
	go startHttp(10, sio)
}

func startHttp(retryCount int, sio *socketio.SocketIOServer) {
	for i := 0; i < retryCount; i++ {
		fmt.Println(http.ListenAndServe(":8082", sio))
		time.Sleep(1000 * time.Millisecond)
	}
	os.Exit(1)
}

func monitorProjectChanges(sio *socketio.SocketIOServer) {
	p, _ := materials.CurrentUserProjects()

	for _, project := range p.Projects() {
		startProjectWatcher(project, sio)
	}
}

func startProjectWatcher(project materials.Project, sio *socketio.SocketIOServer) {
	go func() {
		watcher, err := materials.NewRecursiveWatcher(project.Path)
		if err != nil {
			fmt.Println(err)
			return
		}
		watcher.Run()
		defer watcher.Close()

		for {
			select {
			case e := <-watcher.Events:
				pfs := &ProjectFileStatus{
					Project:  project.Name,
					FilePath: e.Name,
					Status:   eventStatus(e),
				}
				sio.Broadcast("file", pfs)
			}
		}
	}()
}

func eventStatus(event materials.Event) string {
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
