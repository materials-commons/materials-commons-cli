package wsmaterials

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/materials-commons/materials"
	"net/http"
	"os"
	"time"
)

// Start starts up all the webservices and the webserver.
func Start() {
	setupSocketIO()
	addr := setupSite()
	fmt.Println(http.ListenAndServe(addr, nil))
}

type Example struct {
	Name  string
	Stuff string
}

func setupSocketIO() {
	sio := socketio.NewSocketIOServer(&socketio.Config{})
	sio.On("connect", func(ns *socketio.NameSpace) {
		fmt.Println("Connected:", ns.Id())
		sio.Broadcast("connected", ns.Id())
		ns.Emit("file", &Example{"Hello", "World"})
		sio.Broadcast("file", &Example{"From", "Broadcast"})
	})
	sio.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

	})

	go func() {
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
				sio.Broadcast("file", &ProjectFileStatus{file, "File Changed"})
			case folder := <-watcher.Folders:
				fmt.Printf("Folder changed: %s\n", folder)
			}
		}
	}()
	go http.ListenAndServe(":8082", sio)
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
	addr := fmt.Sprintf("%s:%d", materials.Config.Server.Address, materials.Config.Server.Port)
	return addr
}
