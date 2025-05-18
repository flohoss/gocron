ARG V_GOLANG=1.24
ARG V_DEBIAN=bookworm
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

FROM golang:${V_GOLANG}-${V_DEBIAN} AS final

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

# air
RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download
