# AGENTS.md

## Goal

Keep `spec-cli` a thin, predictable Go CLI that scaffolds OpenSpec + Superpowers workflows into AI coding projects.

Good changes usually fit one of these goals:

- Improve setup UX: clearer prompts, safer defaults, better help text, and actionable errors.
- Improve agent reliability: stable `--json` output, deterministic non-interactive flows, and stderr/stdout separation.
- Improve installation correctness: platform detection, skill copy, schema install, and update/doctor behavior.
- Improve maintainability: small command wrappers, testable internal packages, VFS-backed filesystem access, and embedded assets that stay in sync.

Do not turn `spec-cli` into a workflow runner. OpenSpec owns lifecycle execution through `openspec ... --schema superpowers-bridge`; this project only installs and refreshes the scaffolding that makes that flow available.

## Build & Test

```bash
make build              # Build spec-cli with Version/Date/CommitHash ldflags
make unit-test          # Run go test -race -count=1 ./...
make test               # Required before PR: fmt-check + vet + unit-test
make lint               # golangci-lint, including forbidigo VFS rules
make check-schema-sync  # Verify embedded schema assets stay synchronized
```

If dependencies change, run `go mod tidy` and confirm `go.mod` / `go.sum` only contain intentional changes.

## Pre-PR Checks

1. `make test`
2. `make lint`
3. `make check-schema-sync`
4. `go mod tidy` and inspect any module file diff
5. For release or npm installer changes: `make release` and verify generated archive/checksum names

## Who Uses This CLI

The primary users are humans working with AI coding tools, and AI agents automating setup inside user projects. Treat command output as part of the product API:

- stdout is data. JSON envelopes and machine-readable results go there.
- stderr is for progress, prompts, warnings, hints, and diagnostics.
- `--json` must suppress decorative/progress text and return structured output.
- Error messages should be specific enough for an agent to choose a next action.
- Non-interactive mode (`--yes`) must avoid prompts and pick deterministic defaults.

## Architecture Overview

`spec-cli` is a single Cobra binary. Assets are embedded at compile time, business logic lives in `internal/`, and all internal filesystem access goes through `internal/vfs`.

| Layer | Tech | Responsibility |
|-------|------|----------------|
| `cmd/` | Cobra + huh | Parse flags, prompt users, print results, call internal packages |
| `internal/platform/` | Go stdlib | Detect AI coding platforms and map them to OpenSpec tool IDs |
| `internal/skill/` | embed.FS + VFS | Copy the thin `opsx:super` skill into platform skill directories |
| `internal/schema/` | embed.FS + VFS | Install schema bundles into `openspec/schemas/` and append adopters |
| `internal/openspec/` | Runner interface | Wrap `openspec` CLI calls behind testable command execution |
| `embed/` | go:embed | Own all embedded skills, schema bundles, templates, and adopters |

## Source Layout

| Path | What it does |
|------|--------------|
| `main.go` | Calls `os.Exit(cmd.Execute())` |
| `cmd/root.go` | Root Cobra command, version template, subcommand registration |
| `cmd/helpers.go` | Shared path resolution, JSON printing, conditional printer, constants |
| `cmd/init.go` | Interactive/non-interactive scaffolding flow |
| `cmd/status.go` | `openspec list --json` wrapper for active changes |
| `cmd/update.go` | Refresh embedded skills and schema bundles |
| `cmd/doctor.go` | Diagnose OpenSpec CLI, schemas, and installed skill files |
| `cmd/completion.go` | Shell completion for bash/zsh/fish/powershell |
| `internal/vfs/` | Filesystem interface, OS implementation, and package-level helpers |
| `internal/platform/` | Platform registry and project detection |
| `internal/openspec/` | `Runner`, `ExecRunner`, install/list/version helpers |
| `internal/schema/` | Schema list/install/version/adopter fragment helpers |
| `internal/skill/` | Embedded skill copy logic |
| `internal/build/` | Build-time version/date/commit metadata |
| `embed/assets/skills*/` | English and Simplified Chinese `opsx:super` skill assets |
| `embed/assets/schemas/` | Embedded OpenSpec schema bundles |
| `scripts/install.js` | npm postinstall binary downloader with mirror fallback |
| `scripts/run.js` | npm launcher for the downloaded binary |

