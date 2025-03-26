FROM ubuntu:22.04

# Set environment untuk non-interactive
ENV DEBIAN_FRONTEND=noninteractive

# Install dependencies
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

# Install Go dalam satu layer
RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz -O /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz

# Set environment variables untuk Go
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOROOT/bin:$GOPATH/bin:$PATH

# Verifikasi instalasi Go
RUN go version

WORKDIR /app

# Install Air sebelum copy source code
RUN go install github.com/air-verse/air@latest

# Copy source code
COPY . .

# Initialize air config (jika belum ada)
RUN [ -f .air.toml ] || air init

# Prepare runtime environment
RUN mkdir -p /app/out && \
    chown -R nobody:nogroup /app && \
    chmod -R 750 /app && \
    chmod 770 /app/out

USER nobody

ENTRYPOINT ["air"]