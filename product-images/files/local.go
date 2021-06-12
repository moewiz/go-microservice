package files

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

// LocalStorage is an implementation of the Storage interface which works with
// the local disk on the current machine
type LocalStorage struct {
	basePath string
	maxSize  int // maximum number of bytes for files
}

// NewLocalStorage creates a new Local filesystem with the given base path
// basePath is the base didrectory to save files to
// maxSize is the max number of  bytes that  a file can be
func NewLocalStorage(basePath string, maxSize int) (*LocalStorage, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	return &LocalStorage{basePath: p, maxSize: maxSize}, nil
}

// returns the absolute path
func (l *LocalStorage) fullPath(path string) string {
	// append the given path to the base path
	return filepath.Join(l.basePath, path)
}

// Make sure the directory exists,
// if the file exists delete it
func checkLocation(fp string) error {
	// get the directory and make sure it exists
	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("Unable to create directory: %w", err)
	}

	// check the file
	_, err = os.Stat(fp)
	if err == nil {
		// the file is existing, delete it
		err = os.Remove(fp)
		if err != nil {
			return xerrors.Errorf("Unable to delete the file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		// if this is anything other than a not exists error
		return xerrors.Errorf("Unable to get file info: %w", err)
	}

	return nil
}

// Save the contents of the Writer to the given path
// path  is a relative path, basePath will be appended
func (l *LocalStorage) Save(path string, src io.Reader) error {
	// get the full path for the file
	fp := l.fullPath(path)

	if err := checkLocation(fp); err != nil {
		return err
	}

	// create a new file at the full path
	dst, err := os.Create(fp)
	if err != nil {
		return xerrors.Errorf("Unable to create file: %w", err)
	}
	defer dst.Close()

	// write the contents to the new file
	// ensure that we are not writing greather than max bytes
	_, err = io.Copy(dst, src)
	if err != nil {
		return xerrors.Errorf("Unable to write to file: %w", err)
	}

	return nil
}
