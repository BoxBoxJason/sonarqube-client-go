# SonarQube Go Client

A Go client library for SonarQube API, with automatic code generation from the SonarQube API specification.

## Installation

```bash
go get github.com/boxboxjason/sonarqube-client-go
```

## Usage

Import the generated client in your Go code:

```go
package main

import (
    "fmt"
    "log"

    "github.com/boxboxjason/sonarqube-client-go/sonar"
)

func main() {
    // Create a new SonarQube client
    client, err := sonargo.NewClient("http://your-sonarqube-instance/api", "username", "password")
    if err != nil {
        log.Fatal(err)
    }

    // Use the client to interact with SonarQube API
    // Example: List projects
    projects, resp, err := client.Projects.Search(nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d projects\n", len(projects.Components))
}
```

## Generating the Client Code

If you want to regenerate the client code from your SonarQube instance:

### 1. Get the API specification

```bash
curl -u username:password "http://your-sonarqube-instance:9000/api/webservices/list?include_internals=true" -o ./assets/api.json
```

### 2. Configure the Makefile

Edit the `Makefile` to set your SonarQube endpoint and credentials:

```makefile
endpoint := http://your-sonarqube-instance:9000/api
username := admin
password := admin
```

### 3. Generate the code

```bash
make update
```

This will generate client code in the `sonar/` directory based on your SonarQube API specification.

## Features

- ✅ Automatically generated from SonarQube API specification
- ✅ Type-safe Go structs for all API responses
- ✅ Support for all SonarQube API endpoints
- ✅ Handle different response types (JSON, Protocol Buffers, text)
- ✅ Works with modern Go modules (no GOPATH required)

## Requirements

- Go 1.21 or higher
- Access to a SonarQube instance (for code generation)

## License

See LICENSE file for details.
