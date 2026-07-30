<div align="center">

<img src="web/public/static/logo.webp" height="250px">

[![CI](https://github.com/flohoss/gocron/actions/workflows/ci.yaml/badge.svg)](https://github.com/flohoss/gocron/actions/workflows/ci.yaml)
[![Coverage](https://raw.githubusercontent.com/wiki/flohoss/gocron/coverage.svg)](https://github.com/flohoss/gocron/wiki/coverage.html)

[![RepoRanker](https://reporanker.com/badge/flohoss/gocron)](https://reporanker.com/repos/flohoss/gocron)

A task scheduler built with Go and Vue.js that allows users to specify recurring jobs via a simple YAML configuration file. The scheduler reads job definitions, executes commands at specified times using cron expressions, and passes in environment variables for each job.

Tagged GitHub releases include downloadable Linux binaries. Run `./gocron_<version>_linux_<arch> --version` to inspect the embedded version metadata of a downloaded release binary.

The server supports `--config /path/to/config.yaml` for a non-default configuration file.

</div>

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Features](#features)
- [How It Works](#how-it-works)
- [Docker](#docker)
  - [run command](#run-command)
  - [compose file](#compose-file)
- [Screenshots](#screenshots)
  - [Dark mode](#dark-mode)
  - [Light mode](#light-mode)
  - [API Docs](#api-docs)
- [Configuration File](#configuration-file)
  - [YAML Configuration](#yaml-configuration)
  - [Software](#software)
- [Safety & Limitations](#safety--limitations)
- [Star History](#star-history)
- [License](#license)
- [Development setup](#development-setup)
  - [Automatic rebuild and reload](#automatic-rebuild-and-reload)

## Features

- Simple Configuration: Easily define jobs, cron schedules, and environment variables in a YAML config file.
- Cron Scheduling: Supports cron expressions for precise scheduling.
- Environment Variables: Define environment variables specific to each job.
- Easy Job Management: Add and remove jobs quickly with simple configuration.
- Pre-installed backup-software for an easy backup solution

## How It Works

- Defaults Section: This section defines default values that are applied to all jobs. You can specify a default cron expression and environment variables to be inherited by each job.
- Jobs Section: Here, you define multiple jobs. Each job can have its own cron expression, environment variables, and commands to execute.
- Environment Variables: Define environment variables for each job to customize its runtime environment.
- Commands: Each job can have multiple commands, which will be executed in sequence.

## Docker

### run command

```sh
docker run -it --rm \
  --name gocron \
  --hostname gocron \
  -p 8156:8156 \
  -v ./config/:/app/config/ \
  ghcr.io/flohoss/gocron:latest
```

### compose file

```yml
services:
  gocron:
    image: ghcr.io/flohoss/gocron:latest
    restart: always
    container_name: gocron
    hostname: gocron
    volumes:
      - ./config/:/app/config/
    ports:
      - "8156:8156"
```

By default, gocron reads `./config/config.yaml`. You can optionally override the config file path with `--config /path/to/config.yaml`. SQLite data is stored in the same folder by default, and can be overridden with `db.location` inside the config file. Relative `db.location` values are resolved from the config file directory. You can also set the SQLite file name with `db.name` (default: `db.sqlite`).

### Environment overrides (`GC_`)

Config values can be overridden via environment variables using the `GC_` prefix.

- `GC_LOG_LEVEL=debug` overrides `log_level`
- `GC_SERVER_PORT=9000` overrides `server.port`
- `GC_HEALTHCHECK_TYPE=GET` overrides `healthcheck.type`

Rule: dots become underscores and keys are uppercased (`server.address` -> `GC_SERVER_ADDRESS`).

## Screenshots

### Dark mode

<p align="center">
  <img src="screenshots/jobs-dark.webp" width="500" />
  <img src="screenshots/job-dark.webp" width="500" />
  <img src="screenshots/terminal-dark.webp" width="500" />
  <img src="screenshots/filter-dark.webp" width="500" />
</p>

### Light mode

<p align="center">
  <img src="screenshots/jobs-light.webp" width="500" />
  <img src="screenshots/job-light.webp" width="500" />
  <img src="screenshots/terminal-light.webp" width="500" />
  <img src="screenshots/filter-light.webp" width="500" />
</p>

### API Docs

<p align="center">
  <img src="screenshots/api-dark.webp" width="500" />
  <img src="screenshots/api-light.webp" width="500" />
</p>

## Configuration File

### YAML Configuration

The entire configuration is managed via the YAML file, including settings for the timezone, logging, server, and jobs.

For a complete and working configuration example, please refer to the [`config/config.yaml`](config/config.yaml) file in the repository. A minimal excerpt covering the most important fields:

```yaml
time_zone: 'UTC'                   # Sets the TZ environment variable for the process
log_level: 'info'                  # debug | info | warn | error | off
delete_runs_after_days: 7         # Delete run history after N days (0 disables)
db:
  location: '.'                    # Absolute or relative to the config file
  name: 'db.sqlite'

job_defaults:
  cron: '0 3 * * 0'               # Inherited by jobs without their own cron
  timeout: '30s'                  # Optional: inherited by jobs without their own timeout
  retries: 2                      # Optional: inherited by jobs without their own retries
  envs:
    - key: SLEEP_TIME
      value: '5'
  pre_commands:
    - echo "Starting backup..."
  post_commands:
    - echo "Backup finished!"

jobs:
  - name: 'Nightly Backup'
    cron: '0 5 * * 0'             # Overrides job_defaults.cron
    timeout: '30s'                # Abort a single command after this duration
    retries: 2                    # Retry each failing command up to N times
    disable_fail_fast: false      # Stop on first failing command (default)
    envs:
      - key: RESTIC_REPOSITORY
        value: '/backups/nightly'
    commands:
      - restic backup /data
      - restic forget --keep-daily 7
```

See the [Safety & Limitations](#safety--limitations) section below for guidance on secrets, concurrency, and failure semantics.

### Software

You can specify the software you want to install and the version you want to use directly in the configuration file.
Available software packages include: apprise, borgbackup, docker, git, podman, rclone, rdiff-backup, restic, rsync, logrotate, sqlite3, and kopia.

The version format depends on the installation method:

- **apprise** (via pipx): version format like `1.2.0`
- **docker** (via apt): version format like `5:24.0.5-1~debian.11~bullseye`
- **Others** (via apt): standard apt version format

Here is an example of how to set up specific software versions:

```yaml
software:
  - name: "apprise"
    version: "1.2.0"
  - name: "borgbackup"
    version: "1.2.0"
  - name: "docker"
    version: "5:24.0.5-1~debian.11~bullseye"
  - name: "git"
  - name: "podman"
  - name: "rclone"
  - name: "rdiff-backup"
  - name: "restic"
    version: "0.14.0"
  - name: "rsync"
  - name: "logrotate"
  - name: "sqlite3"
  - name: "kopia"
```

## Safety & Limitations

Because GoCron executes shell commands on a schedule, operators should understand the trust boundaries, failure semantics, and security considerations before entrusting it with critical tasks.

### Trust boundaries and permissions

- Commands run inside the container as the process user (root by default in the published image). Use the least-privileged user for jobs that touch sensitive data.
- The web UI terminal is gated by an allow-list (`terminal.allowed_commands` in the config). Do **not** set `allow_all_commands: true` in production.
- Working directory is the container's `/app`. Mount only the directories a job needs.

### Secrets

Do not store passwords, API tokens, or repository credentials in plaintext inside `config.yaml` or `envs` blocks. Prefer one of:

- **Docker secrets / Compose secrets** — mount a secrets file and reference it from the command, e.g. `restic -r $(cat /run/secrets/repo) backup /data`.
- **Environment files** — pass sensitive values via `docker run --env-file` or `environment:` in Compose, and reference them with `${VAR}` expansion in commands.

### Concurrency and overlapping runs

- GoCron is **single-flight**: while any job is running, no other job (or triggered run) will start. Scheduled runs that arrive while a job is active are **skipped**, not queued.
- There is no per-job concurrency limit. If you need parallel jobs, run multiple GoCron instances with separate config files.

### Missed schedules

- GoCron uses standard cron semantics. If the container is down when a schedule fires, that run is **not** caught up on restart. There is no `@reboot` or missed-run replay.

### Time zones and DST

- Set `time_zone` in the config (it sets the `TZ` environment variable). All schedules run in that timezone.
- During a DST "spring forward" gap, a cron expression targeting the skipped hour will not fire. During "fall back", a target in the repeated hour may fire once or twice depending on the cron library. Test schedules around DST transitions.

### Command timeouts

- Set a per-job `timeout` (e.g. `30s`, `5m`) to abort any single command that runs longer than the given duration. A timed-out command counts as a failure.
- A `timeout` set in `job_defaults` is inherited by jobs that don't set their own. If neither is set, commands run without a deadline.

### Retries

- Set `retries` on a job to retry each failing command up to N additional times (N+1 total attempts). Retries apply per-command, not per-job.
- A `retries` value in `job_defaults` is inherited by jobs that don't set their own. If neither is set, commands are not retried.
- Retries are immediate (no backoff). Combine with `timeout` to bound the worst-case runtime.

### Fail-fast behavior

- `disable_fail_fast: false` (default) stops the job on the first failing command.
- `disable_fail_fast: true` continues running remaining commands even if one fails. The run is still marked as failed.

### Log retention

- `delete_runs_after_days` controls how long run history (logs, exit status) is kept in SQLite. Set to `0` to retain forever. A daily cleanup job runs at midnight to prune old records.

### Failure notifications

- Configure the `healthcheck` section to send HTTP callbacks at three phases: `start`, `end`, and `failure`. Each phase supports a URL, query params, and a JSON body.
- Use services like [Healthchecks](https://healthchecks.io) or [Uptime Kuma](https://github.com/louislam/uptime-kuma) to receive alerts when a job fails or doesn't report on schedule.
- Alternatively, install `apprise` via the `software` list to push notifications from within a command.

### Process termination

- On container shutdown, GoCron attempts a graceful shutdown. A running command receives a SIGKILL via Go's `os/exec` context cancellation. Long-running commands may not clean up gracefully — test shutdown behavior for your workloads.

### Supply-chain considerations

- Pre-installing backup tools (`restic`, `borgbackup`, `docker`, `podman`, etc.) increases the image's attack surface. Only list the software you actually need in the `software` section.
- Pin software versions explicitly (see the example above) to avoid surprise upgrades on image rebuild.

## Star History

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=flohoss/gocron&type=Date&theme=dark" />
  <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=flohoss/gocron&type=Date" />
  <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=flohoss/gocron&type=Date" />
</picture>

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/flohoss/gocron/blob/main/LICENSE) file for details.

## Development setup

### Run tests

```bash
# Run Go tests
docker compose run --rm go test ./...

# Install e2e dependencies
docker compose run --rm npm-e2e install --no-audit --no-fund

# Run e2e tests in Docker (uses a separate e2e config via compose override)
docker compose -f compose.yml -f compose.e2e.yml --profile test run --rm e2e
docker compose -f compose.yml -f compose.e2e.yml --profile test down
```

### Update Dependencies

```bash
# Node packages
docker compose run --rm npm install
docker compose run --rm --entrypoint npx npm npm-check-updates -u && docker compose run --rm npm install

# E2E packages
docker compose run --rm npm-e2e install
docker compose run --rm --entrypoint npx npm-e2e npm-check-updates -u && docker compose run --rm npm install

# Go packages
docker compose run --rm go get -u ./...
docker compose run --rm go mod tidy
```

### Automatic rebuild and reload

```sh
docker compose up
```
