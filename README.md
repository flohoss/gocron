# GoCron

<img src="https://github.com/flohoss/gocron/blob/main/web/public/static/logo.png?raw=true" width="250px" height="340px">

A task scheduler built with Go and Vue.js that allows users to specify recurring jobs via a simple YAML configuration file. The scheduler reads job definitions, executes commands at specified times using cron expressions, and passes in environment variables for each job.

## Features

- Simple Configuration: Easily define jobs, cron schedules, and environment variables in a YAML config file.
- Cron Scheduling: Supports cron expressions for precise scheduling.
- Environment Variables: Define environment variables specific to each job.
- Easy Job Management: Add and remove jobs quickly with simple configuration.

# Example Configuration

The following is an example of a valid YAML configuration:

```yml
defaults:
  cron: '0 3 * * *'
  envs:
    - key: RESTIC_PASSWORD_FILE
      value: '/secrets/.resticpwd'
    - key: BASE_REPOSITORY
      value: 'rclone:pcloud:Server/Backups'
    - key: APPDATA_PATH
      value: '/mnt/user/appdata'

jobs:
  - name: Example
    cron: '0 5 * * 0'
    envs:
      - key: RESTIC_POLICY
        value: '--keep-daily 7 --keep-weekly 5 --keep-monthly 12 --keep-yearly 75'
    commands:
      - command: ls -la
      - command: sleep 1
      - command: echo "Done!"
      - command: sleep 1
  - name: Another
    cron: '0 5 * * 0'
    commands:
      - command: docker ps
      - command: sleep 2
      - command: echo "Done!"
```

## How It Works

- Defaults Section: This section defines default values that are applied to all jobs. You can specify a default cron expression and environment variables to be inherited by each job.
- Jobs Section: Here, you define multiple jobs. Each job can have its own cron expression, environment variables, and commands to execute.
- Environment Variables: Define environment variables for each job to customize its runtime environment.
- Commands: Each job can have multiple commands, which will be executed in sequence.

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/flohoss/gocron/blob/main/LICENSE) file for details.

## Development setup

### Automatic rebuild and reload

`docker compose up`

### Rebuild types while docker is running

`docker compose run --rm types`
