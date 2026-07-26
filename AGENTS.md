# Agent Guidance

Read this before making changes. It is rule-oriented and self-contained.

## Code style

- **No comments.** Use descriptive function or service names instead.
- **No code markers** like `// ... existing code ...` in edits.
- Go imports: stdlib, then external, then internal (`github.com/flohoss/mittagskarte/...`), each block alphabetical.

## Git

**Commit message format** — title only, no body:

- `[fix]` — fixes a bug
- `[feature]` — adds new functionality
- `[improve]` — improves existing functionality
- `[meta]` — changes outside the codebase (deployment, CI)
- `[docs]` — documentation
- `[refactor]` — formatting, renaming, structural-only

Capitalize the first letter after the prefix.

## Verification

Before committing, always run:

- **Backend:** `docker compose run --rm go fmt ./...`
- **Frontend:** `docker compose run --rm npm run format`

Only commit if all pass.
