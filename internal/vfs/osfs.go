package vfs

import (
	"io/fs"
	"os"
)

// OsFs delegates every method to the os standard library.
type OsFs struct{}

func (OsFs) Stat(name string) (fs.FileInfo, error) { return os.Stat(name) }
func (OsFs) UserHomeDir() (string, error)          { return os.UserHomeDir() }
func (OsFs) ReadFile(name string) ([]byte, error)  { return os.ReadFile(name) }
func (OsFs) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}
func (OsFs) MkdirAll(path string, perm fs.FileMode) error { return os.MkdirAll(path, perm) }
func (OsFs) ReadDir(name string) ([]os.DirEntry, error)   { return os.ReadDir(name) }
func (OsFs) Remove(name string) error                     { return os.Remove(name) }
func (OsFs) RemoveAll(path string) error                  { return os.RemoveAll(path) }
