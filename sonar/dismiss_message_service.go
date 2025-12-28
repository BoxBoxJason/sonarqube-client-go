// Manage message dismissal.
package sonargo

import "net/http"

type DismissMessageService struct {
	client *Client
}

type DismissMessageCheckObject struct {
	Dismissed bool `json:"dismissed,omitempty"`
}

type DismissMessageCheckOption struct {
	MessageType string `url:"messageType,omitempty"` // Description:"The type of the message dismissed",ExampleValue:""
	ProjectKey  string `url:"projectKey,omitempty"`  // Description:"The project key",ExampleValue:""
}

// Check Check if a message has been dismissed.
func (s *DismissMessageService) Check(opt *DismissMessageCheckOption) (v *DismissMessageCheckObject, resp *http.Response, err error) {
	err = s.ValidateCheckOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "dismiss_message/check", opt)
	if err != nil {
		return
	}
	v = new(DismissMessageCheckObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type DismissMessageDismissOption struct {
	MessageType string `url:"messageType,omitempty"` // Description:"The type of the message dismissed",ExampleValue:""
	ProjectKey  string `url:"projectKey,omitempty"`  // Description:"The project key",ExampleValue:""
}

// Dismiss Dismiss a message.
func (s *DismissMessageService) Dismiss(opt *DismissMessageDismissOption) (resp *http.Response, err error) {
	err = s.ValidateDismissOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "dismiss_message/dismiss", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}
