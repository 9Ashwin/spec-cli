# Comet Schema-Driven Refactor — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor Comet from a hardcoded 5-phase state machine into a pure scaffolding tool that installs OpenSpec-compatible schema bundles, delegating workflow execution to OpenSpec's native `--schema` mechanism.

**Architecture:** Align comet-schema repo structure with openspec-schemas conventions (self-contained schema bundles with schema.yaml + templates + adopters). Extend comet CLI to clone and install schemas into `openspec/schemas/`. Remove all phase skills, shell scripts, and .comet.yaml logic. Simplify `/comet` entry skill to a thin guide pointing to `openspec` commands.

**Tech Stack:** TypeScript (commander), bash (shell scripts, to be removed), vitest (testing), npm (distribution)

---

### Task 1: Restructure comet-schema repo to openspec-schemas convention

**Files:**
- Modify: `schemas/full/schema.yaml` → `superpowers-bridge/schema.yaml`
- Create: `superpowers-bridge/VERSION`
- Create: `superpowers-bridge/templates/brainstorm.md`
- Create: `superpowers-bridge/templates/proposal.md`
- Create: `superpowers-bridge/templates/design.md`
- Create: `superpowers-bridge/templates/spec.md`
- Create: `superpowers-bridge/templates/tasks.md`
- Create: `superpowers-bridge/templates/plan.md`
- Create: `superpowers-bridge/templates/verify.md`
- Create: `superpowers-bridge/templates/retrospective.md`
- Create: `superpowers-bridge/adopters/CLAUDE.md.fragment.md`
- Create: `superpowers-bridge/adopters/CLAUDE.md.fragment.zh.md`
- Create: `superpowers-bridge/README.md`
- Delete: `schemas/` directory

- [ ] **Step 1: Rewrite schema.yaml in OpenSpec artifact-driven format**

Write the new `superpowers-bridge/schema.yaml` following the same format as openspec-schemas, with 8 artifacts (brainstorm, proposal, design, specs, tasks, plan, verify, retrospective) plus an apply section.

Content for `superpowers-bridge/schema.yaml`:

