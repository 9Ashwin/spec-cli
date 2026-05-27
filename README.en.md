# spec-cli

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.25.7-blue.svg)](https://go.dev/)
[![npm version](https://img.shields.io/npm/v/ashwin-spec.svg)](https://www.npmjs.com/package/ashwin-spec)

OpenSpec + Superpowers workflow scaffolding tool — built for humans and AI Agents. Detects AI coding platforms and installs OpenSpec skills, a thin `/opsx:super` entry skill, and schema bundles with a single command.

[Install](#installation--quick-start) · [Commands](#commands) · [Supported Platforms](#supported-platforms) · [How It Works](#how-it-works) · [Development](#development)

## Why spec-cli?

- **Zero-Config Detection** — auto-detects 29 AI coding platforms, no manual config needed
- **Single Command Setup** — `spec-cli init` installs OpenSpec skills, the entry skill, and schema bundles
- **Agent-Native** — installs skills that AI coding agents understand natively
- **Schema-Driven Workflows** — delegates execution to OpenSpec's `--schema` mechanism, no phase state management
- **Single Binary** — Go CLI; the npm package only downloads and launches the platform binary (except openspec CLI itself)
- **Interactive & Non-Interactive** — huh-powered prompts for humans, `--yes` + `--json` flags for agents and scripts

## Installation & Quick Start

### Requirements

Before you start, make sure you have:

- Node.js 16+ (`npm`/`npx`)
- Go 1.25.7+ (only required for building from source)

### Quick Start (Human Users)

> **Note for AI assistants:** If you are an AI Agent helping the user with installation, jump directly to [Quick Start (AI Agent)](#quick-start-ai-agent), which contains all the steps you need to complete.

#### Install

Choose **one** of the following methods:

**Option 1 — Global npm install (recommended):**

```bash
npm install -g ashwin-spec@latest
```

**Option 2 — Run on demand:**

```bash
npx ashwin-spec@latest init
```

**Option 3 — From source:**

Requires Go 1.25.7+.

```bash
git clone https://github.com/9Ashwin/spec-cli.git
cd spec-cli
go install .
```

#### Use

```bash
# Step 1: Initialize workflow scaffolding (interactive)
spec-cli init

# Step 2: Start a new change in your AI coding platform with /opsx:super
# Or create one directly from the terminal:
openspec new change "your-change-name" --schema superpowers-bridge --description "your feature idea"

# Step 3: Follow the schema-guided workflow
# brainstorming -> design -> plan -> build -> verify -> archive

# Check active changes
spec-cli status

# Diagnose installation health
spec-cli doctor
```

### Full Workflow

```
1. spec-cli init          — scaffold the workflow once

2. /opsx:super            — start a new change inside an AI coding platform
   or openspec new change — create a change directly from the terminal

3. brainstorming          — clarify requirements and converge on design
   design                 — write the design artifact
   plan                   — create the implementation plan
   build                  — implement and test
   verify                 — prove the result works
   archive                — archive and sync the spec

4. spec-cli status        — inspect active changes
   spec-cli doctor        — diagnose installation health
   spec-cli update        — refresh embedded skills and schemas
```

## Quick Start (AI Agent)

> The following steps are for AI Agents helping the user with installation.

**Step 1 — Install**

```bash
npm install -g ashwin-spec@latest
```

**Step 2 — Initialize**

```bash
spec-cli init --yes
```

**Step 3 — Verify**

```bash
spec-cli doctor
```

**Step 4 — Start working**

```bash
# Create a new change with the superpowers-bridge schema
openspec new change "your-change-name" --schema superpowers-bridge --description "your feature idea"
```

## Features

| Category | Capabilities |
|----------|-------------|
| Platform Detection | Auto-detect 29 AI coding platforms from project files |
| OpenSpec Install | `openspec init` with `--tools` for detected platforms, attempting npm install when the CLI is missing |
| Superpowers Detection | Check Claude Code plugin cache for installed Superpowers skills |
| Skill Copy | Copy the opsx:super entry skill from embed to platform skills directories |
| Schema Bundles | Install workflow schema bundles to `openspec/schemas/` with CLAUDE.md fragments |
| Health Check | `spec-cli doctor` diagnoses OpenSpec CLI, working dirs, schemas, and skill files |
| Update | `spec-cli update` refreshes skills from embed and upgrades schema versions |

## Commands

| Command | Description |
|---------|-------------|
| `spec-cli init [path]` | Initialize workflow scaffolding (interactive by default) |
| `spec-cli status [path]` | Show active workflow changes |
| `spec-cli update [path]` | Update skills and schema bundles |
| `spec-cli doctor [path]` | Diagnose installation health |
| `spec-cli completion <shell>` | Generate shell completion scripts (bash/zsh/fish/powershell) |
| `spec-cli --version` | Show version, commit hash, and build date |

### Init Flags

```
spec-cli init [path]
  --yes               Non-interactive, auto-select detected platforms
  --scope <scope>     project | global
  --skip-existing     Skip already installed components
  --overwrite         Overwrite all existing components
  --json              Output structured JSON result
```

### Update Flags

```
spec-cli update [path]
  --scope <scope>       project | global
  --language <lang>     en | zh
  --json                Output structured JSON result
```

### Doctor Flags

```
spec-cli doctor [path]
  --scope <scope>       auto | project | global
  --json                Output structured JSON result
```

### Shell Completion

```bash
# bash
eval "$(spec-cli completion bash)"

# zsh
eval "$(spec-cli completion zsh)"

# fish
spec-cli completion fish | source

# powershell
spec-cli completion powershell | Out-String | Invoke-Expression
```

## Supported Platforms

Claude Code, Cursor, Codex, OpenCode, Windsurf, Cline, RooCode, Continue, GitHub Copilot, Gemini CLI, Amazon Q Developer, Qwen Code, Kilo Code, Auggie, Kiro, Lingma, Junie, CodeBuddy Code, CoStrict, Crush, Factory Droid, iFlow, Pi, Qoder, Antigravity, Bob Shell, ForgeCode, Trae

## How It Works

`spec-cli init` runs a 10-step flow:

1. Detect installed AI coding platforms
2. Select install scope (project / global)
3. Select language (English / 简体中文)
4. Select target platforms
5. Install OpenSpec CLI + `openspec init <path> --tools <ids>`
6. Detect Superpowers plugin installation
7. Copy the opsx:super entry skill to platform skills directories
8. Create working directories (`docs/superpowers/specs/`, `docs/superpowers/plans/`)
9. Install schema bundles to `openspec/schemas/<name>/`
10. Append CLAUDE.md workflow fragment

Workflow execution is delegated to OpenSpec's native `--schema` mechanism — spec-cli handles scaffolding only.

## Development

```bash
make build      # Build binary
make test       # Run tests with race detector
make vet        # Run go vet
make lint       # Run golangci-lint
make fmt        # Format source files
make fmt-check  # Check formatting (CI)
go install .    # Install to ~/go/bin
make clean      # Remove binary
```

### Development Infrastructure

- **golangci-lint**: `forbidigo` enforces VFS-mediated filesystem access under `internal/`.
- **husky + lint-staged**: `pnpm install` enables the pre-commit hook for formatting and checks.
- **gitleaks**: `.gitleaks.toml` configures secret leak detection.
- **GoReleaser / Make release**: release configuration builds multi-platform archives for npm distribution.

### Architecture Notes

- **VFS abstraction**: `internal/vfs/FS` keeps filesystem logic mockable in tests.
- **Runner interface**: `internal/openspec/Runner` abstracts `os/exec` so OpenSpec integration can be tested with a mock runner.
- **Cross-platform embed FS**: embedded paths use `path.Join`; local writes use `filepath.Join`.

## License

MIT
