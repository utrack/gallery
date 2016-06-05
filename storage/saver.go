package storage

import (
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

type saver struct {
	path string
}

// NewSaver returns a Saver that uploads
// files to the specified directory.
func NewSaver(dirPath string) (Saver, error) {
	// Check if directory exists
	dirStat, err := os.Stat(dirPath)
	if err != nil {
		return nil, errors.Wrap(err, "directory not found")
	}
	if !dirStat.IsDir() {
		return nil, errors.New("path is not a directory")
	}
	return &saver{path: dirPath}, nil
}

func (u *saver) Upload(name string, r io.Reader) error {
	// Get the filename without directory part,
	// return ErrBadPath if no filename specified
	_, name = filepath.Split(name)
	if len(name) == 0 {
		return ErrBadFilename
	}

	// Open the file and write to it.
	// Should return error if there's pending write to the file
	// outside the app
	f, err := os.OpenFile(filepath.Join(u.path, name), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)

	return err
}
