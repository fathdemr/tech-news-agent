APP_NAME    := tech-news-agent
DOCKER_IMG  := fathdemr/$(APP_NAME)
PLATFORMS   := linux/amd64,linux/arm64

# ─── Local Development ────────────────────────────────────────────────────────
.PHONY: run
run:                        ## Run locally (requires .env)
	go run ./cmd/server

.PHONY: test-run
test-run:                   ## Run once in test mode
	go run ./cmd/server --test

.PHONY: test-connection
test-connection:            ## Test API connections only
	go run ./cmd/server --test-connection

# ─── Build ────────────────────────────────────────────────────────────────────
.PHONY: build
build:                      ## Build Go binary
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(APP_NAME) ./cmd/server

.PHONY: clean
clean:                      ## Remove built binary
	rm -f $(APP_NAME)

# ─── Docker ───────────────────────────────────────────────────────────────────
.PHONY: docker-build
docker-build:               ## Build Docker image locally (current platform)
	docker build -t $(DOCKER_IMG):latest .

.PHONY: docker-push
docker-push:                ## Build multi-platform image & push to Docker Hub
	@echo "==> Building for $(PLATFORMS) and pushing $(DOCKER_IMG):latest..."
	docker buildx build --platform $(PLATFORMS) \
		-t $(DOCKER_IMG):latest \
		--push .
	@echo "==> Done! Your server can now pull the correct platform."

# ─── Server Deploy ────────────────────────────────────────────────────────────
# Usage: make deploy SERVER=user@your.server.ip
.PHONY: deploy
deploy: docker-push         ## Push image & trigger server to pull & restart
	@echo "==> Deploying to $(SERVER)..."
	ssh $(SERVER) "\
		cd ~/tech-news-agent && \
		docker compose pull && \
		docker compose up -d && \
		docker image prune -f"
	@echo "==> Deploy complete!"

# ─── Help ────────────────────────────────────────────────────────────────────
.PHONY: help
help:                       ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
