package hub

import (
	"github.com/utrack/gallery/client"
)

// ConnectionAcceptor accepts incoming connections.
type ConnectionAcceptor interface {
	Accept(client.Connection) error
}

type Hub interface {
	ConnectionAcceptor
}
