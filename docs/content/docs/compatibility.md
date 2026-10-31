---
title: Compatibility
description: "SonarQube Server versions supported by sonarqube-client-go, Go version requirements, and how the SDK is tested against real SonarQube instances."
weight: 60
---

## SonarQube versions

The SDK targets the current SonarQube Server releases: version 25 or newer is recommended. Endpoints and fields that only exist on recent versions are documented as such in the [API reference](https://pkg.go.dev/github.com/boxboxjason/sonarqube-client-go/v2/sonar).

Every change is tested against real SonarQube instances, with the versions pinned in the [Makefile](https://github.com/BoxBoxJason/sonarqube-client-go/blob/main/Makefile) (`sonarqube_version` and `sonarqube_enterprise_version`), which Renovate keeps up to date.

## Editions

- **Community Build**: every service available on that edition is covered, and tested by the integration tests.
- **Developer, Enterprise and Data Center editions**: the services specific to those editions (applications, portfolios through `Views`, audit logs, ...) are covered too, and tested against an Enterprise Edition instance. Calling them on an edition that does not have them returns an error.

## Go versions

The module requires the Go version declared in its [go.mod](https://github.com/BoxBoxJason/sonarqube-client-go/blob/main/go.mod).
