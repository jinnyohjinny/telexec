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

# Install Go dan copy binary ke /app
RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz -O /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz && \
    mkdir -p /app && \
    cp /usr/local/go/bin/go /app/go

# Set PATH untuk mencari binary go di /app terlebih dahulu
ENV PATH="/app:/usr/local/go/bin:${PATH}"
ENV GOROOT=/usr/local/go

WORKDIR /app
COPY . .

# Verifikasi go command
RUN go version && \
    which go && \
    ls -la /app/go

# Build aplikasi
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o telexec .

# Prepare runtime environment
RUN mkdir -p /app/out && \
    chown -R nobody:nogroup /app && \
    chmod -R 750 /app && \
    chmod 770 /app/out

USER nobody

ENTRYPOINT ["/app/telexec"]