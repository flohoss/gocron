#!/bin/sh

cat logo.txt
CMD="./gocron"

if [ ! -d /app/config ]; then
    mkdir -p ./config
fi

if [ ! -e /app/config/config.yml ]; then
    cp /tmp/config.yml ./config/config.yml
fi

if [ -n "$PUID" ] || [ -n "$PGID" ]; then
    USER=appuser
    HOME=/app

    if ! grep -q "$USER" /etc/passwd; then
        addgroup -g "$PGID" "$USER"
        adduser -h "$HOME" -g "" -G "$USER" -D -H -u "$PUID" "$USER"
    fi

    chown "$USER":"$USER" "$HOME" -R
    printf "\nUID: %s GID: %s\n\n" "$PUID" "$PGID"

    # Use `exec su-exec` to correctly switch user and run the process
    exec su-exec "$USER" $CMD
else
    printf "\nWARNING: Running docker as root\n\n"
    exec $CMD
fi
