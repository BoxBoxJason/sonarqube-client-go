package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("System Service", Ordered, func() {
	var (
		client *sonargo.Client
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	Describe("Health Check Endpoints", func() {
		Describe("Health", func() {
			It("should return health status with GREEN/YELLOW/RED", func() {
				health, resp, err := client.System.Health()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(health).NotTo(BeNil())
				Expect(health.Health).To(BeElementOf("GREEN", "YELLOW", "RED"))
			})

			It("should return nodes information for clustered setup", func() {
				health, _, err := client.System.Health()
				Expect(err).NotTo(HaveOccurred())
				Expect(health).NotTo(BeNil())
				// Nodes may be empty for non-clustered setups
				if len(health.Nodes) > 0 {
					for _, node := range health.Nodes {
						Expect(node.Name).NotTo(BeEmpty())
						Expect(node.Health).To(BeElementOf("GREEN", "YELLOW", "RED"))
						Expect(node.Type).NotTo(BeEmpty())
					}
				}
			})
		})

		Describe("Liveness", func() {
			It("should return 204 when system is alive", func() {
				_, resp, err := client.System.Liveness()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})

		Describe("Ping", func() {
			It("should return 'pong' response", func() {
				pong, resp, err := client.System.Ping()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(pong).NotTo(BeNil())
				Expect(*pong).To(Equal("pong"))
			})
		})

		Describe("Status", func() {
			It("should return current system status", func() {
				status, resp, err := client.System.Status()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(status).NotTo(BeNil())
				Expect(status.Status).To(BeElementOf(
					"STARTING",
					"UP",
					"DOWN",
					"RESTARTING",
					"DB_MIGRATION_NEEDED",
					"DB_MIGRATION_RUNNING",
				))
			})

			It("should include version information when system is UP", func() {
				status, _, err := client.System.Status()
				Expect(err).NotTo(HaveOccurred())
				Expect(status).NotTo(BeNil())
				if status.Status == "UP" {
					Expect(status.Version).NotTo(BeEmpty())
				}
			})
		})
	})

	Describe("System Information Endpoints", func() {
		Describe("Info", func() {
			It("should return detailed system information", func() {
				info, resp, err := client.System.Info()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(info).NotTo(BeNil())
			})

			It("should contain sections for system details", func() {
				info, _, err := client.System.Info()
				Expect(err).NotTo(HaveOccurred())
				Expect(info).NotTo(BeNil())
				// Check for common sections that should exist
				// Info contains nested structs for various sections
				Expect(info.System).NotTo(BeNil())
			})
		})
	})

	Describe("Database Migration", func() {
		Describe("DbMigrationStatus", func() {
			It("should return database migration status", func() {
				status, resp, err := client.System.DbMigrationStatus()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(status).NotTo(BeNil())
				Expect(status.State).To(BeElementOf(
					"NO_MIGRATION",
					"NOT_SUPPORTED",
					"MIGRATION_RUNNING",
					"MIGRATION_SUCCEEDED",
					"MIGRATION_FAILED",
					"MIGRATION_REQUIRED",
				))
			})
		})

		// MigrateDb is dangerous and should NOT be run in e2e tests
		// as it can corrupt the database or cause data loss
		Describe("MigrateDb", func() {
			It("should NOT be called in e2e tests (dangerous operation)", func() {
				Skip("MigrateDb is a dangerous operation that can corrupt the database")
			})
		})
	})

	Describe("Upgrades", func() {
		Describe("Upgrades", func() {
			It("should return upgrade information", func() {
				upgrades, resp, err := client.System.Upgrades()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(upgrades).NotTo(BeNil())
			})
		})
	})

	Describe("Log Management", func() {
		Describe("ChangeLogLevel", func() {
			var originalLevel string

			BeforeEach(func() {
				// Store the original level before changing
				// Default log level is INFO
				originalLevel = "INFO"
			})

			AfterEach(func() {
				// Restore original log level
				_, err := client.System.ChangeLogLevel(&sonargo.SystemChangeLogLevelOption{
					Level: originalLevel,
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("should change log level to DEBUG", func() {
				resp, err := client.System.ChangeLogLevel(&sonargo.SystemChangeLogLevelOption{
					Level: "DEBUG",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})

			It("should change log level to TRACE", func() {
				resp, err := client.System.ChangeLogLevel(&sonargo.SystemChangeLogLevelOption{
					Level: "TRACE",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})

			It("should change log level to INFO", func() {
				resp, err := client.System.ChangeLogLevel(&sonargo.SystemChangeLogLevelOption{
					Level: "INFO",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})

			It("should reject invalid log levels", func() {
				resp, err := client.System.ChangeLogLevel(&sonargo.SystemChangeLogLevelOption{
					Level: "INVALID_LEVEL",
				})
				// This should fail with validation error
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Describe("Logs", func() {
			It("should retrieve app logs", func() {
				logs, resp, err := client.System.Logs(&sonargo.SystemLogsOption{
					Name: "app",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				// Logs can be empty if recently rotated
				Expect(logs).NotTo(BeNil())
			})

			It("should retrieve web logs", func() {
				logs, resp, err := client.System.Logs(&sonargo.SystemLogsOption{
					Name: "web",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(logs).NotTo(BeNil())
			})

			It("should retrieve ce logs", func() {
				logs, resp, err := client.System.Logs(&sonargo.SystemLogsOption{
					Name: "ce",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(logs).NotTo(BeNil())
			})

			It("should reject invalid log names", func() {
				_, resp, err := client.System.Logs(&sonargo.SystemLogsOption{
					Name: "invalid_log_name",
				})
				// This should fail with validation error
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("Restart (DANGEROUS)", func() {
		It("should NOT be called in e2e tests (dangerous operation)", func() {
			Skip("Restart is a dangerous operation that would disrupt the test server")
		})
	})
})
