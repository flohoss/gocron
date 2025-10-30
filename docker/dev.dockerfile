ARG V_GOLANG=1.25.3
ARG V_DEBIAN=trixie
FROM debian:${V_DEBIAN}-slim AS final

# Keep this block the same as in the release Dockerfile
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl wget tar ca-certificates tzdata unzip dumb-init \
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

# Install air
RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download
