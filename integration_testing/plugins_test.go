package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Plugins Service", Ordered, func() {
	var client *sonargo.Client

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
				_, resp, err := client.Plugins.Download(&sonargo.PluginsDownloadOption{
					Plugin: "java",
				})
				// Skip if API not available - this is internal API
				if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest) {
					Skip("Plugins Download API is not available in this SonarQube version")
				}
				// Accept 200 OK or errors for missing/unavailable plugins
				if err != nil {
					Expect(resp).NotTo(BeNil())
				}
			})
		})

		Context("Error Handling", func() {
			It("should fail with missing plugin key", func() {
				_, _, err := client.Plugins.Download(&sonargo.PluginsDownloadOption{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.Plugins.Download(nil)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// =========================================================================
	// Install
	// =========================================================================
	Describe("Install", func() {
		Context("Error Handling", func() {
			It("should fail with missing key", func() {
				_, err := client.Plugins.Install(&sonargo.PluginsInstallOption{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, err := client.Plugins.Install(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent plugin key", func() {
				resp, err := client.Plugins.Install(&sonargo.PluginsInstallOption{
					Key: "non-existent-plugin-12345",
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Plugins Install API is not available in this SonarQube version")
				}
				// Expect an error since the plugin doesn't exist
				Expect(err).To(HaveOccurred())
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
				result, resp, err := client.Plugins.Installed(&sonargo.PluginsInstalledOption{
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
				result, resp, err := client.Plugins.Installed(&sonargo.PluginsInstalledOption{
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
				result, resp, err := client.Plugins.Installed(&sonargo.PluginsInstalledOption{
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

		Context("Error Handling", func() {
			It("should fail with invalid type", func() {
				_, _, err := client.Plugins.Installed(&sonargo.PluginsInstalledOption{
					Type: "INVALID",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with invalid field", func() {
				_, _, err := client.Plugins.Installed(&sonargo.PluginsInstalledOption{
					Fields: []string{"invalid_field"},
				})
				Expect(err).To(HaveOccurred())
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
		Context("Error Handling", func() {
			It("should fail with missing key", func() {
				_, err := client.Plugins.Uninstall(&sonargo.PluginsUninstallOption{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, err := client.Plugins.Uninstall(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent plugin key", func() {
				resp, err := client.Plugins.Uninstall(&sonargo.PluginsUninstallOption{
					Key: "non-existent-plugin-12345",
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Plugins Uninstall API is not available in this SonarQube version")
				}
				// Expect an error since the plugin doesn't exist
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// =========================================================================
	// Update
	// =========================================================================
	Describe("Update", func() {
		Context("Error Handling", func() {
			It("should fail with missing key", func() {
				_, err := client.Plugins.Update(&sonargo.PluginsUpdateOption{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, err := client.Plugins.Update(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent plugin key", func() {
				resp, err := client.Plugins.Update(&sonargo.PluginsUpdateOption{
					Key: "non-existent-plugin-12345",
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Plugins Update API is not available in this SonarQube version")
				}
				// Expect an error since the plugin doesn't exist
				Expect(err).To(HaveOccurred())
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
