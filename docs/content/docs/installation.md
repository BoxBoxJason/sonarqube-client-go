---
title: Installation
description: "Install sonar-cli with go install or a release binary, and add the sonarqube-client-go Go SDK to your module with go get."
weight: 10
---

## CLI

Install the latest release with Go:

```bash
go install github.com/boxboxjason/sonarqube-client-go/v2/cmd/sonar-cli@latest
```

Or download the binary for your platform from the [releases](https://github.com/BoxBoxJason/sonarqube-client-go/releases) page.

To build it from source:

```bash
git clone https://github.com/BoxBoxJason/sonarqube-client-go.git
cd sonarqube-client-go
make build
./bin/sonar-cli --version
```

`make build version=1.2.3` stamps a specific version into the binary.

Continue with the [CLI guide]({{< relref "cli" >}}).

## Go SDK

Add the module to your project:

```bash
go get github.com/boxboxjason/sonarqube-client-go/v2/sonar
```

and import the `sonar` package:

```go
import "github.com/boxboxjason/sonarqube-client-go/v2/sonar"
```

The required Go version is the one declared in [go.mod](https://github.com/BoxBoxJason/sonarqube-client-go/blob/main/go.mod).

Continue with the [Go SDK guide]({{< relref "sdk" >}}).
