ARG V_DOCKER=27.3.1
ARG V_GOLANG=1.24
ARG V_NODE=lts
ARG V_ALPINE=3
ARG V_RCLONE=1
ARG V_RESTIC=0.17.3
FROM rclone/rclone:${V_RCLONE} AS rclone
FROM restic/restic:${V_RESTIC} AS restic
FROM alpine:${V_ALPINE} AS logo
WORKDIR /app
RUN apk add figlet
RUN figlet GoCron > logo.txt

FROM node:${V_NODE}-alpine AS node-builder
WORKDIR /app

COPY ./web/package.json ./web/yarn.lock ./
RUN yarn install --frozen-lockfile

COPY ./web/ ./
RUN yarn build

FROM docker:${V_DOCKER}-cli AS final
RUN apk add --update --no-cache \
    python3 py3-pip \
    su-exec dumb-init \
    zip tzdata borgbackup rsync curl rdiff-backup && \
    rm -rf /tmp/* /var/tmp/* /usr/share/man /var/cache/apk/*

# rclone
COPY --from=rclone --chmod=0755 \
    /usr/local/bin/rclone /usr/bin/rclone

# restic
COPY --from=restic --chmod=0755 \
    /usr/bin/restic /usr/bin/restic

# apprise
RUN apk add --no-cache py3-pip && \
    python3 -m venv /venv && \
    /venv/bin/pip install --no-cache-dir apprise==1.9.3 && \
    apk del py3-pip && \
    find /venv -name '*.pyc' -delete && \
    find /venv -type d -name '__pycache__' -exec rm -r {} + && \
    rm -rf /root/.cache /tmp/*
ENV PATH="/venv/bin:$PATH"

WORKDIR /app

# goreleaser
COPY gocron ./gocron

ARG APP_VERSION
ENV APP_VERSION=$APP_VERSION
ARG BUILD_TIME
ENV BUILD_TIME=$BUILD_TIME

COPY ./config/config.yml /tmp/config.yml
COPY --from=logo /app/logo.txt .
COPY --from=node-builder /app/dist/ ./web/
COPY ./docker/entrypoint.sh .

EXPOSE 8156

ENTRYPOINT ["dumb-init", "--", "/app/entrypoint.sh"]
