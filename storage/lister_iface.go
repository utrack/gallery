package storage

import (
	"github.com/utrack/gallery/messages"
)

// Lister returns the directory's listing.
// Only files are listed.
type Lister interface {
	GetList() (messages.FilesInfo, error)
}