## Init Flow

`spec-cli init` runs this flow:

1. Detect installed AI coding platforms with `platform.DetectPlatforms(projectPath)`.
2. Select scope: `project` or `global`.
3. Select skill language: `en` or `zh`.
4. Select target platforms; `--yes` uses detected platforms, or all platforms if none are detected.
5. Resolve the base directory for project/global skill installation.
6. Ensure the OpenSpec CLI is available, then run `openspec init <path> --tools <ids>` through `internal/openspec`.
7. Detect Superpowers in the Claude plugin cache.
8. Copy the embedded `opsx:super` skill into each selected platform.
9. Install embedded schema bundles and append locale-specific `CLAUDE.md` adopter fragments.

## Code Rules

### VFS

- Internal packages must use `internal/vfs` helpers or accept a `vfs.FS`.
- Do not call `os.ReadFile`, `os.WriteFile`, `os.MkdirAll`, `os.Stat`, `os.Remove`, or `os.RemoveAll` directly in `internal/`.
- Tests can replace `vfs.DefaultFS`; package helpers such as `vfs.ReadFile()` delegate to it.

### Embedded Assets

- `//go:embed` directives live in the root `embed/` package.
- Embedded paths use forward slashes and `path.Join`; OS destination paths use `filepath.Join`.
- Embedded files must live under `embed/assets/`.
- When schema source assets change, keep generated/synchronized schema files aligned and run `make check-schema-sync`.

### Platform Registry

- Platform definitions live in `internal/platform/platform.go`.
- Keep `ID`, `Name`, `SkillsDir`, optional `GlobalSkillsDir`, `DetectionPaths`, and `OpenSpecToolID` aligned with OpenSpec's tool registry.
- `DetectionPaths` takes precedence; if absent, detection falls back to `SkillsDir`.
- Adding a platform requires tests for registry validity and detection behavior.

### Commands

- Keep Cobra commands thin: parse flags, resolve paths, call internal packages, print results.
- Commands return errors from `RunE`; `cmd.Execute()` prints `Error: <msg>` and returns exit code `1`.
- Add `--json` to agent-facing commands and keep the schema stable.
- Shared literals such as scope/language values belong in `cmd/helpers.go`.

### OpenSpec Runner

- External `openspec` calls go through `internal/openspec.Runner`.
- Production uses `ExecRunner`; tests should replace `DefaultRunner`.
- Non-critical list/status failures may return empty results when that keeps diagnostics useful.

## Release & Distribution

- Go module: `github.com/9Ashwin/spec-cli`
- npm package: `ashwin-spec`
- Binary name: `spec-cli`
- Release artifacts are built by `make release` / GoReleaser for darwin, linux, and windows.
- `scripts/install.js` downloads the matching archive and verifies checksums; keep installer changes conservative and cross-platform.
- For npm releases: push a `v*` tag. The release workflow runs GoReleaser (build + GitHub release) then publishes to npm via `NPM_TOKEN`. The token is stored as a GitHub Actions secret — no local npm authentication is needed.
- Do not push the tag until the matching commit is on `main` and `make test` + `make lint` pass, because postinstall may fall back to GitHub release downloads when packaged archives are absent.

## Reference Projects

| Project | Path | Useful patterns |
|---------|------|-----------------|
| lark-cli | `/Users/mervyn/workspaces/github/cli` | Agent-first CLI output, VFS conventions, Cobra structure, release/install patterns |
| Memoh | `/Users/mervyn/workspaces/github/Memoh` | Rich AGENTS architecture docs, workspace/runtime boundaries, agent harness ideas to evaluate separately |
