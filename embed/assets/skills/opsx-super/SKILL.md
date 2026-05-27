---
name: opsx:super
description: "Use when the user wants to start, route, or continue OpenSpec + Superpowers work through /opsx:super."
---

# opsx:super

Use this as the front door for OpenSpec + Superpowers work. Its job is routing: decide whether this needs a schema change, avoid duplicate active changes, then hand execution to OpenSpec's `superpowers-bridge` schema.

## Instruction Priority

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

## Create A New Change

Convert the user's request to a short kebab-case name, then run:

```bash
openspec new change "<kebab-case-name>" --schema superpowers-bridge --description "<original request>"
```

After creation, follow the schema artifact instructions:

brainstorm -> proposal -> design -> specs -> tasks -> plan -> apply -> verify -> retrospective/archive.

## Respect Bridge Routing

Do not write to `docs/superpowers/specs/`.
Do not write to `docs/superpowers/plans/`.
Schema artifacts belong under `openspec/changes/<name>/`.

Do not silently fall back to raw `brainstorming`, `writing-plans`, or hand-written artifacts if a required schema step or Superpowers skill is missing. STOP and tell the user what is missing.

## Red Flags

Stop and re-route when you notice:

- "This is a feature, but it is small enough to skip the schema."
- "I'll run brainstorming directly and move the file later."
- "There is already an active change, but a new one is faster."
- "OpenSpec failed, so I'll create the folders by hand."
- "The schema step is unclear, so I'll improvise the next artifact."
