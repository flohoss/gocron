ARG V_GOLANG=1.25
ARG V_DEBIAN=bookworm
FROM golang:${V_GOLANG}-${V_DEBIAN} AS final

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl wget ca-certificates tzdata unzip \
    python3 python3-pip python3-venv pipx \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/*

# Add pipx binary location to PATH
ENV PATH="/root/.local/bin:$PATH"

# air
RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download
