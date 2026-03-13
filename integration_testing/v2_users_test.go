package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("V2 Users Management Service", Ordered, func() {
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
	// Search
	// =========================================================================
	Describe("Search", func() {
		Context("without options", func() {
			It("should return a list of users", func() {
				result, resp, err := client.V2.UsersManagement.Search(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Users).NotTo(BeEmpty())
				Expect(result.Page.PageSize).To(BeNumerically(">", 0))
			})
		})

		Context("with query filter", func() {
			It("should filter users by query", func() {
				result, resp, err := client.V2.UsersManagement.Search(&sonar.UsersSearchOptionV2{
					Query: "admin",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				found := false
				for _, u := range result.Users {
					if u.Login == "admin" {
						found = true
						break
					}
				}
				Expect(found).To(BeTrue())
			})
		})

		Context("with pagination", func() {
			It("should respect page size", func() {
				result, resp, err := client.V2.UsersManagement.Search(&sonar.UsersSearchOptionV2{
					PaginationParamsV2: sonar.PaginationParamsV2{
						PageSize: 1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(len(result.Users)).To(BeNumerically("<=", 1))
				Expect(result.Page.PageSize).To(BeNumerically("==", 1))
			})
		})

		Context("with active filter", func() {
			It("should filter by active status", func() {
				active := true
				result, resp, err := client.V2.UsersManagement.Search(&sonar.UsersSearchOptionV2{
					Active: &active,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				for _, u := range result.Users {
					Expect(u.Active).To(BeTrue())
				}
			})
		})
	})

	// =========================================================================
	// Create
	// =========================================================================
	Describe("Create", func() {
		Context("with valid parameters", func() {
			It("should create a local user", func() {
				login := helpers.UniqueResourceName("v2user")
				user, resp, err := client.V2.UsersManagement.Create(&sonar.UsersCreateOptionsV2{
					Login:    login,
					Name:     "V2 Test User",
					Password: "testPassword123!",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(user).NotTo(BeNil())
				Expect(user.Login).To(Equal(login))
				Expect(user.Name).To(Equal("V2 Test User"))
				Expect(user.Active).To(BeTrue())
				Expect(user.Local).To(BeTrue())
				Expect(user.Id).NotTo(BeEmpty())

				cleanup.RegisterCleanup("v2-user", user.Id, func() error {
					_, cleanupErr := client.V2.UsersManagement.Deactivate(&sonar.UsersDeactivateOptionsV2{
						Id:        user.Id,
						Anonymize: true,
					})
					return helpers.IgnoreNotFoundError(cleanupErr)
				})
			})

			It("should create a user with email and SCM accounts", func() {
				login := helpers.UniqueResourceName("v2usr-full")
				user, resp, err := client.V2.UsersManagement.Create(&sonar.UsersCreateOptionsV2{
					Login:       login,
					Name:        "V2 Full User",
					Password:    "testPassword123!",
					Email:       "v2test@example.com",
					ScmAccounts: []string{"github-account", "gitlab-account"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(user).NotTo(BeNil())
				Expect(user.Email).To(Equal("v2test@example.com"))
				Expect(user.ScmAccounts).To(ConsistOf("github-account", "gitlab-account"))

				cleanup.RegisterCleanup("v2-user", user.Id, func() error {
					_, cleanupErr := client.V2.UsersManagement.Deactivate(&sonar.UsersDeactivateOptionsV2{
						Id:        user.Id,
						Anonymize: true,
					})
					return helpers.IgnoreNotFoundError(cleanupErr)
				})
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil request", func() {
				user, resp, err := client.V2.UsersManagement.Create(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(user).To(BeNil())
			})

			It("should fail with empty login", func() {
				user, resp, err := client.V2.UsersManagement.Create(&sonar.UsersCreateOptionsV2{
					Name:     "No Login User",
					Password: "testPassword123!",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(user).To(BeNil())
			})

			It("should fail with empty name", func() {
				user, resp, err := client.V2.UsersManagement.Create(&sonar.UsersCreateOptionsV2{
					Login:    helpers.UniqueResourceName("v2user"),
					Password: "testPassword123!",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(user).To(BeNil())
			})

			It("should fail with login too short", func() {
				user, resp, err := client.V2.UsersManagement.Create(&sonar.UsersCreateOptionsV2{
					Login:    "a",
					Name:     "Short Login User",
					Password: "testPassword123!",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(user).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Fetch
	// =========================================================================
	Describe("Fetch", func() {
		var createdUser *sonar.UserV2

		BeforeAll(func() {
			login := helpers.UniqueResourceName("v2fetch")
			var resp *http.Response
			var err error
			createdUser, resp, err = client.V2.UsersManagement.Create(&sonar.UsersCreateOptionsV2{
				Login:    login,
				Name:     "V2 Fetch Test User",
				Password: "testPassword123!",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			cleanup.RegisterCleanup("v2-user-fetch", createdUser.Id, func() error {
				_, cleanupErr := client.V2.UsersManagement.Deactivate(&sonar.UsersDeactivateOptionsV2{
					Id:        createdUser.Id,
					Anonymize: true,
				})
				return helpers.IgnoreNotFoundError(cleanupErr)
			})
		})

		Context("with valid user ID", func() {
			It("should fetch the user by ID", func() {
				fetched, resp, err := client.V2.UsersManagement.Fetch(createdUser.Id)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(fetched).NotTo(BeNil())
				Expect(fetched.Id).To(Equal(createdUser.Id))
				Expect(fetched.Login).To(Equal(createdUser.Login))
				Expect(fetched.Name).To(Equal("V2 Fetch Test User"))
				Expect(fetched.Active).To(BeTrue())
			})
		})

		Context("parameter validation", func() {
			It("should fail with empty user ID", func() {
				user, resp, err := client.V2.UsersManagement.Fetch("")
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(user).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Update
	// =========================================================================
	Describe("Update", func() {
		var createdUser *sonar.UserV2

		BeforeAll(func() {
			login := helpers.UniqueResourceName("v2upd")
			var resp *http.Response
			var err error
			createdUser, resp, err = client.V2.UsersManagement.Create(&sonar.UsersCreateOptionsV2{
				Login:    login,
				Name:     "V2 Update Original",
				Password: "testPassword123!",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			cleanup.RegisterCleanup("v2-user-update", createdUser.Id, func() error {
				_, cleanupErr := client.V2.UsersManagement.Deactivate(&sonar.UsersDeactivateOptionsV2{
					Id:        createdUser.Id,
					Anonymize: true,
				})
				return helpers.IgnoreNotFoundError(cleanupErr)
			})
		})

		Context("with valid parameters", func() {
			It("should update user name", func() {
				updated, resp, err := client.V2.UsersManagement.Update(createdUser.Id, &sonar.UsersUpdateOptionsV2{
					Name: "V2 Updated Name",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(updated).NotTo(BeNil())
				Expect(updated.Name).To(Equal("V2 Updated Name"))
			})

			It("should update user email", func() {
				updated, resp, err := client.V2.UsersManagement.Update(createdUser.Id, &sonar.UsersUpdateOptionsV2{
					Email: "v2-updated@example.com",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(updated).NotTo(BeNil())
				Expect(updated.Email).To(Equal("v2-updated@example.com"))
			})

			It("should update SCM accounts", func() {
				updated, resp, err := client.V2.UsersManagement.Update(createdUser.Id, &sonar.UsersUpdateOptionsV2{
					ScmAccounts: &sonar.UpdateFieldListStringV2{
						Value:   []string{"new-github", "new-gitlab"},
						Defined: true,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(updated).NotTo(BeNil())
				Expect(updated.ScmAccounts).To(ConsistOf("new-github", "new-gitlab"))
			})
		})

		Context("parameter validation", func() {
			It("should fail with empty user ID", func() {
				user, resp, err := client.V2.UsersManagement.Update("", &sonar.UsersUpdateOptionsV2{
					Name: "Updated Name",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(user).To(BeNil())
			})

			It("should fail with nil request", func() {
				user, resp, err := client.V2.UsersManagement.Update(createdUser.Id, nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(user).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Deactivate
	// =========================================================================
	Describe("Deactivate", func() {
		Context("with valid parameters", func() {
			It("should deactivate a user", func() {
				login := helpers.UniqueResourceName("v2deact")
				created, resp, err := client.V2.UsersManagement.Create(&sonar.UsersCreateOptionsV2{
					Login:    login,
					Name:     "V2 Deactivate Test",
					Password: "testPassword123!",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				cleanup.RegisterCleanup("v2-user-deact", created.Id, func() error {
					_, cleanupErr := client.V2.UsersManagement.Deactivate(&sonar.UsersDeactivateOptionsV2{
						Id:        created.Id,
						Anonymize: true,
					})
					return helpers.IgnoreNotFoundError(cleanupErr)
				})

				resp, err = client.V2.UsersManagement.Deactivate(&sonar.UsersDeactivateOptionsV2{
					Id: created.Id,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				fetched, resp, err := client.V2.UsersManagement.Fetch(created.Id)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(fetched).NotTo(BeNil())
				Expect(fetched.Active).To(BeFalse())
			})

			It("should deactivate a user with anonymize", func() {
				login := helpers.UniqueResourceName("v2anon")
				created, resp, err := client.V2.UsersManagement.Create(&sonar.UsersCreateOptionsV2{
					Login:    login,
					Name:     "V2 Anonymize Test",
					Password: "testPassword123!",
					Email:    "v2anon@example.com",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				resp, err = client.V2.UsersManagement.Deactivate(&sonar.UsersDeactivateOptionsV2{
					Id:        created.Id,
					Anonymize: true,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.V2.UsersManagement.Deactivate(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with empty user ID", func() {
				resp, err := client.V2.UsersManagement.Deactivate(&sonar.UsersDeactivateOptionsV2{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})
})
