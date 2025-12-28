// Manage GitHub provisioning.
package sonargo

import "net/http"

type GithubProvisioningService struct {
	client *Client
}

type GithubProvisioningCheckObject struct {
	Application   GithubProvisioningCheckObject_sub3   `json:"application,omitempty"`
	Installations []GithubProvisioningCheckObject_sub4 `json:"installations,omitempty"`
}

type GithubProvisioningCheckObject_sub4 struct {
	AutoProvisioning GithubProvisioningCheckObject_sub1 `json:"autoProvisioning,omitempty"`
	Jit              GithubProvisioningCheckObject_sub2 `json:"jit,omitempty"`
	Organization     string                             `json:"organization,omitempty"`
}

type GithubProvisioningCheckObject_sub3 struct {
	AutoProvisioning GithubProvisioningCheckObject_sub1 `json:"autoProvisioning,omitempty"`
	Jit              GithubProvisioningCheckObject_sub2 `json:"jit,omitempty"`
}

type GithubProvisioningCheckObject_sub1 struct {
	ErrorMessage string `json:"errorMessage,omitempty"`
	Status       string `json:"status,omitempty"`
}

type GithubProvisioningCheckObject_sub2 struct {
	Status string `json:"status,omitempty"`
}

// Check Validate Github provisioning configuration.
func (s *GithubProvisioningService) Check() (v *GithubProvisioningCheckObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("POST", "github_provisioning/check", nil)
	if err != nil {
		return
	}
	v = new(GithubProvisioningCheckObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}
