package integration_testing_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = BeforeSuite(func() {
	client, err := helpers.NewDefaultClient()
	Expect(err).NotTo(HaveOccurred())
	Expect(client).NotTo(BeNil())

	err = helpers.CleanupOrphanedResources(client, 0*time.Second)
	Expect(err).NotTo(HaveOccurred())
})

func TestIntegrationTesting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SonarQube SDK E2E Test Suite")
}