```yaml
name: superpowers-bridge
version: 1
description: >
  Spec-driven workflow integrated with Superpowers skills.
  Requirements: Superpowers plugin installed, providing skills:
  brainstorming, writing-plans, using-git-worktrees,
  subagent-driven-development, finishing-a-development-branch.
  Requires a subagent-capable platform (Claude Code, Codex, etc.).
  brainstorm → proposal → specs → tasks → plan → verify → retrospective.
  design is required (reorganizes raw brainstorm output into structured
  Context / Goals / Decisions / Risks / Migration).

artifacts:
  - id: brainstorm
    generates: brainstorm.md
    description: Collaborative design exploration using Superpowers brainstorming
    template: brainstorm.md
    instruction: |
      PRECHECK — required skill availability:
      Before invoking, confirm `superpowers:brainstorming` appears in
      your available skills list. If missing, STOP and inform the user
      that the Superpowers plugin must be installed. Do NOT silently fall back.

      Use the Skill tool to invoke **superpowers:brainstorming**.

      IMPORTANT output redirection:
      - Do NOT write to `docs/superpowers/specs/`. Instead, write the
        brainstorming output directly to this change's `brainstorm.md`.
      - brainstorm.md is a RAW CAPTURE — write the brainstorming
        skill's full output as-is.
      - Do NOT pre-fill design.md during this step.

      The brainstorming skill will:
      1. Explore project context
      2. Ask clarifying questions (one at a time)
      3. Propose 2-3 approaches with trade-offs
      4. Present design sections for approval
      5. Output the validated design
    requires: []

  - id: proposal
    generates: proposal.md
    description: Initial proposal document outlining the change
    template: proposal.md
    instruction: |
      Create the proposal document based on the brainstorming output.

      Read brainstorm.md for the validated design. Extract:
      - **Why**: The motivation
      - **What Changes**: Specific changes agreed upon
      - **Capabilities**: Identify which specs will be created or modified
      - **Impact**: Affected code, APIs, dependencies, or systems.

      Keep it concise (1-2 pages). Focus on the "why" not the "how".
    requires:
      - brainstorm

  - id: design
    generates: design.md
    description: Technical design document with implementation details
    template: design.md
    instruction: |
      Create the design document explaining HOW to implement the change.

      Read brainstorm.md as input. Reorganize into structured sections:
      - **Context**: Background, current state, constraints, stakeholders
      - **Goals / Non-Goals**: What this design achieves and excludes
      - **Decisions**: Key technical choices with rationale
      - **Risks / Trade-offs**: Known limitations
      - **Migration Plan**: Steps to deploy, rollback strategy
      - **Open Questions**: Outstanding decisions

      Focus on architecture and approach, not line-by-line implementation.
    requires:
      - brainstorm

  - id: specs
    generates: "specs/**/*.md"
    description: Detailed specifications for the change
    template: spec.md
    instruction: |
      Create specification files that define WHAT the system should do.

      One spec file per capability listed in the proposal's Capabilities section.
      Delta operations (use ## headers):
      - **ADDED Requirements**: New capabilities
      - **MODIFIED Requirements**: Changed behavior
      - **REMOVED Requirements**: Deprecated features
      - **RENAMED Requirements**: Name changes only

      Format requirements:
      - Each requirement: `### Requirement: <name>`
      - Use SHALL/MUST for normative requirements
      - Each scenario: `#### Scenario: <name>` with WHEN/THEN format
      - Every requirement MUST have at least one scenario.
    requires:
      - proposal

  - id: tasks
    generates: tasks.md
    description: Implementation checklist with trackable tasks
    template: tasks.md
    instruction: |
      Create the task list that breaks down the implementation work.

      **IMPORTANT: Follow the template below exactly.** The apply phase parses
      checkbox format to track progress.

      Guidelines:
      - Group related tasks under ## numbered headings
      - Each task MUST be a checkbox: `- [ ] X.Y Task description`
      - Tasks should be small enough to complete in one session
      - Order tasks by dependency

      Reference specs for what needs to be built and design.md for how.
    requires:
      - specs
      - design

  - id: plan
    generates: plan.md
    description: Micro-task implementation plan using Superpowers writing-plans
    template: plan.md
    instruction: |
      PRECHECK — required skill availability:
      Before invoking, confirm `superpowers:writing-plans` appears in
      your available skills list. If missing, STOP. Do NOT silently fall back.

      Use the Skill tool to invoke **superpowers:writing-plans**.

      IMPORTANT output redirection:
      - Write the plan directly to this change's `plan.md`.

      The writing-plans skill will:
      1. Read the tasks.md and design.md for context
      2. Break each task into 2-5 minute micro-steps (TDD style)
      3. Include exact file paths, code snippets, test commands
      4. Add commit points after each task
    requires:
      - tasks

  - id: verify
    generates: verify.md
    description: Post-implementation verification against specs, design, and tasks
    template: verify.md
    instruction: |
      PRECHECK — implementation evidence:
      Before producing verify.md, run BOTH commands. If either
      returns 0, STOP.

      1. Commit evidence:
         git log --oneline $(git merge-base HEAD origin/main 2>/dev/null || git merge-base HEAD origin/master 2>/dev/null)..HEAD | wc -l
      2. Task progress:
         grep -c '^- \[x\]' openspec/changes/<change-name>/tasks.md

      Only after BOTH return positive numbers, proceed.

      Use the Skill tool to invoke **openspec-verify-change**.

      The verify step MUST perform:
      1. Structural validation: `openspec validate --all --json`
      2. Task completion: every checkbox in tasks.md is `- [x]`
      3. Delta spec sync state
      4. Design/specs coherence
      5. Implementation signal: all code committed
      6. Front-door routing leak detector (warning, non-blocking)
    requires:
      - plan

  - id: retrospective
    generates: retrospective.md
    description: Evidence-first retrospective of completed change
    template: retrospective.md
    instruction: |
      PRECHECK — verify completion evidence:
      Before producing retrospective.md, run:
      1. test -f openspec/changes/<change-name>/verify.md
      2. ! grep -q '^- \[x\] ❌ FAIL' openspec/changes/<change-name>/verify.md

      Only after both succeed, proceed.

      Write a retrospective with evidence-first analysis:
      §0) Evidence — quantitative front-matter
      §1) Wins — what worked well
      §2) Misses — what didn't, with severity markers
      §3) Plan deviations — tasks whose scope changed
      §4) Skill / workflow compliance — each apply-phase skill
      §5) Surprises — assumptions that turned out wrong
      §6) Promote candidates → long-term learning

      Write output to retrospective.md using the template structure.
    requires:
      - verify

apply:
  requires: [plan]
  tracks: tasks.md
  instruction: |
    Before implementing, set up an isolated workspace and executor:

    0. **Pre-flight — verify required Superpowers skills**:
       - superpowers:using-git-worktrees
       - superpowers:subagent-driven-development
       - superpowers:finishing-a-development-branch
       If any required skill is missing, STOP.

    1. **Workspace**: Use **superpowers:using-git-worktrees** to create
       an isolated git worktree.

    2. **Executor**: Use **superpowers:subagent-driven-development** to
       execute the plan.md micro-tasks.

    3. **Verification**: Produce the `verify` artifact.

    4. **Retrospective**: Produce the `retrospective` artifact BEFORE PR.

    5. **Archive**: Run `openspec archive -y` to sync delta specs and
       move the change folder to archive.

    6. **Completion**: Use **superpowers:finishing-a-development-branch**.
```

- [ ] **Step 2: Create VERSION file**

Write `superpowers-bridge/VERSION`:

```
1
```

- [ ] **Step 3: Create artifact templates**

Write `superpowers-bridge/templates/brainstorm.md`:

```markdown
# Brainstorming: {{change_name}}

## Background

## Decision Chain

### Q1: [question]

**Decision:** [answer]

**Rationale:** [why]

### Q2: [question]

**Decision:** [answer]

**Rationale:** [why]

## Design Trade-offs

