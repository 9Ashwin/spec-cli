# AGENTS.md

## Project Overview

spec-cli is a Go CLI tool that scaffolds OpenSpec + Superpowers development workflows into AI coding platform projects. A clean-room rewrite of [Comet](https://github.com/rpamis/comet), it detects installed AI coding platforms and installs OpenSpec skills, Superpowers skills, a thin `/comet` entry skill, and schema bundles into `openspec/schemas/`. Workflow execution is delegated to OpenSpec's native `--schema` mechanism — spec-cli does not manage phase state, guard scripts, or handoff/archive automation.

## Architecture Overview

Single binary (`spec-cli`) built with cobra. Assets (skills, schema bundles) are embedded via `//go:embed` at compile time. All filesystem operations go through a `vfs.FS` interface for testability. Interactive prompts use charmbracelet/huh.

| Layer | Tech | Description |
|-------|------|-------------|
| **cmd/** | cobra | CLI commands — thin wrappers that parse flags and call internal packages |
| **internal/** | Go stdlib + exec | Business logic: platform detection, skill copy, schema install, openspec integration |
| **embed/** | go:embed | Embedded assets: skill files (en/zh) and schema bundles |
| **assets/** | Markdown + YAML | Embedded content: comet SKILL.md, superpowers-bridge schema |

## Tech Stack

- **Language**: Go 1.23
- **CLI Framework**: cobra
- **Interactive Prompts**: charmbracelet/huh + charmbracelet/lipgloss
- **Asset Embedding**: go:embed
- **Build**: Makefile with ldflags version injection
- **Testing**: go test -race

## Reference Projects

| Project | Path | Purpose |
|---------|------|---------|
| lark-cli | `~/workspaces/github/cli` | Go patterns: cobra root, vfs interface, build ldflags, embed |
| comet | `~/workspaces/github/comet` | Original TypeScript codebase: 29 platforms, init flow, openspec/superpowers install logic |
| openspec-schemas | `~/workspaces/github/openspec-schemas` | Schema.yaml format: artifact-driven schema + templates + adopters convention |

## Project Structure

```
spec-cli/
├── main.go                          # os.Exit(cmd.Execute())
├── cmd/                             # cobra commands (thin — parse args, call internal)
│   ├── root.go                      #   Root command: version, help, Execute()
│   ├── init.go                      #   spec-cli init: interactive scaffolding
│   ├── status.go                    #   spec-cli status: openspec list --json wrapper
│   ├── update.go                    #   spec-cli update: refresh skills + schemas
│   └── doctor.go                    #   spec-cli doctor: diagnose install health
├── internal/                        # Go business logic
│   ├── build/                       #   Version/Date ldflags injection (from lark-cli)
│   │   └── build.go                 #     var Version = "DEV"; overridden via -ldflags
│   ├── vfs/                         #   Filesystem abstraction (from lark-cli)
│   │   └── vfs.go                   #     FS interface + OsFs + DefaultFS + helpers
│   ├── platform/                    #   29 AI coding platform definitions
│   │   ├── platform.go              #     Platform struct + AllPlatforms slice
│   │   └── detect.go                #     DetectPlatforms() via DetectionPaths/SkillsDirs
│   ├── skill/                       #   Skill file management
│   │   └── skill.go                 #     CopySkills() from embed to target dirs
│   ├── openspec/                    #   OpenSpec CLI integration
│   │   └── openspec.go              #     InitOpenSpec(), ListChanges(), Version()
│   └── schema/                      #   Schema bundle management
│       └── schema.go                #     ListSchemas(), InstallSchema(), AppendClaudeMdFragment()
├── embed/                           # Go embed targets (root-level for //go:embed paths)
│   ├── skills.go                    #   //go:embed assets/skills/* + assets/skills-zh/*
│   └── schemas.go                   #   //go:embed assets/schemas/*
├── assets/                          # Embedded content
│   ├── skills/comet/SKILL.md        #   English comet entry skill
│   ├── skills-zh/comet/SKILL.md     #   Chinese comet entry skill
│   └── schemas/superpowers-bridge/  #   Schema bundle
│       ├── schema.yaml              #     Artifact definitions + apply flow
│       ├── VERSION                  #     Schema version
│       ├── templates/               #     Artifact templates (proposal, design, tasks, etc.)
│       └── adopters/                #     CLAUDE.md fragments (en + zh)
├── docs/                            # Project documentation
│   ├── specs/                       #   Design specs and implementation plans
│   ├── reference/                   #   Reference documentation
│   └── spec-cli-handoff.md          #   Session handoff document
├── Makefile                         # build, test, vet, fmt targets
├── go.mod                           # github.com/9Ashwin/spec-cli
└── go.sum
```

## Development Guide

### Prerequisites

1. Go 1.23+
2. (optional) `openspec` CLI for integration testing

### Common Commands

| Command | Description |
|---------|-------------|
| `make build` | Build binary with ldflags version injection |
| `make test` | Run all tests with race detector |
| `make vet` | Run go vet |
| `make fmt` | Format all Go source files |
| `go build ./...` | Check all packages compile |
| `./spec-cli --help` | Verify binary works |
| `./spec-cli init` | Run interactive init (test in a temp dir) |

### Build Flags

```bash
go build -trimpath -ldflags "-s -w -X github.com/9Ashwin/spec-cli/internal/build.Version=$(git describe --tags --always --dirty) -X github.com/9Ashwin/spec-cli/internal/build.Date=$(date +%Y-%m-%d)" -o spec-cli .
```

## Key Development Rules

### VFS Interface

- All filesystem operations in `internal/` packages must use `vfs.DefaultFS` (or accept an `vfs.FS` parameter).
- Never call `os.ReadFile`, `os.WriteFile`, `os.MkdirAll`, `os.Stat`, `os.Remove`, `os.RemoveAll` directly in internal packages.
- In tests, replace `vfs.DefaultFS` with an in-memory implementation.
- `cmd/` layer is exempt — it can use `os.` directly for one-off stat checks.

### Embed Pattern

- `//go:embed` directives live in the `embed/` package at module root because Go does not allow `..` in embed paths.
- `embed/skills.go` and `embed/schemas.go` each scope their embed to a specific `assets/` subdirectory.
- Internal packages import `specfs "github.com/9Ashwin/spec-cli/embed"` and read from the exported `embed.FS` variables.
- When adding new skill files or schema bundles, update the corresponding embed directive's glob pattern.

### Platform Definitions

- The 29-platform list in `internal/platform/platform.go` mirrors Comet's `src/core/platforms.ts`.
- Each `Platform` struct: `ID`, `Name`, `DetectionPaths` (paths that indicate the platform is active), `SkillsDirs` (where skills get installed), `OpenSpecToolID` (passed to `openspec init --tools`).
- `DetectPlatforms()` checks `DetectionPaths` first, falls back to checking `SkillsDirs` existence.
- Adding a new platform requires: add to `AllPlatforms`, verify `OpenSpecToolID` matches OpenSpec's registry.

### Schema Bundles

- Schema bundles live under `assets/schemas/<name>/` and are embedded via `embed.SchemasFS`.
- Each bundle: `schema.yaml` (artifact definitions + apply flow), `VERSION`, `templates/` (artifact templates), `adopters/` (CLAUDE.md fragments), `README.md`.
- `schema.InstallSchema()` copies the entire bundle directory from embed to `openspec/schemas/<name>/`.
- Schema version comparison in `spec-cli update` reads the embedded `VERSION` against the installed one.

### CLI Commands

- Root command in `cmd/root.go`: `func Execute() int`, called by `main.go` via `os.Exit(cmd.Execute())`.
- Each subcommand is a `var ...Cmd = &cobra.Command{...}` registered in `init()`.
- Commands use `RunE` for error returns. Flags use `cmd.Flags().BoolVar()` / `cmd.Flags().StringVar()`.
- `--json` flag on all commands outputs structured JSON to stdout (matching Comet's TypeScript JSON format).
- Interactive prompts use huh (`huh.NewSelect`, `huh.NewMultiSelect`); non-interactive mode via `--yes` flag.
