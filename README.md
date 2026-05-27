# spec-cli

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://go.dev/)

OpenSpec + Superpowers workflow scaffolding tool — built for humans and AI Agents. Detects AI coding platforms and installs OpenSpec skills, Superpowers skills, a thin `/comet` entry skill, and schema bundles with a single command.

[Install](#install) · [Commands](#commands) · [Supported Platforms](#supported-platforms) · [Development](#development)

## Why spec-cli?

- **Zero-Config Detection** — auto-detects 29 AI coding platforms, no manual config needed
- **Single Command Setup** — `spec-cli init` scaffolds the complete OpenSpec + Superpowers workflow
- **Agent-Native** — installs skills that AI coding agents understand natively
- **Schema-Driven Workflows** — delegates execution to OpenSpec's `--schema` mechanism, no phase state management
- **Clean-Room Go Rewrite** — single binary, zero npm runtime dependencies (except openspec CLI itself)
- **Interactive & Non-Interactive** — huh-powered prompts for humans, `--yes` + `--json` flags for agents and scripts

## Features

| Category | Capabilities |
|----------|-------------|
| Platform Detection | Auto-detect 29 AI coding platforms from project files |
| OpenSpec Install | `openspec init` with `--tools` for detected platforms, auto-install CLI if missing |
| Superpowers Detection | Check Claude Code plugin cache for installed Superpowers skills |
| Skill Copy | Copy Comet entry skill from embed to platform skills directories |
| Schema Bundles | Install workflow schema bundles to `openspec/schemas/` with CLAUDE.md fragments |
| Health Check | `spec-cli doctor` diagnoses OpenSpec CLI, working dirs, schemas, and skill files |
| Update | `spec-cli update` refreshes skills from embed and upgrades schema versions |

## Install

### Requirements

- Node.js 16+ (`npm`/`npx`)
- Go 1.23+ (build from source only)

### Quick Start (Human Users)

```bash
# Option 1 — From npm (recommended):
npx spec-cli@latest init

# Option 2 — From source:
git clone https://github.com/9Ashwin/spec-cli.git
cd spec-cli
make install
```

### Quick Start (AI Agent)

> The following steps are for AI Agents helping the user with installation.

```bash
# Option 1 — From npm (recommended):
npx spec-cli@latest init --yes

# Option 2 — From source:
git clone https://github.com/9Ashwin/spec-cli.git /tmp/spec-cli
cd /tmp/spec-cli
make install
cd /path/to/user/project
spec-cli init --yes
```

## Commands

| Command | Description |
|---------|-------------|
| `spec-cli init [path]` | Initialize workflow scaffolding (interactive by default) |
| `spec-cli status [path]` | Show active workflow changes |
| `spec-cli update [path]` | Update skills and schema bundles |
| `spec-cli doctor [path]` | Diagnose installation health |

### Init Flags

```
spec-cli init [path]
  --yes               Non-interactive, auto-select detected platforms
  --scope <scope>     project | global
  --skip-existing     Skip already installed components
  --overwrite         Overwrite all existing components
  --json              Output structured JSON result
```

### Shell Completion

```bash
# bash
source <(spec-cli completion bash)

# zsh
source <(spec-cli completion zsh)

# fish
spec-cli completion fish | source
```

## Supported Platforms

Claude Code, Cursor, Codex, OpenCode, Windsurf, Cline, RooCode, Continue, GitHub Copilot, Gemini CLI, Amazon Q Developer, Qwen Code, Kilo Code, Auggie, Kiro, Lingma, Junie, CodeBuddy Code, CoStrict, Crush, Factory Droid, iFlow, Pi, Qoder, Antigravity, Bob Shell, ForgeCode, Trae

## How It Works

`speed-cli init` runs a 10-step flow:

1. Detect installed AI coding platforms
2. Select install scope (project / global)
3. Select language (English / 简体中文)
4. Select target platforms
5. Install OpenSpec CLI + `openspec init <path> --tools <ids>`
6. Detect Superpowers plugin installation
7. Copy Comet entry skill to platform skills directories
8. Create working directories (`docs/superpowers/specs/`, `docs/superpowers/plans/`)
9. Install schema bundles to `openspec/schemas/<name>/`
10. Append CLAUDE.md workflow fragment

Workflow execution is delegated to OpenSpec's native `--schema` mechanism — spec-cli handles scaffolding only.

## Development

```bash
make build      # Build binary
make test       # Run tests with race detector
make vet        # Run go vet
make fmt        # Format source files
make fmt-check  # Check formatting (CI)
make install    # Install to /usr/local/bin
make clean      # Remove binary
```

## License

MIT
