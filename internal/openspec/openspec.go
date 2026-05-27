package openspec

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/9Ashwin/spec-cli/internal/vfs"
)

// InitOpenSpec runs openspec init with the given tool IDs.
// Auto-installs the openspec CLI via npm if not found on PATH.
func InitOpenSpec(projectPath string, toolIDs []string, scope string) error {
	if !commandAvailable("openspec") {
		if err := installOpenSpecCLI(projectPath); err != nil {
			return fmt.Errorf("openspec CLI not available and install failed: %w", err)
		}
		if !commandAvailable("openspec") {
			return fmt.Errorf("openspec CLI install completed but still not found on PATH")
		}
	}

	targetPath := projectPath
	if scope == "global" {
		home, err := vfs.UserHomeDir()
		if err != nil {
			return err
		}
		targetPath = home
	}

	args := []string{"init", targetPath, "--tools", strings.Join(toolIDs, ",")}
	cmd := exec.Command("openspec", args...)
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("openspec init failed: %w\n%s", err, string(output))
	}

	return nil
}

// ChangeInfo is a parsed entry from openspec list --json.
type ChangeInfo struct {
	Name      string            `json:"name"`
	Schema    string            `json:"schema,omitempty"`
	Artifacts map[string]string `json:"artifacts,omitempty"`
}

// ListChanges returns active changes from openspec list --json.
// Returns nil slice if openspec is not installed or the command fails.
func ListChanges(projectPath string) ([]ChangeInfo, error) {
	if !commandAvailable("openspec") {
		return nil, nil
	}

	cmd := exec.Command("openspec", "list", "--json")
	cmd.Dir = projectPath
	output, err := cmd.Output()
	if err != nil {
		return nil, nil
	}

	// Try envelope format first: {"changes": [...]}
	var result struct {
		Changes []ChangeInfo `json:"changes"`
	}
	if err := json.Unmarshal(output, &result); err == nil {
		return result.Changes, nil
	}

	// Fall back to flat array: [...]
	var changes []ChangeInfo
	if err := json.Unmarshal(output, &changes); err != nil {
		return nil, nil
	}
	return changes, nil
}

// Version returns the installed openspec CLI version string.
func Version() (string, error) {
	if !commandAvailable("openspec") {
		return "", fmt.Errorf("openspec not installed")
	}
	cmd := exec.Command("openspec", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func commandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func installOpenSpecCLI(projectPath string) error {
	fmt.Fprintln(os.Stderr, "    Installing OpenSpec CLI...")

	args := []string{"install", "-g", "@fission-ai/openspec@latest"}
	cmd := exec.Command(npmCmd(), args...)
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm install failed: %w\n%s", err, string(output))
	}
	return nil
}

func npmCmd() string {
	if runtime.GOOS == "windows" {
		return "npm.cmd"
	}
	return "npm"
}
