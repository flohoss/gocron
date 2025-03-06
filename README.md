# GoCron

<div align="center">

<img src="web/public/static/logo.webp" height="250px">

[![goreleaser](https://github.com/flohoss/gocron/actions/workflows/release.yaml/badge.svg?branch=main)](https://github.com/flohoss/gocron/actions/workflows/release.yaml)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/flohoss/gocron)

A task scheduler built with Go and Vue.js that allows users to specify recurring jobs via a simple YAML configuration file. The scheduler reads job definitions, executes commands at specified times using cron expressions, and passes in environment variables for each job.

</div>

# Table of Contents

- [GoCron](#gocron)
- [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [How It Works](#how-it-works)
  - [Docker](#docker)
    - [run command](#run-command)
    - [compose file](#compose-file)
  - [Screenshots](#screenshots)
    - [Home](#home)
    - [Job](#job)
    - [Installed software](#installed-software)
    - [OpenAPI Specification (/api/docs)](#openapi-specification-apidocs)
  - [Example Configuration](#example-configuration)
- [Preinstalled Software](#preinstalled-software)
  - [License](#license)
  - [Development setup](#development-setup)
    - [Automatic rebuild and reload](#automatic-rebuild-and-reload)
    - [Rebuild types while docker is running](#rebuild-types-while-docker-is-running)

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
  -e TZ=Europe/Berlin \
  -v ./config/:/app/config/ \
  # -e SEND_ON_SUCCESS=true \
  # Log level can be one of: debug info warn error
  # -e LOG_LEVEL=info \
  # -e PORT=8156 \
  # Uncomment and replace <token> if using ntfy notifications
  # -e NTFY_URL=https://ntfy.hoss.it/ \
  # -e NTFY_TOPIC=Backup \
  # -e NTFY_TOKEN=<token> \
  # Uncomment if using Restic with a password file
  # -v ./.resticpwd:/secrets/.resticpwd \
  # Uncomment if using a preconfigured rclone config
  # -v ./.rclone.conf:/root/.config/rclone/rclone.conf \
  # Uncomment to allow running Docker commands inside the container
  # -v /var/run/docker.sock:/var/run/docker.sock \
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
    environment:
      - TZ=Europe/Berlin
      # - SEND_ON_SUCCESS=true
      # Log level can be one of: debug info warn error
      # - LOG_LEVEL=info
      # - PORT=8156
      # Uncomment and replace <token> if using ntfy notifications
      # - NTFY_URL=https://ntfy.hoss.it/
      # - NTFY_TOPIC=Backup
      # - NTFY_TOKEN=<token>
    volumes:
      - ./config/:/app/config/
      # Uncomment if using Restic with a password file
      # - ./.resticpwd:/secrets/.resticpwd
      # Uncomment if using a preconfigured rclone config
      # - ./.rclone.conf:/root/.config/rclone/rclone.conf
      # Uncomment to allow running Docker commands inside the container
      # - /var/run/docker.sock:/var/run/docker.sock
    port:
      - '8156:8156'
```

## Screenshots

### Home

<img src="img/home.webp" width="500px">

<img src="img/home_light.webp" width="500px">

### Job

<img src="img/job.webp" width="500px">

<img src="img/job_light.webp" width="500px">

### Installed software

<img src="img/software.webp" width="500px">

<img src="img/software_light.webp" width="500px">

### OpenAPI Specification (/api/docs)

<img src="img/api.webp" width="500px">

<img src="img/api_light.webp" width="500px">

## Example Configuration

The following is an example of a valid YAML configuration creating backups with restic every 3 am in the morning and cleaning the repo every Sunday at 5 am:

```yml
defaults:
  # every job will be appended to this cron and the jobs will run sequentially
  cron: '0 3 * * *'
  # global envs to use in all jobs
  envs:
    - key: RESTIC_PASSWORD_FILE
      value: '/secrets/.resticpwd'
    - key: BASE_REPOSITORY
      value: 'rclone:pcloud:Server/Backups'
    - key: APPDATA_PATH
      value: '/mnt/user/appdata'

jobs:
  - name: Cleanup
    # override the default cron
    cron: '0 5 * * 0'
    # envs just for the job, overwriting exiting defaults
    envs:
      - key: RESTIC_POLICY
        value: '--keep-daily 7 --keep-weekly 5 --keep-monthly 12 --keep-yearly 75'
      - key: RESTIC_POLICY_SHORT
        value: '--keep-last 7'
    commands:
      - command: restic -r ${BASE_REPOSITORY}/forgejo forget ${RESTIC_POLICY} --prune
      - command: restic -r ${BASE_REPOSITORY}/paperless forget ${RESTIC_POLICY} --prune
  - name: Forgejo
    envs:
      - key: RESTIC_REPOSITORY
        value: ${BASE_REPOSITORY}/forgejo
    commands:
      - command: docker exec -e PASSWORD=password forgejo-db pg_dump db --username=user
        file_output: ${APPDATA_PATH}/forgejo/.dbBackup.sql
      - command: restic backup ${APPDATA_PATH}/forgejo
  - name: Paperless
    envs:
      - key: RESTIC_REPOSITORY
        value: ${BASE_REPOSITORY}/paperless
    commands:
      - command: docker exec paperless document_exporter ${APPDATA_PATH}/paperless/export
        file_output: ${APPDATA_PATH}/paperless/.export.log
      - command: restic backup ${APPDATA_PATH}/paperless
  - name: Show files
    # only a command per job is required
    commands:
      - command: ls -la
```

# Preinstalled Software

These tools are preinstalled and ready to be used for various operations within your jobs:

1. [BorgBackup](https://www.borgbackup.org/)

> BorgBackup is a fast, secure, and space-efficient backup tool. It deduplicates data and can be used for both local and remote backups. It is widely known for its encryption and compression capabilities, which ensures that backups are both secure and compact.

2. [Restic](https://restic.net/)

> Restic is a fast and secure backup program that supports various backends, including local storage and cloud providers. Restic is optimized for simplicity and speed, offering encrypted backups with minimal configuration. It integrates perfectly with the task scheduler for managing secure backups.

3. [rclone](https://rclone.org/)

> rclone is a command-line program used to manage and transfer files to and from various cloud storage services. It supports numerous cloud providers, including Google Drive, Dropbox, and Amazon S3, making it an excellent choice for managing backups on remote storage solutions. With rclone, you can efficiently sync, move, and manage your data in the cloud.

<picture>
  <source
    media="(prefers-color-scheme: dark)"
    srcset="
      https://api.star-history.com/svg?repos=flohoss/gocron&type=Date&theme=dark
    "
  />
  <source
    media="(prefers-color-scheme: light)"
    srcset="
      https://api.star-history.com/svg?repos=flohoss/gocron&type=Date
    "
  />
  <img
    alt="Star History Chart"
    src="https://api.star-history.com/svg?repos=flohoss/gocron&type=Date"
  />
</picture>

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/flohoss/gocron/blob/main/LICENSE) file for details.

## Development setup

### Automatic rebuild and reload

`docker compose up`

### Rebuild types while docker is running

`docker compose run --rm types`
