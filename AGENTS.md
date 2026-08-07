# Agent Guidance

Read before making changes. Rule-oriented and self-contained.

## Principles

- **CLEAN code.** Small functions, single responsibility, descriptive names, no dead code, no overengineering.
- **No comments.** Use descriptive function or service names instead.
- **No code markers** like `// ... existing code ...` in edits.
- Go imports: stdlib, then external, then internal (`github.com/flohoss/gocron/...`), each block alphabetical.
- Never edit generated files (`services/jobs/*.sql.go`, `models.go`, `web/src/client/`).

## Tooling — always via Docker Compose, never on the host

- **Code generation:** `docker compose run --rm backend sqlc generate`
- **Backend format:** `docker compose run --rm go fmt ./...`
- **Format (all non-Go files):** `docker compose run --rm format`
- **Tests:** `docker compose run --rm go test ./services/ ./config/ ./internal/...`

Only commit if format and tests pass.

## Git

Commit message — title only, no body, capitalize first letter:

- `[fix]` bug fix
- `[feature]` new functionality
- `[improve]` improvement to existing functionality
- `[refactor]` formatting, renaming, structural-only
- `[meta]` changes outside the codebase (deployment, CI)
- `[docs]` documentation
