package platform

import (
	"path/filepath"

	"github.com/9Ashwin/spec-cli/internal/vfs"
)

// DetectPlatforms detects which AI coding platforms are active in the given path.
// If a platform has DetectionPaths, checks those. Otherwise falls back to
// checking if the platform's SkillsDir exists.
func DetectPlatforms(projectPath string) []Platform {
	var detected []Platform

	for _, p := range AllPlatforms {
		if isPlatformDetected(projectPath, p) {
			detected = append(detected, p)
		}
	}

	return detected
}

func isPlatformDetected(projectPath string, p Platform) bool {
	if len(p.DetectionPaths) > 0 {
		for _, dp := range p.DetectionPaths {
			if _, err := vfs.Stat(filepath.Join(projectPath, dp)); err == nil {
				return true
			}
		}
		return false
	}

	// Fall back to checking if the skills directory exists.
	if _, err := vfs.Stat(filepath.Join(projectPath, p.SkillsDir)); err == nil {
		return true
	}
	return false
}
