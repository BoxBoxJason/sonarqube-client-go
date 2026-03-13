package integration_testing_test

import (
	"bytes"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("V2 Analysis Service", Ordered, func() {
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
		errors := cleanup.Cleanup()
		for _, err := range errors {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// GetVersion
	// =========================================================================
	Describe("GetVersion", func() {
		It("should return the scanner engine version", func() {
			version, resp, err := client.V2.Analysis.GetVersion()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(version).NotTo(BeNil())
			Expect(*version).NotTo(BeEmpty())
		})
	})

	// =========================================================================
	// GetJresMetadata
	// =========================================================================
	Describe("GetJresMetadata", func() {
		Context("without options", func() {
			It("should return a list of available JREs", func() {
				jres, resp, err := client.V2.Analysis.GetJresMetadata(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(jres).NotTo(BeEmpty())
				for _, jre := range jres {
					Expect(jre.Id).NotTo(BeEmpty())
					Expect(jre.Filename).NotTo(BeEmpty())
					Expect(jre.Os).NotTo(BeEmpty())
					Expect(jre.Arch).NotTo(BeEmpty())
					Expect(jre.Sha256).NotTo(BeEmpty())
				}
			})
		})

		Context("with OS filter", func() {
			It("should filter JREs by operating system", func() {
				jres, resp, err := client.V2.Analysis.GetJresMetadata(&sonar.AnalysisJresOptions{
					Os: "linux",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				for _, jre := range jres {
					Expect(jre.Os).To(Equal("linux"))
				}
			})
		})

		Context("with OS and arch filter", func() {
			It("should filter JREs by OS and architecture", func() {
				jres, resp, err := client.V2.Analysis.GetJresMetadata(&sonar.AnalysisJresOptions{
					Os:   "linux",
					Arch: "x64",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				for _, jre := range jres {
					Expect(jre.Os).To(Equal("linux"))
					Expect(jre.Arch).To(Equal("x64"))
				}
			})
		})
	})

	// =========================================================================
	// GetJreMetadata
	// =========================================================================
	Describe("GetJreMetadata", func() {
		Context("with valid JRE ID", func() {
			It("should return metadata for a specific JRE", func() {
				jres, resp, err := client.V2.Analysis.GetJresMetadata(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(jres).NotTo(BeEmpty())

				jreID := jres[0].Id
				jre, resp, err := client.V2.Analysis.GetJreMetadata(jreID)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(jre).NotTo(BeNil())
				Expect(jre.Id).To(Equal(jreID))
				Expect(jre.Filename).NotTo(BeEmpty())
				Expect(jre.Sha256).NotTo(BeEmpty())
			})
		})

		Context("parameter validation", func() {
			It("should fail with empty JRE ID", func() {
				jre, resp, err := client.V2.Analysis.GetJreMetadata("")
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(jre).To(BeNil())
			})
		})
	})

	// =========================================================================
	// DownloadJre
	// =========================================================================
	Describe("DownloadJre", func() {
		Context("with valid JRE ID", func() {
			It("should download a JRE binary", func() {
				jres, resp, err := client.V2.Analysis.GetJresMetadata(&sonar.AnalysisJresOptions{
					Os:   "linux",
					Arch: "x64",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				if len(jres) == 0 {
					Skip("No linux/x64 JRE available for download testing")
				}

				var buf bytes.Buffer
				resp, err = client.V2.Analysis.DownloadJre(jres[0].Id, &buf)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(buf.Len()).To(BeNumerically(">", 0))
			})
		})

		Context("parameter validation", func() {
			It("should fail with empty JRE ID", func() {
				var buf bytes.Buffer
				resp, err := client.V2.Analysis.DownloadJre("", &buf)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// GetScannerEngineMetadata
	// =========================================================================
	Describe("GetScannerEngineMetadata", func() {
		It("should return scanner engine metadata", func() {
			engine, resp, err := client.V2.Analysis.GetScannerEngineMetadata()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(engine).NotTo(BeNil())
			Expect(engine.Filename).NotTo(BeEmpty())
			Expect(engine.Sha256).NotTo(BeEmpty())
		})
	})

	// =========================================================================
	// DownloadScannerEngine
	// =========================================================================
	Describe("DownloadScannerEngine", func() {
		It("should download the scanner engine binary", func() {
			var buf bytes.Buffer
			resp, err := client.V2.Analysis.DownloadScannerEngine(&buf)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(buf.Len()).To(BeNumerically(">", 0))
		})
	})

	// =========================================================================
	// GetActiveRules
	// =========================================================================
	Describe("GetActiveRules", func() {
		var projectKey string

		BeforeAll(func() {
			projectKey = helpers.UniqueResourceName("v2arproj")
			_, resp, err := client.Projects.Create(&sonar.ProjectsCreateOption{
				Name:    "V2 Active Rules Test",
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, cleanupErr := client.Projects.Delete(&sonar.ProjectsDeleteOption{
					Project: projectKey,
				})
				return helpers.IgnoreNotFoundError(cleanupErr)
			})
		})

		Context("with valid project key", func() {
			It("should return active rules for the project", func() {
				rules, resp, err := client.V2.Analysis.GetActiveRules(&sonar.AnalysisActiveRuleV2sOptions{
					ProjectKey: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(rules).NotTo(BeNil())
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				rules, resp, err := client.V2.Analysis.GetActiveRules(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(rules).To(BeNil())
			})

			It("should fail with empty project key", func() {
				rules, resp, err := client.V2.Analysis.GetActiveRules(&sonar.AnalysisActiveRuleV2sOptions{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(rules).To(BeNil())
			})
		})
	})
})
