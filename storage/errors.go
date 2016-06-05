package storage

import (
	"errors"
)

// ErrBadFilename is returned by the Saver
// in response to Upload() if filename
// is empty.
var ErrBadFilename = errors.New("bad filename")
