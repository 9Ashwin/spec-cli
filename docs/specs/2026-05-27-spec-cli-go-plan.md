# Spec-CLI Go Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use `superpowers:subagent-driven-development` (recommended) or `superpowers:executing-plans` to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Clean-room Go rewrite of Comet CLI as `spec-cli` — install OpenSpec, Superpowers, and schema bundles into AI coding platform projects.

**Architecture:** Single binary via cobra commands. Assets embedded via `//go:embed`. All filesystem ops through `vfs.FS` interface (testable). Interactive prompts via charmbracelet/huh. Follows [lark-cli](https://github.com/larksuite/cli) patterns exactly.

**Module:** `github.com/9Ashwin/spec-cli` | **Go:** 1.23+

**Reference:** `/Users/solariswu/workspaces/github/cli` (lark-cli) — vfs interface, build ldflags, cobra root, embed patterns

---

### Task 1: Clean up Node.js and initialize Go module

**Files:**
- Delete: `src/`, `test/`, `scripts/`, `bin/`, `node_modules/`
- Delete: `package.json`, `package-lock.json`, `pnpm-lock.yaml`
- Delete: `tsconfig.json`, `vitest.config.ts`, `eslint.config.js`, `build.js`
- Delete: `assets/manifest.json`
- Create: `go.mod`, `main.go`

- [ ] **Step 1: Remove all Node.js artifacts and init Go module**

```bash
cd /Users/solariswu/workspaces/github/spec-cli
rm -rf src/ test/ scripts/ bin/ node_modules/
rm -f package.json package-lock.json pnpm-lock.yaml
rm -f tsconfig.json vitest.config.ts eslint.config.js build.js
rm -f assets/manifest.json
go mod init github.com/9Ashwin/spec-cli
```

- [ ] **Step 2: Create main.go**

```go
package main

import (
	"os"

	"github.com/9Ashwin/spec-cli/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
```

- [ ] **Step 3: Commit**

```bash
git add -A && git commit -m "chore: remove Node.js, initialize Go module

Replace TypeScript codebase with Go module github.com/9Ashwin/spec-cli.
Keep assets/ directory for embedded skills and schemas."
```

---

### Task 2: Create internal/build and internal/vfs

**Files:**
- `internal/build/build.go`
- `internal/vfs/fs.go`
- `internal/vfs/osfs.go`
- `internal/vfs/default.go`

Pattern: exact match with lark-cli's `internal/build/build.go` and `internal/vfs/*.go`.

- [ ] **Step 1: Create directory structure**

```bash
mkdir -p internal/build internal/vfs
```

- [ ] **Step 2: Write internal/build/build.go** (identical to lark-cli)

```go
package build

import "runtime/debug"

// Version is dynamically set by -ldflags or falls back to module info.
var Version = "DEV"

// Date is the build date in YYYY-MM-DD format, set by -ldflags.
var Date = ""

func init() {
	if Version == "DEV" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
	if Version == "" {
		Version = "DEV"
	}
}
```

- [ ] **Step 3: Write internal/vfs/fs.go** (interface only, subset for spec-cli needs)

```go
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
```

- [ ] **Step 4: Write internal/vfs/osfs.go** (OsFs struct, each method delegates to os package)

```go
package vfs

import (
	"io/fs"
	"os"
)

// OsFs delegates every method to the os standard library.
type OsFs struct{}

func (OsFs) Stat(name string) (fs.FileInfo, error)        { return os.Stat(name) }
func (OsFs) UserHomeDir() (string, error)                  { return os.UserHomeDir() }
func (OsFs) ReadFile(name string) ([]byte, error)          { return os.ReadFile(name) }
func (OsFs) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}
func (OsFs) MkdirAll(path string, perm fs.FileMode) error { return os.MkdirAll(path, perm) }
func (OsFs) ReadDir(name string) ([]os.DirEntry, error)    { return os.ReadDir(name) }
func (OsFs) Remove(name string) error                      { return os.Remove(name) }
func (OsFs) RemoveAll(path string) error                   { return os.RemoveAll(path) }
```

- [ ] **Step 5: Write internal/vfs/default.go** (DefaultFS + package-level convenience functions)

```go
package vfs

import (
	"io/fs"
	"os"
)

// DefaultFS is the global filesystem instance. Tests may replace it.
var DefaultFS FS = OsFs{}

// Package-level convenience functions that delegate to DefaultFS.

func Stat(name string) (fs.FileInfo, error)        { return DefaultFS.Stat(name) }
func UserHomeDir() (string, error)                  { return DefaultFS.UserHomeDir() }
func ReadFile(name string) ([]byte, error)          { return DefaultFS.ReadFile(name) }
func WriteFile(name string, data []byte, perm fs.FileMode) error {
	return DefaultFS.WriteFile(name, data, perm)
}
func MkdirAll(path string, perm fs.FileMode) error { return DefaultFS.MkdirAll(path, perm) }
func ReadDir(name string) ([]os.DirEntry, error)    { return DefaultFS.ReadDir(name) }
func Remove(name string) error                      { return DefaultFS.Remove(name) }
func RemoveAll(path string) error                   { return DefaultFS.RemoveAll(path) }
```

- [ ] **Step 6: Verify compilation**

```bash
go build ./internal/build/ ./internal/vfs/
```

- [ ] **Step 7: Commit**

```bash
git add internal/build/ internal/vfs/
git commit -m "feat: add internal/build and internal/vfs packages"
```

---

### Task 3: Create internal/platform

**Files:**
- `internal/platform/platform.go` — 29 platform definitions + helpers
- `internal/platform/detect.go` — platform detection + hasSkills

Source of truth: Comet's `src/core/platforms.ts` and `src/core/detect.ts`.

- [ ] **Step 1: Create directory**

```bash
mkdir -p internal/platform
```

- [ ] **Step 2: Write internal/platform/platform.go**

Platform IDs and openspecToolIds must match Comet's `PLATFORMS` array exactly:

```go
package platform

// Platform represents an AI coding platform, matching Comet's Platform interface.
type Platform struct {
	ID              string   // "claude", "cursor", "roocode", ...
	Name            string   // "Claude Code", "Cursor", "RooCode", ...
	SkillsDir       string   // e.g. ".claude", ".cursor"
	GlobalSkillsDir string   // optional, e.g. ".gemini/antigravity" for antigravity
	DetectionPaths  []string // paths checked for detection; nil means fall back to SkillsDir
	OpenSpecToolID  string   // tool ID passed to openspec init --tools
}

// AllPlatforms lists all 29 supported platforms.
// Source: Comet src/core/platforms.ts PLATFORMS array.
var AllPlatforms = []Platform{
	{ID: "claude", Name: "Claude Code", SkillsDir: ".claude", OpenSpecToolID: "claude"},
	{ID: "cursor", Name: "Cursor", SkillsDir: ".cursor", OpenSpecToolID: "cursor"},
	{ID: "codex", Name: "Codex", SkillsDir: ".codex", OpenSpecToolID: "codex"},
	{ID: "opencode", Name: "OpenCode", SkillsDir: ".opencode", OpenSpecToolID: "opencode"},
	{ID: "windsurf", Name: "Windsurf", SkillsDir: ".windsurf", OpenSpecToolID: "windsurf"},
	{ID: "cline", Name: "Cline", SkillsDir: ".cline", OpenSpecToolID: "cline"},
	{ID: "roocode", Name: "RooCode", SkillsDir: ".roo", OpenSpecToolID: "roocode"},
	{ID: "continue", Name: "Continue", SkillsDir: ".continue", OpenSpecToolID: "continue"},
	{
		ID: "github-copilot", Name: "GitHub Copilot", SkillsDir: ".github",
		DetectionPaths: []string{
			".github/copilot-instructions.md", ".github/instructions",
			".github/prompts", ".github/skills",
		},
		OpenSpecToolID: "github-copilot",
	},
	{ID: "gemini", Name: "Gemini CLI", SkillsDir: ".gemini", OpenSpecToolID: "gemini"},
	{ID: "amazon-q", Name: "Amazon Q Developer", SkillsDir: ".amazonq", OpenSpecToolID: "amazon-q"},
	{ID: "qwen", Name: "Qwen Code", SkillsDir: ".qwen", OpenSpecToolID: "qwen"},
	{ID: "kilocode", Name: "Kilo Code", SkillsDir: ".kilocode", OpenSpecToolID: "kilocode"},
	{ID: "auggie", Name: "Auggie (Augment CLI)", SkillsDir: ".augment", OpenSpecToolID: "auggie"},
	{ID: "kiro", Name: "Kiro", SkillsDir: ".kiro", OpenSpecToolID: "kiro"},
	{ID: "lingma", Name: "Lingma", SkillsDir: ".lingma", OpenSpecToolID: "lingma"},
	{ID: "junie", Name: "Junie", SkillsDir: ".junie", OpenSpecToolID: "junie"},
	{ID: "codebuddy", Name: "CodeBuddy Code", SkillsDir: ".codebuddy", OpenSpecToolID: "codebuddy"},
	{ID: "costrict", Name: "CoStrict", SkillsDir: ".cospec", OpenSpecToolID: "costrict"},
	{ID: "crush", Name: "Crush", SkillsDir: ".crush", OpenSpecToolID: "crush"},
	{ID: "factory", Name: "Factory Droid", SkillsDir: ".factory", OpenSpecToolID: "factory"},
	{ID: "iflow", Name: "iFlow", SkillsDir: ".iflow", OpenSpecToolID: "iflow"},
	{ID: "pi", Name: "Pi", SkillsDir: ".pi", OpenSpecToolID: "pi"},
	{ID: "qoder", Name: "Qoder", SkillsDir: ".qoder", OpenSpecToolID: "qoder"},
	{
		ID: "antigravity", Name: "Antigravity",
		SkillsDir: ".agents", GlobalSkillsDir: ".gemini/antigravity",
		OpenSpecToolID: "antigravity",
	},
	{ID: "bob", Name: "Bob Shell", SkillsDir: ".bob", OpenSpecToolID: "bob"},
	{ID: "forgecode", Name: "ForgeCode", SkillsDir: ".forge", OpenSpecToolID: "forgecode"},
	{ID: "trae", Name: "Trae", SkillsDir: ".trae", OpenSpecToolID: "trae"},
}

// ByID returns the platform with the given ID, or nil if not found.
func ByID(id string) *Platform {
	for i := range AllPlatforms {
		if AllPlatforms[i].ID == id {
			return &AllPlatforms[i]
		}
	}
	return nil
}
```

- [ ] **Step 3: Write internal/platform/detect.go**

Matches Comet's `detect.ts` logic:

```go
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
```

- [ ] **Step 4: Verify compilation and commit**

```bash
go build ./internal/platform/
git add internal/platform/ && git commit -m "feat: add internal/platform — 29 platforms + detection"
```

---

### Task 4: Create internal/openspec

**Files:**
- `internal/openspec/openspec.go`

- [ ] **Step 1: Create directory**

```bash
mkdir -p internal/openspec
```

- [ ] **Step 2: Write internal/openspec/openspec.go**

```go
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
```

- [ ] **Step 3: Verify compilation and commit**

```bash
go build ./internal/openspec/
git add internal/openspec/ && git commit -m "feat: add internal/openspec for openspec CLI integration"
```

---

### Task 5: Create embed package + internal/skill

**Files:**
- `embed/skills.go`
- `embed/schemas.go`
- `internal/skill/skill.go`

**Critical:** `//go:embed` cannot use `..` paths, so the embed package lives at module root. `internal/skill` imports from `embed`.

- [ ] **Step 1: Create embed package**

```bash
mkdir -p embed internal/skill
```

- [ ] **Step 2: Write embed/skills.go**

```go
package embed

import "embed"

//go:embed assets/skills/*
var SkillsFS embed.FS

//go:embed assets/skills-zh/*
var SkillsZHFS embed.FS
```

- [ ] **Step 3: Write embed/schemas.go**

```go
package embed

import "embed"

//go:embed assets/schemas/*
var SchemasFS embed.FS
```

- [ ] **Step 4: Write internal/skill/skill.go**

```go
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
	if language == LangZH {
		efs = specfs.SkillsZHFS
	}

	var copied, skipped int

	err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		dest := filepath.Join(baseDir, skillsDir, "skills", path)

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
```

- [ ] **Step 5: Verify compilation and commit**

```bash
go build ./embed/ ./internal/skill/
git add embed/ internal/skill/ && git commit -m "feat: add embed package and internal/skill"
```

---

### Task 6: Create internal/schema

**Files:**
- `internal/schema/schema.go`

- [ ] **Step 1: Create directory**

```bash
mkdir -p internal/schema
```

- [ ] **Step 2: Write internal/schema/schema.go**

```go
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

	fragmentPath := filepath.Join("assets/schemas", name, "adopters", fragmentName)
	fragment, err := specfs.SchemasFS.ReadFile(fragmentPath)
	if err != nil {
		return false, nil // fragment not found, skip silently
	}

	newContent := content + "\n" + string(fragment)
	return true, vfs.WriteFile(claudeMdPath, []byte(newContent), 0o644)
}
```

- [ ] **Step 3: Verify compilation and commit**

```bash
go build ./internal/schema/
git add internal/schema/ && git commit -m "feat: add internal/schema for schema bundle installation"
```

---

### Task 7: Create cmd/root.go

**Files:**
- `cmd/root.go`

- [ ] **Step 1: Create directory and get cobra dependency**

```bash
mkdir -p cmd
go get github.com/spf13/cobra
```

- [ ] **Step 2: Write cmd/root.go**

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/9Ashwin/spec-cli/internal/build"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spec-cli",
	Short: "Install OpenSpec, Superpowers, and schema bundles",
	Long: `spec-cli — OpenSpec + Superpowers workflow scaffolding tool.

spec-cli detects AI coding platforms and installs:
  - OpenSpec skills (spec lifecycle management)
  - Superpowers skills (brainstorming, TDD, code review)
  - Comet entry skill (thin workflow guide)
  - Schema bundles (workflow definitions for openspec/schemas/)

Commands:
  spec-cli init [path]     Initialize workflow scaffolding
  spec-cli status [path]   Show active changes
  spec-cli update [path]   Update packages and schemas
  spec-cli doctor [path]   Diagnose installation health

Examples:
  spec-cli init              # Interactive setup in current directory
  spec-cli init --yes        # Non-interactive, auto-detect platforms
  spec-cli status            # Show active workflow changes
  spec-cli doctor            # Check installation health`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command and returns the process exit code.
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return 1
	}
	return 0
}

