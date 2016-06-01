package filewatch

import (
	"github.com/ansel1/merry"
	"github.com/fsnotify/fsnotify"
	"github.com/utrack/gallery/messages"
)

type watcher struct {
	doneChan chan bool

	notificationChan chan messages.FileChangeNotification

	w *fsnotify.Watcher
}

// New returns the initiated Notifier that
// watches the specified directory.
func NewNotifier(path string) (Notifier, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, merry.Wrap(err)
	}

	err = fsWatcher.Add(path)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	ret := &watcher{
		doneChan:         make(chan bool, 2),
		notificationChan: make(chan messages.FileChangeNotification, 30),
		w:                fsWatcher,
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
