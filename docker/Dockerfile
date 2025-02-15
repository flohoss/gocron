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

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY . .
RUN sqlc generate
RUN swag init --dir /app/cmd,/app/handlers,/app/services --parseDependency --requiredByDefault
RUN go build -ldflags="-s -w" cmd/main.go

FROM node:${V_NODE}-alpine AS node-builder
WORKDIR /app

COPY --from=go-builder /app/docs/ /docs/
COPY ./web/package.json ./web/yarn.lock ./
RUN yarn install --frozen-lockfile

COPY ./web/ ./
RUN yarn types
RUN yarn build

FROM docker:${V_DOCKER}-cli AS final
RUN apk add --update zip tzdata borgbackup dumb-init && \
    rm -rf /tmp/* /var/tmp/* /usr/share/man /var/cache/apk/*

# rclone
COPY --from=rclone --chmod=0755 \
    /usr/local/bin/rclone /usr/bin/rclone

# restic
COPY --from=restic --chmod=0755 \
    /usr/bin/restic /usr/bin/restic

WORKDIR /app
COPY --from=node-builder /app/dist/ ./web/
COPY --from=go-builder /app/main .
COPY --from=go-builder /app/config/config.yml ./config/config.yml

# Envs
ARG APP_VERSION
ENV APP_VERSION=$APP_VERSION
ARG BUILD_TIME
ENV BUILD_TIME=$BUILD_TIME

EXPOSE 8156

ENTRYPOINT ["dumb-init", "--"]
CMD ["/app/main"]