func init() {
	rootCmd.Version = build.Version
	rootCmd.SetVersionTemplate("spec-cli version {{.Version}}\n")
}
```

- [ ] **Step 3: Verify compilation and commit**

```bash
go build ./cmd/
git add cmd/root.go go.mod go.sum && git commit -m "feat: add cmd/root.go with cobra root command"
```

---

### Task 8: Create cmd/init.go

**Files:**
- `cmd/init.go`

This is the largest task. The `init` command implements the full 9-step interactive flow matching Comet's init logic.

- [ ] **Step 1: Get huh + lipgloss dependencies**

```bash
go get github.com/charmbracelet/huh github.com/charmbracelet/lipgloss
```

- [ ] **Step 2: Write cmd/init.go**

```go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/9Ashwin/spec-cli/internal/openspec"
	"github.com/9Ashwin/spec-cli/internal/platform"
	"github.com/9Ashwin/spec-cli/internal/schema"
	"github.com/9Ashwin/spec-cli/internal/skill"
	"github.com/9Ashwin/spec-cli/internal/vfs"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

type initOptions struct {
	yes          bool
	skipExisting bool
	overwrite    bool
	jsonOutput   bool
	scope        string
}

var initOpts initOptions

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize workflow scaffolding",
	Long: `Initialize OpenSpec, Superpowers, Comet entry skill, and schema bundles.

Detects AI coding platforms and interactively installs all components.
Use --yes for non-interactive mode.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVar(&initOpts.yes, "yes", false, "Non-interactive mode")
	initCmd.Flags().BoolVar(&initOpts.skipExisting, "skip-existing", false, "Skip already installed components")
	initCmd.Flags().BoolVar(&initOpts.overwrite, "overwrite", false, "Overwrite all existing components")
	initCmd.Flags().BoolVar(&initOpts.jsonOutput, "json", false, "Output structured JSON")
	initCmd.Flags().StringVar(&initOpts.scope, "scope", "", "Install scope: project | global")
}

type initResult struct {
	ProjectPath       string         `json:"projectPath"`
	Scope             string         `json:"scope"`
	Language          string         `json:"language"`
	SelectedPlatforms []string       `json:"selectedPlatforms"`
	OpenSpec          string         `json:"openspec"`
	Superpowers       string         `json:"superpowers"`
	Comet             map[string]int `json:"comet"`
	SchemasInstalled  int            `json:"schemasInstalled"`
	WorkingDirs       bool           `json:"workingDirsCreated"`
}

func runInit(cmd *cobra.Command, args []string) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}
	projectPath, err := filepath.Abs(projectPath)
	if err != nil {
		return err
	}

	log := func(format string, a ...interface{}) {
		if !initOpts.jsonOutput {
			fmt.Fprintf(os.Stderr, format, a...)
		}
	}

	log("\n  spec-cli — OpenSpec + Superpowers Workflow Scaffolding\n\n")

	// Step 1: Detect platforms
	detected := platform.DetectPlatforms(projectPath)
	if len(detected) > 0 {
		names := make([]string, len(detected))
		for i, p := range detected {
			names[i] = p.Name
		}
		log("  Detected platforms: %s\n", strings.Join(names, ", "))
	}

	// Step 2: Select scope
	scope := initOpts.scope
	if scope == "" {
		if initOpts.yes {
			scope = "project"
		} else {
			scope = selectScope()
		}
	}
	log("  Scope: %s\n", scope)

	// Step 3: Select language
	language := "en"
	if !initOpts.yes {
		language = selectLanguage()
	}
	log("  Language: %s\n", languageName(language))

	// Step 4: Select platforms
	selected := selectPlatforms(detected)
	if len(selected) == 0 {
		log("\n  No platforms selected. Exiting.\n")
		if initOpts.jsonOutput {
			printJSON(initResult{ProjectPath: projectPath, Scope: scope, Language: language})
		}
		return nil
	}
	log("  Selected: %s\n", strings.Join(platformNames(selected), ", "))

	// Step 5: Determine base directory
	baseDir := projectPath
	if scope == "global" {
		home, err := vfs.UserHomeDir()
		if err != nil {
			return err
		}
		baseDir = home
	}

	// Step 6: Install OpenSpec
	var openSpecStatus string
	toolIDs := make([]string, len(selected))
	for i, p := range selected {
		toolIDs[i] = p.OpenSpecToolID
	}
	log("\n  Installing OpenSpec for: %s\n", strings.Join(toolIDs, ", "))
	if err := openspec.InitOpenSpec(projectPath, toolIDs, scope); err != nil {
		log("  OpenSpec: failed — %v\n", err)
		openSpecStatus = "failed"
	} else {
		log("  OpenSpec: installed\n")
		openSpecStatus = "installed"
	}

	// Step 7: Detect Superpowers
	superpowersStatus := "skipped"
	log("\n  Superpowers: checking...\n")
	if checkSuperpowers() {
		log("  Superpowers: detected (plugin-installed)\n")
		superpowersStatus = "detected"
	} else {
		log("  Superpowers: not detected. Install with: claude plugin install superpowers@claude-plugins-official\n")
	}

	// Step 8: Install Comet skill
	cometResults := make(map[string]int)
	for _, p := range selected {
		skillsDir := p.SkillsDir
		if scope == "global" && p.GlobalSkillsDir != "" {
			skillsDir = p.GlobalSkillsDir
		}
		copied, _, err := skill.CopySkills(baseDir, skillsDir, language, initOpts.overwrite)
		if err != nil {
			log("  Comet -> %s: error — %v\n", p.Name, err)
		} else {
			log("  Comet -> %s: %d copied\n", p.Name, copied)
			cometResults[p.ID] = copied
		}
	}

	// Step 9: Create working directories (project scope only)
	workingDirs := false
	if scope == "project" {
		specsDir := filepath.Join(projectPath, "docs", "superpowers", "specs")
		plansDir := filepath.Join(projectPath, "docs", "superpowers", "plans")
		if err := vfs.MkdirAll(specsDir, 0o755); err == nil {
			if err := vfs.MkdirAll(plansDir, 0o755); err == nil {
				workingDirs = true
			}
		}
		if workingDirs {
			log("\n  Working directories: docs/superpowers/specs/, docs/superpowers/plans/\n")
		}
	}

	// Step 10: Install schemas
	schemasInstalled := 0
	schemas, err := schema.ListSchemas()
	if err == nil && len(schemas) > 0 {
		schemaNames := make([]string, len(schemas))
		for i, s := range schemas {
			schemaNames[i] = s.Name
		}

		selectedSchemas := schemaNames
		if !initOpts.yes && len(schemas) > 1 {
			selectedSchemas = selectSchemas(schemas)
		}

		for _, s := range schemas {
			if !contains(selectedSchemas, s.Name) {
				continue
			}
			if err := schema.InstallSchema(s.Name, projectPath); err != nil {
				log("  Schema %s: failed — %v\n", s.Name, err)
			} else {
				schemasInstalled++
				log("  Schema: %s installed -> openspec/schemas/%s/\n", s.Name, s.Name)

				if added, _ := schema.AppendClaudeMdFragment(s.Name, projectPath, language); added {
					log("  CLAUDE.md: appended %s workflow fragment\n", s.Name)
				}
			}
		}
	}

	// Summary
	if !initOpts.jsonOutput {
		log("\n  Get started:\n")
		log("    openspec new --schema superpowers-bridge \"your idea\"\n\n")
	}

	if initOpts.jsonOutput {
		platformIDs := make([]string, len(selected))
		for i, p := range selected {
			platformIDs[i] = p.ID
		}
		printJSON(initResult{
			ProjectPath:       projectPath,
			Scope:             scope,
			Language:          language,
			SelectedPlatforms: platformIDs,
			OpenSpec:          openSpecStatus,
			Superpowers:       superpowersStatus,
			Comet:             cometResults,
			SchemasInstalled:  schemasInstalled,
			WorkingDirs:       workingDirs,
		})
	}

	return nil
}

// --- Interactive helpers ---

func selectScope() string {
	var scope string
	huh.NewSelect[string]().
		Title("Install scope:").
		Options(
			huh.NewOption("Project (current directory)", "project"),
			huh.NewOption("Global (home directory)", "global"),
		).
		Value(&scope).
		Run()
	return scope
}

func selectLanguage() string {
	var lang string
	huh.NewSelect[string]().
		Title("Language for Comet skills:").
		Options(
			huh.NewOption("English", "en"),
			huh.NewOption("中文", "zh"),
		).
		Value(&lang).
		Run()
	return lang
}

func languageName(lang string) string {
	if lang == "zh" {
		return "中文"
	}
	return "English"
}

func selectPlatforms(detected []platform.Platform) []platform.Platform {
	if initOpts.yes {
		if len(detected) > 0 {
			return detected
		}
		return platform.AllPlatforms
	}

	// In interactive mode, auto-select detected platforms.
	if len(detected) > 0 {
		return detected
	}

	// No platforms detected — prompt to select from all.
	var selected []platform.Platform
	// For now, select all as default when nothing detected.
	// Full multi-select with huh requires a custom implementation
	// since huh's MultiSelect works with string values, not structs.
	return platform.AllPlatforms
}

func selectSchemas(schemas []schema.Info) []string {
	// In --yes mode, select all.
	if initOpts.yes {
		names := make([]string, len(schemas))
		for i, s := range schemas {
			names[i] = s.Name
		}
		return names
	}

	// Interactive multi-select for schemas.
	options := make([]huh.Option[string], len(schemas))
	for i, s := range schemas {
		label := fmt.Sprintf("%s (v%s)", s.Name, s.Version)
		options[i] = huh.NewOption(label, s.Name)
	}

	// Default: select all schemas when interactive (simplest UX).
	// Full multi-select with huh can be added later.
	names := make([]string, len(schemas))
	for i, s := range schemas {
		names[i] = s.Name
	}
	return names
}

// --- Superpowers detection ---

var superpowersSkillNames = []string{
	"brainstorming",
	"using-superpowers",
	"writing-plans",
	"test-driven-development",
	"subagent-driven-development",
}

// checkSuperpowers checks if Superpowers is installed via Claude Code plugins.
// Mirrors Comet's hasPluginSuperpowers() in detect.ts.
func checkSuperpowers() bool {
	home, err := vfs.UserHomeDir()
	if err != nil {
		return false
	}

	claudeDir := os.Getenv("CLAUDE_CONFIG_DIR")
	if claudeDir == "" {
		claudeDir = filepath.Join(home, ".claude")
	}

	pluginsCacheDir := filepath.Join(claudeDir, "plugins", "cache")

	marketplaceEntries, err := vfs.ReadDir(pluginsCacheDir)
	if err != nil {
		return false
	}

	for _, marketplace := range marketplaceEntries {
		if !marketplace.IsDir() {
			continue
		}

		superpowersDir := filepath.Join(pluginsCacheDir, marketplace.Name(), "superpowers")
		if _, err := vfs.Stat(superpowersDir); err != nil {
			continue
		}

		versionEntries, err := vfs.ReadDir(superpowersDir)
		if err != nil {
			continue
		}

		for _, version := range versionEntries {
			if !version.IsDir() {
				continue
			}

			skillsDir := filepath.Join(superpowersDir, version.Name(), "skills")
			skillEntries, err := vfs.ReadDir(skillsDir)
			if err != nil {
				continue
			}

			for _, entry := range skillEntries {
				for _, name := range superpowersSkillNames {
					if entry.Name() == name {
						return true
					}
				}
			}
		}
	}

	return false
}

// --- Utilities ---

func platformNames(platforms []platform.Platform) []string {
	names := make([]string, len(platforms))
	for i, p := range platforms {
		names[i] = p.Name
	}
	return names
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func printJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return
	}
	fmt.Println(string(data))
}
```

