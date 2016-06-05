package storage

import (
	"github.com/fsnotify/fsnotify"
	"sync"
	"time"
)

const (
	tickPeriod = time.Second * 2
)

// eventThrottler throttles the modification
// events and sends out only first notification
// in a row.
type eventThrottler struct {
	// map of filename -> throttle
	files   map[string]*throttle
	filesMu sync.RWMutex

	outChan chan fsnotify.Event
}

type throttle struct {
	timer *time.Timer
	event fsnotify.Event
}

func newEventThrottler() *eventThrottler {
	return &eventThrottler{
		files:   map[string]*throttle{},
		outChan: make(chan fsnotify.Event, 10),
	}
}

func (e *eventThrottler) pushEvent(ev fsnotify.Event) {
	if ev.Op != fsnotify.Write {
		e.dumpEventIfExists(ev.Name)
		e.deleteThrottle(ev.Name)
	}

	e.createThrottleIfNotExists(ev.Name)

	e.filesMu.RLock()
	defer e.filesMu.RUnlock()

	thr := e.files[ev.Name]
	thr.timer.Stop()

	// Reset the event if not modification
	// If event isn't empty
	if ev.Op != fsnotify.Write || thr.event.Name == `` {
		thr.event = ev
	}
	thr.timer.Reset(tickPeriod)
}

func (e *eventThrottler) dumpEventIfExists(fname string) {
	e.filesMu.RLock()
	defer e.filesMu.RUnlock()

	thr, ok := e.files[fname]
	if !ok {
		return
	}

	thr.timer.Stop()
	e.outChan <- thr.event
}

func (e *eventThrottler) deleteThrottle(fname string) {
	e.filesMu.Lock()
	defer e.filesMu.Unlock()
	delete(e.files, fname)
}

func (e *eventThrottler) createThrottleIfNotExists(fname string) {
	e.filesMu.Lock()
	defer e.filesMu.Unlock()
	// Skip creation if not exists
	_, ok := e.files[fname]
	if ok {
		return
	}
	// Create new throttle that will dump the event after tickPeriod
	thr := &throttle{}
	thr.timer = time.AfterFunc(tickPeriod, func() {
		e.outChan <- thr.event
		e.deleteThrottle(fname)
	})
	e.files[fname] = thr
}
