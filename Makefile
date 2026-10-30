# Sentinel AI — Makefile
# ======================
# Local development commands for the Autonomous SRE Engineer project.

.PHONY: help setup infra-up infra-down kind-up kind-down run-api run-worker \
        run-orchestrator trigger-incident build-api docker-up docker-down \
        lint test clean

GO      ?= go
PYTHON  ?= python3
DOCKER  ?= docker
KIND    ?= kind
KUBECTL ?= kubectl

PROJECT := sentinel-ai
API_DIR  := cmd/sre-api
ORCH_DIR := agent-orchestrator

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}'

setup: ## Initialize the repo (go mod, python venv)
	@echo "==> Initializing Go module"
	cd $(API_DIR) && $(GO) mod init sentinel-ai/sre-api || true
	cd $(API_DIR) && $(GO) mod tidy || true
	@echo "==> Setting up Python virtual environment"
	cd $(ORCH_DIR) && $(PYTHON) -m venv .venv || true
	cd $(ORCH_DIR) && . .venv/bin/activate && pip install --upgrade pip || true
	@echo "==> Setup complete"

# ---------------------------------------------------------------------------
# Local infrastructure (Postgres, Redis via docker-compose)
# ---------------------------------------------------------------------------
infra-up: ## Start Postgres + Redis with docker-compose
	@echo "==> Starting local infrastructure"
	$(DOCKER) compose up -d postgres redis
	@echo "==> Infrastructure ready"

infra-down: ## Stop local infrastructure
	$(DOCKER) compose down

# ---------------------------------------------------------------------------
# Kind Kubernetes cluster
# ---------------------------------------------------------------------------
kind-up: ## Create Kind cluster and deploy monitoring + sample app
	@echo "==> Creating Kind cluster"
	$(KIND) create cluster --config deployments/kind/cluster.yaml --name sentinel || true
	@echo "==> Deploying Prometheus & Loki"
	$(KUBECTL) apply -f deployments/monitoring/prometheus/
	$(KUBECTL) apply -f deployments/monitoring/loki/
	@echo "==> Deploying sample application"
	$(KUBECTL) apply -f deployments/application/
	@echo "==> Kind cluster ready"

kind-down: ## Delete Kind cluster
	$(KIND) delete cluster --name sentinel || true

# ---------------------------------------------------------------------------
# Go services
# ---------------------------------------------------------------------------
run-api: ## Run the Go SRE API
	cd $(API_DIR) && $(GO) run main.go

run-worker: ## Run the Go SRE worker
	cd cmd/sre-worker && $(GO) run main.go

build-api: ## Build the Go SRE API binary
	cd $(API_DIR) && $(GO) build -o ../../bin/sre-api .

# ---------------------------------------------------------------------------
# Python orchestrator
# ---------------------------------------------------------------------------
run-orchestrator: ## Run the Python LangGraph orchestrator
	cd $(ORCH_DIR) && . .venv/bin/activate && $(PYTHON) main.py

# ---------------------------------------------------------------------------
# Incident demo
# ---------------------------------------------------------------------------
trigger-incident: ## Trigger a sample incident (CPU spike / error rate)
	@echo "==> Triggering incident scenario"
	bash test-infra/scenarios/high-error-rate.sh || true
	@echo "==> Incident triggered. Watch the orchestrator."

# ---------------------------------------------------------------------------
# Docker
# ---------------------------------------------------------------------------
docker-up: ## Build and start all Docker services
	$(DOCKER) compose up -d --build

docker-down: ## Stop all Docker services
	$(DOCKER) compose down

# ---------------------------------------------------------------------------
# Lint / test / clean
# ---------------------------------------------------------------------------
lint: ## Run Go lint + Python format checks
	cd $(API_DIR) && $(GO) vet ./... || true
	cd $(ORCH_DIR) && . .venv/bin/activate && $(PYTHON) -m py_compile *.py agents/*.py orchestration/*.py llm/*.py tools/*.py schemas/*.py || true

test: ## Run Go unit tests
	cd $(API_DIR) && $(GO) test ./... || true

clean: ## Clean build artifacts
	rm -rf bin/ $(API_DIR)/.venv $(ORCH_DIR)/.venv || true