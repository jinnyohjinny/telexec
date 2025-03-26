IMAGE_NAME = telexec-ubuntu
CONTAINER_NAME = telexec-bot
OUT_DIR = ./out
ENV_FILE = .env

.PHONY: build run stop start clean logs help

build:
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME) .

run:
	@echo "Starting container in detached mode..."
	docker run -d \
		--restart unless-stopped \
		--name $(CONTAINER_NAME) \
		-v $(OUT_DIR):/app/out \
		--env-file $(ENV_FILE) \
		$(IMAGE_NAME)

stop:
	@echo "Stopping container..."
	docker stop $(CONTAINER_NAME)

start:
	@echo "Starting existing container..."
	docker start $(CONTAINER_NAME)

clean:
	@echo "Removing container..."
	docker rm -f $(CONTAINER_NAME) || true

logs:
	@echo "Showing logs..."
	docker logs -f $(CONTAINER_NAME)

shell:
	docker exec -it $(CONTAINER_NAME) /bin/bash

deploy: build clean run

help:
	@echo "Penggunaan:"
	@echo "  make build     - Build Docker image"
	@echo "  make run       - Jalankan container"
	@echo "  make stop      - Stop container"
	@echo "  make start     - Start container yang sudah ada"
	@echo "  make logs      - Tampilkan logs container"
	@echo "  make shell     - Masuk ke container"
	@echo "  make clean     - Hapus container"
	@echo "  make deploy    - Rebuild dan restart container"
	@echo "  make help      - Tampilkan bantuan ini"