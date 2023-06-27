ARG DOCKER_VERSION
ARG GOLANG_VERSION
ARG NODE_VERSION
ARG ALPINE_VERSION
ARG RCLONE_VERSION
ARG RESTIC_VERSION
FROM rclone/rclone:${RCLONE_VERSION} AS rclone
FROM restic/restic:${RESTIC_VERSION} AS restic
FROM golang:${GOLANG_VERSION}-alpine AS goBuilder
WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" cmd/gobackup/gobackup.go

FROM node:${NODE_VERSION}-alpine AS nodeBuilder

COPY ./package.json .
COPY ./yarn.lock .
RUN yarn install --frozen-lockfile

COPY ./web ./web
COPY ./tailwind.config.js .
RUN yarn run tailwind:build

FROM alpine:${ALPINE_VERSION} AS logo
RUN apk add figlet
RUN figlet GoBackup > logo.txt

FROM docker:${DOCKER_VERSION}-cli AS final
RUN apk add tzdata

# Rclone
COPY --from=rclone --chmod=0755 \
    /usr/local/bin/rclone /usr/bin/rclone

# Restic
COPY --from=restic --chmod=0755 \
    /usr/bin/restic /usr/bin/restic

# GoBackup
WORKDIR /app
COPY scripts/entrypoint.sh .

COPY --from=logo /logo.txt .
COPY --from=nodeBuilder /web/static/ ./web/static/
COPY --from=nodeBuilder /web/templates ./web/templates
COPY --from=goBuilder /app/gobackup .

# Envs
ARG VERSION
ENV VERSION=$VERSION
ARG BUILD_TIME
ENV BUILD_TIME=$BUILD_TIME

ENTRYPOINT [ "/app/entrypoint.sh" ]
