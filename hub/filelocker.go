package hub

import (
	"sync"
)

type filelocker struct {
	files map[string]interface{}
	mu    sync.Mutex
}

func newFileLocker() fileLocker {
	return &filelocker{
		files: make(map[string]interface{}),
	}
}

func (l *filelocker) LockIfNotExists(name string) (fileLock, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Lock the filename
	_, ok := l.files[name]
	if ok {
		return nil, ErrFileExists
	}
	l.files[name] = nil

	// Return FileLock with callback set
	return filelock{
		f:    l.unlock,
		name: name,
	}, nil
}

// unlock is called by filelock on filelock.Unlock().
func (l *filelocker) unlock(name string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.files, name)
}

// filelock represents obtained lock on the file.
type filelock struct {
	// f is a function that should be called on Unlock().
	f func(string)
	// name is this file's name.
	name string
}

func (l *filelock) Unlock() {
	l.f(l.name)
}
