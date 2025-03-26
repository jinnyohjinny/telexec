FROM ubuntu:22.04 AS builder

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    git \
    build-essential \
    golang-go \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o telexec .

FROM ubuntu:22.04

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/telexec /app/telexec
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/out /app/out

RUN chown -R nobody:nogroup /app && \
    chmod -R 750 /app && \
    mkdir -p /app/out && \
    chmod 770 /app/out

USER nobody

HEALTHCHECK --interval=30s --timeout=3s \
  CMD ps aux | grep '[t]elexec' || exit 1

ENTRYPOINT ["/app/telexec"]