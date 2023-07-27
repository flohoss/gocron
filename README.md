# GoBackup

[![pipeline status](https://gitlab.unjx.de/flohoss/gobackup/badges/main/pipeline.svg)](https://gitlab.unjx.de/flohoss/gobackup/-/commits/main)

## A docker-compose example

```yml
services:
  gobackup:
    image: git.unjx.de/flohoss/gobackup:latest
    container_name: gobackup
    hostname: gobackup
    restart: always
    environment:
      - TZ=Europe/Berlin
      - BACKUP_CRON=0 2 * * *
      - CLEANUP_CRON=0 3 * * 0
      - HEALTHCHECK_URL=https://health.example.de/ping/
      - HEALTHCHECK_UUID=xxx-xxx-xxx-xxx-xxx
      - NOTIFICATION_URL=telegram://xxx:xxx@telegram?chats=xxx
    ports:
      - '127.0.0.1:8080:8080'
    volumes:
      - ./gobackup/:/app/storage/
      - /var/run/docker.sock:/var/run/docker.sock
      - /opt/docker/:/opt/docker/
      - ./.resticpwd:/secrets/.resticpwd:ro
      - ./.rclone.conf:/root/.config/rclone/rclone.conf:ro
```
