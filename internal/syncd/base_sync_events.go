package syncd

import (
	"context"
	"log"
	"strings"
	"sync"
	"syncer/internal/fsutils"
	"syncer/internal/logging"

	"github.com/fsnotify/fsnotify"
)

type SyncWatcher struct {
	sync.Mutex
	locked bool

	ctx context.Context
	ch  chan *SyncAction
	src string
	dst string
}

func (w *SyncWatcher) Disable() {
	w.Lock()
	defer w.Unlock()

	w.locked = true
}

func (w *SyncWatcher) Enable() {
	w.Lock()
	defer w.Unlock()

	w.locked = false
}

func (w *SyncWatcher) isLocked() bool {
	w.Lock()
	defer w.Unlock()

	return w.locked
}

func (w *SyncWatcher) Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer func() {
			logging.Log.Info("event listen stop")
			wg.Done()
		}()
		logging.Log.Info("event listen...")

	eventLoop:
		for {
			select {
			case <-w.ctx.Done():
				break eventLoop
			case event, ok := <-watcher.Events:
				if !ok {
					logging.Log.Error("event receive error")
					break
				}

				if w.isLocked() {
					break
				}

				var action SyncActionType

				if event.Has(fsnotify.Write) {
					action = MODIFY
				} else if event.Has(fsnotify.Create) {
					action = COPY
				} else if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					action = REMOVE
				}
				if isDir := fsutils.CheckDir(event.Name); !isDir {
					if action != "" {
						action := &SyncAction{
							src:      w.src,
							dst:      w.dst,
							filename: w.formatName(event.Name),
							action:   action,
						}
						w.send(action)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					logging.Log.Errorf("watcher error: %v", err)
					return
				}
			}
		}
	}()

	for _, dir := range w.getSubdirs() {
		err = watcher.Add(dir)
		logging.Log.Infof("subscribe to events on %v", dir)

		if err != nil {
			logging.Log.Errorf("add watcher error: %v", err)
		}
	}

	wg.Wait()
}

func (w *SyncWatcher) send(action *SyncAction) {
	select {
	case w.ch <- action:
		// message sent
	default:
		// message dropped
	}
}

func (w *SyncWatcher) formatName(filename string) string {
	return strings.Replace(filename, w.src, "", 1)
}

func (w *SyncWatcher) getSubdirs() []string {
	dirs, err := fsutils.GetSubdirs(w.src)
	if err != nil {
		logging.Log.Errorf("get subdirs error: %v", err)
	}
	return dirs
}
