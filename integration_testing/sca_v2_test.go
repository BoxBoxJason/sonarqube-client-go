package integration_testing_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("SCA V2 Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	Describe("ListClis", func() {
		Context("Functional Tests", func() {
			It("should list CLIs or return an expected error", func() {
				result, resp, err := client.V2.Sca.ListClis(context.Background(), nil)
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	Describe("GetEnablement", func() {
		Context("Functional Tests", func() {
			It("should return enablement or an expected error", func() {
				result, resp, err := client.V2.Sca.GetEnablement(context.Background())
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	Describe("SetEnablement", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Sca.SetEnablement(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})
		Context("Functional Tests", func() {
			It("should set enablement or return an expected error", func() {
				result, resp, err := client.V2.Sca.SetEnablement(context.Background(), &sonar.ScaSetEnablementOptions{
					Enablement: false,
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

	Describe("SearchDependencyRisks", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Sca.SearchDependencyRisks(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})
		Context("Functional Tests", func() {
			It("should search or return an expected error", func() {
				result, resp, err := client.V2.Sca.SearchDependencyRisks(context.Background(), &sonar.ScaDependencyRisksSearchOptions{
					ProjectKey: "nonexistent-project",
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

	Describe("SearchReleases", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Sca.SearchReleases(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})
		Context("Functional Tests", func() {
			It("should search or return an expected error", func() {
				result, resp, err := client.V2.Sca.SearchReleases(context.Background(), &sonar.ScaReleasesSearchOptions{
					ProjectKey: "nonexistent-project",
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

	Describe("GetSbomReport", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Sca.GetSbomReport(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})
		Context("Functional Tests", func() {
			It("should return SBOM or an expected error", func() {
				result, resp, err := client.V2.Sca.GetSbomReport(context.Background(), &sonar.ScaSbomReportOptions{
					ProjectKey: "nonexistent-project",
					Type:       sonar.ScaSbomReportTypeCycloneDX,
					Format:     sonar.ScaSbomReportFormatJSON,
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
