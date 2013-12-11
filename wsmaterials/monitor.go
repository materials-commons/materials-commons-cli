package wsmaterials

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/materials-commons/materials"
	"net/http"
	"os"
	"time"
)

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
			//fmt.Printf("File changed: %s\n", file)
			pfs := ProjectFileStatus{
				FilePath: file,
				Status:   "File Changed",
			}
			sio.Broadcast("file", &pfs)
		case folder := <-watcher.Folders:
			//fmt.Printf("Folder changed: %s\n", folder)
			pfs := ProjectFileStatus{
				FilePath: folder,
				Status:   "Directory Changed",
			}
			sio.Broadcast("dir", &pfs)
		}
	}
}
