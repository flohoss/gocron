ARG DOCKER_VERSION
ARG GOLANG_VERSION
ARG NODE_VERSION
ARG ALPINE_VERSION
ARG RCLONE_VERSION
FROM rclone/rclone:${RCLONE_VERSION} AS rclone
FROM golang:${GOLANG_VERSION}-alpine AS goBuilder
WORKDIR /app
RUN apk add gcc

COPY ./scripts/swagger.sh .
RUN ./swagger.sh install

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY . .
RUN ./swagger.sh init
RUN go build -ldflags="-s -w" cmd/gobackup/gobackup.go

FROM node:${NODE_VERSION}-alpine AS nodeBuilder
WORKDIR /app

COPY ./web/package.json .
COPY ./web/yarn.lock .
RUN yarn install --frozen-lockfile

COPY --from=goBuilder /app/docs/swagger.json ../docs/swagger.json
COPY ./web/ .
RUN yarn run types:openapi
RUN yarn run build

FROM alpine:${ALPINE_VERSION} AS logo
RUN apk add figlet
RUN figlet GoBackup > logo.txt

FROM golang:${GOLANG_VERSION}-alpine AS resticBuilder
ARG RESTIC_VERSION
RUN apk add wget make
RUN wget https://github.com/restic/restic/archive/refs/tags/v${RESTIC_VERSION}.zip -q
RUN unzip v${RESTIC_VERSION}.zip -q
RUN cd restic-${RESTIC_VERSION} && make restic

FROM docker:${DOCKER_VERSION}-cli AS final
ARG RESTIC_VERSION
RUN apk add tzdata

# Rclone
COPY --from=rclone --chmod=0755 \
    /usr/local/bin/rclone /usr/bin/rclone

# Restic
COPY --from=resticBuilder --chmod=0755 \
    /go/restic-${RESTIC_VERSION}/restic /usr/bin/restic

# GoBackup
WORKDIR /app
COPY ./scripts/entrypoint.sh .

COPY --from=logo /logo.txt .
COPY --from=nodeBuilder /app/dist/ ./web/
COPY --from=goBuilder /app/gobackup .

# Envs
ARG APP_VERSION
ENV APP_VERSION=$APP_VERSION
ARG BUILD_TIME
ENV BUILD_TIME=$BUILD_TIME

ENTRYPOINT [ "/app/entrypoint.sh" ]
