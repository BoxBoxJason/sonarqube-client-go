package_name := sonargo
target_dir := sonar
endpoint := http://127.0.0.1:9000/api
username := admin
password := admin
container_engine := docker
sonarqube_version := 25.12.0.117093-community

.PHONY: setup.sonar clean generate test lint no-diff coverage

clean:
	rm -f ${target_dir}/*.go
	rm -rf integration_testing

generate: setup.sonar
	go generate ./...

# Run all unit tests
test:
	rm -rf integration_testing
	@command -v gotestsum >/dev/null 2>&1 || { echo "Installing gotestsum..."; go install gotest.tools/gotestsum@latest; }
	@mkdir -p codequality
	gotestsum --junitfile codequality/unit-tests.xml --format-icons octicons -- ./...
	rm -rf integration_testing

# Run tests with coverage report
coverage:
	rm -rf integration_testing
	@command -v gotestsum >/dev/null 2>&1 || { echo "Installing gotestsum..."; go install gotest.tools/gotestsum@latest; }
	@mkdir -p codequality
	gotestsum --junitfile codequality/unit-tests.xml --format-icons octicons -- -coverprofile=codequality/coverage.out -covermode=atomic ./...
	rm -rf integration_testing
	@echo "Coverage report generated: codequality/coverage.html"

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
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest; }
	@mkdir -p codequality
	golangci-lint run ./...

# Check for uncommitted changes after generation (useful for CI)
no-diff: generate
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Error: There are uncommitted changes after running 'make generate'"; \
		echo "Please run 'make generate' locally and commit the changes."; \
		git status --porcelain; \
		git diff --stat; \
		exit 1; \
	else \
		echo "Success: No uncommitted changes detected."; \
	fi

# Setup SonarQube instance for generation / integration testing
# If SonarQube API is already reachable, skip setup
# Else use container engine to start a SonarQube instance with a port mapping
setup.sonar:
	@command -v curl >/dev/null 2>&1 || { echo "curl is required but not installed. Please install curl."; exit 1; }
	@if curl -s -u ${username}:${password} ${endpoint}/system/health | grep -q "GREEN"; then \
		echo "SonarQube API is reachable at ${endpoint}. Skipping setup."; \
	else \
		if [ -n "$$GITHUB_ACTIONS" ] || [ -n "$$CI" ]; then \
			echo "Detected CI environment; not starting container. Waiting for SonarQube at ${endpoint}..."; \
		else \
			echo "Starting SonarQube instance using ${container_engine}..."; \
			${container_engine} run -d --name sonargo-sonarqube -p 9000:9000 docker.io/library/sonarqube:${sonarqube_version}; \
			echo "Waiting for SonarQube to be ready..."; \
		fi; \
		until curl -s -u ${username}:${password} ${endpoint}/system/health | grep -q "GREEN"; do \
			printf "."; \
			sleep 5; \
		done; \
		echo "\nSonarQube is ready at ${endpoint}."; \
	fi
