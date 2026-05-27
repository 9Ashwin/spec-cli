# Spec-CLI Go Implementation Design

## Context

spec-cli is a Go CLI tool that installs OpenSpec, Superpowers, and schema bundles into AI coding platform projects. It is a clean-room rewrite of the [Comet](https://github.com/rpamis/comet) concept, using the [lark-cli](https://github.com/larksuite/cli) project as a reference for Go patterns (cobra, charmbracelet/huh, go:embed, vfs).

**What it does:** `spec-cli init` detects installed AI coding platforms, interactively installs OpenSpec skills, Superpowers skills, a thin `/comet` entry skill, and schema bundles into `openspec/schemas/`. The workflow execution itself is delegated to OpenSpec's native `--schema` mechanism.

**What it does NOT do:** Phase state management, guard scripts, handoff/archive automation — those are handled by OpenSpec schema instructions.

## Architecture

```
spec-cli/
├── main.go                          # os.Exit(cmd.Execute())
├── cmd/
│   ├── root.go                      # cobra root + Execute()
│   ├── init.go                      # spec-cli init
│   ├── status.go                    # spec-cli status
│   ├── update.go                    # spec-cli update
│   └── doctor.go                    # spec-cli doctor
├── internal/
│   ├── build/
│   │   └── build.go                 # Version / Date (ldflags injected)
│   ├── vfs/
│   │   └── vfs.go                   # FS interface (testable)
│   ├── platform/
│   │   ├── platform.go              # 29 platform definitions
│   │   └── detect.go                # platform auto-detection
│   ├── skill/
│   │   └── skill.go                 # copy skills from embed to target dirs
│   ├── openspec/
│   │   └── openspec.go              # openspec init --tools invocation
│   └── schema/
│       └── schema.go                # schema bundle install from embed
├── assets/
│   ├── skills/comet/SKILL.md        # embedded: English entry skill
│   ├── skills-zh/comet/SKILL.md     # embedded: Chinese entry skill
│   └── schemas/superpowers-bridge/  # embedded: schema bundle
│       ├── schema.yaml
│       ├── VERSION
│       ├── templates/
│       └── adopters/
├── Makefile
├── go.mod                           # github.com/9Ashwin/spec-cli
└── go.sum
```

## Dependencies

```
cmd/init.go
  ├── internal/platform   (DetectPlatforms)
  ├── internal/skill      (CopySkills)
  ├── internal/schema     (ListSchemas, InstallSchema)
  └── internal/openspec   (InitOpenSpec)

cmd/status.go
  └── internal/openspec   (exec openspec list --json)

cmd/update.go
  ├── internal/skill      (CopySkills with overwrite)
  └── internal/schema     (version comparison + re-install)

cmd/doctor.go
  ├── internal/openspec   (CLI version check)
  ├── internal/platform   (skills presence check)
  └── internal/schema     (schema presence check)
```

## Embed Strategy

Two embed entry points, each scoped to its own assets subdirectory:

```go
// internal/skill/skill.go
//go:embed assets/skills/*
var skillsFS embed.FS

func CopySkills(targetDir string, platform Platform, language string, overwrite bool) (int, int, error)
```

```go
// internal/schema/schema.go
//go:embed assets/schemas/*
var schemasFS embed.FS

func ListSchemas() ([]SchemaInfo, error)
func InstallSchema(name, targetDir string) error
func GetSchemaVersion(name string) (string, error)
```

skill and schema packages use `internal/vfs.FS` for all filesystem writes, enabling in-memory testing.

## VFS Interface

```go
// internal/vfs/vfs.go
type FS interface {
    MkdirAll(path string, perm os.FileMode) error
    WriteFile(name string, data []byte, perm os.FileMode) error
    ReadFile(name string) ([]byte, error)
    Stat(name string) (os.FileInfo, error)
    RemoveAll(path string) error
}

type OsFS struct{}

var DefaultFS FS = OsFS{}  // replace in tests
```

## Platform Detection

```go
// internal/platform/platform.go
type Platform struct {
    ID             string
    Name           string
    DetectionPaths []string   // e.g. [".claude", ".claude.json"]
    SkillsDirs     []string   // e.g. [".claude/"]
    OpenSpecToolID string     // e.g. "claude"
}

var AllPlatforms = []Platform{...}  // 29 platforms

// internal/platform/detect.go
func DetectPlatforms(projectPath string) []Platform
func SkillsDir(p Platform, scope Scope) string
```

Detection checks `DetectionPaths` against the filesystem. For platforms without `DetectionPaths`, falls back to checking `SkillsDirs` directory existence.

Keep the 29-platform list from Comet's `src/core/platforms.ts`. Each platform added to the source list requires:
1. Add to `AllPlatforms` in `platform.go`
2. Map `OpenSpecToolID` to OpenSpec's tool ID

## Interactive Init Flow (charmbracelet/huh)

```
spec-cli init [path]
  │
  ├─ Step 1: Platform detection (automatic)
  │    DetectPlatforms() marks results "(detected)" in selection list
  │
  ├─ Step 2: Scope
  │    huh.NewSelect[Scope]() → project / global
  │
  ├─ Step 3: Language
  │    huh.NewSelect[string]() → English / 中文
  │
  ├─ Step 4: Platform selection (multi-select)
  │    huh.NewMultiSelect[Platform]() — detected pre-checked
  │    --yes: auto-select detected, or all if none detected
  │
  ├─ Step 5: Existing component handling
  │    Per platform: check openspec/superpowers/comet presence
  │    Multiple existing: bulk choice (overwrite-all/skip-all/per-component)
  │    Single existing: overwrite/skip
  │    --yes: skip existing
  │    --overwrite: overwrite all
  │    --skip-existing: skip all
  │
  ├─ Step 6: Install OpenSpec
  │    exec: openspec init <path> --tools <tool-ids>
  │    Auto-install CLI if missing: npm install -g @fission-ai/openspec@latest
  │
  ├─ Step 7: Install Superpowers
  │    Check ~/.claude/plugins/cache/*/superpowers/*/skills/
  │    Warn if not found, point to plugin install docs
  │
  ├─ Step 8: Install Comet skill
  │    skill.CopySkills() → writes SKILL.md from embed
  │    Project scope: create docs/superpowers/specs/ and plans/
  │
  └─ Step 9: Install schemas
       schema.ListSchemas() → available bundles (multi-select)
       schema.InstallSchema() → openspec/schemas/<name>/
       Append CLAUDE.md fragment if CLAUDE.md exists
```

## CLI Flags

```
spec-cli init [path]
  --yes               Non-interactive, auto-select detected platforms
  --scope <scope>     project | global
  --skip-existing     Skip already installed components
  --overwrite         Overwrite all existing components
  --json              Output structured JSON result

spec-cli status [path]
  --json              Output JSON

spec-cli update [path]
  --json              Output JSON
  --language <lang>   Override detected language
  --scope <scope>     Update only global or project

spec-cli doctor [path]
  --json              Output JSON
  --scope <scope>     auto | project | global
```

## Status Command

Reads `openspec list --json` output, parses active changes, displays them with schema name and artifact progress. No `.comet.yaml` parsing — state is tracked by OpenSpec natively.

## Update Command

1. Check for newer npm version (if installed via npm)
2. Re-copy skill files from embed to installed targets (overwrite)
3. Compare embedded schema versions against installed `VERSION` files
4. Update schemas when versions differ

## Doctor Command

Checks:
- `openspec` CLI installed and on PATH
- Working directories (docs/superpowers/specs/, plans/)
- Schema bundles present in `openspec/schemas/`
- Skill files present for each detected platform

## Entry Skill

`/comet` skill (written by `skill.CopySkills`) is a thin guide:

```markdown
---
name: comet
description: "Comet — OpenSpec + Superpowers development workflow."
---

# Comet — OpenSpec + Superpowers Development Workflow

## Workflow

openspec new --schema superpowers-bridge <name>   # start a new change
openspec status --change "<name>"                  # check progress
openspec archive --change "<name>" -y              # archive when done

## Active Changes

openspec list --json
```

## What Gets Deleted

All Node.js/TypeScript code from the original comet copy:
- `src/`, `test/`, `scripts/`, `bin/`
- `package.json`, `package-lock.json`, `pnpm-lock.yaml`
- `tsconfig.json`, `vitest.config.ts`, `eslint.config.js`
- `build.js`, `node_modules/`
- `assets/manifest.json` (replaced by embed)
