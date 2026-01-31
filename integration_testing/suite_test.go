package integration_testing_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIntegrationTesting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SonarQube SDK E2E Test Suite")
}
