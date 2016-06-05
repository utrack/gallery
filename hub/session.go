package hub

import (
	"encoding/json"
	"github.com/utrack/gallery/client"
)

// disconMsg is a message about client getting disconnected.
// Contains the client's ID and discon reason.
type disconMsg struct {
	sessid uint64
	reason string
}

// session is an internal wrapper around the client's connection.
// It serves as a mediator between the hub and connection.
type session struct {
	// sessid is this session's ID.
	sessid uint64
	// conn is the underlying client connection.
	conn client.Connection
	// disconChan is the hub's discon msg channel.
	disconChan chan<- disconMsg
}

// newSession returns the session.
// Remember to execute runPump() for the session
// to become operational.
func newSession(id uint64, conn client.Connection, dChan chan<- disconMsg) *session {
	return &session{
		sessid:     id,
		conn:       conn,
		disconChan: dChan,
	}
}

// runPump starts the session's message pump.
func (s *session) runPump() {
	go s.disconPump()
}

// send forwards the message to the connection.
func (s *session) send(msg json.RawMessage) {
	s.conn.Send(msg)
}

// disconPump forwards discon messages from the client,
// joining them with the sessid.
func (s *session) disconPump() {
	err := <-s.conn.DisconChan()
	s.disconChan <- disconMsg{
		sessid: s.sessid,
		reason: err.Error(),
	}
	return
}
