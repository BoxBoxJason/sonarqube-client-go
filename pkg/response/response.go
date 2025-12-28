package response

import (
	"errors"

	glog "github.com/magicsong/color-glog"

	"github.com/boxboxjason/sonarqube-client-go/pkg/api"
)

type ExampleFetcher struct {
	endpoint, username, password string
}

func NewExampleFetcher(endpoint, username, password string) *ExampleFetcher {
	return &ExampleFetcher{endpoint: endpoint, username: username, password: password}
}
func (e *ExampleFetcher) GetResponseExample(service *api.WebService) (examples []*WebservicesResponseExampleResp, err error) {
	if service == nil || len(service.Actions) == 0 {
		return nil, errors.New("service cannot be empty")
	}
	c, err := NewClient(e.endpoint, e.username, e.password)
	if err != nil {
		return
	}
	for index, action := range service.Actions {
		opt := &WebservicesResponseExampleOption{
			Action:     action.Key,
			Controller: service.Path,
		}
		if !action.HasResponseExample {
			glog.V(3).Infof("%s of service %s does not have examples", action.Key, service.Path)
			continue
		}
		glog.V(3).Infof("%s of service %s HAVE examples", action.Key, service.Path)
		resp, err := c.Webservices.ResponseExample(opt)
		if err != nil {
			return nil, err
		}
		examples = append(examples, resp)
		service.Actions[index].ResponseType = resp.Format
	}
	return
}
