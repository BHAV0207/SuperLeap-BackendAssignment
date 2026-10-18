CLUSTER   ?= lead-crm
NAMESPACE ?= lead-crm
RELEASE   ?= lead-crm
CHART     := deploy/helm/lead-crm
IMAGE     ?= ghcr.io/bhav0207/lead-crm
# image tag is always a git commit sha so the running version maps to a commit
TAG       ?= $(shell git rev-parse --short=7 HEAD)
APP_URL   ?= http://localhost:8088

.PHONY: help up cluster secrets monitoring image deploy smoke rollback status test lint down

help:
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-12s %s\n", $$1, $$2}'

up: cluster secrets monitoring image deploy smoke ## everything from scratch on a local kind cluster

cluster: ## create the kind cluster
	@kind get clusters | grep -qx $(CLUSTER) || kind create cluster --config deploy/kind-config.yaml
	kubectl config use-context kind-$(CLUSTER)

secrets: ## create the db secret with a random password (only if it does not exist)
	kubectl create namespace $(NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	@kubectl -n $(NAMESPACE) get secret lead-crm-db >/dev/null 2>&1 || \
		kubectl -n $(NAMESPACE) create secret generic lead-crm-db \
			--from-literal=password="$$(openssl rand -base64 24)"

monitoring: ## install prometheus, grafana, alertmanager and jaeger
	kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -
	@kubectl -n monitoring get secret grafana-admin >/dev/null 2>&1 || \
		kubectl -n monitoring create secret generic grafana-admin \
			--from-literal=admin-user=admin \
			--from-literal=admin-password="$$(openssl rand -base64 18)"
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts >/dev/null 2>&1 || true
	helm repo update prometheus-community >/dev/null
	helm upgrade --install kube-prometheus-stack prometheus-community/kube-prometheus-stack \
		--version 92.10.1 -n monitoring -f deploy/monitoring/kube-prometheus-stack.yaml --wait --timeout 10m
	kubectl apply -f deploy/monitoring/jaeger.yaml

image: ## build the image locally and load it into kind (CI builds the same image and pushes to GHCR)
	docker build -t $(IMAGE):$(TAG) \
		--build-arg VERSION=$(TAG) \
		--build-arg COMMIT=$$(git rev-parse HEAD) \
		--build-arg BUILD_TIME=$$(date -u +%Y-%m-%dT%H:%M:%SZ) .
	kind load docker-image $(IMAGE):$(TAG) --name $(CLUSTER)

# --atomic rolls back automatically if the new pods never become ready
deploy: ## helm upgrade to TAG (rolling update, auto rollback on failure)
	helm upgrade --install $(RELEASE) $(CHART) -n $(NAMESPACE) \
		--set app.image.repository=$(IMAGE) \
		--set app.image.tag=$(TAG) \
		--atomic --wait --timeout 3m
	helm -n $(NAMESPACE) history $(RELEASE) --max 5

smoke: ## check the deployed version and run a create/read user journey
	APP_URL=$(APP_URL) EXPECTED_VERSION=$(TAG) ./scripts/smoke-test.sh

rollback: ## roll back to the previous helm revision (or REVISION=n)
	helm -n $(NAMESPACE) rollback $(RELEASE) $(REVISION) --wait
	helm -n $(NAMESPACE) history $(RELEASE) --max 5

status: ## pods, release history and running version
	kubectl -n $(NAMESPACE) get pods,svc,pvc
	helm -n $(NAMESPACE) history $(RELEASE) --max 5
	@curl -s $(APP_URL)/version; echo

test: ## go tests
	go test -race ./...

lint: ## gofmt, vet and chart lint
	test -z "$$(gofmt -l .)"
	go vet ./...
	helm lint $(CHART) --set app.image.tag=lint

down: ## delete the kind cluster
	kind delete cluster --name $(CLUSTER)
