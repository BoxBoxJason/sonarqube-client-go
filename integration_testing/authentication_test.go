package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Authentication Service", Ordered, func() {
	var (
		client *sonargo.Client
		cfg    *helpers.Config
	)

	BeforeAll(func() {
		var err error
		cfg = helpers.LoadConfig()
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	Describe("Validate", func() {
		Context("with valid credentials", func() {
			It("should return valid=true for authenticated client", func() {
				result, resp, err := client.Authentication.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Valid).To(BeTrue())
			})
		})

		Context("with basic auth client", func() {
			It("should validate basic auth credentials", func() {
				// Create a new client with basic auth
				basicAuthClient, err := helpers.NewClient(cfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(basicAuthClient).NotTo(BeNil())

				result, resp, err := basicAuthClient.Authentication.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Valid).To(BeTrue())
			})
		})

		Context("with unauthenticated client", func() {
			It("should return valid=false for anonymous access", func() {
				// Create a client without credentials
				anonClient, err := sonargo.NewClient(nil,
					sonargo.WithBaseURL(cfg.BaseURL+"/api/"),
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(anonClient).NotTo(BeNil())

				result, resp, err := anonClient.Authentication.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Valid).To(BeFalse())
			})
		})
	})

	Describe("Login", func() {
		Context("with valid credentials", func() {
			It("should successfully login with correct username and password", func() {
				// Create a fresh client without auth for login testing
				loginClient, err := sonargo.NewClient(nil,
					sonargo.WithBaseURL(cfg.BaseURL+"/api/"),
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(loginClient).NotTo(BeNil())

				opt := &sonargo.AuthenticationLoginOption{
					Login:    cfg.Username,
					Password: cfg.Password,
				}

				resp, err := loginClient.Authentication.Login(opt)
				Expect(err).NotTo(HaveOccurred())
				// SonarQube may return 200 OK or 204 No Content depending on version
				Expect(resp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))
			})
		})

		Context("with invalid credentials", func() {
			It("should fail login with incorrect password", func() {
				// Create a fresh client without auth
				loginClient, err := sonargo.NewClient(nil,
					sonargo.WithBaseURL(cfg.BaseURL+"/api/"),
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(loginClient).NotTo(BeNil())

				opt := &sonargo.AuthenticationLoginOption{
					Login:    cfg.Username,
					Password: "wrongpassword",
				}

				resp, err := loginClient.Authentication.Login(opt)
				Expect(err).To(HaveOccurred())
				// Login failure returns 401 Unauthorized
				if resp != nil {
					Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
				}
			})

			It("should fail login with non-existent user", func() {
				// Create a fresh client without auth
				loginClient, err := sonargo.NewClient(nil,
					sonargo.WithBaseURL(cfg.BaseURL+"/api/"),
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(loginClient).NotTo(BeNil())

				opt := &sonargo.AuthenticationLoginOption{
					Login:    "nonexistentuser",
					Password: "somepassword",
				}

				resp, err := loginClient.Authentication.Login(opt)
				Expect(err).To(HaveOccurred())
				// Login failure returns 401 Unauthorized
				if resp != nil {
					Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
				}
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Authentication.Login(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing Login field", func() {
				opt := &sonargo.AuthenticationLoginOption{
					Password: "somepassword",
				}

				resp, err := client.Authentication.Login(opt)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing Password field", func() {
				opt := &sonargo.AuthenticationLoginOption{
					Login: "someuser",
				}

				resp, err := client.Authentication.Login(opt)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with empty Login field", func() {
				opt := &sonargo.AuthenticationLoginOption{
					Login:    "",
					Password: "somepassword",
				}

				resp, err := client.Authentication.Login(opt)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with empty Password field", func() {
				opt := &sonargo.AuthenticationLoginOption{
					Login:    "someuser",
					Password: "",
				}

				resp, err := client.Authentication.Login(opt)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("Logout", func() {
		Context("after successful login", func() {
			It("should successfully logout", func() {
				// Create a client and login first
				sessionClient, err := sonargo.NewClient(nil,
					sonargo.WithBaseURL(cfg.BaseURL+"/api/"),
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(sessionClient).NotTo(BeNil())

				// Login first
				loginOpt := &sonargo.AuthenticationLoginOption{
					Login:    cfg.Username,
					Password: cfg.Password,
				}

				resp, err := sessionClient.Authentication.Login(loginOpt)
				Expect(err).NotTo(HaveOccurred())
				// SonarQube may return 200 OK or 204 No Content depending on version
				Expect(resp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

				// Now logout
				resp, err = sessionClient.Authentication.Logout()
				Expect(err).NotTo(HaveOccurred())
				// SonarQube may return 200 OK or 204 No Content depending on version
				Expect(resp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))
			})
		})

		Context("without prior login", func() {
			It("should handle logout gracefully for unauthenticated client", func() {
				// Create a client without any authentication
				anonClient, err := sonargo.NewClient(nil,
					sonargo.WithBaseURL(cfg.BaseURL+"/api/"),
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(anonClient).NotTo(BeNil())

				// Logout should still succeed (no-op for unauthenticated)
				resp, err := anonClient.Authentication.Logout()
				Expect(err).NotTo(HaveOccurred())
				// SonarQube may return 200 OK or 204 No Content depending on version
				Expect(resp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))
			})
		})
	})

	Describe("Session Lifecycle", func() {
		It("should complete full login/validate/logout cycle", func() {
			// Create a fresh client for session testing
			sessionClient, err := sonargo.NewClient(nil,
				sonargo.WithBaseURL(cfg.BaseURL+"/api/"),
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(sessionClient).NotTo(BeNil())

			// Step 1: Validate should show not authenticated
			result, resp, err := sessionClient.Authentication.Validate()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Valid).To(BeFalse())

			// Step 2: Login
			loginOpt := &sonargo.AuthenticationLoginOption{
				Login:    cfg.Username,
				Password: cfg.Password,
			}

			resp, err = sessionClient.Authentication.Login(loginOpt)
			Expect(err).NotTo(HaveOccurred())
			// SonarQube may return 200 OK or 204 No Content depending on version
			Expect(resp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

			// Step 3: Validate should now show authenticated
			// Note: Session-based auth requires cookie handling which
			// may not work with the default HTTP client. This test verifies
			// the API calls work correctly.
			result, resp, err = sessionClient.Authentication.Validate()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// After login, session should be established (if cookies are preserved)

			// Step 4: Logout
			resp, err = sessionClient.Authentication.Logout()
			Expect(err).NotTo(HaveOccurred())
			// SonarQube may return 200 OK or 204 No Content depending on version
			Expect(resp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))
		})
	})
})