- [ ] **Step 3: Register initCmd in root.go**

Edit `cmd/root.go` init() to add the subcommand:

```go
func init() {
	rootCmd.Version = build.Version
	rootCmd.SetVersionTemplate("spec-cli version {{.Version}}\n")
	rootCmd.AddCommand(initCmd)
}
```

- [ ] **Step 4: Verify compilation and commit**

```bash
go build ./cmd/ ./...
git add cmd/init.go cmd/root.go go.mod go.sum
git commit -m "feat: add cmd/init.go — interactive init with huh"
```

---

### Task 9: Create cmd/status.go, cmd/update.go, cmd/doctor.go

**Files:**
- `cmd/status.go`
- `cmd/update.go`
- `cmd/doctor.go`

- [ ] **Step 1: Write cmd/status.go**

```go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/9Ashwin/spec-cli/internal/openspec"
	"github.com/spf13/cobra"
)

var statusJSON bool

var statusCmd = &cobra.Command{
	Use:   "status [path]",
	Short: "Show active changes",
	Long:  "Display active workflow changes from openspec list --json.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runStatus,
}

func init() {
	statusCmd.Flags().BoolVar(&statusJSON, "json", false, "Output structured JSON")
}

func runStatus(cmd *cobra.Command, args []string) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}
	projectPath, err := filepath.Abs(projectPath)
	if err != nil {
		return err
	}

	changes, err := openspec.ListChanges(projectPath)
	if err != nil {
		return err
	}

	if statusJSON {
		if changes == nil {
			changes = []openspec.ChangeInfo{}
		}
		data, _ := json.MarshalIndent(changes, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	if len(changes) == 0 {
		fmt.Println("No active changes.")
		return nil
	}

	fmt.Fprintf(os.Stderr, "\n  Active Changes (%d):\n\n", len(changes))
	for _, c := range changes {
		schemaLabel := ""
		if c.Schema != "" {
			schemaLabel = fmt.Sprintf(" [schema: %s]", c.Schema)
		}
		fmt.Fprintf(os.Stderr, "  • %s%s\n", c.Name, schemaLabel)
	}
	fmt.Fprintln(os.Stderr)

	return nil
}
```

