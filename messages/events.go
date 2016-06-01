package messages

//go:generate stringer -type=FileChangeType

// FileChangeNotification is sent to the client
// when a change to the fileset had occurred.
type FileChangeNotification struct {
	// Filename is the file's name.
	Filename string `json:"filename"`
	// Action marks the operation that was performed
	// for the file.
	Action FileChangeType `json:"change_type"`
}

// FileChangeType marks the action type.
type FileChangeType uint8

const (
	// ChangeUnknown is the default value of FileChangeType.
	ChangeUnknown FileChangeType = iota
	// ChangeAddition is sent when a new file was added.
	ChangeAddition
	// ChangeRemoval is sent when the file was deleted.
	ChangeRemoval
	// ChangeModification is sent when the file was modified.
	ChangeModification
)
