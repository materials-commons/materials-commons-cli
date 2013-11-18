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
	folders := subfolders(path)
	if len(folders) == 0 {
		return nil, errors.New("No folders to watch")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	rWatcher := &RecursiveWatcher{Watcher: watcher}
	rWatcher.Files = make(chan string, 10)
	rWatcher.Folders = make(chan string, len(folders))

	for _, folder := range folders {
		rWatcher.addFolder(folder)
	}

	return rWatcher, nil
}

func NewRecursiveWatcherPaths(paths []string) (*RecursiveWatcher, error) {
	if len(paths) == 0 {
		return nil, errors.New("No paths to watch")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	rWatcher := &RecursiveWatcher{Watcher: watcher}
	rWatcher.Files = make(chan string, 10)
	rWatcher.Folders = make(chan string, 100)

	for _, path := range paths {
		folders := subfolders(path)
		for _, folder := range folders {
			rWatcher.addFolder(folder)
		}
	}

	return rWatcher, nil
}

func (watcher *RecursiveWatcher) addFolder(folder string) {
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
						log.Printf("Error on stat for %s: %s\n", event.Name, err.Error())
					} else if finfo.IsDir() {
						watcher.addFolder(event.Name)
					} else {
						watcher.Files <- event.Name
					}

				case event.IsModify():
					finfo, err := os.Stat(event.Name)
					if err != nil {
						log.Printf("Error on stat for %s: %s\n", event.Name, err.Error())
					} else if !finfo.IsDir() {
						watcher.Files <- event.Name
					}
					log.Println("IsModify")

				case event.IsDelete():
					finfo, err := os.Stat(event.Name)
					// Don't do Stat because file no longer exists!
					if err != nil {
						log.Printf("Error on stat for %s: %s\n", event.Name, err.Error())
					} else if finfo.IsDir() {
						watcher.Folders <- event.Name
					} else {
						watcher.Files <- event.Name
					}

				case event.IsRename():
					finfo, err := os.Stat(event.Name)
					if err != nil {
						log.Printf("Error on stat for %s: %s\n", event.Name, err.Error())
					} else if finfo.IsDir() {
						watcher.Folders <- event.Name
					} else {
						watcher.Files <- event.Name
					}
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()
}

func subfolders(path string) (paths []string) {
	filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
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
