# AGENTS.md

## Project Overview

spec-cli is a Go CLI tool that scaffolds OpenSpec + Superpowers development workflows into AI coding platform projects. A clean-room rewrite of [Comet](https://github.com/rpamis/comet), it detects installed AI coding platforms and installs OpenSpec skills, Superpowers skills, a thin `/comet` entry skill, and schema bundles into `openspec/schemas/`. Workflow execution is delegated to OpenSpec's native `--schema` mechanism — spec-cli does not manage phase state, guard scripts, or handoff/archive automation.

## Architecture Overview

Single binary (`spec-cli`) built with cobra. Assets (skills, schema bundles) are embedded via `//go:embed` at compile time. All filesystem operations go through a `vfs.FS` interface. Interactive prompts use charmbracelet/huh.

| Layer | Tech | Description |
|-------|------|-------------|
| **cmd/** | cobra | CLI commands — thin wrappers that parse flags and call internal packages |
| **internal/** | Go stdlib + exec | Business logic: platform detection, skill copy, schema install, openspec integration |
| **embed/** | go:embed + Markdown + YAML | Embedded assets: skill files (en/zh), schema bundles, templates, adopters |

## Tech Stack

- **Language**: Go 1.23
- **CLI Framework**: [cobra](https://github.com/spf13/cobra)
- **Interactive Prompts**: [charmbracelet/huh](https://github.com/charmbracelet/huh) + [lipgloss](https://github.com/charmbracelet/lipgloss)
- **Asset Embedding**: go:embed
- **Build**: Makefile with ldflags version injection
- **Testing**: go test -race

## Reference Projects

| Project | Path | Purpose |
|---------|------|---------|
| lark-cli | `/Users/solariswu/workspaces/github/cli` | Go patterns: cobra root, vfs interface, build ldflags, embed |
| comet | `/Users/solariswu/workspaces/github/comet` | Original TypeScript codebase: 29 platforms, init flow, openspec/superpowers install logic |
| openspec-schemas | `/Users/solariswu/workspaces/github/openspec-schemas` | Schema.yaml format: artifact-driven schema + templates + adopters convention |

## Project Structure

```
spec-cli/
├── main.go                             # os.Exit(cmd.Execute())
├── cmd/                                # cobra commands (thin — parse args, call internal)
│   ├── root.go                         #   Root command: version, help, Execute()
│   ├── init.go                         #   spec-cli init: interactive scaffolding (9-step flow)
│   ├── status.go                       #   spec-cli status: openspec list --json wrapper
│   ├── update.go                       #   spec-cli update: refresh skills + schemas
│   └── doctor.go                       #   spec-cli doctor: diagnose install health
├── internal/                           # Go business logic
│   ├── build/                          #   Version/Date ldflags injection (from lark-cli)
│   │   └── build.go                    #     var Version = "DEV"; overridden via -ldflags
│   ├── vfs/                            #   Filesystem abstraction (from lark-cli)
│   │   ├── fs.go                       #     FS interface
│   │   ├── osfs.go                     #     OsFs struct (delegates to os package)
│   │   └── default.go                  #     DefaultFS + package-level convenience functions
│   ├── platform/                       #   29 AI coding platform definitions
│   │   ├── platform.go                 #     Platform struct + AllPlatforms slice
│   │   └── detect.go                   #     DetectPlatforms() via DetectionPaths/SkillsDir
│   ├── skill/                          #   Skill file management
│   │   └── skill.go                    #     CopySkills() from embed to target dirs
│   ├── openspec/                       #   OpenSpec CLI integration
│   │   └── openspec.go                 #     InitOpenSpec(), ListChanges(), Version()
│   └── schema/                         #   Schema bundle management
│       └── schema.go                   #     ListSchemas(), InstallSchema(), AppendClaudeMdFragment()
├── embed/                              # Go embed targets + embedded assets
│   ├── skills.go                       #   //go:embed all:assets/skills + all:assets/skills-zh
│   ├── schemas.go                      #   //go:embed all:assets/schemas
│   └── assets/                         #   Embedded content (must be in same dir as embed .go files)
│       ├── skills/comet/SKILL.md       #     English comet entry skill
│       ├── skills-zh/comet/SKILL.md    #     简体中文 comet entry skill
│       └── schemas/superpowers-bridge/ #     Schema bundle
│           ├── schema.yaml             #       Artifact definitions + apply flow
│           ├── VERSION                 #       Schema version
│           ├── templates/              #       Artifact templates (proposal, design, tasks, etc.)
│           └── templates/adopters/     #       CLAUDE.md fragments (en + zh)
├── docs/                               # Project documentation
│   ├── specs/                          #   Design specs and implementation plans
│   │   ├── 2026-05-27-spec-cli-go-design.md
│   │   └── 2026-05-27-spec-cli-go-plan.md
│   └── spec-cli-handoff.md             #   Session handoff document
├── Makefile                            # build, test, vet, fmt targets
├── go.mod                              # github.com/9Ashwin/spec-cli
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
| `./spec-cli --version` | Show build version |
| `./spec-cli init` | Run interactive init (test in a temp dir) |
| `./spec-cli init --yes` | Non-interactive init, auto-detect platforms |
| `./spec-cli status` | Show active workflow changes |
| `./spec-cli doctor` | Diagnose installation health |

### Build Flags

```bash
make build
# Equivalent to:
go build -trimpath -ldflags "-s -w -X github.com/9Ashwin/spec-cli/internal/build.Version=$(git describe --tags --always --dirty) -X github.com/9Ashwin/spec-cli/internal/build.Date=$(date +%Y-%m-%d)" -o spec-cli .
```

## Key Development Rules

### VFS Interface

- All filesystem operations in `internal/` packages must use `vfs.DefaultFS` or accept an `vfs.FS` parameter.
- Never call `os.ReadFile`, `os.WriteFile`, `os.MkdirAll`, `os.Stat`, `os.Remove`, `os.RemoveAll` directly in internal packages.
- In tests, replace `vfs.DefaultFS` with an in-memory implementation.
- Package-level convenience functions (`vfs.Stat()`, `vfs.ReadFile()`, etc.) delegate to `DefaultFS` — tests can swap the global.

### Embed Pattern

- `//go:embed` directives live in the `embed/` package at module root. Go does not allow `..` in embed paths, so patterns are relative to the source file directory.
- `embed/skills.go` uses `//go:embed all:assets/skills` for recursive embedding. The `all:` prefix (Go 1.22+) includes nested directories.
- Embedded assets must live under `embed/assets/` — the same directory tree as the Go source files.
- Internal packages import `specfs "github.com/9Ashwin/spec-cli/embed"` and read from the exported `embed.FS` variables.

### Platform Definitions

- The 29-platform list in `internal/platform/platform.go` mirrors Comet's `src/core/platforms.ts`.
- Each `Platform` struct: `ID`, `Name`, `SkillsDir` (where skills are installed), `GlobalSkillsDir` (optional, for global scope), `DetectionPaths` (paths that indicate the platform is active, nil falls back to SkillsDir), `OpenSpecToolID` (passed to `openspec init --tools`).
- `DetectPlatforms()` checks `DetectionPaths` first; if nil, falls back to checking `SkillsDir` existence.
- Adding a new platform: add to `AllPlatforms` slice, verify `OpenSpecToolID` matches OpenSpec's registry.

### Schema Bundles

- Schema bundles live under `embed/assets/schemas/<name>/` and are embedded via `embed.SchemasFS`.
- Each bundle: `schema.yaml` (artifact definitions + apply flow), `VERSION`, `templates/` (artifact templates), `templates/adopters/` (CLAUDE.md fragments en + zh), `README.md`.
- `schema.InstallSchema()` copies the entire bundle directory from embed to `openspec/schemas/<name>/`.
- `schema.AppendClaudeMdFragment()` appends locale-specific fragment to `CLAUDE.md` if present and not already added.
- Schema version comparison in `spec-cli update` reads the embedded `VERSION` against the installed one.

### CLI Commands

- Root command in `cmd/root.go`: `func Execute() int`, called by `main.go` via `os.Exit(cmd.Execute())`.
- Each subcommand is a `var ...Cmd = &cobra.Command{...}` with `RunE` for error returns.
- Flags use `cmd.Flags().BoolVar()` / `cmd.Flags().StringVar()`, registered in `init()`.
- Subcommands registered in `root.go`'s `init()` via `rootCmd.AddCommand(...)`.
- `--json` flag on all commands outputs structured JSON to stdout (matching Comet's TypeScript JSON format).
- Interactive prompts use huh (`huh.NewSelect`); non-interactive mode via `--yes` flag.

### Init Flow

`speed-cli init` executes a 10-step flow:

1. Detect platforms (auto, via `platform.DetectPlatforms()`)
2. Select scope (project / global) — interactive or `--scope`
3. Select language (en / zh) — interactive or defaults to en
4. Select platforms — interactive multi-select, auto in `--yes` mode
5. Determine base directory (project path or home)
6. Install OpenSpec (`openspec init <path> --tools <ids>`, auto-installs CLI if missing)
7. Detect Superpowers (check `~/.claude/plugins/cache/*/superpowers/*/skills/`)
8. Install Comet skill (`skill.CopySkills()` from embed)
9. Create working directories (`docs/superpowers/specs/`, `docs/superpowers/plans/`)
10. Install schema bundles (`schema.InstallSchema()` + `AppendClaudeMdFragment()`)

### Error Handling

- Commands return errors via `RunE`; cobra prints to stderr.
- `SilenceUsage: true` and `SilenceErrors: true` on root to suppress cobra's default error output.
- `Execute()` in `cmd/root.go` catches errors and prints `Error: <msg>` to stderr, returns exit code 1.
- Internal packages wrap errors with `fmt.Errorf("...: %w", err)` for context.
- Non-critical failures (e.g., `openspec list --json` fails) return nil results rather than errors.
