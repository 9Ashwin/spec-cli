package skill

import (
	"io/fs"
	"path/filepath"

	specfs "github.com/9Ashwin/spec-cli/embed"
	"github.com/9Ashwin/spec-cli/internal/vfs"
)

const (
	LangEN = "en"
	LangZH = "zh"
)

// CopySkills copies the Comet entry skill from embed to the target platform's
// skills directory. Returns (copied, skipped, error).
func CopySkills(baseDir, skillsDir, language string, overwrite bool) (int, int, error) {
	efs := specfs.SkillsFS
	sourceDir := "assets/skills"
	if language == LangZH {
		efs = specfs.SkillsZHFS
		sourceDir = "assets/skills-zh"
	}

	var copied, skipped int

	err := fs.WalkDir(efs, sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		dest := filepath.Join(baseDir, skillsDir, "skills", relPath)

		if !overwrite {
			if _, statErr := vfs.Stat(dest); statErr == nil {
				skipped++
				return nil
			}
		}

		data, readErr := efs.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		if mkdirErr := vfs.MkdirAll(filepath.Dir(dest), 0o755); mkdirErr != nil {
			return mkdirErr
		}

		if writeErr := vfs.WriteFile(dest, data, 0o644); writeErr != nil {
			return writeErr
		}

		copied++
		return nil
	})

	return copied, skipped, err
}
