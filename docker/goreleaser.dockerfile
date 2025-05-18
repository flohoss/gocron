ARG V_NODE=lts
ARG V_ALPINE=3
ARG V_DEBIAN=bookworm
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

FROM debian:${V_DEBIAN}-slim AS tools

RUN apt-get update && apt-get install -y \
    curl unzip zip gnupg ca-certificates bzip2 python3 python3-venv \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Install restic
RUN curl -L https://github.com/restic/restic/releases/download/v0.18.0/restic_0.18.0_linux_amd64.bz2 \
    | bunzip2 -c > /usr/local/bin/restic && chmod +x /usr/local/bin/restic

# Install rclone
RUN curl https://rclone.org/install.sh | bash

# Install docker
RUN apt-get install ca-certificates 
RUN install -m 0755 -d /etc/apt/keyrings
RUN curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
RUN chmod a+r /etc/apt/keyrings/docker.asc
RUN echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian \
    $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
RUN apt-get update && apt-get install -y \
    docker-ce-cli docker-compose-plugin \
    && rm -rf /var/lib/apt/lists/*

# Install apprise
RUN python3 -m venv /venv && \
    /venv/bin/pip install --no-cache-dir apprise && \
    rm -rf /root/.cache /tmp/*

FROM debian:${V_DEBIAN}-slim AS final

RUN apt-get update && apt-get install -y \
    curl wget dumb-init python3 rsync tzdata borgbackup rdiff-backup \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN rm -rf /usr/share/doc /usr/share/man /usr/share/locale /var/cache/* /tmp/*

# Copy tools from tools stage
COPY --from=tools /usr/local/bin/restic /usr/local/bin/restic
COPY --from=tools /usr/bin/docker /usr/local/bin/docker
COPY --from=tools /usr/libexec/docker/cli-plugins/docker-compose /usr/libexec/docker/cli-plugins/docker-compose
COPY --from=tools /usr/bin/rclone /usr/local/bin/rclone
COPY --from=tools /venv /venv
ENV PATH="/venv/bin:$PATH"

WORKDIR /app

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

ENTRYPOINT ["/app/entrypoint.sh"]
