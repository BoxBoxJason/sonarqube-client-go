package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Metrics Service", Ordered, func() {
	var (
		client *sonargo.Client
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// Search
	// =========================================================================
	Describe("Search", func() {
		Context("Functional Tests", func() {
			It("should search all metrics with nil options", func() {
				result, resp, err := client.Metrics.Search(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Metrics).NotTo(BeEmpty())
			})

			It("should search metrics with empty options", func() {
				result, resp, err := client.Metrics.Search(&sonargo.MetricsSearchOption{})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Metrics).NotTo(BeEmpty())
			})

			It("should return metrics with valid properties", func() {
				result, resp, err := client.Metrics.Search(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				for _, metric := range result.Metrics {
					Expect(metric.Key).NotTo(BeEmpty())
					Expect(metric.Name).NotTo(BeEmpty())
					Expect(metric.Type).NotTo(BeEmpty())
				}
			})

			It("should find common metrics", func() {
				result, resp, err := client.Metrics.Search(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				// Look for common metrics
				metricKeys := make(map[string]bool)
				for _, metric := range result.Metrics {
					metricKeys[metric.Key] = true
				}

				// At least some common metrics should exist
				commonMetrics := []string{"bugs", "code_smells", "vulnerabilities", "ncloc"}
				foundCount := 0
				for _, key := range commonMetrics {
					if metricKeys[key] {
						foundCount++
					}
				}
				Expect(foundCount).To(BeNumerically(">", 0), "Should find at least one common metric")
			})

			It("should support pagination", func() {
				result, resp, err := client.Metrics.Search(&sonargo.MetricsSearchOption{
					PaginationArgs: sonargo.PaginationArgs{
						PageSize: 5,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(len(result.Metrics)).To(BeNumerically("<=", 5))
			})
		})
	})

	// =========================================================================
	// Types
	// =========================================================================
	Describe("Types", func() {
		Context("Functional Tests", func() {
			It("should list metric types", func() {
				result, resp, err := client.Metrics.Types()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Types).NotTo(BeEmpty())
			})

			It("should return common metric types", func() {
				result, resp, err := client.Metrics.Types()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				// Check for common types
				typeSet := make(map[string]bool)
				for _, t := range result.Types {
					typeSet[t] = true
				}

				// Common metric types: INT, FLOAT, PERCENT, BOOL, STRING, etc.
				Expect(typeSet["INT"]).To(BeTrue(), "Should have INT type")
				Expect(typeSet["FLOAT"]).To(BeTrue(), "Should have FLOAT type")
				Expect(typeSet["PERCENT"]).To(BeTrue(), "Should have PERCENT type")
			})

			It("should return consistent results on multiple calls", func() {
				result1, resp1, err := client.Metrics.Types()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp1.StatusCode).To(Equal(http.StatusOK))

				result2, resp2, err := client.Metrics.Types()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp2.StatusCode).To(Equal(http.StatusOK))

				Expect(len(result1.Types)).To(Equal(len(result2.Types)))
			})
		})
	})
})
