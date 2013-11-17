package materials

import (
	"errors"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"path/filepath"
)

type RecursiveWatcher struct {
	*fsnotify.Watcher
	Files   chan string
	Folders chan string
}

func NewRecursiveWatcher(path string) (*RecursiveWatcher, error) {
	folders := Subfolders(path)
	if len(folders) == 0 {
		return nil, errors.New("No folders to watch")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	rWatcher := &RecursiveWatcher{Watcher: watcher}
	rWatcher.Files = make(chan string, 10)
	rWatcher.Folders = make(chan string, len(folder))

	for _, folder := range folders {
		rWatcher.AddFolder(folder)
	}

	return rWatcher, nil
}

func (watcher *RecursiveWatcher) AddFolder(folder string) {
	err := watcher.WatchFlags(folder, fsnotify.FSN_ALL)
	if err != nil {
		log.Println("Error watching folder: ", folder, err)
	}
	watcher.Folders <- folder
}

func (watcher *RecursiveWatcher) Run() {
	go func() {
		for {
			select {
			case event := <-watcher.Event:
				switch {
				case event.IsCreate():
					finfo, err := os.Stat(event.Name)
					if err != nil {
						// do something
					} else if finfo.IsDir() {
						watcher.AddFolder(event.Name)
					} else {
						watcher.Files <- event.Name
					}

				case event.IsModify():
					log.Println("IsModify")
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()
}

func Subfolders(path string) (paths []string) {
	filepath.Walk(path, func(subpath, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			hidden := filepath.HasPrefix(name, ".") && name != "." && name != ".."
			if hidden {
				return filepath.SkipDir
			} else {
				paths = append(paths, subpath)
			}
		}
		return nil
	})
	return paths
}
