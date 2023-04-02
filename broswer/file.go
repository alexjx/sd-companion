package broswer

import "time"

type File struct {
	// Path is the path of the file, include file name
	Path string `json:"path"`

	// Size is the size of the file
	Size int64 `json:"size"`

	// ModifiedAt is the time when the file is modified
	ModifiedAt time.Time `json:"modified_at"`

	// Metadata is the metadata of the file
	Metadata map[string]string `json:"metadata,omitempty"`
}
