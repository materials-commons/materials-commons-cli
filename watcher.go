/*
* Based on code at: https://github.com/gophertown/looper/blob/master/watch.go
 */

package materials

import (
	"errors"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"path/filepath"
)

type Event struct {
	*fsnotify.FileEvent
}

type RecursiveWatcher struct {
	*fsnotify.Watcher
	Events chan Event
}

func NewRecursiveWatcher(path string) (*RecursiveWatcher, error) {
	directories := subdirs(path)
	if len(directories) == 0 {
		return nil, errors.New("No directories to watch")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	rWatcher := &RecursiveWatcher{Watcher: watcher}
	rWatcher.Events = make(chan Event, 10)

	for _, dir := range directories {
		rWatcher.addDirectory(dir)
	}

	return rWatcher, nil
}

func (watcher *RecursiveWatcher) addDirectory(dir string) {
	err := watcher.WatchFlags(dir, fsnotify.FSN_ALL)
	if err != nil {
		log.Println("Error watching directory: ", dir, err)
	}
}

func (watcher *RecursiveWatcher) Run() {
	go func() {
		for {
			select {
			case event := <-watcher.Event:
				watcher.handleEvent(event)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()
}

func (watcher *RecursiveWatcher) handleEvent(event *fsnotify.FileEvent) {
	e := Event{
		FileEvent: event,
	}
	if event.IsCreate() {
		finfo, err := os.Stat(event.Name)
		if err != nil {
			log.Printf("Error on stat for %s: %s\n", event.Name, err.Error())
		} else if finfo.IsDir() {
			watcher.addDirectory(event.Name)
		}
	}
	watcher.Events <- e
}

func subdirs(path string) (paths []string) {
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