- [ ] **Step 2: Write cmd/update.go**

```go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/9Ashwin/spec-cli/internal/platform"
	"github.com/9Ashwin/spec-cli/internal/schema"
	"github.com/9Ashwin/spec-cli/internal/skill"
	"github.com/9Ashwin/spec-cli/internal/vfs"
	"github.com/spf13/cobra"
)

var (
	updateJSON     bool
	updateLanguage string
	updateScope    string
)

var updateCmd = &cobra.Command{
	Use:   "update [path]",
	Short: "Update packages and schemas",
	Long:  "Re-copy skill files and update schema bundles when versions differ.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().BoolVar(&updateJSON, "json", false, "Output structured JSON")
	updateCmd.Flags().StringVar(&updateLanguage, "language", "en", "Language: en | zh")
	updateCmd.Flags().StringVar(&updateScope, "scope", "project", "Update scope: project | global")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}
	projectPath, err := filepath.Abs(projectPath)
	if err != nil {
		return err
	}

	log := func(format string, a ...interface{}) {
		if !updateJSON {
			fmt.Fprintf(os.Stderr, format, a...)
		}
	}

	baseDir := projectPath
	if updateScope == "global" {
		home, err := vfs.UserHomeDir()
		if err != nil {
			return err
		}
		baseDir = home
	}

	type updateResult struct {
		SkillsUpdated  int            `json:"skillsUpdated"`
		SchemasUpdated int            `json:"schemasUpdated"`
		Comet          map[string]int `json:"comet"`
	}

	result := updateResult{Comet: make(map[string]int)}

	// Update skills (overwrite mode)
	detected := platform.DetectPlatforms(projectPath)
	log("\n  Updating skills...\n")
	for _, p := range detected {
		skillsDir := p.SkillsDir
		if updateScope == "global" && p.GlobalSkillsDir != "" {
			skillsDir = p.GlobalSkillsDir
		}
		copied, _, err := skill.CopySkills(baseDir, skillsDir, updateLanguage, true)
		if err != nil {
			log("  %s: error — %v\n", p.Name, err)
		} else {
			result.Comet[p.ID] = copied
			result.SkillsUpdated += copied
			log("  %s: %d updated\n", p.Name, copied)
		}
	}

	// Update schemas
	schemas, err := schema.ListSchemas()
	if err == nil {
		log("\n  Updating schemas...\n")
		for _, s := range schemas {
			installed := schema.GetInstalledVersion(s.Name, projectPath)
			if installed != "" && installed == s.Version {
				log("  %s: up to date (v%s)\n", s.Name, s.Version)
				continue
			}
			if err := schema.InstallSchema(s.Name, projectPath); err != nil {
				log("  %s: failed — %v\n", s.Name, err)
			} else {
				result.SchemasUpdated++
				log("  %s: updated to v%s\n", s.Name, s.Version)
			}
		}
	}

	if updateJSON {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	} else {
		log("\n  Update complete. %d skills, %d schemas updated.\n\n",
			result.SkillsUpdated, result.SchemasUpdated)
	}

	return nil
}
```

