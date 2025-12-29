package response

import (
	"errors"
	"fmt"

	"github.com/boxboxjason/sonarqube-client-go/pkg/api"
	glog "github.com/magicsong/color-glog"
)

const glogLevel = 3

// ExampleFetcher fetches response examples for web services.
type ExampleFetcher struct {
	endpoint, username, password string
}

// NewExampleFetcher creates a new ExampleFetcher.
func NewExampleFetcher(endpoint, username, password string) *ExampleFetcher {
	return &ExampleFetcher{endpoint: endpoint, username: username, password: password}
}

// GetResponseExample fetches response examples for the given service.
func (e *ExampleFetcher) GetResponseExample(service *api.WebService) ([]*WebservicesResponseExampleResp, error) {
	if service == nil || len(service.Actions) == 0 {
		return nil, errors.New("service cannot be empty")
	}

	client, err := NewClient(e.endpoint, e.username, e.password)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	examples := make([]*WebservicesResponseExampleResp, 0, len(service.Actions))

	for index, action := range service.Actions {
		opt := &WebservicesResponseExampleOption{
			Action:     action.Key,
			Controller: service.Path,
		}

		if !action.HasResponseExample {
			glog.V(glogLevel).Infof("%s of service %s does not have examples", action.Key, service.Path)

			continue
		}

		glog.V(glogLevel).Infof("%s of service %s HAVE examples", action.Key, service.Path)

		resp, err := client.Webservices.ResponseExample(opt)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch response example: %w", err)
		}

		resp.Name = action.Key
		examples = append(examples, resp)
		service.Actions[index].ResponseType = resp.Format
	}

	return examples, nil
}