## Approved Design
```

Write `superpowers-bridge/templates/proposal.md`:

```markdown
# Proposal: {{change_name}}

## Why

## What Changes

## Capabilities

### New Capabilities
- `capability-name`: [description]

### Modified Capabilities
- `existing-capability`: [what changes]

## Impact
```

Write `superpowers-bridge/templates/design.md`:

```markdown
# Design: {{change_name}}

## Context

## Goals / Non-Goals

### Goals

### Non-Goals

## Decisions

## Risks / Trade-offs

## Migration Plan

## Open Questions
```

Write `superpowers-bridge/templates/spec.md`:

```markdown
# Spec: {{capability_name}}

## ADDED Requirements

### Requirement: [name]

[Description using SHALL/MUST]

#### Scenario: [name]

- **WHEN** [condition]
- **THEN** [expected outcome]

## MODIFIED Requirements

### Requirement: [name]

[Full updated requirement text]

#### Scenario: [name]

- **WHEN** [condition]
- **THEN** [expected outcome]
```

Write `superpowers-bridge/templates/tasks.md`:

```markdown
# Tasks: {{change_name}}

## 1. [Phase Name]

- [ ] 1.1 [Task description]
- [ ] 1.2 [Task description]

## 2. [Phase Name]

- [ ] 2.1 [Task description]
- [ ] 2.2 [Task description]
```

Write `superpowers-bridge/templates/plan.md`:

```markdown
# Plan: {{change_name}}

> change: {{change_name}}
> design-doc: {{design_doc_path}}

## Micro-Tasks

### Task 1: [Name]

- [ ] Step 1: [description]
- [ ] Step 2: [description]
- [ ] Commit: `git commit -m "[message]"`
```

Write `superpowers-bridge/templates/verify.md`:

```markdown
# Verification: {{change_name}}

## Overall Decision

- [ ] ✅ PASS
- [ ] ⚠️ PASS WITH WARNINGS
- [ ] ❌ FAIL

## 1. Structural Validation

## 2. Task Completion

## 3. Delta Spec Sync State

## 4. Design/Specs Coherence

## 5. Implementation Signal

## 6. Front-Door Routing Leak Detector

## 7. Deferred Dogfood vs Automated Test Equivalence
```

Write `superpowers-bridge/templates/retrospective.md`:

```markdown
# Retrospective: {{change_name}}

## §0 Evidence

- **Commits:** [count]
- **Diff:** +[added] -[removed] ([files] files)
- **Tasks:** [done]/[total] done
- **Active hours:** [hours]
- **Subagent dispatches:** [count]
- **New dependencies:** [list or none]
- **Post-merge bugs:** [count]
- **OpenSpec validate:** [state]
- **Test coverage:** [signal]
- **Commit chain:** [one-line summary]

## §1 Wins

## §2 Misses

## §3 Plan Deviations

## §4 Skill / Workflow Compliance

## §5 Surprises

## §6 Promote Candidates → Long-Term Learning

- [ ] [emoji] [one-sentence learning]
  → **Promote to** [destination]
  > **Why**: [reason]
  > **How to apply**: [trigger condition]
```

- [ ] **Step 4: Create CLAUDE.md fragments**

Write `superpowers-bridge/adopters/CLAUDE.md.fragment.md`:

```markdown
## Comet + Superpowers Bridge Workflow

This project uses the Comet workflow with the `superpowers-bridge` schema.

### Starting a new change

```
/openspec:new --schema superpowers-bridge <change-name>
```

### Checking status

```
/openspec:status --change "<change-name>"
```

### Workflow

The schema defines 8 artifacts: brainstorm → proposal → design → specs → tasks → plan → verify → retrospective.
After plan completion, apply phase executes: workspace setup → subagent-driven-development → verify → retrospective → archive → PR.
```

Write `superpowers-bridge/adopters/CLAUDE.md.fragment.zh.md`:

```markdown
## Comet + Superpowers Bridge 工作流

本项目使用 Comet 工作流和 `superpowers-bridge` schema。

### 新建变更

```
/openspec:new --schema superpowers-bridge <change-name>
```

### 查看状态

```
/openspec:status --change "<change-name>"
```

### 工作流

Schema 定义了 8 个 artifact：brainstorm → proposal → design → specs → tasks → plan → verify → retrospective。
Plan 完成后进入 apply 阶段：工作空间设置 → subagent-driven-development → verify → retrospective → archive → PR。
```

- [ ] **Step 5: Create schema README**

Write `superpowers-bridge/README.md`:

```markdown
# superpowers-bridge Schema

Bridges OpenSpec's artifact governance (the **what**) with Superpowers execution skills (the **how**) into a single workflow.

Schema version: v1.

## Install

Run `comet init` in your project root. Comet will install this schema into `openspec/schemas/superpowers-bridge/`.

## Artifacts

