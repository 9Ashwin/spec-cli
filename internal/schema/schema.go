package schema

import (
	"io/fs"
	"path/filepath"
	"strings"

	specfs "github.com/9Ashwin/spec-cli/embed"
	"github.com/9Ashwin/spec-cli/internal/vfs"
)

// Info describes an available schema bundle.
type Info struct {
	Name    string
	Version string
}

// ListSchemas returns available schema bundles from embed.
func ListSchemas() ([]Info, error) {
	entries, err := fs.ReadDir(specfs.SchemasFS, "assets/schemas")
	if err != nil {
		return nil, err
	}

	var schemas []Info
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		name := entry.Name()
		version := "unknown"

		versionPath := filepath.Join("assets/schemas", name, "VERSION")
		if data, err := specfs.SchemasFS.ReadFile(versionPath); err == nil {
			version = strings.TrimSpace(string(data))
		}

		schemas = append(schemas, Info{Name: name, Version: version})
	}

	return schemas, nil
}

// InstallSchema copies a schema bundle from embed to openspec/schemas/<name>/.
// Removes any existing installation first.
func InstallSchema(name, projectPath string) error {
	sourceDir := filepath.Join("assets/schemas", name)
	targetDir := filepath.Join(projectPath, "openspec", "schemas", name)

	_ = vfs.RemoveAll(targetDir)

	return fs.WalkDir(specfs.SchemasFS, sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		dest := filepath.Join(targetDir, relPath)

		if d.IsDir() {
			return vfs.MkdirAll(dest, 0o755)
		}

		data, err := specfs.SchemasFS.ReadFile(path)
		if err != nil {
			return err
		}

		if err := vfs.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}

		return vfs.WriteFile(dest, data, 0o644)
	})
}

// IsInstalled checks whether a schema is already installed.
func IsInstalled(name, projectPath string) bool {
	versionPath := filepath.Join(projectPath, "openspec", "schemas", name, "VERSION")
	_, err := vfs.Stat(versionPath)
	return err == nil
}

// GetInstalledVersion returns the installed schema version, or empty string.
func GetInstalledVersion(name, projectPath string) string {
	versionPath := filepath.Join(projectPath, "openspec", "schemas", name, "VERSION")
	data, err := vfs.ReadFile(versionPath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// AppendClaudeMdFragment appends the schema's CLAUDE.md fragment if not already present.
func AppendClaudeMdFragment(name, projectPath, locale string) (bool, error) {
	claudeMdPath := filepath.Join(projectPath, "CLAUDE.md")

	data, err := vfs.ReadFile(claudeMdPath)
	if err != nil {
		return false, nil // no CLAUDE.md, skip
	}

	content := string(data)
	if strings.Contains(content, name) {
		return false, nil // already present
	}

	var fragmentName string
	if locale == "zh" {
		fragmentName = "CLAUDE.md.fragment.zh.md"
	} else {
		fragmentName = "CLAUDE.md.fragment.md"
	}

	fragmentPath := filepath.Join("assets/schemas", name, "templates", "adopters", fragmentName)
	fragment, err := specfs.SchemasFS.ReadFile(fragmentPath)
	if err != nil {
		return false, nil // fragment not found, skip silently
	}

	newContent := content + "\n" + string(fragment)
	return true, vfs.WriteFile(claudeMdPath, []byte(newContent), 0o644)
}
