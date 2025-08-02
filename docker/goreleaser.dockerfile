ARG V_NODE=lts
ARG V_ALPINE=3
ARG V_DEBIAN=bookworm
FROM alpine:${V_ALPINE} AS logo
WORKDIR /app
RUN apk add --no-cache figlet
RUN figlet GoCron > logo.txt

FROM node:${V_NODE}-alpine AS node-builder
WORKDIR /app

COPY ./web/package.json ./web/yarn.lock ./
RUN yarn install --frozen-lockfile

COPY ./web/ ./
RUN yarn build

FROM debian:${V_DEBIAN}-slim AS final

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl wget ca-certificates tzdata unzip dumb-init \
    python3 python3-pip python3-venv pipx \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/*

# Add pipx binary location to PATH
ENV PATH="/root/.local/bin:$PATH"

WORKDIR /app

COPY gocron ./gocron

ARG APP_VERSION
ENV APP_VERSION=$APP_VERSION
ARG BUILD_TIME
ENV BUILD_TIME=$BUILD_TIME

COPY --from=logo /app/logo.txt .
COPY --from=node-builder /app/dist/ ./web/

EXPOSE 8156

ENTRYPOINT ["dumb-init", "--", "/app/gocron"]
