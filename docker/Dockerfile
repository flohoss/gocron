ARG V_DOCKER=27.3.1
ARG V_GOLANG=1.23
ARG V_NODE=lts
ARG V_ALPINE=3
ARG V_RCLONE=1
ARG V_RESTIC=0.17.3
FROM rclone/rclone:${V_RCLONE} AS rclone
FROM restic/restic:${V_RESTIC} AS restic
FROM golang:${V_GOLANG}-alpine AS go-builder
RUN apk add cmake g++ gcc
WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-s -w" cmd/gobackup/gobackup.go

FROM node:${V_NODE}-alpine AS node-builder
WORKDIR /web

COPY /docs/swagger.json ../docs/swagger.json
COPY ./web/ .
RUN yarn install --frozen-lockfile
RUN yarn run types:openapi
RUN yarn run build

FROM alpine:${V_ALPINE} AS logo
RUN apk add figlet
RUN figlet GoBackup > logo.txt

FROM docker:${V_DOCKER}-cli AS final
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
COPY --from=node-builder /web/dist/ ./web/
COPY --from=go-builder /app/gobackup .

# Envs
ARG APP_VERSION
ENV APP_VERSION=$APP_VERSION
ARG BUILD_TIME
ENV BUILD_TIME=$BUILD_TIME

EXPOSE 8080

ENTRYPOINT [ "/app/entrypoint.sh" ]
