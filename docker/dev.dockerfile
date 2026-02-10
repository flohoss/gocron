ARG V_GOLANG=1.25.6
ARG V_DEBIAN=trixie
ARG V_AIR=1.64.5
ARG V_SQLC=1.30.0
FROM debian:${V_DEBIAN}-slim AS final

# Keep this block the same as in the release Dockerfile
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl git wget tar ca-certificates tzdata unzip dumb-init \
    python3 python3-pip python3-venv pipx gnupg \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/*

# Add pipx binary location to PATH
ENV PATH="/root/.local/bin:${PATH}"

# Install go
ARG V_GOLANG
ARG TARGETARCH
ENV GO_VERSION=$V_GOLANG
RUN curl -fsSLo /tmp/go.tgz "https://go.dev/dl/go${GO_VERSION}.linux-${TARGETARCH}.tar.gz" \
    && tar -C /usr/local -xzf /tmp/go.tgz \
    && rm /tmp/go.tgz
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"

ARG V_AIR
RUN go install github.com/air-verse/air@v${V_AIR}
ARG V_SQLC
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v${V_SQLC}

WORKDIR /app

ENV APP_VERSION=v0.0.0.0-dev

COPY ./go.mod ./go.sum ./
RUN go mod download
