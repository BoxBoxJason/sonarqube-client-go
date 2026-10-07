package enterprise_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("Hotspots migration", Ordered, func() {
	var (
		client     *sonar.Client
		cleanup    *helpers.CleanupManager
		projectKey string
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		cleanup = helpers.NewCleanupManager(client)

		projectKey = helpers.UniqueResourceName("hotspots-migration")
		_, _, err = client.Projects.Create(context.Background(), &sonar.ProjectsCreateOptions{Name: projectKey, Project: projectKey})
		Expect(err).NotTo(HaveOccurred())
		cleanup.RegisterCleanup("project", projectKey, func() error {
			_, err := client.Projects.Delete(context.Background(), &sonar.ProjectsDeleteOptions{Project: projectKey})
			return err
		})
	})

	AfterAll(func() {
		for _, err := range cleanup.Cleanup() {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	It("should report the migration status of the instance", func() {
		result, resp, err := client.Hotspots.MigrationStatus(context.Background(), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(result).NotTo(BeNil())
		Expect(result.RemainingHotspots).To(BeNumerically(">=", 0))
	})

	It("should report a complete migration for a project without hotspots", func() {
		result, resp, err := client.Hotspots.MigrationStatus(context.Background(), &sonar.HotspotsMigrationStatusOptions{Project: projectKey})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(result.RemainingHotspots).To(BeZero())
		Expect(result.Complete).To(BeTrue())
	})

	It("should dry-run a migration without writing anything", func() {
		result, resp, err := client.Hotspots.MigrateToIssues(context.Background(), &sonar.HotspotsMigrateToIssuesOptions{
			Project: projectKey,
			DryRun:  true,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(result.DryRun).To(BeTrue())
		Expect(result.Skipped).To(BeZero())
		Expect(result.Projects).To(BeEmpty())
	})
})

var _ = Describe("Notifications group subscriptions", Ordered, func() {
	const notificationType = "security-alert-raised"

	var (
		client    *sonar.Client
		groupUUID string
		groupName string
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())

		groupName = helpers.UniqueResourceName("notif-group")
		group, _, err := client.V2.Authorizations.CreateGroup(context.Background(), &sonar.AuthorizationsCreateGroupOptions{Name: groupName})
		Expect(err).NotTo(HaveOccurred())
		Expect(group).NotTo(BeNil())
		groupUUID = group.Id
		DeferCleanup(func() {
			_, _ = client.V2.Authorizations.DeleteGroup(context.Background(), groupUUID)
		})
	})

	It("should subscribe a group to a notification type", func() {
		resp, err := client.Notifications.AddGroup(context.Background(), &sonar.NotificationsGroupOptions{GroupUuid: groupUUID, Type: notificationType})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
	})

	It("should list the group subscription", func() {
		result, resp, err := client.Notifications.ListGroups(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(result.Subscriptions).To(ContainElement(sonar.NotificationsGroupSubscription{
			GroupUuid:        groupUUID,
			GroupName:        groupName,
			NotificationType: notificationType,
			ChannelKey:       "EmailNotificationChannel",
		}))
	})

	It("should list the current user's group subscriptions", func() {
		result, resp, err := client.Notifications.ListGroupSubscriptions(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(result).NotTo(BeNil())
	})

	It("should list group-subscription notification types", func() {
		result, resp, err := client.Notifications.List(context.Background(), &sonar.NotificationsListOptions{Filter: "groupSubscription"})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(result.GlobalTypes).To(ContainElement(notificationType))
	})

	It("should remove the group subscription", func() {
		resp, err := client.Notifications.RemoveGroup(context.Background(), &sonar.NotificationsGroupOptions{GroupUuid: groupUUID, Type: notificationType})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

		result, _, err := client.Notifications.ListGroups(context.Background())
		Expect(err).NotTo(HaveOccurred())
		for _, sub := range result.Subscriptions {
			Expect(sub.GroupUuid).NotTo(Equal(groupUUID))
		}
	})
})
