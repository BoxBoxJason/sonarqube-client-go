package_name := sonargo
sdk_dir := sonar
cli_dirs := ./cmd/... ./internal/...
endpoint := http://127.0.0.1:9000
enterprise_endpoint := http://127.0.0.1:9001
username := admin
password := admin
sonarqube_version := 26.7.0.124771-community
sonarqube_enterprise_version := 2026.3.1-enterprise
version := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
build_time := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# target can be: all (default), sdk, cli
target := all

# Compute the lint/test paths based on target
ifeq ($(target),sdk)
  target_paths := ./${sdk_dir}/...
else ifeq ($(target),cli)
  target_paths := ${cli_dirs}
else
  target_paths := ./${sdk_dir}/... ${cli_dirs}
endif

# Automatically detect container engine (docker or podman)
ifeq ($(shell command -v docker 2>/dev/null),)
  ifeq ($(shell command -v podman 2>/dev/null),)
    $(error Neither docker nor podman is installed. Please install one of them.)
  endif
  container_engine := podman
else
  container_engine := docker
endif

.PHONY: setup.sonar setup.sonar.enterprise test lint vuln coverage api api.enterprise build

# Run all unit tests (use target=sdk|cli|all to filter)
test:
	@command -v gotestsum >/dev/null 2>&1 || { echo "Installing gotestsum..."; go install gotest.tools/gotestsum@v1.13.0; }
	@mkdir -p codequality
	CGO_ENABLED=1 gotestsum --junitfile codequality/unit-tests.xml --format-icons octicons -- -race ${target_paths}

# Run tests with coverage report (use target=sdk|cli|all to filter)
coverage:
	@command -v gotestsum >/dev/null 2>&1 || { echo "Installing gotestsum..."; go install gotest.tools/gotestsum@v1.13.0; }
	@mkdir -p codequality
	gotestsum --junitfile codequality/unit-tests.xml --format-icons octicons -- -coverprofile=codequality/coverage.out -covermode=atomic ${target_paths}
	@echo "Coverage report generated: codequality/coverage.html"

# Run integration tests
e2e: setup.sonar
	@command -v ginkgo >/dev/null 2>&1 || { echo "Installing ginkgo..."; go install github.com/onsi/ginkgo/v2/ginkgo@v2.30.0; }
	SONAR_TOKEN= SONAR_URL=${endpoint} SONAR_USERNAME=${username} SONAR_PASSWORD=${password} ginkgo -r integration_testing

# Run enterprise edition integration tests
e2e.enterprise: setup.sonar.enterprise
	@command -v ginkgo >/dev/null 2>&1 || { echo "Installing ginkgo..."; go install github.com/onsi/ginkgo/v2/ginkgo@v2.30.0; }
	SONAR_TOKEN= SONAR_URL=${enterprise_endpoint} SONAR_USERNAME=${username} SONAR_PASSWORD=${password} ginkgo -r integration_testing

# Build the CLI binary to ./bin/sonar-cli.
# Version defaults to the latest git tag/commit. Override with: make build version=1.2.3
# Build time is stamped automatically at build time.
build:
	go build -o bin/sonar-cli -ldflags "-X github.com/boxboxjason/sonarqube-client-go/v2/internal/cli.version=$(version) -X github.com/boxboxjason/sonarqube-client-go/v2/internal/cli.buildTime=$(build_time)" ./cmd/sonar-cli

# Generate changelog using git-cliff
changelog:
	@command -v git-cliff >/dev/null 2>&1 || { echo "Please install git-cliff: https://github.com/orhun/git-cliff/releases"; exit 1; }
	git-cliff -c cliff.toml -o CHANGELOG.md

# Verify changelog is up-to-date (CI-friendly)
changelog-check:
	@command -v git-cliff >/dev/null 2>&1 || { echo "Please install git-cliff: https://github.com/orhun/git-cliff/releases"; exit 1; }
	@git-cliff -c cliff.toml -o /tmp/CHANGELOG.md
	@if [ ! -f CHANGELOG.md ]; then \
		echo "CHANGELOG.md does not exist, generating one with 'make changelog'"; \
		rm -f /tmp/CHANGELOG.md; \
		exit 1; \
	fi
	@if ! cmp -s CHANGELOG.md /tmp/CHANGELOG.md; then \
		echo "CHANGELOG.md is out of date. Run 'make changelog' and commit the changes."; \
		rm -f /tmp/CHANGELOG.md; \
		exit 1; \
	else \
		echo "CHANGELOG.md is up to date."; \
		rm -f /tmp/CHANGELOG.md; \
	fi

# Run golangci-lint (use target=sdk|cli|all to filter)
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2; }
	@mkdir -p codequality
	golangci-lint run ${target_paths}

# Scan dependencies and stdlib for known vulnerabilities (govulncheck).
vuln:
	@command -v govulncheck >/dev/null 2>&1 || { echo "Installing govulncheck..."; go install golang.org/x/vuln/cmd/govulncheck@v1.3.0; }
	govulncheck ./...

