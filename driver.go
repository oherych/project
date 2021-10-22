package project

import (
	"os"
	"path"
)

// DriverInterface describes communication interface with file system
type DriverInterface interface {
	CreateDir(filePath string, perm os.FileMode) error
	CreateFile(filePath string, data []byte, perm os.FileMode) error
}

// FileSystem is driver for work with local file system
// Root - basic path
type FileSystem struct {
	Root string
}

// CreateDir create new folder located on `Root` path and file path from arguments
// new folder will be created with provided permission
// If there is an error, it will be of type *PathError.
func (driver FileSystem) CreateDir(filePath string, perm os.FileMode) error {
	return os.Mkdir(path.Join(driver.Root, filePath), perm)
}

// CreateFile create new file located on `Root` path and file path from arguments
// new file will be created with provided permission
// If there is an error, it will be of type *PathError.
func (driver FileSystem) CreateFile(filePath string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(path.Join(driver.Root, filePath), os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}

	return err
}
