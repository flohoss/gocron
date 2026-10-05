# Agent Guidance

Read before making changes. Rule-oriented and self-contained.

## Principles

- **CLEAN code.** Small functions, single responsibility, descriptive names, no dead code, no overengineering.
- **No comments.** Use descriptive function or service names instead.
- **Test comments are allowed** when the tested behaviour is not obvious from the test name: a short block above the test explaining what it asserts and why it exists.
- **No code markers** like `// ... existing code ...` in edits.
- Go imports: stdlib, then external, then internal (`github.com/flohoss/gocron/...`), each block alphabetical.
- Never edit generated files (`services/jobs/*.sql.go`, `models.go`, `web/src/client/`).

## Tooling — always via Docker Compose, never on the host

- **Code generation:** `docker compose run --rm backend sqlc generate`
- **Backend format:** `docker compose run --rm go fmt ./...`
- **Format (all non-Go files):** `docker compose run --rm format`
- **Tests:** `docker compose run --rm go test ./...`
- **Coverage:** `docker compose run --rm --entrypoint sh go scripts/coverage.sh` — excludes generated files (detected by the `// Code generated ... DO NOT EDIT.` marker)
- **Coverage badge:** `docker compose run --rm --entrypoint sh go scripts/coverage-badge.sh` — writes `assets/coverage.svg`, which CI commits to `main` when the value changes
- **E2E tests (Headless):**
  ```sh
  docker compose -f compose.yml -f compose.e2e.yml --profile test run --rm e2e
  # teardown when done:
  docker compose -f compose.yml -f compose.e2e.yml --profile dev --profile test down
  ```
- **E2E interactive (Cypress GUI):**
  ```sh
  docker compose -f compose.yml -f compose.e2e.yml --profile dev --profile test up -d backend frontend
  npm --prefix web/e2e run open
  # teardown when done:
  docker compose -f compose.yml -f compose.e2e.yml --profile dev --profile test down
  ```
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

### Compose pinned versions

`compose.yml` pins tool versions (`V_GOLANG`, `V_AIR`, `V_SQLC`, `V_NODE`, `V_CYPRESS`, `V_DEBIAN`). Check for new versions when updating deps:

- `V_GOLANG` — `go.dev/dl/?mode=json`
- `V_AIR` — `go list -m -versions github.com/air-verse/air`
- `V_SQLC` — `go list -m -versions github.com/sqlc-dev/sqlc`
- `V_NODE` — `nodejs.org/en/about/previous-releases`
- `V_CYPRESS` — `hub.docker.com/r/cypress/included/tags`
- `V_DEBIAN` — Debian stable codename

Rebuild images after bumping: `docker compose build backend go format`

### TypeScript major

`vue-tsc` breaks on TS majors it doesn't support yet (`ERR_PACKAGE_PATH_NOT_EXPORTED` for `./lib/tsc`). After `npm-check-updates -u`, if `typescript` was bumped to a new major, check if `vue-tsc` supports it (`docker compose run --rm npm run build`). If not, revert `typescript` in `web/package.json` to the previous major before installing.

### mapstructure time-zone decoding

Never reintroduce `mapstructure.StringToTimeLocationHookFunc` for `time_zone`. It leaks the raw `time.LoadLocation` error, which is platform-dependent: a sandboxed darwin build returns `operation not permitted` (EACCES) for an unknown zone instead of `unknown time zone`, while Linux hits ENOENT and falls through to the embedded `time/tzdata`. Keep the project-owned `timeLocationDecodeHookFunc` in `config/config.go` — it guarantees the error message on all platforms. Check this when bumping `github.com/go-viper/mapstructure` or `github.com/spf13/viper`.

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
