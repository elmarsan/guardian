package files

import (
	"io"
	"mime"
	"os"
	"path"
)

// Storage defines file operations.
type Storage interface {
	// Save creates file in fpath from r.
	Save(fpath string, f io.Reader) error

	// GetAllInfo return all StorageFiles.
	GetAllInfo() (*[]StorageFile, error)

	// Write writes StorageFile into given w.
	Write(fpath string, w io.Writer) (*StorageFile, error)
}

// StorageFile represents file located in Storage.
type StorageFile struct {
	Path string
	Name string
	Mime string
	Size int64
	ETag string
}

// NewStorageFile creates StorageFile from file descriptor f.
func NewStorageFile(f *os.File) (*StorageFile, error) {
	// Figure out Content-type header from file extension
	ext := path.Ext(f.Name())
	mime := mime.TypeByExtension(ext)

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return &StorageFile{
		Name: f.Name(),
		Mime: mime,
		Path: f.Name(),
		Size: stat.Size(),
	}, nil
}
