#!/bin/sh

cat logo.txt
CMD=./gocron
# copy example config if it does not exist
cp -n /tmp/config.json /app/config/config.json

if [ -n "$PUID" ] || [ -n "$PGID" ]; then
    USER=appuser
    HOME=/app

    if ! grep -q "$USER" /etc/passwd; then
        addgroup -g "$PGID" "$USER"
        adduser -h "$HOME" -g "" -G "$USER" -D -H -u "$PUID" "$USER"
    fi

    chown "$USER":"$USER" "$HOME" -R
    printf "\nUID: %s GID: %s\n\n" "$PUID" "$PGID"
    exec su -c - $USER "$CMD"
else
    printf "\nWARNING: Running docker as root\n\n"
    exec "$CMD"
fi
