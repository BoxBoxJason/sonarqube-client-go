package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("V2 Authorizations Service", Ordered, func() {
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
	// Groups
	// =========================================================================
	Describe("SearchGroups", func() {
		Context("without options", func() {
			It("should return a list of groups", func() {
				result, resp, err := client.V2.Authorizations.SearchGroups(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Groups).NotTo(BeEmpty())
				Expect(result.Page.PageSize).To(BeNumerically(">", 0))
			})
		})

		Context("with query filter", func() {
			It("should filter groups by query", func() {
				groupName := helpers.UniqueResourceName("v2srchg")
				group, resp, err := client.V2.Authorizations.CreateGroup(&sonar.AuthorizationsCreateGroupOptions{
					Name:        groupName,
					Description: "Search test group",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				cleanup.RegisterCleanup("v2-group-search", group.Id, func() error {
					_, cleanupErr := client.V2.Authorizations.DeleteGroup(group.Id)
					return helpers.IgnoreNotFoundError(cleanupErr)
				})

				result, resp, err := client.V2.Authorizations.SearchGroups(&sonar.AuthorizationsSearchGroupsOptions{
					Query: groupName,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				found := false
				for _, g := range result.Groups {
					if g.Name == groupName {
						found = true
						break
					}
				}
				Expect(found).To(BeTrue())
			})
		})

		Context("with pagination", func() {
			It("should respect page size", func() {
				result, resp, err := client.V2.Authorizations.SearchGroups(&sonar.AuthorizationsSearchGroupsOptions{
					PaginationParamsV2: sonar.PaginationParamsV2{
						PageSize: 1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(len(result.Groups)).To(BeNumerically("<=", 1))
				Expect(result.Page.PageSize).To(BeNumerically("==", 1))
			})
		})
	})

	Describe("CreateGroup", func() {
		Context("with valid parameters", func() {
			It("should create a group", func() {
				groupName := helpers.UniqueResourceName("v2grp")
				group, resp, err := client.V2.Authorizations.CreateGroup(&sonar.AuthorizationsCreateGroupOptions{
					Name:        groupName,
					Description: "V2 integration test group",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				Expect(group).NotTo(BeNil())
				Expect(group.Name).To(Equal(groupName))
				Expect(group.Description).To(Equal("V2 integration test group"))
				Expect(group.Id).NotTo(BeEmpty())

				cleanup.RegisterCleanup("v2-group-create", group.Id, func() error {
					_, cleanupErr := client.V2.Authorizations.DeleteGroup(group.Id)
					return helpers.IgnoreNotFoundError(cleanupErr)
				})
			})

			It("should create a group with name only", func() {
				groupName := helpers.UniqueResourceName("v2grpmin")
				group, resp, err := client.V2.Authorizations.CreateGroup(&sonar.AuthorizationsCreateGroupOptions{
					Name: groupName,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				Expect(group).NotTo(BeNil())
				Expect(group.Name).To(Equal(groupName))
				Expect(group.Id).NotTo(BeEmpty())

				cleanup.RegisterCleanup("v2-group-min", group.Id, func() error {
					_, cleanupErr := client.V2.Authorizations.DeleteGroup(group.Id)
					return helpers.IgnoreNotFoundError(cleanupErr)
				})
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil request", func() {
				group, resp, err := client.V2.Authorizations.CreateGroup(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(group).To(BeNil())
			})

			It("should fail with empty name", func() {
				group, resp, err := client.V2.Authorizations.CreateGroup(&sonar.AuthorizationsCreateGroupOptions{
					Description: "No Name Group",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(group).To(BeNil())
			})
		})
	})

	Describe("FetchGroup", func() {
		var createdGroup *sonar.Group

		BeforeAll(func() {
			groupName := helpers.UniqueResourceName("v2gfetch")
			var resp *http.Response
			var err error
			createdGroup, resp, err = client.V2.Authorizations.CreateGroup(&sonar.AuthorizationsCreateGroupOptions{
				Name:        groupName,
				Description: "Fetch test group",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			cleanup.RegisterCleanup("v2-group-fetch", createdGroup.Id, func() error {
				_, cleanupErr := client.V2.Authorizations.DeleteGroup(createdGroup.Id)
				return helpers.IgnoreNotFoundError(cleanupErr)
			})
		})

		Context("with valid group ID", func() {
			It("should fetch the group by ID", func() {
				fetched, resp, err := client.V2.Authorizations.FetchGroup(createdGroup.Id)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(fetched).NotTo(BeNil())
				Expect(fetched.Id).To(Equal(createdGroup.Id))
				Expect(fetched.Name).To(Equal(createdGroup.Name))
				Expect(fetched.Description).To(Equal("Fetch test group"))
			})
		})

		Context("parameter validation", func() {
			It("should fail with empty group ID", func() {
				group, resp, err := client.V2.Authorizations.FetchGroup("")
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(group).To(BeNil())
			})
		})
	})

	Describe("UpdateGroup", func() {
		var createdGroup *sonar.Group

		BeforeAll(func() {
			groupName := helpers.UniqueResourceName("v2gupd")
			var resp *http.Response
			var err error
			createdGroup, resp, err = client.V2.Authorizations.CreateGroup(&sonar.AuthorizationsCreateGroupOptions{
				Name:        groupName,
				Description: "Original description",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			cleanup.RegisterCleanup("v2-group-update", createdGroup.Id, func() error {
				_, cleanupErr := client.V2.Authorizations.DeleteGroup(createdGroup.Id)
				return helpers.IgnoreNotFoundError(cleanupErr)
			})
		})

		Context("with valid parameters", func() {
			It("should update group name", func() {
				newName := helpers.UniqueResourceName("v2gnew")
				updated, resp, err := client.V2.Authorizations.UpdateGroup(createdGroup.Id, &sonar.AuthorizationsUpdateGroupOptions{
					Name: newName,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(updated).NotTo(BeNil())
				Expect(updated.Name).To(Equal(newName))
			})

			It("should update group description", func() {
				updated, resp, err := client.V2.Authorizations.UpdateGroup(createdGroup.Id, &sonar.AuthorizationsUpdateGroupOptions{
					Description: "Updated description",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(updated).NotTo(BeNil())
				Expect(updated.Description).To(Equal("Updated description"))
			})
		})

		Context("parameter validation", func() {
			It("should fail with empty group ID", func() {
				group, resp, err := client.V2.Authorizations.UpdateGroup("", &sonar.AuthorizationsUpdateGroupOptions{
					Name: "Updated",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(group).To(BeNil())
			})

			It("should fail with nil request", func() {
				group, resp, err := client.V2.Authorizations.UpdateGroup(createdGroup.Id, nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(group).To(BeNil())
			})
		})
	})

	Describe("DeleteGroup", func() {
		Context("with valid group ID", func() {
			It("should delete a group", func() {
				groupName := helpers.UniqueResourceName("v2gdel")
				group, resp, err := client.V2.Authorizations.CreateGroup(&sonar.AuthorizationsCreateGroupOptions{
					Name: groupName,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				cleanup.RegisterCleanup("v2-group-delete", group.Id, func() error {
					_, cleanupErr := client.V2.Authorizations.DeleteGroup(group.Id)
					return helpers.IgnoreNotFoundError(cleanupErr)
				})

				resp, err = client.V2.Authorizations.DeleteGroup(group.Id)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(BeNumerically(">=", 200))
				Expect(resp.StatusCode).To(BeNumerically("<", 300))
			})
		})

		Context("parameter validation", func() {
			It("should fail with empty group ID", func() {
				resp, err := client.V2.Authorizations.DeleteGroup("")
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Group Memberships
	// =========================================================================
	Describe("SearchGroupMemberships", func() {
		Context("without options", func() {
			It("should return a list of group memberships", func() {
				result, resp, err := client.V2.Authorizations.SearchGroupMemberships(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.GroupMemberships).NotTo(BeEmpty())
			})
		})

		Context("with group filter", func() {
			It("should filter memberships by group ID", func() {
				groups, resp, err := client.V2.Authorizations.SearchGroups(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(groups.Groups).NotTo(BeEmpty())

				groupID := groups.Groups[0].Id
				result, resp, err := client.V2.Authorizations.SearchGroupMemberships(&sonar.AuthorizationsSearchGroupMembershipsOptions{
					GroupId: groupID,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				for _, m := range result.GroupMemberships {
					Expect(m.GroupId).To(Equal(groupID))
				}
			})
		})
	})

	Describe("CreateGroupMembership", func() {
		var (
			createdGroup *sonar.Group
			createdUser  *sonar.UserV2
		)

		BeforeAll(func() {
			groupName := helpers.UniqueResourceName("v2gmem")
			var resp *http.Response
			var err error
			createdGroup, resp, err = client.V2.Authorizations.CreateGroup(&sonar.AuthorizationsCreateGroupOptions{
				Name: groupName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			cleanup.RegisterCleanup("v2-group-membership", createdGroup.Id, func() error {
				_, cleanupErr := client.V2.Authorizations.DeleteGroup(createdGroup.Id)
				return helpers.IgnoreNotFoundError(cleanupErr)
			})

			login := helpers.UniqueResourceName("v2umem")
			createdUser, resp, err = client.V2.UsersManagement.Create(&sonar.UsersCreateOptionsV2{
				Login:    login,
				Name:     "V2 Membership User",
				Password: "testPassword123!",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			cleanup.RegisterCleanup("v2-user-membership", createdUser.Id, func() error {
				_, cleanupErr := client.V2.UsersManagement.Deactivate(&sonar.UsersDeactivateOptionsV2{
					Id:        createdUser.Id,
					Anonymize: true,
				})
				return helpers.IgnoreNotFoundError(cleanupErr)
			})
		})

		Context("with valid parameters", func() {
			It("should create a group membership", func() {
				membership, resp, err := client.V2.Authorizations.CreateGroupMembership(&sonar.AuthorizationsCreateGroupMembershipOptions{
					GroupId: createdGroup.Id,
					UserId:  createdUser.Id,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				Expect(membership).NotTo(BeNil())
				Expect(membership.GroupId).To(Equal(createdGroup.Id))
				Expect(membership.UserId).To(Equal(createdUser.Id))
				Expect(membership.Id).NotTo(BeEmpty())

				cleanup.RegisterCleanup("v2-membership", membership.Id, func() error {
					_, cleanupErr := client.V2.Authorizations.DeleteGroupMembership(membership.Id)
					return helpers.IgnoreNotFoundError(cleanupErr)
				})

				result, resp, err := client.V2.Authorizations.SearchGroupMemberships(&sonar.AuthorizationsSearchGroupMembershipsOptions{
					GroupId: createdGroup.Id,
					UserId:  createdUser.Id,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result.GroupMemberships).NotTo(BeEmpty())
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil request", func() {
				membership, resp, err := client.V2.Authorizations.CreateGroupMembership(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(membership).To(BeNil())
			})

			It("should fail with empty group ID", func() {
				membership, resp, err := client.V2.Authorizations.CreateGroupMembership(&sonar.AuthorizationsCreateGroupMembershipOptions{
					UserId: createdUser.Id,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(membership).To(BeNil())
			})

			It("should fail with empty user ID", func() {
				membership, resp, err := client.V2.Authorizations.CreateGroupMembership(&sonar.AuthorizationsCreateGroupMembershipOptions{
					GroupId: createdGroup.Id,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(membership).To(BeNil())
			})
		})
	})

	Describe("DeleteGroupMembership", func() {
		Context("with valid membership ID", func() {
			It("should delete a group membership", func() {
				groupName := helpers.UniqueResourceName("v2gmdel")
				group, resp, err := client.V2.Authorizations.CreateGroup(&sonar.AuthorizationsCreateGroupOptions{
					Name: groupName,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				cleanup.RegisterCleanup("v2-group-delm", group.Id, func() error {
					_, cleanupErr := client.V2.Authorizations.DeleteGroup(group.Id)
					return helpers.IgnoreNotFoundError(cleanupErr)
				})

				login := helpers.UniqueResourceName("v2umdel")
				user, resp, err := client.V2.UsersManagement.Create(&sonar.UsersCreateOptionsV2{
					Login:    login,
					Name:     "V2 Del Membership User",
					Password: "testPassword123!",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				cleanup.RegisterCleanup("v2-user-delm", user.Id, func() error {
					_, cleanupErr := client.V2.UsersManagement.Deactivate(&sonar.UsersDeactivateOptionsV2{
						Id:        user.Id,
						Anonymize: true,
					})
					return helpers.IgnoreNotFoundError(cleanupErr)
				})

				membership, resp, err := client.V2.Authorizations.CreateGroupMembership(&sonar.AuthorizationsCreateGroupMembershipOptions{
					GroupId: group.Id,
					UserId:  user.Id,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				cleanup.RegisterCleanup("v2-membership-del", membership.Id, func() error {
					_, cleanupErr := client.V2.Authorizations.DeleteGroupMembership(membership.Id)
					return helpers.IgnoreNotFoundError(cleanupErr)
				})

				resp, err = client.V2.Authorizations.DeleteGroupMembership(membership.Id)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(BeNumerically(">=", 200))
				Expect(resp.StatusCode).To(BeNumerically("<", 300))
			})
		})

		Context("parameter validation", func() {
			It("should fail with empty membership ID", func() {
				resp, err := client.V2.Authorizations.DeleteGroupMembership("")
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})
})
