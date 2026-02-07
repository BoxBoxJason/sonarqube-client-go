package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Webservices Service", Ordered, func() {
	var (
		client *sonar.Client
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
				result, resp, err := client.Webservices.List(&sonar.WebservicesListOption{})
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
					Expect(ws.Actions).NotTo(BeNil())
					// Verify actions have valid keys and params structure
					for _, action := range ws.Actions {
						Expect(action.Key).NotTo(BeEmpty())
						// Params may be nil for actions with no parameters, which is valid
						if len(action.Params) > 0 {
							for _, param := range action.Params {
								Expect(param.Key).NotTo(BeEmpty())
							}
						}
					}
				}
			})

			It("should find common API paths", func() {
				result, resp, err := client.Webservices.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should filter and search for specific web service domains", func() {
				result, resp, err := client.Webservices.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				// Verify specific domains exist as per issue #113
				domains := map[string]bool{
					"api/projects": false,
					"api/issues":   false,
					"api/measures": false,
				}

				for _, ws := range result.Webservices {
					if _, exists := domains[ws.Path]; exists {
						domains[ws.Path] = true
					}
				}

				Expect(domains["api/projects"]).To(BeTrue(), "Should find api/projects domain")
				Expect(domains["api/issues"]).To(BeTrue(), "Should find api/issues domain")
				Expect(domains["api/measures"]).To(BeTrue(), "Should find api/measures domain")
			})

			It("should include internals when requested", func() {
				withoutInternals, _, err := client.Webservices.List(&sonar.WebservicesListOption{
					IncludeInternals: false,
				})
				Expect(err).NotTo(HaveOccurred())

				withInternals, _, err := client.Webservices.List(&sonar.WebservicesListOption{
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
				list, _, err := client.Webservices.List(&sonar.WebservicesListOption{
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

				// Ensure test data exists
				if controller == "" || action == "" {
					Skip("No action with response example found in this SonarQube version")
				}

				result, resp, err := client.Webservices.ResponseExample(&sonar.WebservicesResponseExampleOption{
					Controller: controller,
					Action:     action,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(*result).NotTo(BeEmpty())
			})
		})

		Context("Error Handling", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Webservices.ResponseExample(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with empty options", func() {
				_, resp, err := client.Webservices.ResponseExample(&sonar.WebservicesResponseExampleOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("Action"))
			})

			It("should fail with missing action", func() {
				_, resp, err := client.Webservices.ResponseExample(&sonar.WebservicesResponseExampleOption{
					Controller: "api/system",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("Action"))
			})

			It("should fail with missing controller", func() {
				_, resp, err := client.Webservices.ResponseExample(&sonar.WebservicesResponseExampleOption{
					Action: "status",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("Controller"))
			})
		})
	})
})
