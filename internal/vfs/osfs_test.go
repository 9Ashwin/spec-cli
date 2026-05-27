package vfs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOsFs_MkdirAll_and_Stat(t *testing.T) {
	tmp := t.TempDir()
	fs := OsFs{}

	dir := filepath.Join(tmp, "a", "b", "c")
	if err := fs.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	info, err := fs.Stat(dir)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected directory")
	}
}

func TestOsFs_WriteFile_and_ReadFile(t *testing.T) {
	tmp := t.TempDir()
	fs := OsFs{}

	p := filepath.Join(tmp, "test.txt")
	content := []byte("hello spec-cli")

	if err := fs.WriteFile(p, content, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got, err := fs.ReadFile(p)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("got %q, want %q", got, content)
	}
}

func TestOsFs_RemoveAll(t *testing.T) {
	tmp := t.TempDir()
	fs := OsFs{}

	dir := filepath.Join(tmp, "removeme")
	_ = fs.MkdirAll(dir, 0o755)
	_ = fs.WriteFile(filepath.Join(dir, "f.txt"), []byte("x"), 0o644)

	if err := fs.RemoveAll(dir); err != nil {
		t.Fatalf("RemoveAll: %v", err)
	}

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Error("directory should not exist after RemoveAll")
	}
}

func TestDefaultFS_isOsFs(t *testing.T) {
	_, ok := DefaultFS.(OsFs)
	if !ok {
		t.Errorf("DefaultFS is %T, expected OsFs", DefaultFS)
	}
}

func TestPackageLevelFunctions(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "pkg.txt")

	if err := WriteFile(p, []byte("pkg"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	data, err := ReadFile(p)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != "pkg" {
		t.Errorf("got %q, want %q", data, "pkg")
	}

	info, err := Stat(p)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Size() != 3 {
		t.Errorf("size=%d, want 3", info.Size())
	}
}
