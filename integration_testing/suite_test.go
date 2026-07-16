package integration_testing_test

import (
	"fmt"
	"strings"
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

	err = helpers.CleanupOrphanedResources(client, 0*time.Second)
	Expect(err).NotTo(HaveOccurred())
})

// AfterSuite performs strict schema validation across every API response
// observed during the run: every e2e client is wired to record any JSON field
// with no corresponding field on its destination Go struct (see
// helpers.RecordSchemaMismatches). Failing here, rather than in the
// individual spec that happened to trigger it, gives one consolidated report
// of every struct that has drifted from the real SonarQube API.
var _ = AfterSuite(func() {
	mismatches := helpers.SchemaMismatches()
	if len(mismatches) == 0 {
		return
	}

	lines := make([]string, 0, len(mismatches))
	for _, mismatch := range mismatches {
		lines = append(lines, mismatch.String())
	}

	Fail(fmt.Sprintf(
		"strict schema validation found %d field(s) in API responses with no match in their Go struct:\n%s",
		len(mismatches), strings.Join(lines, "\n"),
	))
})

func TestIntegrationTesting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SonarQube SDK E2E Test Suite")
}
