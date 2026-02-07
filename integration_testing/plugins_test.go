package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Plugins Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// Available
	// =========================================================================
	Describe("Available", func() {
		Context("Functional Tests", func() {
			It("should list available plugins", func() {
				result, resp, err := client.Plugins.Available()
				// Skip if API not available or requires marketplace configuration
				if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest) {
					Skip("Plugins Available API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// CancelAll
	// =========================================================================
	Describe("CancelAll", func() {
		Context("Functional Tests", func() {
			It("should cancel all pending plugin operations", func() {
				resp, err := client.Plugins.CancelAll()
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Plugins CancelAll API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(SatisfyAny(Equal(http.StatusOK), Equal(http.StatusNoContent)))
			})
		})
	})

	// =========================================================================
	// Download
	// =========================================================================
	Describe("Download", func() {
		Context("Functional Tests", func() {
			It("should attempt to download a plugin", func() {
				_, resp, err := client.Plugins.Download(&sonar.PluginsDownloadOption{
					Plugin: "java",
				})
				// Skip if API not available - this is internal API
				if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest) {
					Skip("Plugins Download API is not available in this SonarQube version")
				}
				// Assert successful download path: no error and 200 status
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})

		Context("Parameter Validation", func() {
			It("should fail with missing plugin key", func() {
				result, resp, err := client.Plugins.Download(&sonar.PluginsDownloadOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with nil options", func() {
				result, resp, err := client.Plugins.Download(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Install
	// =========================================================================
	Describe("Install", func() {
		Context("Parameter Validation", func() {
			It("should fail with missing key", func() {
				resp, err := client.Plugins.Install(&sonar.PluginsInstallOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with nil options", func() {
				resp, err := client.Plugins.Install(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with non-existent plugin key", func() {
				resp, err := client.Plugins.Install(&sonar.PluginsInstallOption{
					Key: "non-existent-plugin-12345",
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Plugins Install API is not available in this SonarQube version")
				}
				// Expect an error since the plugin doesn't exist
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(BeNumerically(">", 399))
			})
		})
	})

	// =========================================================================
	// Installed
	// =========================================================================
	Describe("Installed", func() {
		Context("Functional Tests", func() {
			It("should list installed plugins with no options", func() {
				result, resp, err := client.Plugins.Installed(nil)
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Plugins Installed API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Plugins).NotTo(BeNil())
			})

			It("should list installed plugins with category field", func() {
				result, resp, err := client.Plugins.Installed(&sonar.PluginsInstalledOption{
					Fields: []string{"category"},
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Plugins Installed API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list bundled plugins only", func() {
				result, resp, err := client.Plugins.Installed(&sonar.PluginsInstalledOption{
					Type: "BUNDLED",
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Plugins Installed API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list external plugins only", func() {
				result, resp, err := client.Plugins.Installed(&sonar.PluginsInstalledOption{
					Type: "EXTERNAL",
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Plugins Installed API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Parameter Validation", func() {
			It("should fail with invalid type", func() {
				result, resp, err := client.Plugins.Installed(&sonar.PluginsInstalledOption{
					Type: "INVALID",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with invalid field", func() {
				result, resp, err := client.Plugins.Installed(&sonar.PluginsInstalledOption{
					Fields: []string{"invalid_field"},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Pending
	// =========================================================================
	Describe("Pending", func() {
		Context("Functional Tests", func() {
			It("should list pending plugin operations", func() {
				result, resp, err := client.Plugins.Pending()
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Plugins Pending API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// Uninstall
	// =========================================================================
	Describe("Uninstall", func() {
		Context("Parameter Validation", func() {
			It("should fail with missing key", func() {
				resp, err := client.Plugins.Uninstall(&sonar.PluginsUninstallOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with nil options", func() {
				resp, err := client.Plugins.Uninstall(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with non-existent plugin key", func() {
				resp, err := client.Plugins.Uninstall(&sonar.PluginsUninstallOption{
					Key: "non-existent-plugin-12345",
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Plugins Uninstall API is not available in this SonarQube version")
				}
				// Expect an error since the plugin doesn't exist
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(BeNumerically(">", 399))
			})
		})
	})

	// =========================================================================
	// Update
	// =========================================================================
	Describe("Update", func() {
		Context("Parameter Validation", func() {
			It("should fail with missing key", func() {
				resp, err := client.Plugins.Update(&sonar.PluginsUpdateOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with nil options", func() {
				resp, err := client.Plugins.Update(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with non-existent plugin key", func() {
				resp, err := client.Plugins.Update(&sonar.PluginsUpdateOption{
					Key: "non-existent-plugin-12345",
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Plugins Update API is not available in this SonarQube version")
				}
				// Expect an error since the plugin doesn't exist
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(BeNumerically(">", 399))
			})
		})
	})

	// =========================================================================
	// Updates
	// =========================================================================
	Describe("Updates", func() {
		Context("Functional Tests", func() {
			It("should list plugins with available updates", func() {
				result, resp, err := client.Plugins.Updates()
				// Skip if API not available or requires marketplace configuration
				if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest) {
					Skip("Plugins Updates API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})
})
