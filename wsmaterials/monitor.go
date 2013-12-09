package wsmaterials

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/materials-commons/materials"
)

func (p *ProjectResource) monitorEventLoop() {
	var projectPaths []string
	for _, project := range p.Projects() {
		projectPaths = append(projectPaths, project.Path)
	}

	sio := socketio.NewSocketIOServer(&socketio.Config{})

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
			p.events[0] = ProjectFileStatus{
				FilePath: file,
				Status:   "File Changed",
			}
			sio.Broadcast("file", file)
		case folder := <-watcher.Folders:
			fmt.Printf("Folder changed: %s\n", folder)
			p.events[0] = ProjectFileStatus{
				FilePath: folder,
				Status:   "Directory Changed",
			}
		}
	}
}