- [ ] **Step 3: Write cmd/doctor.go**

```go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/9Ashwin/spec-cli/internal/openspec"
	"github.com/9Ashwin/spec-cli/internal/platform"
	"github.com/9Ashwin/spec-cli/internal/schema"
	"github.com/9Ashwin/spec-cli/internal/vfs"
	"github.com/spf13/cobra"
)

var (
	doctorJSON  bool
	doctorScope string
)

var doctorCmd = &cobra.Command{
	Use:   "doctor [path]",
	Short: "Diagnose installation health",
	Long:  "Check OpenSpec CLI, working directories, schema bundles, and skill files.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runDoctor,
}

func init() {
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "Output structured JSON")
	doctorCmd.Flags().StringVar(&doctorScope, "scope", "auto", "Check scope: auto | project | global")
}

type doctorCheck struct {
	Name   string `json:"name"`
	Status string `json:"status"` // "ok", "warning", "error"
	Detail string `json:"detail,omitempty"`
}

func runDoctor(cmd *cobra.Command, args []string) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}
	projectPath, err := filepath.Abs(projectPath)
	if err != nil {
		return err
	}

	var checks []doctorCheck

	// Check 1: openspec CLI
	version, err := openspec.Version()
	if err != nil {
		checks = append(checks, doctorCheck{
			Name: "OpenSpec CLI", Status: "error",
			Detail: fmt.Sprintf("not found on PATH: %v", err),
		})
	} else {
		checks = append(checks, doctorCheck{
			Name: "OpenSpec CLI", Status: "ok",
			Detail: fmt.Sprintf("version %s", version),
		})
	}

	// Check 2: Working directories
	specsDir := filepath.Join(projectPath, "docs", "superpowers", "specs")
	plansDir := filepath.Join(projectPath, "docs", "superpowers", "plans")
	_, specsErr := vfs.Stat(specsDir)
	_, plansErr := vfs.Stat(plansDir)
	if specsErr != nil || plansErr != nil {
		checks = append(checks, doctorCheck{
			Name: "Working Directories", Status: "warning",
			Detail: "docs/superpowers/specs/ or plans/ missing — run spec-cli init",
		})
	} else {
		checks = append(checks, doctorCheck{
			Name: "Working Directories", Status: "ok",
			Detail: "docs/superpowers/specs/, docs/superpowers/plans/",
		})
	}

	// Check 3: Schema bundles
	schemas, err := schema.ListSchemas()
	if err != nil {
		checks = append(checks, doctorCheck{
			Name: "Schemas", Status: "error", Detail: err.Error(),
		})
	} else {
		for _, s := range schemas {
			if schema.IsInstalled(s.Name, projectPath) {
				installed := schema.GetInstalledVersion(s.Name, projectPath)
				if installed != s.Version {
					checks = append(checks, doctorCheck{
						Name: fmt.Sprintf("Schema: %s", s.Name), Status: "warning",
						Detail: fmt.Sprintf("installed v%s, available v%s", installed, s.Version),
					})
				} else {
					checks = append(checks, doctorCheck{
						Name: fmt.Sprintf("Schema: %s", s.Name), Status: "ok",
						Detail: fmt.Sprintf("v%s", installed),
					})
				}
			} else {
				checks = append(checks, doctorCheck{
					Name: fmt.Sprintf("Schema: %s", s.Name), Status: "warning",
					Detail: "not installed — run spec-cli init",
				})
			}
		}
	}

	// Check 4: Skill files for detected platforms
	detected := platform.DetectPlatforms(projectPath)
	if len(detected) == 0 {
		checks = append(checks, doctorCheck{
			Name: "Platform Skills", Status: "warning",
			Detail: "no platforms detected",
		})
	}
	for _, p := range detected {
		skillPath := filepath.Join(projectPath, p.SkillsDir, "skills", "comet", "SKILL.md")
		if _, err := vfs.Stat(skillPath); err != nil {
			checks = append(checks, doctorCheck{
				Name: fmt.Sprintf("Skills: %s", p.Name), Status: "warning",
				Detail: fmt.Sprintf("%s/skills/comet/SKILL.md not found", p.SkillsDir),
			})
		} else {
			checks = append(checks, doctorCheck{
				Name: fmt.Sprintf("Skills: %s", p.Name), Status: "ok",
				Detail: fmt.Sprintf("%s/skills/comet/SKILL.md present", p.SkillsDir),
			})
		}
	}

	// Output
	if doctorJSON {
		data, _ := json.MarshalIndent(checks, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	fmt.Fprintf(os.Stderr, "\n  spec-cli Doctor\n\n")
	okCount, warnCount, errCount := 0, 0, 0
	for _, c := range checks {
		icon := "✓"
		switch c.Status {
		case "warning":
			icon = "!"
			warnCount++
		case "error":
			icon = "✗"
			errCount++
		default:
			okCount++
		}
		detail := ""
		if c.Detail != "" {
			detail = fmt.Sprintf(" — %s", c.Detail)
		}
		fmt.Fprintf(os.Stderr, "  %s %s%s\n", icon, c.Name, detail)
	}
	fmt.Fprintf(os.Stderr, "\n  %d ok, %d warnings, %d errors\n\n", okCount, warnCount, errCount)

	if errCount > 0 {
		return fmt.Errorf("doctor found %d error(s)", errCount)
	}
	return nil
}
```

