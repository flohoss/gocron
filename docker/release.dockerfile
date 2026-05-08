ARG V_DEBIAN
ARG V_GOLANG
ARG V_NODE

FROM node:${V_NODE}-alpine AS node-builder
WORKDIR /app

COPY ./web/package.json ./web/yarn.lock ./
RUN yarn install --frozen-lockfile --network-timeout 30000 --silent

COPY ./web/openapi.json ./web/openapi-ts.config.ts ./
RUN yarn types

COPY ./web/ ./
RUN yarn build

FROM golang:${V_GOLANG} AS golang-builder
WORKDIR /app

ARG APP_VERSION
ARG BUILD_TIME
ARG REPO_URL

COPY ./go.mod ./go.sum ./
RUN go mod download > /dev/null 2>&1

COPY . .
COPY --from=node-builder /app/dist ./internal/webui/dist
RUN go build -ldflags="-s -w -X github.com/flohoss/gocron/internal/buildinfo.Version=${APP_VERSION} -X github.com/flohoss/gocron/internal/buildinfo.BuildTime=${BUILD_TIME} -X github.com/flohoss/gocron/internal/buildinfo.RepoURL=${REPO_URL}" -o gocron .

FROM debian:${V_DEBIAN}-slim AS final

# Keep this block the same as in the dev Dockerfile
RUN apt-get update > /dev/null 2>&1 && apt-get install -y --no-install-recommends \
    curl wget tar ca-certificates tzdata unzip dumb-init \
    python3 python3-pip python3-venv pipx gnupg > /dev/null 2>&1 \
    && apt-get clean > /dev/null 2>&1 && rm -rf /var/lib/apt/lists/* /tmp/*

# Add pipx binary location to PATH
ENV PATH="/root/.local/bin:${PATH}"

WORKDIR /app

COPY --from=golang-builder /app/gocron .

EXPOSE 8156

ENTRYPOINT ["dumb-init", "--", "/app/gocron"]
