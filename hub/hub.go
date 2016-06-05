package hub

import (
	"github.com/utrack/gallery/client"
	"sync"
)

type hub struct {
	locker fileLocker
	path   string

	// conns is a map of connid -> connection
	conns      map[uint64]*session
	connMu     sync.RWMutex
	lastConnId uint64

	disconMsgs chan disconMsg
}

func NewHub(path string) interface{} {
	return &hub{
		locker: newFileLocker(),
		path:   path,
	}
}

func (h *hub) Accept(c client.Connection) error {
	h.connMu.Lock()
	defer h.connMu.Unlock()
	h.lastConnId++

	sess := newSession(h.lastConnId, c, h.disconMsgs)
	sess.runPump()

	h.conns[sess.sessid] = sess
	return nil
}
