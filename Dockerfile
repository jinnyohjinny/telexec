# Stage 1: Builder
FROM ubuntu:22.04 AS builder

# Install Go
RUN apt-get update && \
    apt-get install -y wget && \
    wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz -O /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# Install build dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    git \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o telexec .

# Stage 2: Runtime
FROM ubuntu:22.04

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    git \
    && rm -rf /var/lib/apt/lists/*

# Install Go di runtime juga jika ingin menjalankan 'go install'
ENV GOROOT=/usr/local/go
ENV PATH="${GOROOT}/bin:${PATH}"
COPY --from=builder /usr/local/go /usr/local/go

WORKDIR /app
COPY --from=builder /app/telexec /app/telexec
COPY --from=builder /app/.env /app/.env

RUN mkdir -p /app/out && \
    chown -R nobody:nogroup /app && \
    chmod -R 750 /app && \
    chmod 770 /app/out

USER nobody

HEALTHCHECK --interval=30s --timeout=3s \
  CMD ps aux | grep '[t]elexec' || exit 1

ENTRYPOINT ["/app/telexec"]