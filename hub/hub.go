package hub

import (
	"github.com/pkg/errors"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/utrack/gallery/client"
	"github.com/utrack/gallery/messages"
	"github.com/utrack/gallery/storage"
	"log"
	"sync"
)

type hub struct {
	// Interfaces to the storage.
	notifier storage.Notifier
	lister   storage.Lister

	// Locker locks filenames before the upload.
	// It isn't allowed to upload two files with
	// the same filename at once.
	locker fileLocker

	// conns is a map of connid -> session.
	// session wraps the client's connection.
	conns map[uint64]*session
	// connMu protects the conns.
	connMu sync.RWMutex
	// lastConnId equals last ID given to the connection.
	lastConnId uint64

	// disconMsgs stores the discon messages coming from the
	// clients' sessions.
	disconMsgs chan disconMsg
}

// NewHub returns the Hub that uses passed Lister, Notifier and Saver.
// It is fully initiated and ready to accept connections.
func NewHub(list storage.Lister, notif storage.Notifier) Hub {
	ret := &hub{
		notifier: notif,
		lister:   list,

		locker: newFileLocker(),

		conns: make(map[uint64]*session),

		disconMsgs: make(chan disconMsg, 20),
	}
	go ret.pump()
	return ret
}

func (h *hub) Accept(c client.Connection) error {
	// retrieve the directory listing, err if failed
	list, err := h.lister.GetList()
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve the directory listing")
	}
	buf, err := ffjson.Marshal(list)
	if err != nil {
		// programmer's error, shouldn't happen
		panic(err)
	}
	defer ffjson.Pool(buf)

	// finally register the connection
	h.connMu.Lock()
	defer h.connMu.Unlock()
	h.lastConnId++

	// init session, send the listing
	sess := newSession(h.lastConnId, c, h.disconMsgs)
	sess.runPump()
	sess.send(buf)

	h.conns[sess.sessid] = sess
	return nil
}

// pump processes files' events and discon notifications.
func (h *hub) pump() {
	notifChan := h.notifier.GetNotificationChan()
	defer h.notifier.Close()
	for {
		select {
		case notif := <-notifChan:
			h.fanoutNotification(notif)
		case dMsg := <-h.disconMsgs:
			h.disconSession(dMsg.sessid)
			log.Printf("Connection %v dropped: %v", dMsg.sessid, dMsg.reason)
		}
	}
}

// fanoutNotification sends the notification to every
// client connected.
func (h *hub) fanoutNotification(n messages.FileChangeNotification) {
	// Marshal the notification's JSON
	buf, err := ffjson.Marshal(n)

	// that's totally unexpected, programmer's error
	// happens only with unmarshallable types (chan)
	if err != nil {
		panic(err)
	}

	// Pool the []byte buffer
	defer ffjson.Pool(buf)

	h.connMu.RLock()
	defer h.connMu.RUnlock()
	// fanout
	for _, sess := range h.conns {
		sess.send(buf)
	}
}

// disconSession removes a connection from the conns' dict.
// It is assumed that the connection had stopped itself
// (naturally, because the discon msg was rcvd).
func (h *hub) disconSession(sessId uint64) {
	h.connMu.Lock()
	defer h.connMu.Unlock()
	delete(h.conns, sessId)
}
