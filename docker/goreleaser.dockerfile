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
RUN yarn install --frozen-lockfile --production=true

COPY ./web/ ./
RUN yarn build

FROM debian:${V_DEBIAN}-slim AS tools

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential wget autoconf automake libtool pkg-config \
    curl unzip zip gnupg ca-certificates bzip2 \
    # Python dependencies
    python3 python3-dev python3-pip python3-venv python3-wheel \
    # borgbackup & rdiff-backup dependencies
    libssl-dev liblz4-dev libzstd-dev libxxhash-dev libacl1-dev librsync-dev \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/*

# Install restic
RUN curl -L https://github.com/restic/restic/releases/download/v0.18.0/restic_0.18.0_linux_amd64.bz2 \
    | bunzip2 -c > /usr/local/bin/restic && chmod +x /usr/local/bin/restic

# Install rclone
RUN curl https://rclone.org/install.sh | bash

# Install docker-ce-cli and docker-compose-plugin
RUN install -m 0755 -d /etc/apt/keyrings
RUN curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
RUN chmod a+r /etc/apt/keyrings/docker.asc
RUN echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian \
    $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
RUN apt-get update && apt-get install -y --no-install-recommends \
    docker-ce-cli docker-compose-plugin \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/*

# Install apprise
RUN python3 -m venv /venv_apprise && \
    /venv_apprise/bin/pip install --no-cache-dir apprise && \
    rm -rf /root/.cache /tmp/*

# Install borgbackup
RUN python3 -m venv /venv_borg && \
    /venv_borg/bin/pip install --no-cache-dir borgbackup && \
    rm -rf /root/.cache /tmp/*

# Install rdiff-backup
RUN python3 -m venv /venv_rdiffbackup && \
    /venv_rdiffbackup/bin/pip install --no-cache-dir rdiff-backup && \
    rm -rf /root/.cache /tmp/*

# Install inotify-tools from source
RUN wget -q https://github.com/inotify-tools/inotify-tools/archive/refs/tags/4.23.9.0.tar.gz -O /tmp/inotify-tools.tar.gz && \
    tar -xzf /tmp/inotify-tools.tar.gz -C /tmp && \
    cd /tmp/inotify-tools-4.23.9.0 && \
    autoreconf -i && \
    ./configure --disable-shared && \
    make && \
    make install && \
    rm -rf /tmp/inotify-tools*

FROM debian:${V_DEBIAN}-slim AS final

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl wget ca-certificates rsync tzdata \
    podman podman-compose \
    python3-minimal openssh-client librsync2 \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/*

# Clean up common documentation and cache directories to reduce image size
RUN rm -rf /usr/share/doc /usr/share/man /usr/share/locale /var/cache/*

# Copy tools from tools stage
COPY --from=tools /usr/local/bin/inotifywait /usr/local/bin/inotifywait
COPY --from=tools /usr/local/bin/restic /usr/local/bin/restic
COPY --from=tools /usr/bin/rclone /usr/local/bin/rclone

# Docker
COPY --from=tools /usr/bin/docker /usr/local/bin/docker
COPY --from=tools /usr/libexec/docker/cli-plugins/docker-compose /usr/libexec/docker/cli-plugins/docker-compose

# Apprise
COPY --from=tools /venv_apprise /venv_apprise
# Borgbackup
COPY --from=tools /venv_borg /venv_borg
# Rdiff-backup
COPY --from=tools /venv_rdiffbackup /venv_rdiffbackup

# Set PATH to include all venv bins. Order matters for potential name clashes.
ENV PATH="/venv_apprise/bin:/venv_borg/bin:/venv_rdiffbackup/bin:$PATH"

WORKDIR /app

COPY gocron ./gocron

ARG APP_VERSION
ENV APP_VERSION=$APP_VERSION
ARG BUILD_TIME
ENV BUILD_TIME=$BUILD_TIME

COPY --from=logo /app/logo.txt .
COPY --from=node-builder /app/dist/ ./web/
COPY ./docker/entrypoint.sh .

EXPOSE 8156

ENTRYPOINT ["/app/entrypoint.sh"]