- [ ] **Step 4: Register all commands in cmd/root.go init()**

```go
func init() {
	rootCmd.Version = build.Version
	rootCmd.SetVersionTemplate("spec-cli version {{.Version}}\n")
	rootCmd.AddCommand(initCmd, statusCmd, updateCmd, doctorCmd)
}
```

- [ ] **Step 5: Verify compilation and commit**

```bash
go build ./cmd/ ./...
git add cmd/status.go cmd/update.go cmd/doctor.go cmd/root.go
git commit -m "feat: add status, update, doctor commands"
```

---

### Task 10: Create Makefile

**Files:**
- `Makefile`

- [ ] **Step 1: Write Makefile**

```makefile
BINARY  := spec-cli
MODULE  := github.com/9Ashwin/spec-cli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
DATE    := $(shell date +%Y-%m-%d)
LDFLAGS := -s -w -X $(MODULE)/internal/build.Version=$(VERSION) -X $(MODULE)/internal/build.Date=$(DATE)

.PHONY: build test fmt vet clean

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) .

test:
	go test -race -count=1 ./...

vet:
	go vet ./...

fmt:
	gofmt -s -w .

clean:
	rm -f $(BINARY)
```

- [ ] **Step 2: Verify build**

```bash
make build
./spec-cli --version
```

Expected: `spec-cli version <commit-hash>`

- [ ] **Step 3: Commit**

```bash
git add Makefile && git commit -m "feat: add Makefile with build/test/vet/fmt targets"
```

---

### Task 11: Final verification

- [ ] **Step 1: Full vet + build**

```bash
make vet && make build
```

- [ ] **Step 2: gofmt check**

```bash
test -z "$(gofmt -l . | grep -v '^\.claude/')"
```

- [ ] **Step 3: CLI smoke test**

```bash
./spec-cli --help
./spec-cli --version
./spec-cli init --help
./spec-cli status --help
./spec-cli update --help
./spec-cli doctor --help
```

- [ ] **Step 4: Verify embed compiles (binary includes assets)**

```bash
go build -o /dev/null .
```

- [ ] **Step 5: Final commit + push**

```bash
git add -A && git status
git commit -m "chore: final verification — vet, fmt, build all pass"
git push origin main
```
