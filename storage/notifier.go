package storage

import (
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/utrack/gallery/messages"
)

type watcher struct {
	doneChan chan bool

	notificationChan chan messages.FileChangeNotification

	w *fsnotify.Watcher

	throttler *eventThrottler
}

// New returns the initiated Notifier that
// watches the specified directory.
func NewNotifier(path string) (Notifier, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't init watcher")
	}

	err = fsWatcher.Add(path)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't watch directory")
	}

	ret := &watcher{
		doneChan:         make(chan bool, 2),
		notificationChan: make(chan messages.FileChangeNotification, 30),
		w:                fsWatcher,
		throttler:        newEventThrottler(),
	}
	go ret.pump()
	return ret, nil
}

// pump processes the fsnotify events.
func (w *watcher) pump() {
	defer w.w.Close()
	for {
		select {
		case ev := <-w.w.Events:
			w.throttler.pushEvent(ev)
		case ev := <-w.throttler.outChan:
			notification := fsEventToNotification(ev)
			w.notificationChan <- notification
		case <-w.doneChan:
			return
		}
	}
}

func (w *watcher) Close() error {
	w.doneChan <- true
	return nil
}

func (w *watcher) GetNotificationChan() <-chan messages.FileChangeNotification {
	return w.notificationChan
}

// fsEventNotification converts the changes' info from fsnotify format
// to the internal messages.FileChangeNotifications.
func fsEventToNotification(ev fsnotify.Event) messages.FileChangeNotification {
	var ret messages.FileChangeNotification

	ret.Filename = ev.Name
	switch ev.Op {
	case fsnotify.Create:
		ret.Action = messages.ChangeAddition
	case fsnotify.Remove:
		ret.Action = messages.ChangeRemoval
	case fsnotify.Write:
		ret.Action = messages.ChangeModification
	case fsnotify.Rename:
		ret.Action = messages.ChangeRename
	}
	return ret
}
