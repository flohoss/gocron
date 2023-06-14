ARG DOCKER_VERSION
ARG GOLANG_VERSION
ARG NODE_VERSION
ARG ALPINE_VERSION
ARG RCLONE_VERSION
ARG RESTIC_VERSION
FROM rclone/rclone:${RCLONE_VERSION} AS rclone
FROM restic/restic:${RESTIC_VERSION} AS restic
FROM golang:${GOLANG_VERSION}-alpine AS goBuilder

WORKDIR /src
COPY ./src/go.mod .
COPY ./src/go.sum .
RUN go mod download

COPY ./src/ .
RUN go build -ldflags="-s -w"

FROM node:${NODE_VERSION}-alpine AS nodeBuilder

COPY package.json .
COPY yarn.lock .
RUN yarn install

COPY tailwind.config.js .
COPY ./src/templates/ ./src/templates/
COPY ./src/static/ ./src/static/
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
WORKDIR /src
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

COPY --from=logo /logo.txt .
COPY --from=nodeBuilder /src/static/ ./static/
COPY --from=nodeBuilder /src/templates ./templates
COPY --from=goBuilder /src/gobackup .

# Envs
ARG VERSION
ENV VERSION=$VERSION
ARG BUILDTIME
ENV BUILDTIME=$BUILDTIME

ENTRYPOINT [ "/src/entrypoint.sh" ]
