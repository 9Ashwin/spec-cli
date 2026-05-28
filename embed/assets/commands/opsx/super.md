name: "OPSX: Super"
description: Start a new change with Superpowers workflow
category: Workflow
tags: [superpowers, openspec, workflow]
---

Start a new change using the Superpowers workflow schema.

Instruct the user to describe what they want to build, then convert the description to kebab-case and run:

```bash
openspec new change --schema superpowers-bridge "<kebab-case-name>"
```

Then follow the artifact instructions from the schema (brainstorm → design → specs → tasks → plan → apply → verify → archive).
