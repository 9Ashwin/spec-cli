# spec-cli Handoff

## What We're Building

**spec-cli** — a Go CLI tool that installs OpenSpec, Superpowers, and schema bundles into AI coding platform projects. Clean-room rewrite of Comet's scaffolding features, with workflow execution delegated to OpenSpec's native `--schema` mechanism.

- **Repo:** https://github.com/9Ashwin/spec-cli
- **Local:** `/Users/solariswu/workspaces/github/spec-cli`
- **Module:** `github.com/9Ashwin/spec-cli`

## What's Done

1. **Design:** [`docs/specs/2026-05-27-spec-cli-go-design.md`](https://github.com/9Ashwin/spec-cli/blob/main/docs/specs/2026-05-27-spec-cli-go-design.md) — Go project structure, embed strategy, init flow, vfs pattern
2. **Plan:** [`docs/specs/2026-05-27-spec-cli-go-plan.md`](https://github.com/9Ashwin/spec-cli/blob/main/docs/specs/2026-05-27-spec-cli-go-plan.md) — 11 tasks, bite-sized steps
3. **Cleanup:** Node.js code removed, Go module init'd
4. **AGENTS.md:** [`AGENTS.md`](https://github.com/9Ashwin/spec-cli/blob/main/AGENTS.md) — reference project paths, build commands, architecture overview
5. **Assets preserved:** `assets/skills/comet/SKILL.md`, `assets/skills-zh/comet/SKILL.md`, `assets/schemas/superpowers-bridge/` (schema.yaml + templates + adopters)

## What's Next

Execute the Go implementation plan. **11 tasks total**, starting from Task 2 (Task 1 is done):

- **Task 2:** Create `internal/build/build.go` + `internal/vfs/vfs.go`
- **Task 3:** Create `internal/platform/` (29 platforms + detection)
- **Task 4:** Create `internal/openspec/` (exec openspec CLI)
- **Task 5:** Create `embed/` package + `internal/skill/` (embed + copy)
- **Task 6:** Create `internal/schema/` (schema install)
- **Task 7:** Create `cmd/root.go` (cobra root)
- **Task 8:** Create `cmd/init.go` (interactive init with huh)
- **Task 9:** Create `cmd/status.go`, `cmd/update.go`, `cmd/doctor.go`
- **Task 10:** Create `Makefile`
- **Task 11:** Final verification (vet, build, fmt)

## Key Dependencies

```bash
go get github.com/spf13/cobra          # CLI framework
go get github.com/charmbracelet/huh    # Interactive prompts
go get github.com/charmbracelet/lipgloss # Terminal styling
```

## Reference Projects

| Project | Local Path | What To Reference |
|---------|-----------|-------------------|
| lark-cli | `/Users/solariswu/workspaces/github/cli` | Go patterns: vfs interface, build ldflags, cobra root, embed |
| comet | `/Users/solariswu/workspaces/github/comet` | Original TS code: 29 platforms, init flow logic |
| openspec-schemas | `/Users/solariswu/workspaces/github/openspec-schemas` | Schema.yaml format reference |

## Key Design Decisions

- **vfs.FS interface** for all filesystem ops (testable, from lark-cli pattern)
- **embed at module root** (`embed/skills.go`, `embed/schemas.go`) — Go's `//go:embed` can't use `..` paths
- **cobra commands** are thin — parse flags, call internal packages
- **huh for interactive prompts** — scope select, language select, platform multi-select, schema multi-select
- **Single binary** distribution, no npm dependency at runtime (except calling `openspec init` which installs itself)
- **No phase skills, no shell scripts, no .comet.yaml** — all removed, workflow executed by OpenSpec

## Suggested Skills

- `superpowers:subagent-driven-development` — execute the implementation plan task-by-task with fresh subagents
- `superpowers:test-driven-development` — write Go tests before implementation code
- `superpowers:using-git-worktrees` — isolate work in a git worktree if needed
