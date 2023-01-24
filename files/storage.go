package files

import "io"

// FileInfo represents file server information containing abs path and name.
type FileInfo struct {
	Path string
	Name string
	Mime string
	Size int64
	ETag string
}

// Storage defines file operations.
type Storage interface {
	// Save creates file in fpath from r.
	Save(fpath string, f io.Reader) error

	// GetAllInfo return all files contained in the storage.
	GetAllInfo() (*[]FileInfo, error)

	// Write writes file contained in fpath into given w.
	Write(fpath string, w io.Writer) (*FileInfo, error)
}
