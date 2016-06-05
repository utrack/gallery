package storage

import (
	"io"
)

// Saver saves the Reader's contents to the file.
type Saver interface {
	Upload(name string, r io.Reader) error
}
