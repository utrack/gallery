package hub

import (
	"github.com/utrack/gallery/client"
)

// fileLocker is used to obtain the exclusive access
// to a file.
type fileLocker interface {
	// LockIfNotExists locks the file with given name if it does not exist.
	LockIfNotExists(name string) (fileLock, error)
}

// fileLock is used to drop exclusive lock from a file.
type fileLock interface {
	Unlock()
}

// ConnectionAcceptor accepts incoming connections.
type ConnectionAcceptor interface {
	Accept(client.Connection) error
}

type Hub interface {
	ConnectionAcceptor
}
