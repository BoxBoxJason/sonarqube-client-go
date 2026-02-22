package_name := sonargo
target_dir := sonar
endpoint := http://127.0.0.1:9000
username := admin
password := admin
sonarqube_version := 26.2.0.119303-community

# Automatically detect container engine (docker or podman)
ifeq ($(shell command -v docker 2>/dev/null),)
  ifeq ($(shell command -v podman 2>/dev/null),)
    $(error Neither docker nor podman is installed. Please install one of them.)
  endif
  container_engine := podman
else
  container_engine := docker
endif

.PHONY: setup.sonar test lint coverage api

# Run all unit tests
test:
	@command -v gotestsum >/dev/null 2>&1 || { echo "Installing gotestsum..."; go install gotest.tools/gotestsum@v1.13.0; }
	@mkdir -p codequality
	gotestsum --junitfile codequality/unit-tests.xml --format-icons octicons -- ./${target_dir}/...

# Run tests with coverage report
coverage:
	@command -v gotestsum >/dev/null 2>&1 || { echo "Installing gotestsum..."; go install gotest.tools/gotestsum@v1.13.0; }
	@mkdir -p codequality
	gotestsum --junitfile codequality/unit-tests.xml --format-icons octicons -- -coverprofile=codequality/coverage.out -covermode=atomic ./${target_dir}/...
	@echo "Coverage report generated: codequality/coverage.html"

# Run integration tests
e2e: setup.sonar
	@command -v ginkgo >/dev/null 2>&1 || { echo "Installing ginkgo..."; go install github.com/onsi/ginkgo/v2/ginkgo@v2.28.1; }
	SONAR_TOKEN= SONAR_URL=${endpoint} SONAR_USERNAME=${username} SONAR_PASSWORD=${password} ginkgo -r integration_testing

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

# Run golangci-lint
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0; }
	@mkdir -p codequality
	golangci-lint run ./${target_dir}/...

# Fetch SonarQube API specification
api:
	@command -v curl >/dev/null 2>&1 || { echo "curl is required but not installed. Please install curl."; exit 1; }
	@mkdir -p assets
	@echo "Fetching SonarQube API specification from ${endpoint}/api/webservices/list..."
	curl -u ${username}:${password} "${endpoint}/api/webservices/list?include_internals=true" -o ./assets/api.json
	@echo "API specification saved to ./assets/api.json"

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

# Teardown SonarQube instance
teardown.sonar:
	@if ${container_engine} ps -a --format '{{.Names}}' | grep -w sonargo-sonarqube >/dev/null 2>&1; then \
		echo "Stopping SonarQube instance..."; \
		${container_engine} rm -f sonargo-sonarqube >/dev/null 2>&1 || echo "Failed to remove SonarQube container."; \
		echo "SonarQube instance stopped."; \
	else \
		echo "No SonarQube container 'sonargo-sonarqube' found. Nothing to teardown."; \
	fi
