# spec-cli

OpenSpec + Superpowers workflow scaffolding tool. Detects AI coding platforms and installs OpenSpec skills, Superpowers skills, a thin `/comet` entry skill, and schema bundles.

## Install

```bash
go install github.com/9Ashwin/spec-cli@latest
```

Or build from source:

```bash
git clone https://github.com/9Ashwin/spec-cli.git
cd spec-cli
make build
```

## Usage

```bash
# Interactive setup
spec-cli init

# Non-interactive, auto-detect platforms
spec-cli init --yes

# Show active workflow changes
spec-cli status

# Update skills and schemas
spec-cli update

# Diagnose installation health
spec-cli doctor
```

### Init Options

```
spec-cli init [path]
  --yes               Non-interactive, auto-select detected platforms
  --scope <scope>     project | global
  --skip-existing     Skip already installed components
  --overwrite         Overwrite all existing components
  --json              Output structured JSON result
```

## Supported Platforms

29 AI coding platforms: Claude Code, Cursor, Codex, OpenCode, Windsurf, Cline, RooCode, Continue, GitHub Copilot, Gemini CLI, Amazon Q Developer, Qwen Code, Kilo Code, Auggie, Kiro, Lingma, Junie, CodeBuddy Code, CoStrict, Crush, Factory Droid, iFlow, Pi, Qoder, Antigravity, Bob Shell, ForgeCode, Trae

## How It Works

`speed-cli init` runs a 10-step flow:

1. Detect installed AI coding platforms
2. Select install scope (project / global)
3. Select language (English / 简体中文)
4. Select target platforms
5. Install OpenSpec CLI + init with `--tools`
6. Detect Superpowers plugin
7. Copy Comet entry skill to platform skills dirs
8. Create working directories (`docs/superpowers/`)
9. Install schema bundles to `openspec/schemas/`
10. Append CLAUDE.md workflow fragment

Workflow execution is delegated to OpenSpec's native `--schema` mechanism — spec-cli handles scaffolding only.

## License

MIT
