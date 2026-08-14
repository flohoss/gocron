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
- **E2E tests:** `docker compose -f compose.yml -f compose.e2e.yml --profile test run --rm e2e` then `... down`
- **Rebuild & reload:** `docker compose up`

Only commit if format and tests pass.

## Common Commands

```sh
docker compose run --rm npm install
docker compose run --rm --entrypoint npx npm --yes npm-check-updates -u && docker compose run --rm npm install
docker compose run --rm npm-e2e install
docker compose run --rm --entrypoint npx npm-e2e --yes npm-check-updates -u && docker compose run --rm npm-e2e install
docker compose run --rm go get -u ./...
docker compose run --rm go mod tidy
docker compose run --rm go fmt ./...
```

### TypeScript major

`vue-tsc` breaks on TS majors it doesn't support yet (`ERR_PACKAGE_PATH_NOT_EXPORTED` for `./lib/tsc`). After `npm-check-updates -u`, if `typescript` was bumped to a new major, check if `vue-tsc` supports it (`docker compose run --rm npm run build`). If not, revert `typescript` in `web/package.json` to the previous major before installing.

## Git

- Do not commit automatically — wait until explicitly asked.
- One commit per concern — never batch unrelated changes.
- Title only, no body. Capitalize first letter after the prefix:
  - `[fix]` bug fix
  - `[feature]` new functionality
  - `[improve]` improvement to existing functionality
  - `[refactor]` formatting, renaming, structural-only
  - `[meta]` deployment, CI
  - `[docs]` documentation
