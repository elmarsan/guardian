package files

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path"
	"path/filepath"
)

// LocalStorage implements Storage using disk.
type LocalStorage struct {
	basePath string
	// TODO remove path
	path string
}

// NewLocalStorage returns Storage implementation using disk.
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	return &LocalStorage{
		basePath: p,
		path:     basePath,
	}, nil
}

// Save saves new file in path with content contained in r.
func (l *LocalStorage) Save(fpath string, r io.Reader) error {
	// fp := l.fullPath(fpath)
	fp := fpath

	// get the directory and make sure it exists
	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Unable to create directory: %w", err)
	}

	// if the file exists delete it
	_, err = os.Stat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			return fmt.Errorf("Unable to delete file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		// if this is anything other than a not exists error
		return fmt.Errorf("Unable to get file info: %w", err)
	}

	// create a new file
	f, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("Unable to create file: %w", err)
	}
	defer f.Close()

	// write file content
	_, err = io.Copy(f, r)
	if err != nil {
		return fmt.Errorf("Unable to write to file: %w", err)
	}

	return nil
}

// Write writes file contained in path into given w.
func (l *LocalStorage) Write(fpath string, w io.Writer) (*FileInfo, error) {
	// fp := l.fullPath(fpath)
	fp := fpath

	// open the file
	f, err := os.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file: %w", err)
	}
	defer f.Close()

	// copy content to w
	s, err := io.Copy(w, f)
	if err != nil {
		return nil, fmt.Errorf("Unable to write file: %w", err)
	}

	// Figure out Content-type header from file extension
	ext := path.Ext(f.Name())
	mime := mime.TypeByExtension(ext)

	return &FileInfo{
		Name: f.Name(),
		Mime: mime,
		Path: fp,
		Size: s,
	}, nil
}

// GetAllInfo return all files contained in the storage.
func (l *LocalStorage) GetAllInfo() (*[]FileInfo, error) {
	files := []FileInfo{}

	// Walk through base path and get all files.
	err := filepath.Walk(l.path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !f.IsDir() {
			f := FileInfo{
				Path: path,
				Name: f.Name(),
			}

			files = append(files, f)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &files, nil
}

// fullPath returns the absolute path
func (l *LocalStorage) fullPath(path string) string {
	return filepath.Join(l.basePath, path)
}
