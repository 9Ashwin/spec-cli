---
name: opsx-super
description: "Use when the user wants to start, route, or continue OpenSpec + Superpowers work through /opsx:super."
---

# opsx:super

Use this as the front door for OpenSpec + Superpowers work. Its job is routing: decide whether this needs a schema change, avoid duplicate active changes, then hand execution to OpenSpec's `superpowers-bridge` schema.

<EXTREMELY-IMPORTANT>
This is a GATED workflow. Each phase (brainstorm, proposal, design, specs, tasks, plan, apply, verify) has an exit gate defined in the schema artifact instructions. You MUST pause at every gate that requires user input — especially brainstorm approval, design sign-off, scope changes, verification failure handling, and branch/PR decisions.

The first hard gate is brainstorm → proposal. The brainstorm EXIT GATE says: "Stop here until the user has approved the proposed design direction." Do NOT write proposal.md or advance further until the user explicitly approves.

Auto-advance through unambiguous completed phases only. When a schema exit gate says to stop, stop and ask.
</EXTREMELY-IMPORTANT>

## Instruction Priority

User instructions always take precedence. If the user explicitly says not to create a schema change, do not override them; explain the tradeoff and follow their requested path.

When routing work, follow this order:

1. User's explicit instruction
2. Project instructions such as AGENTS.md / CLAUDE.md
3. OpenSpec schema artifact instructions
4. This entry skill
5. Default agent behavior

If these conflict, stop and explain the conflict instead of silently choosing a path.

## The Rule

Before creating a change, always do both:

1. Classify the request as schema change or direct PR.
2. Inspect active changes with:

```bash
openspec list --json
```

Do not silently create duplicate changes.

If `openspec` is missing, the project is not initialized, or `superpowers-bridge` is not installed, STOP and tell the user to run `spec-cli init` or `spec-cli doctor`.

## Route The Request

Use `superpowers-bridge` for:

- New features or capabilities
- Architecture changes
- Breaking changes
- External contract, schema, data model, or cross-system changes

Do not open a schema change for:

- Bug fixes that restore intended behavior without changing contracts
- Typos or docs-only edits
- Test backfills
- Linter/config value tweaks
- Non-breaking dependency updates

For direct PR cases, tell the user this does not need an `opsx:super` change and proceed normally.

## Active Change Handling

- No active changes: create a new change when the request needs one.
- One active change: ask whether to continue it or create a new one, unless the user explicitly asked for a new change.
- Multiple active changes: list them and ask which one to continue, or whether to create a new change.
- User asks to continue: follow OpenSpec status/instructions for the active change instead of creating a new directory.

## Continuous Execution

After creating or selecting a change, enter the Continuous Execution loop below. Auto-advance through phases only when exit gates are satisfied — the first gate (brainstorm approval) is the most critical. Do not skip it.

For every `opsx:super` invocation that enters `superpowers-bridge`, inspect OpenSpec status first:

```bash
openspec status --change "<name>" --json
```

Then read the current schema artifact instructions and continue from the next incomplete schema step. Before advancing to the next phase, verify the current phase's EXIT GATE (defined in the schema artifact instruction) is satisfied. If the gate requires user input, STOP and ask. Advance through unambiguous completed steps automatically:

brainstorm -> proposal -> design -> specs -> tasks -> plan -> apply action -> verify -> retrospective/archive.

For each schema step, the schema instructions decide which Superpowers skill applies. Invoke the relevant Superpowers skill before acting on that schema step. Do not hand-write artifacts from memory when a schema step requires a skill.

Skill output path overrides (schema artifact instructions take precedence over skill defaults):
- **brainstorming**: The skill's default output path is `docs/superpowers/specs/` and its default terminal state is invoking `writing-plans`. When running under superpowers-bridge, override both: write the raw brainstorming output to the change's `brainstorm.md` (per the schema artifact instruction), and after the user approves the design, proceed to the **proposal** artifact — do NOT invoke writing-plans.
- **writing-plans**: The skill's default output path is `docs/superpowers/plans/`. When running under superpowers-bridge, write the plan to the change's `plan.md` instead.

For artifact steps, get the concrete instructions from OpenSpec:

```bash
openspec instructions <artifact-id> --change "<name>" --json
```

Use the returned `instruction`, `template`, `outputPath`, and `dependencies`. Write the artifact to `outputPath`, then re-run status before advancing.

For the apply action, do not look for an `apply` artifact in status. Get the action instructions from OpenSpec:

```bash
openspec instructions apply --change "<name>" --json
```

Use the returned `contextFiles`, `tasks`, `progress`, and `instruction` to run the implementation phase.

Do not delegate Continuous Execution to `openspec-continue-change` or `/opsx:continue`; that skill intentionally stops after one artifact. `opsx:super` owns the status -> instructions -> artifact/action -> status loop until a stop condition below.

Treat this as a gated workflow, not an unconditional run-to-the-end prompt:

- Auto-advance only when the current artifact/action has satisfied its schema instruction and exit gate.
- On resume, trust OpenSpec status/instructions and files on disk, not conversation history.
- When a schema instruction says to stop for a user decision, stop and ask. Do not choose defaults for design approval, scope expansion, verification failure handling, or branch/PR handling.

Stop only when one of these happens:

- The schema workflow is complete.
- The next schema instruction requires an explicit user decision.
- OpenSpec status/instructions are unavailable or conflict with project/user instructions.
- A required schema step, Superpowers skill, command, or artifact is missing.
- Verification fails or the schema instruction says to stop.

Do not use conversation history as the source of truth for progress. Re-read OpenSpec status/instructions on resume, after tool failures, and before advancing to the next phase.

## Create A New Change

Convert the user's request to a short kebab-case name, then run:

```bash
openspec new change "<kebab-case-name>" --schema superpowers-bridge --description "<original request>"
```

After creation, enter Continuous Execution for the new change.

## Respect Bridge Routing

Do not write to `docs/superpowers/specs/`.
Do not write to `docs/superpowers/plans/`.
Schema artifacts belong under `openspec/changes/<name>/`.

Do not silently fall back to raw `brainstorming`, `writing-plans`, or hand-written artifacts if a required schema step or Superpowers skill is missing. STOP and tell the user what is missing.

## Red Flags

These thoughts mean STOP. Re-read OpenSpec status and schema instructions before continuing:

- "This is a feature, but it is small enough to skip the schema."
- "I'll run brainstorming directly and move the file later."
- "There is already an active change, but a new one is faster."
- "OpenSpec failed, so I'll create the folders by hand."
- "The schema step is unclear, so I'll improvise the next artifact."
