package integration_testing_test

import (
	"context"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("Architecture V2 Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// FileGraph
	// =========================================================================
	Describe("FileGraph", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Architecture.FileGraph(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required project key", func() {
				result, resp, err := client.V2.Architecture.FileGraph(context.Background(), &sonar.ArchitectureFileGraphOptions{
					BranchKey: "main",
					Source:    "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ProjectKey"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should return file graph or an enterprise-only error", func() {
				result, resp, err := client.V2.Architecture.FileGraph(context.Background(), &sonar.ArchitectureFileGraphOptions{
					ProjectKey: "nonexistent-project",
					BranchKey:  "main",
					Source:     "java",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
					// Enterprise-only / not-found gate: accept the expected error codes only,
					// so unrelated failures (network errors, 5xx, decoding issues) are not
					// silently swallowed by this test.
					Expect(resp.StatusCode).To(BeElementOf(http.StatusNotFound, http.StatusForbidden, http.StatusPaymentRequired))
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	// =========================================================================
	// Directives
	// =========================================================================
	Describe("ListDirectives", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Architecture.ListDirectives(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without a project id", func() {
				result, resp, err := client.V2.Architecture.ListDirectives(context.Background(), &sonar.ArchitectureListDirectivesOptions{})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should list directives or return an expected error", func() {
				result, resp, err := client.V2.Architecture.ListDirectives(context.Background(), &sonar.ArchitectureListDirectivesOptions{
					ProjectId: "nonexistent-project-arch-v2",
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

	Describe("CreateDirective", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Architecture.CreateDirective(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail with an unknown directive type", func() {
				result, resp, err := client.V2.Architecture.CreateDirective(context.Background(), &sonar.ArchitectureCreateDirectiveOptions{
					ProjectId: "nonexistent-project-arch-v2",
					DirectiveData: sonar.ArchitectureDirectiveData{
						FromNodeKey:    "a",
						GraphType:      "file_graph",
						GraphEcoSystem: "java",
						Type:           "NOT_A_TYPE",
					},
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should create a directive or return an expected error", func() {
				result, resp, err := client.V2.Architecture.CreateDirective(context.Background(), &sonar.ArchitectureCreateDirectiveOptions{
					ProjectId: "nonexistent-project-arch-v2",
					DirectiveData: sonar.ArchitectureDirectiveData{
						FromNodeKey:    "a",
						ToNodeKey:      "b",
						GraphType:      "file_graph",
						GraphEcoSystem: "java",
						Type:           "REMOVE_EDGE",
					},
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

	Describe("GetDirective / DeleteDirective", func() {
		Context("Parameter Validation", func() {
			It("should fail GetDirective with an empty id", func() {
				result, resp, err := client.V2.Architecture.GetDirective(context.Background(), "")
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail DeleteDirective with an empty id", func() {
				resp, err := client.V2.Architecture.DeleteDirective(context.Background(), "")
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should get a directive or return an expected error", func() {
				result, resp, err := client.V2.Architecture.GetDirective(context.Background(), "00000000-0000-0000-0000-000000000000")
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})

			It("should delete a directive or return an expected error", func() {
				resp, err := client.V2.Architecture.DeleteDirective(context.Background(), "00000000-0000-0000-0000-000000000000")
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
				}
			})
		})
	})

	// =========================================================================
	// Intended architecture
	// =========================================================================
	Describe("GetIntendedArchitecture", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Architecture.GetIntendedArchitecture(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should export the intended architecture or return an expected error", func() {
				result, resp, err := client.V2.Architecture.GetIntendedArchitecture(context.Background(), &sonar.ArchitectureIntendedOptions{
					ProjectId: "nonexistent-project-arch-v2",
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
	// Models
	// =========================================================================
	Describe("ListModels / CreateModel", func() {
		Context("Parameter Validation", func() {
			It("should fail ListModels with nil options", func() {
				result, resp, err := client.V2.Architecture.ListModels(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail CreateModel without a project id", func() {
				result, resp, err := client.V2.Architecture.CreateModel(context.Background(), &sonar.ArchitectureCreateModelOptions{})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should list models or return an expected error", func() {
				result, resp, err := client.V2.Architecture.ListModels(context.Background(), &sonar.ArchitectureModelsListOptions{
					ProjectId: "nonexistent-project-arch-v2",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})

			It("should create a model or return an expected error", func() {
				result, resp, err := client.V2.Architecture.CreateModel(context.Background(), &sonar.ArchitectureCreateModelOptions{
					ProjectId: "nonexistent-project-arch-v2",
					Model: &sonar.ArchitectureModelData{
						Perspectives: []sonar.ArchitectureModelPerspective{{Qualifiers: "file", Language: "java"}},
					},
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

	Describe("GetModel / UpdateModel / DeleteModel", func() {
		Context("Parameter Validation", func() {
			It("should fail GetModel with an empty id", func() {
				result, resp, err := client.V2.Architecture.GetModel(context.Background(), "")
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail UpdateModel with a nil body", func() {
				result, resp, err := client.V2.Architecture.UpdateModel(context.Background(), "00000000-0000-0000-0000-000000000000", nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail DeleteModel with an empty id", func() {
				resp, err := client.V2.Architecture.DeleteModel(context.Background(), "")
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should get a model or return an expected error", func() {
				result, resp, err := client.V2.Architecture.GetModel(context.Background(), "00000000-0000-0000-0000-000000000000")
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})

			It("should update a model or return an expected error", func() {
				result, resp, err := client.V2.Architecture.UpdateModel(context.Background(), "00000000-0000-0000-0000-000000000000", &sonar.ArchitectureUpdateModelOptions{
					Data: "updated",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})

			It("should delete a model or return an expected error", func() {
				resp, err := client.V2.Architecture.DeleteModel(context.Background(), "00000000-0000-0000-0000-000000000000")
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
				}
			})
		})
	})

	// =========================================================================
	// Project configurations
	// =========================================================================
	Describe("GetProjectConfigurations", func() {
		Context("Functional Tests", func() {
			It("should return project configurations or an expected error", func() {
				result, resp, err := client.V2.Architecture.GetProjectConfigurations(context.Background(), &sonar.ArchitectureProjectConfigurationsOptions{
					ProjectKey: "nonexistent-project-arch-v2",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})

			It("should return private project configurations or an expected error", func() {
				result, resp, err := client.V2.Architecture.GetPrivateProjectConfigurations(context.Background(), &sonar.ArchitectureProjectConfigurationsOptions{
					ProjectKey: "nonexistent-project-arch-v2",
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
	// Snapshots
	// =========================================================================
	Describe("ListSnapshots", func() {
		Context("Parameter Validation", func() {
			It("should fail without required identifiers", func() {
				result, resp, err := client.V2.Architecture.ListSnapshots(context.Background(), &sonar.ArchitectureListSnapshotsOptions{
					OrganizationId: "o1",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should list snapshots or return an expected error", func() {
				result, resp, err := client.V2.Architecture.ListSnapshots(context.Background(), &sonar.ArchitectureListSnapshotsOptions{
					OrganizationId: "nonexistent-organization",
					BranchId:       "nonexistent-branch",
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
	// Structure
	// =========================================================================
	Describe("GetStructure", func() {
		Context("Parameter Validation", func() {
			It("should fail with an unknown ecosystem", func() {
				result, resp, err := client.V2.Architecture.GetStructure(context.Background(), &sonar.ArchitectureStructureOptions{
					OrganizationId: "o1",
					BranchId:       "b1",
					Ecosystem:      "cobol",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should return a structure graph or an expected error", func() {
				result, resp, err := client.V2.Architecture.GetStructure(context.Background(), &sonar.ArchitectureStructureOptions{
					OrganizationId: "nonexistent-organization",
					BranchId:       "nonexistent-branch",
					Ecosystem:      "java",
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
	// Analysis and scanner data
	// =========================================================================
	Describe("StoreAnalysis", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Architecture.StoreAnalysis(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should store analysis data or return an expected error", func() {
				result, resp, err := client.V2.Architecture.StoreAnalysis(context.Background(), &sonar.ArchitectureStoreAnalysisOptions{
					CeTaskId: "nonexistent-ce-task",
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

	Describe("InitiateScannerDataUpload / UploadScannerData", func() {
		Context("Parameter Validation", func() {
			It("should fail InitiateScannerDataUpload without an analysis id", func() {
				result, resp, err := client.V2.Architecture.InitiateScannerDataUpload(context.Background(), &sonar.ArchitectureScannerDataOptions{})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail UploadScannerData without a payload", func() {
				resp, err := client.V2.Architecture.UploadScannerData(context.Background(), &sonar.ArchitectureUploadScannerDataOptions{
					CeTaskId: "ce-1",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should initiate a scanner-data upload or return an expected error", func() {
				result, resp, err := client.V2.Architecture.InitiateScannerDataUpload(context.Background(), &sonar.ArchitectureScannerDataOptions{
					AnalysisId: "nonexistent-analysis",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})

			It("should upload scanner data or return an expected error", func() {
				resp, err := client.V2.Architecture.UploadScannerData(context.Background(), &sonar.ArchitectureUploadScannerDataOptions{
					CeTaskId: "nonexistent-ce-task",
					Payload:  strings.NewReader("scanner-report"),
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
				}
			})
		})
	})
})
