package integration_testing_test

import (
	"context"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("V2 Fix Suggestions Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		if os.Getenv("SONAR_FIX_SUGGESTIONS_E2E") == "" {
			Skip("set SONAR_FIX_SUGGESTIONS_E2E=1 to run Enterprise fix suggestions tests")
		}

		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	Describe("GetServiceInfo", func() {
		It("should return fix suggestions service info", func() {
			result, resp, err := client.V2.FixSuggestions.GetServiceInfo(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Status).NotTo(BeEmpty())
		})
	})

	Describe("ListSupportedLLMProviders", func() {
		It("should return supported providers", func() {
			result, resp, err := client.V2.FixSuggestions.ListSupportedLLMProviders(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})
	})
})