| # | Artifact | Generates | Requires |
|---|----------|-----------|----------|
| 1 | brainstorm | brainstorm.md | — |
| 2 | proposal | proposal.md | brainstorm |
| 3 | design | design.md | brainstorm |
| 4 | specs | specs/*/spec.md | proposal |
| 5 | tasks | tasks.md | specs, design |
| 6 | plan | plan.md | tasks |
| 7 | verify | verify.md | plan |
| 8 | retrospective | retrospective.md | verify |
| — | apply | (execution) | plan |

## Requirements

- OpenSpec CLI (`npm install -g @fission-ai/openspec@latest`)
- Superpowers plugin for your AI coding platform

## License

MIT
```

- [ ] **Step 6: Remove old schemas/ directory and commit**

```bash
cd /Users/solariswu/workspaces/github/comet-schema
rm -rf schemas/
git add -A
git commit -m "restructure: align with openspec-schemas convention

Replace schemas/full/ with superpowers-bridge/ bundle following the
openspec-schemas directory structure: schema.yaml (artifact-driven),
templates/, adopters/, README.md, VERSION."
```

---

### Task 2: Add schema install to comet CLI (src/core/openspec.ts)

**Files:**
- Modify: `src/core/openspec.ts`

- [ ] **Step 1: Add schema install and utility functions**

Append to `src/core/openspec.ts`:

```typescript
import { promises as fs } from 'fs';
import path from 'path';

const COMET_SCHEMA_REPO = 'https://github.com/9Ashwin/comet-schema.git';

interface SchemaInfo {
  name: string;
  version: string;
  path: string;
}

async function cloneCometSchema(tempDir: string): Promise<string> {
  const repoDir = path.join(tempDir, 'comet-schema');
  execSync(`git clone --depth 1 ${COMET_SCHEMA_REPO} ${repoDir}`, {
    stdio: 'pipe',
    timeout: 60_000,
  });
  return repoDir;
}

async function listAvailableSchemas(repoDir: string): Promise<SchemaInfo[]> {
  const entries = await fs.readdir(repoDir, { withFileTypes: true });
  const schemas: SchemaInfo[] = [];

  for (const entry of entries) {
    if (!entry.isDirectory() || entry.name.startsWith('.')) continue;
    const schemaYaml = path.join(repoDir, entry.name, 'schema.yaml');
    const versionFile = path.join(repoDir, entry.name, 'VERSION');

    try {
      await fs.access(schemaYaml);
      let version = 'unknown';
      try {
        version = (await fs.readFile(versionFile, 'utf-8')).trim();
      } catch {
        // version file optional
      }
      schemas.push({ name: entry.name, version, path: path.join(repoDir, entry.name) });
    } catch {
      // skip directories without schema.yaml
    }
  }

  return schemas;
}

async function installSchema(
  schemaPath: string,
  projectPath: string,
  schemaName: string,
): Promise<boolean> {
  const targetDir = path.join(projectPath, 'openspec', 'schemas', schemaName);

  try {
    await fs.mkdir(path.dirname(targetDir), { recursive: true });
    // Remove existing if present
    await fs.rm(targetDir, { recursive: true, force: true });
    // Copy schema bundle
    await fs.cp(schemaPath, targetDir, { recursive: true });
    return true;
  } catch (error) {
    console.error(`    Failed to install schema ${schemaName}: ${(error as Error).message}`);
    return false;
  }
}

async function appendClaudeMdFragment(
  projectPath: string,
  schemaPath: string,
  locale: string,
): Promise<boolean> {
  const claudeMdPath = path.join(projectPath, 'CLAUDE.md');
  try {
    await fs.access(claudeMdPath);
  } catch {
    return false; // no CLAUDE.md, skip
  }

  const fragmentName = locale === 'zh'
    ? 'CLAUDE.md.fragment.zh.md'
    : 'CLAUDE.md.fragment.md';
  const fragmentPath = path.join(schemaPath, 'adopters', fragmentName);

  try {
    const fragment = await fs.readFile(fragmentPath, 'utf-8');
    const existing = await fs.readFile(claudeMdPath, 'utf-8');

    // Check if fragment already exists
    if (existing.includes('superpowers-bridge')) {
      return false; // already present
    }

    await fs.appendFile(claudeMdPath, `\n${fragment}`);
    return true;
  } catch {
    return false;
  }
}

export {
  installOpenSpec,
  isCommandAvailable,
  buildOpenSpecInitCommand,
  quoteShellArg,
  cloneCometSchema,
  listAvailableSchemas,
  installSchema,
  appendClaudeMdFragment,
  COMET_SCHEMA_REPO,
};
export type { SchemaInfo };
```

- [ ] **Step 2: Run tests to verify no regressions**

```bash
npx vitest run test/ts/openspec.test.ts
```

Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add src/core/openspec.ts
git commit -m "feat(openspec): add schema install, list, and Claude.md fragment functions"
```

---

### Task 3: Extend comet init with schema installation

**Files:**
- Modify: `src/commands/init.ts`

- [ ] **Step 1: Add schema install step to init flow**

After the working directories step (line 322) and before `displaySummary`, add schema installation.

Replace the `displaySummary` function signature and add a new function. In `src/commands/init.ts`:

Replace the `get started` lines in `displaySummary` (lines 176-179):

```typescript
  console.log(`\n  Get started:`);
  console.log(`    /openspec:new --schema superpowers-bridge "your idea"  — Start a new change\n`);
}
```

Add schema install logic after line 323 (`await createWorkingDirs(projectPath);`):

```typescript
  // Schema installation
  const os = await import('os');
  const tempDir = path.join(os.tmpdir(), `comet-schema-${Date.now()}`);
  let schemasInstalled = 0;

  try {
    const { cloneCometSchema, listAvailableSchemas, installSchema, appendClaudeMdFragment } =
      await import('../core/openspec.js');

    log(`\n  Fetching schemas from comet-schema...`);
    const repoDir = await cloneCometSchema(tempDir);
    const availableSchemas = await listAvailableSchemas(repoDir);

    if (availableSchemas.length === 0) {
      log(`  No schemas found.`);
    } else {
      log(`  Available schemas: ${availableSchemas.map((s) => `${s.name} (v${s.version})`).join(', ')}`);

      const schemaNames = options.yes
        ? availableSchemas.map((s) => s.name)
        : await checkbox({
            message: 'Select schemas to install:',
            choices: availableSchemas.map((s) => ({
              name: `${s.name} (v${s.version})`,
              value: s.name,
              checked: true,
            })),
            required: false,
          });

      for (const schema of availableSchemas) {
        if (!schemaNames.includes(schema.name)) continue;

        const installed = await installSchema(schema.path, projectPath, schema.name);
        if (installed) {
          schemasInstalled++;
          log(`  Schema: ${schema.name} installed -> openspec/schemas/${schema.name}/`);

          // Append CLAUDE.md fragment
          const fragmentAdded = await appendClaudeMdFragment(
            projectPath,
            schema.path,
            language.id,
          );
          if (fragmentAdded) {
            log(`  CLAUDE.md: appended ${schema.name} workflow fragment`);
          }
        }
      }
    }
  } catch (error) {
    log(`  Schema install skipped: ${(error as Error).message}`);
  } finally {
    await fs.rm(tempDir, { recursive: true, force: true }).catch(() => undefined);
  }
```

Update JSON output (line 341) to include schema info:

```typescript
          workingDirsCreated: scope === 'project',
          schemasInstalled,
```

Update imports at top — add:

```typescript
import { checkbox } from '@inquirer/prompts';
```

(in addition to the existing `{ checkbox, select }` import on line 3)

- [ ] **Step 2: Run tests**

```bash
npx vitest run test/ts/init.test.ts
```

Expected: PASS (existing tests, schema step skipped in test environment)

- [ ] **Step 3: Commit**

```bash
git add src/commands/init.ts
git commit -m "feat(init): add schema installation step to comet init"
```

---

### Task 4: Simplify comet status to use openspec list

**Files:**
- Modify: `src/commands/status.ts`

- [ ] **Step 1: Rewrite status.ts**

Replace the entire file content. The new status command reads `openspec list --json` instead of parsing `.comet.yaml`:

```typescript
import path from 'path';
import { execSync } from 'child_process';

interface OpenspecChange {
  name: string;
  schema?: string;
  artifacts?: Record<string, string>;
  tasksCompleted?: number;
  tasksTotal?: number;
}

interface StatusOptions {
  json?: boolean;
}

async function getActiveChanges(projectPath: string): Promise<OpenspecChange[]> {
  try {
    const output = execSync('openspec list --json', {
      cwd: projectPath,
      stdio: 'pipe',
      timeout: 15_000,
    }).toString().trim();

    if (!output) return [];

    const parsed = JSON.parse(output);
    if (Array.isArray(parsed)) return parsed;
    if (parsed.changes && Array.isArray(parsed.changes)) return parsed.changes;
    return [];
  } catch {
    return [];
  }
}

function displayStatus(changes: OpenspecChange[]): void {
  if (changes.length === 0) {
    console.log('No active changes.\n');
    console.log('Get started:');
    console.log('  openspec new --schema superpowers-bridge <name>\n');
    return;
  }

  console.log('Active Changes:\n');

  for (let i = 0; i < changes.length; i++) {
    const c = changes[i];
    const schema = c.schema ?? 'unknown';
    const taskStr =
      c.tasksTotal !== undefined && c.tasksTotal > 0
        ? ` [${c.tasksCompleted ?? 0}/${c.tasksTotal} tasks]`
        : '';
    console.log(`  ${i + 1}. ${c.name} [schema: ${schema}${taskStr}]`);
    if (c.artifacts) {
      for (const [artifact, file] of Object.entries(c.artifacts)) {
        if (file) console.log(`     ${artifact}: ${file}`);
      }
    }
    console.log();
  }
}

export async function statusCommand(
  targetPath: string,
  options: StatusOptions = {},
): Promise<void> {
  const projectPath = path.resolve(targetPath);
  const changes = await getActiveChanges(projectPath);

  if (options.json) {
    console.log(JSON.stringify({ changes }, null, 2));
    return;
  }

  displayStatus(changes);
}
```

- [ ] **Step 2: Run tests**

```bash
npx vitest run test/ts/status.test.ts
```

Expected: PASS (tests may need updates — see Task 10)

- [ ] **Step 3: Commit**

```bash
git add src/commands/status.ts
git commit -m "refactor(status): use openspec list --json instead of .comet.yaml"
```

---

### Task 5: Simplify doctor command

**Files:**
- Modify: `src/commands/doctor.ts`

- [ ] **Step 1: Remove .comet.yaml validation, add schema checks**

Remove `checkCometYamlValidity` and `checkScriptsPresent` functions. Add `checkSchemasPresent`:

```typescript
async function checkSchemasPresent(projectPath: string): Promise<CheckResult> {
  const schemasDir = path.join(projectPath, 'openspec', 'schemas');
  if (!(await fileExists(schemasDir))) {
    return {
      check: 'openspec schemas',
      status: 'warn',
      message: 'no schemas installed — run: comet init',
    };
  }

  const entries = await readDir(schemasDir);
  const schemaDirs = [];
  for (const entry of entries) {
    const stat = await fs.stat(path.join(schemasDir, entry));
    if (stat.isDirectory() && !entry.startsWith('.')) {
      schemaDirs.push(entry);
    }
  }

  if (schemaDirs.length === 0) {
    return {
      check: 'openspec schemas',
      status: 'warn',
      message: 'schemas directory exists but is empty',
    };
  }

  return {
    check: 'openspec schemas',
    status: 'pass',
    message: `${schemaDirs.length} installed (${schemaDirs.join(', ')})`,
  };
}
```

Update `collectResults` to call `checkSchemasPresent` instead of `checkScriptsPresent` and `checkCometYamlValidity`:

```typescript
async function collectResults(projectPath: string, scope: DoctorScope): Promise<CheckResult[]> {
  const results: CheckResult[] = [];
  results.push(await checkOpenSpecCli());
  if (scope !== 'global') {
    results.push(await checkWorkingDirs(projectPath));
    results.push(await checkSchemasPresent(projectPath));
  }
  results.push(...(await checkSkillCompleteness(projectPath, scope)));
  return results;
}
```

Remove the `VALID_YAML_FIELDS` constant (lines 19-32) and unused imports.

- [ ] **Step 2: Run tests**

```bash
npx vitest run test/ts/doctor.test.ts
```

Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add src/commands/doctor.ts
git commit -m "refactor(doctor): replace .comet.yaml checks with schema checks"
```

---

### Task 6: Extend comet update with schema upgrade

**Files:**
- Modify: `src/commands/update.ts`

- [ ] **Step 1: Add schema update step after skill update**

After the skill copy loop (around line 260), add schema update logic. Insert:

```typescript
  // Schema update
  log(`\n  Checking schema updates...`);
  const os = await import('os');
  const tempDir = path.join(os.tmpdir(), `comet-schema-update-${Date.now()}`);
  let schemasUpdated = 0;

  try {
    const { cloneCometSchema, listAvailableSchemas, installSchema } =
      await import('../core/openspec.js');

    const repoDir = await cloneCometSchema(tempDir);
    const remoteSchemas = await listAvailableSchemas(repoDir);
    const schemasDir = path.join(projectPath, 'openspec', 'schemas');

    for (const schema of remoteSchemas) {
      const localVersionFile = path.join(schemasDir, schema.name, 'VERSION');
      let localVersion = 'none';
      try {
        localVersion = (await fs.readFile(localVersionFile, 'utf-8')).trim();
      } catch {
        // schema not installed locally
      }

      if (localVersion !== schema.version) {
        log(`  Updating schema: ${schema.name} (v${localVersion} → v${schema.version})`);
        const installed = await installSchema(schema.path, projectPath, schema.name);
        if (installed) schemasUpdated++;
      } else {
        log(`  Schema ${schema.name}: up to date (v${localVersion})`);
      }
    }
  } catch (error) {
    log(`  Schema update skipped: ${(error as Error).message}`);
  } finally {
    await fs.rm(tempDir, { recursive: true, force: true }).catch(() => undefined);
  }
```

Update JSON output to include schema counts.

- [ ] **Step 2: Run tests**

```bash
npx vitest run test/ts/update.test.ts
```

Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add src/commands/update.ts
git commit -m "feat(update): add schema update step to comet update"
```

---

### Task 7: Remove phase skills and shell scripts

**Files:**
- Delete: `assets/skills/comet-open/SKILL.md`
- Delete: `assets/skills/comet-design/SKILL.md`
- Delete: `assets/skills/comet-build/SKILL.md`
- Delete: `assets/skills/comet-verify/SKILL.md`
- Delete: `assets/skills/comet-archive/SKILL.md`
- Delete: `assets/skills/comet-hotfix/SKILL.md`
- Delete: `assets/skills/comet-tweak/SKILL.md`
- Delete: `assets/skills/comet/scripts/comet-guard.sh`
- Delete: `assets/skills/comet/scripts/comet-state.sh`
- Delete: `assets/skills/comet/scripts/comet-handoff.sh`
- Delete: `assets/skills/comet/scripts/comet-archive.sh`
- Delete: `assets/skills/comet/scripts/comet-yaml-validate.sh`
- Delete: `assets/skills/comet/reference/dirty-worktree.md`
- Modify: `assets/manifest.json`
- Also delete Chinese counterparts: `assets/skills-zh/comet-*/SKILL.md` and `assets/skills-zh/comet/scripts/*`

- [ ] **Step 1: Remove skill directories**

```bash
rm -rf assets/skills/comet-open
rm -rf assets/skills/comet-design
rm -rf assets/skills/comet-build
rm -rf assets/skills/comet-verify
rm -rf assets/skills/comet-archive
rm -rf assets/skills/comet-hotfix
rm -rf assets/skills/comet-tweak
rm -rf assets/skills/comet/scripts
rm -rf assets/skills/comet/reference

rm -rf assets/skills-zh/comet-open
rm -rf assets/skills-zh/comet-design
rm -rf assets/skills-zh/comet-build
rm -rf assets/skills-zh/comet-verify
rm -rf assets/skills-zh/comet-archive
rm -rf assets/skills-zh/comet-hotfix
rm -rf assets/skills-zh/comet-tweak
rm -rf assets/skills-zh/comet/scripts
rm -rf assets/skills-zh/comet/reference
```

- [ ] **Step 2: Update manifest.json**

Replace `assets/manifest.json`:

```json
{
  "version": "0.4.0",
  "skills": [
    "comet/SKILL.md"
  ],
  "languages": [
    { "id": "en", "name": "English", "skillsDir": "skills" },
    { "id": "zh", "name": "中文", "skillsDir": "skills-zh" }
  ]
}
```

- [ ] **Step 3: Commit**

```bash
git add -A assets/
git commit -m "refactor: remove phase skills and shell scripts

Remove 7 phase skills (comet-open through comet-archive), 2 preset
skills (hotfix/tweak), 5 shell scripts (guard/state/handoff/archive/
yaml-validate), and dirty-worktree reference. Workflow execution is
now delegated to OpenSpec's native --schema mechanism."
```

---

### Task 8: Simplify /comet entry skill

**Files:**
- Modify: `assets/skills/comet/SKILL.md`
- Modify: `assets/skills-zh/comet/SKILL.md`

- [ ] **Step 1: Rewrite comet SKILL.md as a thin guide**

Replace `assets/skills/comet/SKILL.md`:

```markdown
---
name: comet
description: "Comet — OpenSpec + Superpowers dual-star development workflow. Start with /comet to check active changes and get recommended next commands."
---

# Comet — OpenSpec + Superpowers Development Workflow

OpenSpec handles WHAT (proposals, specs, lifecycle, archive).
Superpowers handles HOW (brainstorming, design, planning, execution).

Comet installs both and provides schema bundles for `openspec/schemas/`.

## Setup Check

Verify the environment:

```bash
openspec --version          # must be installed
openspec schemas            # should list superpowers-bridge
```

If schemas are missing, run `comet init` in the project root.

## Workflow

```
openspec new --schema superpowers-bridge <name>   # start a new change
openspec status --change "<name>"                  # check progress
openspec archive --change "<name>" -y              # archive when done
```

The `superpowers-bridge` schema defines 8 artifacts:

```
brainstorm → proposal → specs → tasks → plan → verify → retrospective
                 ↘ design ↗
```

After plan completion, the apply phase executes: workspace setup → subagent-driven-development → verify → retrospective → archive → PR.

## Active Changes

To see active changes, run:

```bash
openspec list --json
```

Or use `comet status` for a formatted view.

## Quick Reference

| Command | Purpose |
|---------|---------|
| `comet init` | Install Comet + OpenSpec + Superpowers + schemas |
| `comet status` | Show active changes |
| `comet update` | Update Comet packages and schemas |
| `comet doctor` | Diagnose installation health |
```

Write Chinese version to `assets/skills-zh/comet/SKILL.md`:

```markdown
---
name: comet
description: "Comet — OpenSpec + Superpowers 双星开发工作流。使用 /comet 查看活跃变更和推荐命令。"
---

# Comet — OpenSpec + Superpowers 开发工作流

OpenSpec 负责 WHAT（提案、规格、生命周期、归档）。
Superpowers 负责 HOW（头脑风暴、设计、规划、执行）。

Comet 安装两者并为 `openspec/schemas/` 提供 schema bundle。

## 环境检查

```bash
openspec --version          # 必须已安装
openspec schemas            # 应列出 superpowers-bridge
```

如果 schema 缺失，在项目根目录运行 `comet init`。

## 工作流

```
openspec new --schema superpowers-bridge <name>   # 新建变更
openspec status --change "<name>"                  # 查看进度
openspec archive --change "<name>" -y              # 完成后归档
```

`superpowers-bridge` schema 定义了 8 个 artifact：

```
brainstorm → proposal → specs → tasks → plan → verify → retrospective
                 ↘ design ↗
```

Plan 完成后进入 apply 阶段：工作空间设置 → subagent-driven-development → verify → retrospective → archive → PR。

## 查看活跃变更

```bash
openspec list --json
```

或使用 `comet status` 查看格式化输出。

## 快速参考

| 命令 | 用途 |
|------|------|
| `comet init` | 安装 Comet + OpenSpec + Superpowers + schemas |
| `comet status` | 查看活跃变更 |
| `comet update` | 更新 Comet 和 schemas |
| `comet doctor` | 诊断安装健康状态 |
```

- [ ] **Step 2: Commit**

```bash
git add assets/skills/comet/SKILL.md assets/skills-zh/comet/SKILL.md
git commit -m "refactor(comet): simplify entry skill to thin guide"
```

---

### Task 9: Update tests

**Files:**
- Delete: `test/ts/comet-scripts.test.ts`
- Modify: `test/ts/status.test.ts`
- Modify: `test/ts/doctor.test.ts`
- Modify: `test/ts/init.test.ts`

- [ ] **Step 1: Remove shell script tests**

```bash
rm test/ts/comet-scripts.test.ts
```

- [ ] **Step 2: Update status.test.ts**

The status test needs updating since `getActiveChanges` now calls `openspec list --json` instead of reading `.comet.yaml`.

Read the current test file and adjust:

```bash
npx vitest run test/ts/status.test.ts
```

Fix any failing tests. Expected changes:
- Mock `execSync` for `openspec list --json` instead of creating `.comet.yaml` files

- [ ] **Step 3: Update doctor.test.ts**

Remove tests that check for `.comet.yaml` validation. Add tests for schema presence checks.

```bash
npx vitest run test/ts/doctor.test.ts
```

Fix any failing tests.

- [ ] **Step 4: Update init.test.ts**

Add test for schema installation step.

```bash
npx vitest run test/ts/init.test.ts
```

- [ ] **Step 5: Run full test suite**

```bash
npx vitest run
```

Expected: all tests PASS

- [ ] **Step 6: Commit**

```bash
git add test/
git commit -m "test: remove shell script tests, update for schema-driven refactor"
```

---

### Task 10: Update docs and changelog

**Files:**
- Modify: `README.md`
- Modify: `README-zh.md`
- Modify: `CHANGELOG.md`

- [ ] **Step 1: Update README.md**

Replace the workflow section with schema-driven version. Remove references to phase commands, guard scripts, .comet.yaml. Add schema installation description.

Key changes to `README.md`:
- Remove "Five Phases" table and replace with schema artifact flow
- Remove "State Management" section (.comet.yaml)
- Remove "Reliability Features" section (guard/state scripts)
- Remove "Phase Skills" references, keep only `/comet` entry
- Update "Commands" section to reflect simplified workflow
- Update "Project Structure" diagram

- [ ] **Step 2: Update README-zh.md**

Same changes as above, in Chinese.

- [ ] **Step 3: Update CHANGELOG.md**

Prepend new version entry:

```markdown
## What's Changed [0.4.0] - 2026-05-27

### Changed

- **Schema-Driven Architecture**: Refactor Comet from a hardcoded 5-phase state machine into a pure scaffolding tool. Workflow execution is now delegated to OpenSpec's native `--schema` mechanism via schema bundles installed into `openspec/schemas/`.
- **Simplified Workflow**: Removed 7 phase skills (comet-open through comet-archive), 2 preset skills (hotfix/tweak), 5 shell scripts (guard/state/handoff/archive/yaml-validate), and `.comet.yaml` state machine logic.
- **Schema Installation**: `comet init` now installs schema bundles from comet-schema repository into `openspec/schemas/`. Supports interactive selection and non-interactive `--yes` mode.
- **Schema Update**: `comet update` checks remote schema versions and updates local copies.

### Removed

- **Phase Skills**: comet-open, comet-design, comet-build, comet-verify, comet-archive, comet-hotfix, comet-tweak
- **Shell Scripts**: comet-guard.sh, comet-state.sh, comet-handoff.sh, comet-archive.sh, comet-yaml-validate.sh
- **State Machine**: .comet.yaml is no longer created or managed by Comet. State is tracked by OpenSpec natively.

### Tests

- Removed shell script test suite. Updated status, doctor, and init tests to reflect schema-driven architecture.
```

- [ ] **Step 4: Bump version in package.json**

```json
"version": "0.4.0"
```

- [ ] **Step 5: Commit**

```bash
git add README.md README-zh.md CHANGELOG.md package.json
git commit -m "docs: update for schema-driven architecture, bump to 0.4.0"
```

---

### Task 11: Final verification

- [ ] **Step 1: Run full test suite**

```bash
npx vitest run
```

Expected: all tests PASS

- [ ] **Step 2: Run TypeScript type check**

```bash
npx tsc --noEmit
```

Expected: no errors

- [ ] **Step 3: Verify build**

```bash
node build.js
```

Expected: build succeeds

- [ ] **Step 4: Run linter**

```bash
npx eslint src/
```

Expected: no errors

- [ ] **Step 5: Push to comet-schema repo**

```bash
cd /Users/solariswu/workspaces/github/comet-schema
git push origin main
```
