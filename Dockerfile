FROM ubuntu:22.04

# Set environment untuk non-interactive
ENV DEBIAN_FRONTEND=noninteractive

# Install runtime dependencies and build tools tanpa interaksi
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    git \
    wget \
    build-essential \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && ln -fs /usr/share/zoneinfo/UTC /etc/localtime

# Install Go
RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz -O /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz

ENV GOROOT=/usr/local/go
ENV PATH="${GOROOT}/bin:${PATH}"

WORKDIR /app

# Copy source files
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o telexec .

# Prepare runtime environment
RUN mkdir -p /app/out && \
    chown -R nobody:nogroup /app && \
    chmod -R 750 /app && \
    chmod 770 /app/out

USER nobody

HEALTHCHECK --interval=30s --timeout=3s \
  CMD ps aux | grep '[t]elexec' || exit 1

ENTRYPOINT ["/app/telexec"]