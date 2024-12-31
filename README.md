# GoBackup

## Development setup

### automatic rebuild but manual reload

`docker compose up`

### automatic reload

```sh
templ generate --watch --proxy="http://localhost:8156" --cmd="go run cmd/main.go"
docker compose run --rm tailwind
```

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
      - HEALTHCHECK_URL=https://health.example.de/ping/
      - HEALTHCHECK_UUID=xxx-xxx-xxx-xxx-xxx
      - NOTIFICATION_URL=telegram://xxx:xxx@telegram?chats=xxx
    ports:
      - '127.0.0.1:8156:8156'
    volumes:
      - ./gobackup/:/app/config
      - /var/run/docker.sock:/var/run/docker.sock
      - /opt/docker/:/opt/docker/
      - ./.resticpwd:/secrets/.resticpwd:ro
      - ./.rclone.conf:/root/.config/rclone/rclone.conf:ro
```
