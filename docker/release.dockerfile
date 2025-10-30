ARG V_GOLANG=1.25
ARG V_NODE=lts
ARG V_ALPINE=3
ARG V_DEBIAN=trixie
FROM alpine:${V_ALPINE} AS logo
WORKDIR /app
RUN apk add figlet > /dev/null 2>&1
RUN figlet GoCron > logo.txt

FROM golang:${V_GOLANG} AS golang-builder
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download > /dev/null 2>&1

COPY . .
RUN go build -o gocron .

FROM node:${V_NODE}-alpine AS node-builder
WORKDIR /app

COPY ./web/package.json ./web/yarn.lock ./
RUN yarn install --frozen-lockfile --network-timeout 30000 --silent

COPY ./web/openapi.json ./web/openapi-ts.config.ts ./
RUN yarn types

COPY ./web/ ./
RUN yarn build

FROM debian:${V_DEBIAN}-slim AS final

# Keep this block the same as in the dev Dockerfile
RUN apt-get update > /dev/null 2>&1 && apt-get install -y --no-install-recommends \
    curl wget tar ca-certificates tzdata unzip dumb-init \
    python3 python3-pip python3-venv pipx gnupg > /dev/null 2>&1 \
    && apt-get clean > /dev/null 2>&1 && rm -rf /var/lib/apt/lists/* /tmp/*

# Add pipx binary location to PATH
ENV PATH="/root/.local/bin:${PATH}"

WORKDIR /app

ARG APP_VERSION
ENV APP_VERSION=$APP_VERSION
ARG BUILD_TIME
ENV BUILD_TIME=$BUILD_TIME
ARG REPO
ENV REPO=$REPO

COPY --from=logo /app/logo.txt .
COPY --from=golang-builder /app/gocron .
COPY --from=node-builder /app/dist/ ./web/

EXPOSE 8156

ENTRYPOINT ["dumb-init", "--", "/app/gocron"]
