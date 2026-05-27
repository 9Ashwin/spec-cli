# spec-cli

OpenSpec + Superpowers workflow scaffolding tool (Go). Installs OpenSpec skills, Superpowers skills, and schema bundles into AI coding platform projects.

## Reference Projects

| Project | Path | Purpose |
|---------|------|---------|
| lark-cli | `/Users/solariswu/workspaces/github/cli` | Go patterns: cobra, huh, vfs, embed, build ldflags |
| comet | `/Users/solariswu/workspaces/github/comet` | Original TypeScript codebase: 29 platforms, init flow, openspec/superpowers install logic |
| openspec-schemas | `/Users/solariswu/workspaces/github/openspec-schemas` | Schema.yaml format: artifact-driven schema + templates + adopters convention |

## Build

```bash
make build      # build binary with ldflags version injection
make test       # go test -race ./...
make vet        # go vet ./...
make fmt        # gofmt -s -w .
```

## Architecture

```
cmd/            # cobra commands (thin — parse args, call internal)
  root.go       # root command + Execute()
  init.go       # spec-cli init (interactive scaffolding)
  status.go     # spec-cli status (reads openspec list --json)
  update.go     # spec-cli update (refresh skills + schemas)
  doctor.go     # spec-cli doctor (diagnose install health)
internal/
  build/        # Version/Date ldflags (from lark-cli)
  vfs/          # FS interface + OsFs + DefaultFS (from lark-cli)
  platform/     # 29 platform definitions + auto-detection
  skill/        # Copy skills from embed to target dirs
  openspec/     # openspec CLI integration (exec)
  schema/       # Schema bundle install from embed
embed/
  skills.go     # //go:embed assets/skills/* + skills-zh/*
  schemas.go    # //go:embed assets/schemas/*
assets/
  skills/       # English comet SKILL.md
  skills-zh/    # Chinese comet SKILL.md
  schemas/      # Schema bundles (superpowers-bridge/)
```

## Key Patterns (from lark-cli)

- `vfs.DefaultFS` — global filesystem, replace in tests
- `internal/build/build.go` — `var Version = "DEV"` with ldflags override
- `cmd/root.go` — `func Execute() int`, called by `main.go` `os.Exit(cmd.Execute())`
- cobra commands use `RunE`, flags via `cmd.Flags().BoolVar(...)`, registered in `init()`
- embed at module root level (`embed/` package) because `//go:embed` can't use `..` paths
