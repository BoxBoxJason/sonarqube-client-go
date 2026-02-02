package integration_testing_test

import (
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("ProjectBadges Service", Ordered, func() {
	var (
		client     *sonargo.Client
		cleanup    *helpers.CleanupManager
		projectKey string
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		cleanup = helpers.NewCleanupManager(client)

		// Create a test project for badge operations
		projectKey = helpers.UniqueResourceName("badge")
		_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
			Name:    "ProjectBadges Test Project",
			Project: projectKey,
		})
		Expect(err).NotTo(HaveOccurred())

		cleanup.RegisterCleanup("project", projectKey, func() error {
			_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
				Project: projectKey,
			})
			return err
		})
	})

	AfterAll(func() {
		errors := cleanup.Cleanup()
		for _, err := range errors {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// Token
	// =========================================================================
	Describe("Token", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.ProjectBadges.Token(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required project", func() {
				result, resp, err := client.ProjectBadges.Token(&sonargo.ProjectBadgesTokenOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Project"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should get token for a project", func() {
				result, resp, err := client.ProjectBadges.Token(&sonargo.ProjectBadgesTokenOption{
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Token).NotTo(BeEmpty())
			})
		})

		Context("Non-Existent Project", func() {
			It("should fail for non-existent project", func() {
				result, resp, err := client.ProjectBadges.Token(&sonargo.ProjectBadgesTokenOption{
					Project: "non-existent-project",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// RenewToken
	// =========================================================================
	Describe("RenewToken", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.ProjectBadges.RenewToken(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required project", func() {
				resp, err := client.ProjectBadges.RenewToken(&sonargo.ProjectBadgesRenewTokenOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Project"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should renew token for a project", func() {
				// Get current token
				tokenBefore, _, err := client.ProjectBadges.Token(&sonargo.ProjectBadgesTokenOption{
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				// Renew token
				resp, err := client.ProjectBadges.RenewToken(&sonargo.ProjectBadgesRenewTokenOption{
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// Get new token
				tokenAfter, _, err := client.ProjectBadges.Token(&sonargo.ProjectBadgesTokenOption{
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(tokenAfter.Token).NotTo(Equal(tokenBefore.Token))
			})
		})

		Context("Non-Existent Project", func() {
			It("should fail for non-existent project", func() {
				resp, err := client.ProjectBadges.RenewToken(&sonargo.ProjectBadgesRenewTokenOption{
					Project: "non-existent-project",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Measure
	// =========================================================================
	Describe("Measure", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.ProjectBadges.Measure(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required project", func() {
				result, resp, err := client.ProjectBadges.Measure(&sonargo.ProjectBadgesMeasureOption{
					Metric: "coverage",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Project"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required metric", func() {
				result, resp, err := client.ProjectBadges.Measure(&sonargo.ProjectBadgesMeasureOption{
					Project: projectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Metric"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid metric", func() {
				result, resp, err := client.ProjectBadges.Measure(&sonargo.ProjectBadgesMeasureOption{
					Project: projectKey,
					Metric:  "invalid_metric",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Metric"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should get coverage badge", func() {
				result, resp, err := client.ProjectBadges.Measure(&sonargo.ProjectBadgesMeasureOption{
					Project: projectKey,
					Metric:  "coverage",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// Badge is SVG content
				Expect(strings.Contains(*result, "<svg")).To(BeTrue())
			})

			It("should get ncloc badge", func() {
				result, resp, err := client.ProjectBadges.Measure(&sonargo.ProjectBadgesMeasureOption{
					Project: projectKey,
					Metric:  "ncloc",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should get alert_status badge", func() {
				result, resp, err := client.ProjectBadges.Measure(&sonargo.ProjectBadgesMeasureOption{
					Project: projectKey,
					Metric:  "alert_status",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Non-Existent Project", func() {
			It("should return an error badge for non-existent project", func() {
				result, resp, err := client.ProjectBadges.Measure(&sonargo.ProjectBadgesMeasureOption{
					Project: "non-existent-project",
					Metric:  "coverage",
				})
				// Badge API may return 200 with an error badge instead of an error
				if err == nil {
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(result).NotTo(BeNil())
				} else {
					if resp != nil {
						Expect(resp.StatusCode).To(BeNumerically(">=", 400))
					}
				}
			})
		})
	})

	// =========================================================================
	// QualityGate
	// =========================================================================
	Describe("QualityGate", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.ProjectBadges.QualityGate(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required project", func() {
				result, resp, err := client.ProjectBadges.QualityGate(&sonargo.ProjectBadgesQualityGateOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Project"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should get quality gate badge", func() {
				result, resp, err := client.ProjectBadges.QualityGate(&sonargo.ProjectBadgesQualityGateOption{
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// Badge is SVG content
				Expect(strings.Contains(*result, "<svg")).To(BeTrue())
			})
		})

		Context("Non-Existent Project", func() {
			It("should return an error badge for non-existent project", func() {
				result, resp, err := client.ProjectBadges.QualityGate(&sonargo.ProjectBadgesQualityGateOption{
					Project: "non-existent-project",
				})
				// Badge API may return 200 with an error badge instead of an error
				if err == nil {
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(result).NotTo(BeNil())
				} else {
					if resp != nil {
						Expect(resp.StatusCode).To(BeNumerically(">=", 400))
					}
				}
			})
		})
	})

	// =========================================================================
	// Full Workflow
	// =========================================================================
	Describe("Full Workflow", func() {
		It("should get token, renew it, and use it for badge access", func() {
			// Get token
			tokenResult, resp, err := client.ProjectBadges.Token(&sonargo.ProjectBadgesTokenOption{
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(tokenResult.Token).NotTo(BeEmpty())

			// Use token to get badge
			result, resp, err := client.ProjectBadges.Measure(&sonargo.ProjectBadgesMeasureOption{
				Project: projectKey,
				Metric:  "coverage",
				Token:   tokenResult.Token,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())

			// Renew token
			resp, err = client.ProjectBadges.RenewToken(&sonargo.ProjectBadgesRenewTokenOption{
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Get new token
			newTokenResult, resp, err := client.ProjectBadges.Token(&sonargo.ProjectBadgesTokenOption{
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(newTokenResult.Token).NotTo(Equal(tokenResult.Token))
		})
	})
})
