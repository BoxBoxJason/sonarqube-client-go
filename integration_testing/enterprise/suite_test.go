// Package enterprise_test contains integration tests that only make sense
// against a real SonarQube Enterprise Edition (or above) instance.
//
// Unlike the specs in integration_testing/, these specs do not tolerate an
// "enterprise-only error" as a passing outcome: they assert the actual
// enterprise behavior. The whole suite is skipped up front (not failed) if
// the connected server is not running Enterprise Edition or above, so it is
// safe to run in any environment via `make e2e.enterprise` — see the
// project Makefile and README for how to point it at a licensed instance.
package enterprise_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
)

var _ = BeforeSuite(func() {
	client, err := helpers.NewDefaultClient()
	Expect(err).NotTo(HaveOccurred())
	Expect(client).NotTo(BeNil())

	edition, err := helpers.GetEdition(client)
	Expect(err).NotTo(HaveOccurred())

	if !helpers.IsEnterpriseOrAbove(edition) {
		Skip(fmt.Sprintf("this suite requires a SonarQube Enterprise Edition (or above) instance, got edition %q", edition))
	}

	err = helpers.CleanupOrphanedResources(client, 0*time.Second)
	Expect(err).NotTo(HaveOccurred())
})

func TestEnterpriseIntegrationTesting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SonarQube SDK Enterprise Edition E2E Test Suite")
}
