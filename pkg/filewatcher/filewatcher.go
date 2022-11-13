// package filewatcher is to get events of file system in a realtime.
package filewatcher

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/howeyc/fsnotify"
)

type fileWatcherImpl struct {
	path    string
	watcher *fsnotify.Watcher

	stopCh chan struct{}
	stopWg sync.WaitGroup
	// TODO implement single-instance + multi event subsctibers fot the same path.
}

func New() FileWatcher {
	return &fileWatcherImpl{
		stopCh: make(chan struct{}, 1),
	}
}

func (fw *fileWatcherImpl) Start(path string, created, modified, deleted FileEventCallbackFn) error {
	if fw.watcher != nil { // TODO is not concurrently-safe
		return errors.New("create another file watcher instance or stop the watching first")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("fs watcher creation failed^ %w", err)
	}

	if err = watcher.Watch(path); err != nil {
		return fmt.Errorf("watching for '%s' failed^ %w", path, err)
	}

	fw.path = path
	fw.watcher = watcher

	fw.stopWg.Add(1)
	go func() {
		defer fw.stopWg.Done()
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() && created != nil {
					created(ev.Name)
				} else if ev.IsDelete() && deleted != nil {
					deleted(ev.Name)
				} else if ev.IsModify() && deleted != nil {
					modified(ev.Name)
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			case <-fw.stopCh:
				return
			}
		}
	}()
	return nil
}

func (fw *fileWatcherImpl) Stop() {
	fw.stopCh <- struct{}{}
	fw.stopWg.Wait()

	fw.watcher.RemoveWatch(fw.path)
	fw.watcher = nil
	fw.path = ""
}