# Fetch SonarQube community edition API specification
api: setup.sonar
	@command -v curl >/dev/null 2>&1 || { echo "curl is required but not installed. Please install curl."; exit 1; }
	@mkdir -p assets
	@echo "Fetching SonarQube API v1 specification from ${endpoint}/api/webservices/list..."
	curl -u ${username}:${password} "${endpoint}/api/webservices/list?include_internals=true" -o ./assets/api.json
	@echo "Fetching SonarQube API v2 specification from ${endpoint}/api/v2/api-docs..."
	curl -u ${username}:${password} "${endpoint}/api/v2/api-docs" -o ./assets/api.v2.json
	@echo "API specification saved to ./assets/api.json and ./assets/api.v2.json"

# Fetch SonarQube enterprise edition API specification
# Requires a running SonarQube Enterprise Edition instance at enterprise_endpoint.
# Override the endpoint with: make api.enterprise enterprise_endpoint=http://my-sonarqube:9000
api.enterprise: setup.sonar.enterprise
	@command -v curl >/dev/null 2>&1 || { echo "curl is required but not installed. Please install curl."; exit 1; }
	@mkdir -p assets
	@echo "Fetching SonarQube Enterprise API v1 specification from ${enterprise_endpoint}/api/webservices/list..."
	curl -u ${username}:${password} "${enterprise_endpoint}/api/webservices/list?include_internals=true" -o ./assets/api.enterprise.json
	@echo "Fetching SonarQube Enterprise API v2 specification from ${enterprise_endpoint}/api/v2/api-docs..."
	curl -u ${username}:${password} "${enterprise_endpoint}/api/v2/api-docs" -o ./assets/api.enterprise.v2.json
	@echo "Enterprise API specification saved to ./assets/api.enterprise.json and ./assets/api.enterprise.v2.json"


# Setup SonarQube instance for integration testing
# If SonarQube API is already reachable, skip setup
# Else use container engine to start a SonarQube instance with a port mapping
setup.sonar:
	@command -v curl >/dev/null 2>&1 || { echo "curl is required but not installed. Please install curl."; exit 1; }
	@if curl -s -u ${username}:${password} ${endpoint}/api/system/health | grep -q "GREEN"; then \
		echo "SonarQube API is reachable at ${endpoint}/api. Skipping setup."; \
	else \
		if [ -n "$$GITHUB_ACTIONS" ] || [ -n "$$CI" ]; then \
			echo "Detected CI environment; not starting container. Waiting for SonarQube at ${endpoint}/api..."; \
		else \
			echo "Starting SonarQube instance using ${container_engine}..."; \
			${container_engine} run -d --name sonargo-sonarqube -p 9000:9000 docker.io/library/sonarqube:${sonarqube_version}; \
			echo "Waiting for SonarQube to be ready..."; \
		fi; \
		until curl -s -u ${username}:${password} ${endpoint}/api/system/health | grep -q "GREEN"; do \
			printf "."; \
			sleep 5; \
		done; \
		echo "\nSonarQube is ready at ${endpoint}."; \
	fi

# Setup SonarQube Enterprise Edition instance for integration testing.
# Requires a valid SonarQube Enterprise Edition license.
# Override the endpoint with: make setup.sonar.enterprise enterprise_endpoint=http://my-sonarqube:9000
setup.sonar.enterprise:
	@command -v curl >/dev/null 2>&1 || { echo "curl is required but not installed. Please install curl."; exit 1; }
	@if curl -s -u ${username}:${password} ${enterprise_endpoint}/api/system/health | grep -q "GREEN"; then \
		echo "SonarQube Enterprise API is reachable at ${enterprise_endpoint}/api. Skipping setup."; \
	else \
		if [ -n "$$GITHUB_ACTIONS" ] || [ -n "$$CI" ]; then \
			echo "Detected CI environment; not starting container. Waiting for SonarQube Enterprise at ${enterprise_endpoint}/api..."; \
		else \
			echo "Starting SonarQube Enterprise instance using ${container_engine}..."; \
			${container_engine} run -d --name sonargo-sonarqube-enterprise -p 9001:9000 docker.io/library/sonarqube:${sonarqube_enterprise_version}; \
			echo "Waiting for SonarQube Enterprise to be ready..."; \
		fi; \
		until curl -s -u ${username}:${password} ${enterprise_endpoint}/api/system/health | grep -q "GREEN"; do \
			printf "."; \
			sleep 5; \
		done; \
		echo "\nSonarQube Enterprise is ready at ${enterprise_endpoint}."; \
	fi

# Teardown SonarQube instance
teardown.sonar:
	@if ${container_engine} ps -a --format '{{.Names}}' | grep -w sonargo-sonarqube >/dev/null 2>&1; then \
		echo "Stopping SonarQube instance..."; \
		${container_engine} rm -f sonargo-sonarqube >/dev/null 2>&1 || echo "Failed to remove SonarQube container."; \
		echo "SonarQube instance stopped."; \
	else \
		echo "No SonarQube container 'sonargo-sonarqube' found. Nothing to teardown."; \
	fi

# Teardown SonarQube Enterprise Edition instance
teardown.sonar.enterprise:
	@if ${container_engine} ps -a --format '{{.Names}}' | grep -w sonargo-sonarqube-enterprise >/dev/null 2>&1; then \
		echo "Stopping SonarQube Enterprise instance..."; \
		${container_engine} rm -f sonargo-sonarqube-enterprise >/dev/null 2>&1 || echo "Failed to remove SonarQube Enterprise container."; \
		echo "SonarQube Enterprise instance stopped."; \
	else \
		echo "No SonarQube Enterprise container 'sonargo-sonarqube-enterprise' found. Nothing to teardown."; \
	fi
