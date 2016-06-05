/*Package client provides common interface to the
clients and their connections.*/
package client

import (
	"encoding/json"
)

// Sender is able to send JSON messages.
// Errors returned via DisconChan.
type Sender interface {
	Send(json.RawMessage)
}

// Connection is a client's long-lived conn.
type Connection interface {
	Sender

	// DisconChan returns a client's disconnect notification channel,
	// which receives single value (error) at most. Error contains the
	// discon's description which describes why the client had gone away.
	DisconChan() <-chan error
}
