package wsmaterials

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/materials-commons/materials"
	"net/http"
)

func startMonitor() {
	sio := socketio.NewSocketIOServer(&socketio.Config{})
	sio.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// Nothing to do right now
	})
	go monitorProjectChanges(sio)
	go http.ListenAndServe(":8082", sio)
}

func monitorProjectChanges(sio *socketio.SocketIOServer) {
	p, _ := materials.CurrentUserProjects()

	var projectPaths []string
	for _, project := range p.Projects() {
		projectPaths = append(projectPaths, project.Path)
	}

	watcher, err := materials.NewRecursiveWatcherPaths([]string{"/tmp/a", "/tmp/b"})
	if err != nil {
		return
	}
	watcher.Run()
	defer watcher.Close()

	for {
		select {
		case file := <-watcher.Files:
			fmt.Printf("File changed: %s\n", file)
			pfs := ProjectFileStatus{
				FilePath: file,
				Status:   "File Changed",
			}
			sio.Broadcast("file", &pfs)
		case folder := <-watcher.Folders:
			fmt.Printf("Folder changed: %s\n", folder)
			pfs := ProjectFileStatus{
				FilePath: folder,
				Status:   "Directory Changed",
			}
			sio.Broadcast("dir", &pfs)
		}
	}
}
