ARG DOCKER_VERSION
ARG GOLANG_VERSION
ARG NODE_VERSION
ARG ALPINE_VERSION
ARG RCLONE_VERSION
ARG RESTIC_VERSION
FROM rclone/rclone:${RCLONE_VERSION} AS rclone
FROM restic/restic:${RESTIC_VERSION} AS restic
FROM golang:${GOLANG_VERSION}-alpine AS goBuilder
RUN apk add cmake g++ gcc
WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" cmd/gobackup/gobackup.go

FROM node:${NODE_VERSION}-alpine AS nodeBuilder
WORKDIR /web

COPY /docs/swagger.json ../docs/swagger.json
COPY ./web/ .
RUN yarn install --frozen-lockfile
RUN yarn run types:openapi
RUN yarn run build

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
COPY ./scripts/entrypoint.sh .

COPY --from=logo /logo.txt .
COPY --from=nodeBuilder /web/dist/ ./web/
COPY --from=goBuilder /app/gobackup .

# Envs
ARG APP_VERSION
ENV APP_VERSION=$APP_VERSION
ARG BUILD_TIME
ENV BUILD_TIME=$BUILD_TIME

ENTRYPOINT [ "/app/entrypoint.sh" ]
