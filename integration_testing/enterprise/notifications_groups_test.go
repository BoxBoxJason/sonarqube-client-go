package enterprise_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

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
