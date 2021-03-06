package storage

import (
	"github.com/utrack/gallery/messages"
)

// Notifier pipes the filesystem change notifications to the output channel.
// Use GetNotificationChan() to retrieve it.
type Notifier interface {
	GetNotificationChan() <-chan messages.FileChangeNotification
	// Close shuts down the Notifier.
	Close() error
}
