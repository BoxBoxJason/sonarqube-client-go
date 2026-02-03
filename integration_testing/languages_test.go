package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Languages Service", Ordered, func() {
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
	// List
	// =========================================================================
	Describe("List", func() {
		Context("Functional Tests", func() {
			It("should list all languages", func() {
				result, resp, err := client.Languages.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Languages).NotTo(BeEmpty())
			})

			It("should return languages with valid properties", func() {
				result, resp, err := client.Languages.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				for _, lang := range result.Languages {
					Expect(lang.Key).NotTo(BeEmpty())
					Expect(lang.Name).NotTo(BeEmpty())
				}
			})

			It("should list languages with empty options", func() {
				result, resp, err := client.Languages.List(&sonargo.LanguagesListOption{})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Languages).NotTo(BeEmpty())
			})

			It("should filter languages with query", func() {
				result, resp, err := client.Languages.List(&sonargo.LanguagesListOption{
					Query: "java",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// May or may not find results depending on installed plugins
			})

			It("should limit results with page size", func() {
				result, resp, err := client.Languages.List(&sonargo.LanguagesListOption{
					PageSize: 2,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(len(result.Languages)).To(BeNumerically("<=", 2))
			})

			It("should return consistent results on multiple calls", func() {
				result1, resp1, err := client.Languages.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp1.StatusCode).To(Equal(http.StatusOK))

				result2, resp2, err := client.Languages.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp2.StatusCode).To(Equal(http.StatusOK))

				Expect(len(result1.Languages)).To(Equal(len(result2.Languages)))
			})
		})
	})
})
