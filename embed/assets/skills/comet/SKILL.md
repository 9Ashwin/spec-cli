---
name: comet
description: "Comet — OpenSpec + Superpowers development workflow."
---

# Comet — OpenSpec + Superpowers Development Workflow

## Quick Start

```bash
# Start a new change
openspec new --schema superpowers-bridge "your idea"

# Check active changes
openspec list --json

# Archive when done
openspec archive --change "<name>" -y
```

## Workflow

OpenSpec handles the WHAT (proposal, spec lifecycle, archive).
Superpowers handles the HOW (brainstorming, design, plan, build, verify).

The `superpowers-bridge` schema connects them — follow the artifact instructions injected at each step.
