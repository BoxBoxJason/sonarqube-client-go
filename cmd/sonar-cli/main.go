// Package main provides the entry point for the sonar-cli command-line tool.
// sonar-cli is a CLI wrapper around the SonarQube Go SDK, providing access to all
// SonarQube API endpoints from the command line.
//
// Install:
//
//	go install github.com/boxboxjason/sonarqube-client-go/cmd/sonar-cli@latest
//
// Usage:
//
//	sonar-cli [global flags] <service> <method> [flags]
package main

import (
	"os"

	"github.com/boxboxjason/sonarqube-client-go/internal/cli"
)

func main() {
	err := cli.Execute()
	if err != nil {
		os.Exit(1)
	}
}
