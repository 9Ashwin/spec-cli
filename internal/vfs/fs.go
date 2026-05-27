package vfs

import (
	"io/fs"
	"os"
)

// FS abstracts filesystem operations used across the project.
// Implementations must behave identically to the corresponding os package functions.
type FS interface {
	// Query
	Stat(name string) (fs.FileInfo, error)
	UserHomeDir() (string, error)

	// Read/Write
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error

	// Directory/File management
	MkdirAll(path string, perm fs.FileMode) error
	ReadDir(name string) ([]os.DirEntry, error)
	Remove(name string) error
	RemoveAll(path string) error
}
