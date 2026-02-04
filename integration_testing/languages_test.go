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
		Context("Parameter Validation", func() {
			It("should succeed with nil options", func() {
				result, resp, err := client.Languages.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Languages).NotTo(BeEmpty())
			})

			It("should succeed with empty options", func() {
				result, resp, err := client.Languages.List(&sonargo.LanguagesListOption{})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Languages).NotTo(BeEmpty())
			})
		})

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

			It("should verify common languages present", func() {
				result, resp, err := client.Languages.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				// Build a map of language keys for easy lookup
				langKeys := make(map[string]bool)
				for _, lang := range result.Languages {
					langKeys[lang.Key] = true
				}

				// Verify common languages are present (as per issue #110)
				commonLanguages := []string{"java", "js", "py", "go"}
				foundCount := 0
				for _, lang := range commonLanguages {
					if langKeys[lang] {
						foundCount++
					}
				}
				// At least some common languages should be present
				Expect(foundCount).To(BeNumerically(">", 0),
					"At least one common language (java, js, py, go) should be present")
			})

			It("should filter languages with query", func() {
				result, resp, err := client.Languages.List(&sonargo.LanguagesListOption{
					Query: "java",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// If results are returned, verify the query parameter is working
				// The filter may return languages where key or name partially matches
				if len(result.Languages) > 0 {
					// At least one result should contain "java" in the key
					foundJava := false
					for _, lang := range result.Languages {
						if lang.Key == "java" {
							foundJava = true
							break
						}
					}
					Expect(foundJava).To(BeTrue(), "Should find 'java' language when filtering with 'java' query")
				}
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
