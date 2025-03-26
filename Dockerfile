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

# Install Go
RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz -O /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz

# Install Air untuk live reload
RUN go install github.com/cosmtrek/air@latest

ENV PATH="/root/go/bin:/usr/local/go/bin:${PATH}"
ENV GOROOT=/usr/local/go
ENV GOPATH=/root/go

WORKDIR /app

# Copy semua file termasuk konfigurasi Air (jika ada)
COPY . .

RUN air init

# Buat direktori out jika diperlukan
RUN mkdir -p /app/out && \
    chown -R nobody:nogroup /app && \
    chmod -R 750 /app && \
    chmod 770 /app/out

USER nobody

# Gunakan Air sebagai entrypoint
ENTRYPOINT ["air"]