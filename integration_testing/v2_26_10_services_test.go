package integration_testing_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

// Specs for the V2 endpoints covered with SonarQube 26.10 (instance mode, email configurations,
// and the `local` safety check of the user update).
var _ = Describe("SonarQube 26.10 V2 services", Ordered, func() {
	var (
		client *sonar.Client
		ctx    = context.Background()
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Clean code policy mode", func() {
		It("should get the mode and update it to its current value", func() {
			current, resp, err := client.V2.CleanCodePolicy.GetMode(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(current.Mode).To(BeElementOf(sonar.InstanceModeMQR, sonar.InstanceModeStandardExperience))

			updated, _, err := client.V2.CleanCodePolicy.UpdateMode(ctx, &sonar.CleanCodePolicyUpdateModeOptions{Mode: current.Mode})
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Mode).To(Equal(current.Mode))
		})

		It("should reject an invalid mode", func() {
			_, _, err := client.V2.CleanCodePolicy.UpdateMode(ctx, &sonar.CleanCodePolicyUpdateModeOptions{Mode: "bogus"})
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Email configurations", func() {
		It("should create, fetch, update and delete the email configuration", func() {
			created, _, err := client.V2.System.CreateEmailConfiguration(ctx, &sonar.SystemEmailConfigurationCreateOptions{
				Host: "smtp.example.com", Port: "587", SecurityProtocol: "STARTTLS", FromAddress: "sonar@example.com",
				FromName: "SonarQube", SubjectPrefix: "[SonarQube]", AuthMethod: sonar.EmailAuthMethodBasic,
				Username: "user", BasicPassword: "password",
			})
			Expect(err).NotTo(HaveOccurred())
			DeferCleanup(func() {
				_, _ = client.V2.System.DeleteEmailConfiguration(ctx, created.Id)
			})
			Expect(created.IsBasicPasswordSet).To(BeTrue())

			list, _, err := client.V2.System.SearchEmailConfigurations(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(list.EmailConfigurations).NotTo(BeEmpty())

			updated, _, err := client.V2.System.UpdateEmailConfiguration(ctx, created.Id, &sonar.SystemEmailConfigurationUpdateOptions{
				SubjectPrefix: "[SQ]", OauthScope: "https://example.com/.default",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.SubjectPrefix).To(Equal("[SQ]"))
			Expect(updated.OauthScope).NotTo(BeNil())
			Expect(*updated.OauthScope).To(Equal("https://example.com/.default"))

			fetched, _, err := client.V2.System.GetEmailConfiguration(ctx, created.Id)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetched.Host).To(Equal("smtp.example.com"))

			resp, err := client.V2.System.DeleteEmailConfiguration(ctx, created.Id)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})
	})

	Describe("User update local safety check", func() {
		It("should apply the update when the local status matches and reject it otherwise", func() {
			login := helpers.UniqueResourceName("local-user")
			local := true

			user, _, err := client.V2.UsersManagement.Create(ctx, &sonar.UsersCreateOptionsV2{Login: login, Name: login, Password: "Sup3r-Secret-pass!"})
			Expect(err).NotTo(HaveOccurred())
			DeferCleanup(func() {
				_, _ = client.V2.UsersManagement.Deactivate(ctx, &sonar.UsersDeactivateOptionsV2{Id: user.Id})
			})

			updated, _, err := client.V2.UsersManagement.Update(ctx, user.Id, &sonar.UsersUpdateOptionsV2{Name: "renamed", Local: &local})
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Name).To(Equal("renamed"))

			notLocal := false
			_, _, err = client.V2.UsersManagement.Update(ctx, user.Id, &sonar.UsersUpdateOptionsV2{Name: "again", Local: &notLocal})
			Expect(err).To(HaveOccurred())
		})
	})
})
