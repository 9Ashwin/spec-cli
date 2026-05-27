# Comet Schema-Driven Architecture

## Context

Comet currently hardcodes its 5-phase workflow (open → design → build → verify → archive) across 6+ locations: shell scripts, phase skills, the TypeScript CLI, and the main entry skill. Adding a new workflow variant requires touching all of them.

OpenSpec provides a native schema mechanism: `schema.yaml` defines artifacts and their dependencies; `openspec new --schema <name>` selects a schema per change. The [openspec-schemas](https://github.com/JiangWay/openspec-schemas) community repository distributes schema bundles following a standard directory convention.

This design refactors Comet to act purely as a scaffolding tool (installer + entry point) while adopting openspec-schemas' artifact-driven schema pattern for workflow definition. Comet installs schemas; OpenSpec executes them.

## Goals / Non-Goals

**Goals:**
- Align comet-schema repo structure and schema.yaml format with openspec-schemas conventions
- Simplify comet CLI to focus on installation and status display
- Remove the hardcoded phase state machine (shell scripts, phase skills, .comet.yaml)
- Delegate workflow execution to OpenSpec's native `--schema` mechanism

**Non-Goals:**
- Modify OpenSpec CLI behavior
- Drop platform support (29 platforms remain)
- Change the npm distribution model (`npm install -g @rpamis/comet`)

## Architecture

```
comet (scaffolding)              comet-schema (content)
├── src/cli/                     ├── superpowers-bridge/
│   ├── index.ts                 │   ├── schema.yaml      ← artifacts + apply
│   └── commands/                │   ├── VERSION
│       ├── init.ts    ← +schema │   ├── templates/        ← artifact templates
│       ├── status.ts  ← 精简    │   ├── adopters/         ← CLAUDE.md fragments
│       └── update.ts  ← +schema │   └── README.md
├── assets/skills/comet/         └── spec-driven/          ← future
│   └── SKILL.md       ← 精简       └── ...
└── assets/manifest.json

User project (after comet init)
├── openspec/
│   ├── config.yaml
│   ├── schemas/                 ← comet installs here
│   │   └── superpowers-bridge/
│   └── changes/
├── .claude/skills/
│   ├── comet/SKILL.md           ← /comet entry
│   ├── openspec-*/              ← OpenSpec skills
│   └── brainstorming/           ← Superpowers skills
```

## comet-schema Structure

Each schema is a self-contained directory following the openspec-schemas convention:

```
superpowers-bridge/
├── schema.yaml       # artifact definitions + apply flow
├── VERSION            # semantic version (e.g., "1")
├── templates/         # markdown templates for each artifact
│   ├── brainstorm.md
│   ├── proposal.md
│   ├── design.md
│   ├── spec.md
│   ├── tasks.md
│   ├── plan.md
│   ├── verify.md
│   └── retrospective.md
├── adopters/          # CLAUDE.md routing fragments
│   ├── CLAUDE.md.fragment.md
│   └── CLAUDE.md.fragment.zh.md
└── README.md
```

## schema.yaml Format

Follows the OpenSpec artifact-driven format. Two top-level sections:

### artifacts

8 artifacts forming a dependency DAG:

```
brainstorm → proposal → specs → tasks → plan → verify → retrospective
                 ↘ design ↗
```

Each artifact defines:
- `id` — unique identifier
- `generates` — output filename
- `description` — one-line purpose
- `template` — path to template file under templates/
- `instruction` — agent execution instructions (skill invocation, output redirection, PRECHECK)
- `requires` — list of prerequisite artifact ids

### apply

Execution flow after plan completion:
- `requires: [plan]` — plan must exist before apply starts
- `tracks: tasks.md` — checkbox progress tracking
- `instruction` — step-by-step execution: workspace setup → subagent-driven-development → verify → retrospective → archive → PR

Strict PRECHECK: All required Superpowers skills must be available. Missing skills = hard stop.

## What Gets Removed from Comet

Script functions (handoff context, archive automation, guard checks) don't disappear — they move into the schema.yaml `apply.instruction` and individual artifact instructions, which tell the agent to run the equivalent `openspec` commands and file checks at the right time.

| Artifact | Reason |
|----------|--------|
| `assets/skills/comet-open/` ~ `comet-archive/` (7 skills) | Phase execution delegated to OpenSpec schema |
| `assets/skills/comet-hotfix/`, `comet-tweak/` (2 presets) | Become separate schema bundles |
| `assets/skills/comet/scripts/comet-state.sh` | State machine replaced by schema artifact DAG |
| `assets/skills/comet/scripts/comet-guard.sh` | Guards replaced by schema instruction PRECHECKs |
| `assets/skills/comet/scripts/comet-yaml-validate.sh` | Validation handled by `openspec schema validate` |
| `assets/skills/comet/scripts/comet-handoff.sh` | Context passing handled by schema instructions |
| `assets/skills/comet/scripts/comet-archive.sh` | Archive handled by `openspec archive` |
| `.comet.yaml` state machine logic | State tracked by OpenSpec natively |

## What Stays in Comet

| Artifact | Change |
|----------|--------|
| `src/cli/index.ts` | Unchanged |
| `src/core/detect.ts` (platform detection) | Unchanged |
| `src/core/skills.ts` (file copy) | Unchanged |
| `src/core/platforms.ts` (29 platforms) | Unchanged |
| `src/core/openspec.ts` | Extended: add schema install |
| `src/core/superpowers.ts` | Unchanged |
| `src/commands/init.ts` | Extended: add schema clone/copy step |
| `src/commands/status.ts` | Simplified: read `openspec list --json` |
| `src/commands/update.ts` | Extended: add schema update step |
| `src/commands/doctor.ts` | Simplified: check openspec/schemas/ |
| `assets/manifest.json` | Trimmed to comet entry skill only |
| `assets/skills/comet/SKILL.md` | Simplified: guide users to openspec commands |

## /comet Entry Skill (Simplified)

The entry skill no longer contains phase detection logic. It becomes a thin guide:

1. Check `openspec list --json` for active changes
2. Display active changes with their schemas
3. Output the next recommended command:

```
openspec new --schema superpowers-bridge <name>   # start a new change
openspec status --change "<name>"                  # check progress
openspec archive --change "<name>" -y              # archive when done
```

## Installation Flow

`comet init` gains a schema installation step:

```
comet init
  1. Platform detection (existing)
  2. Scope / language selection (existing)
  3. Install OpenSpec skills (existing)
  4. Install Superpowers skills (existing)
  5. Install Comet skill (simplified /comet entry)
  6. Install schemas (NEW)
     - git clone comet-schema to temp dir
     - list available schema bundles
     - user selects which to install
     - cp -R to openspec/schemas/<name>/
     - append CLAUDE.md fragment (if CLAUDE.md exists)
```

`comet update` gains schema upgrade:
- Compare local VERSION against remote
- Show diff, user confirms before overwrite

## Migration Path

1. **Phase 1: comet-schema repo** — Align directory structure and schema.yaml format with openspec-schemas. Add templates, adopters, README. (Already started at github.com/9Ashwin/comet-schema)
2. **Phase 2: comet CLI** — Add schema install to `comet init`/`comet update`. Simplify `comet status` to read OpenSpec output. Remove phase command mapping.
3. **Phase 3: skills cleanup** — Remove 7 phase skills and 5 shell scripts. Simplify `/comet` entry skill. Update manifest.json.
4. **Phase 4: test & release** — Update tests. Update changelog. Release as a new minor version with migration notes for existing users.
