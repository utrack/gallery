package storage

import (
	"github.com/pkg/errors"
	"github.com/utrack/gallery/messages"
	"io/ioutil"
	"os"
)

type lister struct {
	path string
}

// NewLister returns the Lister that returns listings
// for the given path.
func NewLister(path string) (Lister, error) {
	dir, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "error when opening directory")
	}

	// return error if not a directory
	if !dir.IsDir() {
		return nil, errors.New("not a directory")
	}
	return &lister{path: path}, nil
}

func (l *lister) GetList() (messages.FilesInfo, error) {
	files, err := ioutil.ReadDir(l.path)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't retrieve directory listing")
	}
	ret := make(messages.FilesInfo, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		ret = append(ret, messages.FileInfo{Filename: f.Name()})
	}
	return ret, nil
}
