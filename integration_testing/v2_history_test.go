package integration_testing_test

import (
	"context"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("V2 History Service", Ordered, func() {
	var client *sonar.Client

	// A syntactically valid (but non-existent) 36-character entity UUID.
	const dummyEntityID = "00000000-0000-0000-0000-000000000000"

	startDate := time.Now().AddDate(0, -1, 0).UTC().Format(time.RFC3339)

	// expectedErrorCodes are the status codes an unlicensed community instance
	// may return for these portfolio/application-scoped endpoints. Listing them
	// explicitly keeps unrelated failures (decode errors, panics) from being
	// silently swallowed.
	expectedErrorCodes := []int{
		http.StatusBadRequest,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusPaymentRequired,
	}

	// entityHistoryErrorCodes additionally tolerates HTTP 500: on the community
	// edition the measures/issue-count history endpoints return a generic
	// server error (rather than 400/404) when the referenced entity UUID does
	// not resolve, because the backing portfolio/governance features are not
	// available in this edition.
	entityHistoryErrorCodes := append([]int{http.StatusInternalServerError}, expectedErrorCodes...)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// GetIssueCountHistory
	// =========================================================================
	Describe("GetIssueCountHistory", func() {
		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.History.GetIssueCountHistory(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required fields", func() {
				result, resp, err := client.V2.History.GetIssueCountHistory(context.Background(), &sonar.HistoryIssueCountHistoryOptions{
					EntityType: "PROJECT_BRANCH",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("functional", func() {
			It("should return issue count history or an expected error", func() {
				result, resp, err := client.V2.History.GetIssueCountHistory(context.Background(), &sonar.HistoryIssueCountHistoryOptions{
					EntityID:   dummyEntityID,
					EntityType: "PROJECT_BRANCH",
					StartDate:  startDate,
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).To(BeElementOf(entityHistoryErrorCodes))
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	// =========================================================================
	// GetMeasuresHistory
	// =========================================================================
	Describe("GetMeasuresHistory", func() {
		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.History.GetMeasuresHistory(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with an invalid entity type", func() {
				result, resp, err := client.V2.History.GetMeasuresHistory(context.Background(), &sonar.HistoryMeasuresHistoryOptions{
					EntityType: "BRANCH",
					EntityID:   dummyEntityID,
					MetricKeys: []string{"coverage"},
					StartDate:  startDate,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("functional", func() {
			It("should return measures history or an expected error", func() {
				result, resp, err := client.V2.History.GetMeasuresHistory(context.Background(), &sonar.HistoryMeasuresHistoryOptions{
					EntityType: sonar.HistoryEntityTypeProjectBranch,
					EntityID:   dummyEntityID,
					MetricKeys: []string{"coverage", "bugs"},
					StartDate:  startDate,
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).To(BeElementOf(entityHistoryErrorCodes))
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	// =========================================================================
	// GetProjectIssueCounts
	// =========================================================================
	Describe("GetProjectIssueCounts", func() {
		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.History.GetProjectIssueCounts(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail when no selector is provided", func() {
				result, resp, err := client.V2.History.GetProjectIssueCounts(context.Background(), &sonar.HistoryProjectIssueCountsOptions{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail when both selector forms are provided", func() {
				result, resp, err := client.V2.History.GetProjectIssueCounts(context.Background(), &sonar.HistoryProjectIssueCountsOptions{
					PortfolioID: "portfolio-1",
					EntityType:  "PORTFOLIO",
					EntityID:    dummyEntityID,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("functional", func() {
			It("should return project issue counts or an expected error", func() {
				result, resp, err := client.V2.History.GetProjectIssueCounts(context.Background(), &sonar.HistoryProjectIssueCountsOptions{
					PortfolioID: helpers.UniqueResourceName("portfolio"),
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).To(BeElementOf(entityHistoryErrorCodes))
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	// =========================================================================
	// GetProjectMeasures
	// =========================================================================
	Describe("GetProjectMeasures", func() {
		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.History.GetProjectMeasures(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without a metric key", func() {
				result, resp, err := client.V2.History.GetProjectMeasures(context.Background(), &sonar.HistoryProjectMeasuresOptions{
					PortfolioID: "portfolio-1",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("MetricKey"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("functional", func() {
			It("should return project measures or an expected error", func() {
				result, resp, err := client.V2.History.GetProjectMeasures(context.Background(), &sonar.HistoryProjectMeasuresOptions{
					MetricKey:   "coverage",
					PortfolioID: helpers.UniqueResourceName("portfolio"),
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).To(BeElementOf(entityHistoryErrorCodes))
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})
})
