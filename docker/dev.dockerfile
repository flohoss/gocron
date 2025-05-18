ARG V_GOLANG=1.24
FROM golang:${V_GOLANG}-bookworm

# Install required packages
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    python3-venv \
    dumb-init \
    curl \
    zip \
    tzdata \
    restic \
    rsync \
    borgbackup \
    rdiff-backup \
    && rm -rf /var/lib/apt/lists/*

# docker
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

# rclone
RUN curl https://rclone.org/install.sh | bash

# restic
RUN restic self-update

# apprise
RUN python3 -m venv /venv && \
    /venv/bin/pip install --no-cache-dir apprise && \
    rm -rf /root/.cache /tmp/*
ENV PATH="/venv/bin:$PATH"

# air
RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download
