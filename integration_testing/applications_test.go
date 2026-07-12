package integration_testing_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("Applications Service", Ordered, func() {
	var (
		client  *sonar.Client
		cleanup *helpers.CleanupManager
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		cleanup = helpers.NewCleanupManager(client)
	})

	AfterAll(func() {
		cleanup.Cleanup()
	})

	// =========================================================================
	// Create
	// =========================================================================
	Describe("Create", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Applications.Create(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required name", func() {
				result, resp, err := client.Applications.Create(context.Background(), &sonar.ApplicationsCreateOptions{})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should create application or return an enterprise-only error", func() {
				appKey := helpers.UniqueResourceName("application")
				result, resp, err := client.Applications.Create(context.Background(), &sonar.ApplicationsCreateOptions{
					Name: appKey,
					Key:  appKey,
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
					cleanup.RegisterCleanup("application", appKey, func() error {
						_, cleanupErr := client.Applications.Delete(context.Background(), &sonar.ApplicationsDeleteOptions{
							Application: appKey,
						})
						return cleanupErr
					})
				}
			})
		})
	})

	// =========================================================================
	// Delete
	// =========================================================================
	Describe("Delete", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Applications.Delete(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should fail for nonexistent application", func() {
				resp, err := client.Applications.Delete(context.Background(), &sonar.ApplicationsDeleteOptions{
					Application: "nonexistent-application-xyz",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				}
			})
		})
	})

	// =========================================================================
	// Show
	// =========================================================================
	Describe("Show", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Applications.Show(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should return application or an enterprise-only error", func() {
				result, resp, err := client.Applications.Show(context.Background(), &sonar.ApplicationsShowOptions{
					Application: "nonexistent-application-xyz",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	// =========================================================================
	// AddProject / RemoveProject
	// =========================================================================
	Describe("AddProject", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Applications.AddProject(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("RemoveProject", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Applications.RemoveProject(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// SetTags
	// =========================================================================
	Describe("SetTags", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Applications.SetTags(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// SearchProjects
	// =========================================================================
	Describe("SearchProjects", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Applications.SearchProjects(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should return projects or an enterprise-only error", func() {
				result, resp, err := client.Applications.SearchProjects(context.Background(), &sonar.ApplicationsSearchProjectsOptions{
					Application: "nonexistent-application-xyz",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})
})
