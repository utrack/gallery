package messages

// FileInfo represents the file's data.
type FileInfo struct {
	Filename string `json:"filename"`
}

// FilesInfo is a group of FileInfo.
type FilesInfo []FileInfo
