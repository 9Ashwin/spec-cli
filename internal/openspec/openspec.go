package openspec

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/9Ashwin/spec-cli/internal/vfs"
)

// Runner abstracts external command execution for testability.
type Runner interface {
	Run(dir, name string, args ...string) ([]byte, error)
	Available(name string) bool
}

// ExecRunner is the production implementation that shells out.
type ExecRunner struct{}

func (ExecRunner) Run(dir, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

func (ExecRunner) Available(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// DefaultRunner is the package-level runner. Tests may replace it.
var DefaultRunner Runner = ExecRunner{}

// LogWriter is where install progress is written. Tests may replace it.
var LogWriter io.Writer = os.Stderr

// InitOpenSpec runs openspec init with the given tool IDs.
// Auto-installs the openspec CLI via npm if not found on PATH.
func InitOpenSpec(projectPath string, toolIDs []string, scope string) error {
	if !DefaultRunner.Available("openspec") {
		if err := installOpenSpecCLI(projectPath); err != nil {
			return fmt.Errorf("openspec CLI not available and install failed: %w", err)
		}
		if !DefaultRunner.Available("openspec") {
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
	output, err := DefaultRunner.Run(projectPath, "openspec", args...)
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
	if !DefaultRunner.Available("openspec") {
		return nil, nil
	}

	output, err := DefaultRunner.Run(projectPath, "openspec", "list", "--json")
	if err != nil {
		return nil, nil //nolint:nilerr // status is best-effort; callers treat failures as no active changes.
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
		return nil, nil //nolint:nilerr // status is best-effort; invalid output is treated as no active changes.
	}
	return changes, nil
}

// Version returns the installed openspec CLI version string.
func Version() (string, error) {
	if !DefaultRunner.Available("openspec") {
		return "", fmt.Errorf("openspec not installed")
	}
	output, err := DefaultRunner.Run("", "openspec", "--version")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func installOpenSpecCLI(projectPath string) error {
	fmt.Fprintln(LogWriter, "    Installing OpenSpec CLI...")

	args := []string{"install", "-g", "@fission-ai/openspec@latest"}
	output, err := DefaultRunner.Run(projectPath, npmCmd(), args...)
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
