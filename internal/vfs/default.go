package vfs

import (
	"io/fs"
	"os"
)

// DefaultFS is the global filesystem instance. Tests may replace it.
var DefaultFS FS = OsFs{}

// Package-level convenience functions that delegate to DefaultFS.

func Stat(name string) (fs.FileInfo, error) { return DefaultFS.Stat(name) }
func UserHomeDir() (string, error)          { return DefaultFS.UserHomeDir() }
func ReadFile(name string) ([]byte, error)  { return DefaultFS.ReadFile(name) }
func WriteFile(name string, data []byte, perm fs.FileMode) error {
	return DefaultFS.WriteFile(name, data, perm)
}
func MkdirAll(path string, perm fs.FileMode) error { return DefaultFS.MkdirAll(path, perm) }
func ReadDir(name string) ([]os.DirEntry, error)   { return DefaultFS.ReadDir(name) }
func Remove(name string) error                     { return DefaultFS.Remove(name) }
func RemoveAll(path string) error                  { return DefaultFS.RemoveAll(path) }
