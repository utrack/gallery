package hub

import (
	"github.com/utrack/gallery/client"
)

// ConnectionAcceptor accepts incoming connections.
type ConnectionAcceptor interface {
	Accept(client.Connection) error
}

// Hub routes notifications from the storage to clients.
type Hub interface {
	ConnectionAcceptor
}
