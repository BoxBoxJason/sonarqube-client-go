package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Webservices Service", Ordered, func() {
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
			It("should list webservices with nil options", func() {
				result, resp, err := client.Webservices.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Webservices).NotTo(BeEmpty())
			})

			It("should list webservices with empty options", func() {
				result, resp, err := client.Webservices.List(&sonargo.WebservicesListOption{})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Webservices).NotTo(BeEmpty())
			})

			It("should return webservices with valid properties", func() {
				result, resp, err := client.Webservices.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				for _, ws := range result.Webservices {
					Expect(ws.Path).NotTo(BeEmpty())
				}
			})

			It("should find common API paths", func() {
				result, resp, err := client.Webservices.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				paths := make(map[string]bool)
				for _, ws := range result.Webservices {
					paths[ws.Path] = true
				}

				// Common APIs that should exist
				Expect(paths["api/system"]).To(BeTrue(), "Should have api/system")
				Expect(paths["api/projects"]).To(BeTrue(), "Should have api/projects")
			})

			It("should include internals when requested", func() {
				withoutInternals, _, err := client.Webservices.List(&sonargo.WebservicesListOption{
					IncludeInternals: false,
				})
				Expect(err).NotTo(HaveOccurred())

				withInternals, _, err := client.Webservices.List(&sonargo.WebservicesListOption{
					IncludeInternals: true,
				})
				Expect(err).NotTo(HaveOccurred())

				// With internals should have at least as many
				totalActionsWithout := 0
				for _, ws := range withoutInternals.Webservices {
					totalActionsWithout += len(ws.Actions)
				}

				totalActionsWith := 0
				for _, ws := range withInternals.Webservices {
					totalActionsWith += len(ws.Actions)
				}

				Expect(totalActionsWith).To(BeNumerically(">=", totalActionsWithout))
			})
		})
	})

	// =========================================================================
	// ResponseExample
	// =========================================================================
	Describe("ResponseExample", func() {
		Context("Functional Tests", func() {
			It("should get response example for an action with example", func() {
				// First find an action with a response example
				list, _, err := client.Webservices.List(&sonargo.WebservicesListOption{
					IncludeInternals: true,
				})
				Expect(err).NotTo(HaveOccurred())

				var controller, action string
				for _, ws := range list.Webservices {
					for _, a := range ws.Actions {
						if a.HasResponseExample {
							controller = ws.Path
							action = a.Key
							break
						}
					}
					if controller != "" {
						break
					}
				}

				if controller != "" && action != "" {
					result, resp, err := client.Webservices.ResponseExample(&sonargo.WebservicesResponseExampleOption{
						Controller: controller,
						Action:     action,
					})
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(result).NotTo(BeNil())
					Expect(*result).NotTo(BeEmpty())
				}
			})
		})

		Context("Error Handling", func() {
			It("should fail with missing action", func() {
				_, _, err := client.Webservices.ResponseExample(&sonargo.WebservicesResponseExampleOption{
					Controller: "api/system",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with missing controller", func() {
				_, _, err := client.Webservices.ResponseExample(&sonargo.WebservicesResponseExampleOption{
					Action: "status",
				})
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
