// Manage user groups.
package sonargo

import "net/http"

type UserGroupsService struct {
	client *Client
}

type UserGroupsCreateObject struct {
	Group UserGroupsCreateObject_sub1 `json:"group,omitempty"`
}

type UserGroupsCreateObject_sub1 struct {
	Default      bool   `json:"default,omitempty"`
	Description  string `json:"description,omitempty"`
	ID           string `json:"id,omitempty"`
	MembersCount int64  `json:"membersCount,omitempty"`
	Name         string `json:"name,omitempty"`
	Organization string `json:"organization,omitempty"`
}

type UserGroupsSearchObject struct {
	Groups []UserGroupsSearchObject_sub1 `json:"groups,omitempty"`
	Paging UserGroupsSearchObject_sub2   `json:"paging,omitempty"`
}

type UserGroupsSearchObject_sub1 struct {
	Default      bool   `json:"default,omitempty"`
	Description  string `json:"description,omitempty"`
	Managed      bool   `json:"managed,omitempty"`
	MembersCount int64  `json:"membersCount,omitempty"`
	Name         string `json:"name,omitempty"`
}

type UserGroupsSearchObject_sub2 struct {
	PageIndex int64 `json:"pageIndex,omitempty"`
	PageSize  int64 `json:"pageSize,omitempty"`
	Total     int64 `json:"total,omitempty"`
}

type UserGroupsUsersObject struct {
	Paging UserGroupsUsersObject_sub1   `json:"paging,omitempty"`
	Users  []UserGroupsUsersObject_sub2 `json:"users,omitempty"`
}

type UserGroupsUsersObject_sub2 struct {
	Login    string `json:"login,omitempty"`
	Managed  bool   `json:"managed,omitempty"`
	Name     string `json:"name,omitempty"`
	Selected bool   `json:"selected,omitempty"`
}

type UserGroupsUsersObject_sub1 struct {
	PageIndex int64 `json:"pageIndex,omitempty"`
	PageSize  int64 `json:"pageSize,omitempty"`
	Total     int64 `json:"total,omitempty"`
}

type UserGroupsAddUserOption struct {
	Login string `url:"login,omitempty"` // Description:"User login",ExampleValue:"g.hopper"
	Name  string `url:"name,omitempty"`  // Description:"Group name",ExampleValue:"sonar-administrators"
}

// AddUser Add a user to a group.<br />'name' must be provided.<br />Requires the following permission: 'Administer System'.
func (s *UserGroupsService) AddUser(opt *UserGroupsAddUserOption) (resp *http.Response, err error) {
	err = s.ValidateAddUserOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "user_groups/add_user", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type UserGroupsCreateOption struct {
	Description string `url:"description,omitempty"` // Description:"Description for the new group. A group description cannot be larger than 200 characters.",ExampleValue:"Default group for new users"
	Name        string `url:"name,omitempty"`        // Description:"Name for the new group. A group name cannot be larger than 255 characters and must be unique. The value 'anyone' (whatever the case) is reserved and cannot be used.",ExampleValue:"sonar-users"
}

// Create Create a group.<br>Requires the following permission: 'Administer System'.
func (s *UserGroupsService) Create(opt *UserGroupsCreateOption) (v *UserGroupsCreateObject, resp *http.Response, err error) {
	err = s.ValidateCreateOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "user_groups/create", opt)
	if err != nil {
		return
	}
	v = new(UserGroupsCreateObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type UserGroupsDeleteOption struct {
	Name string `url:"name,omitempty"` // Description:"Group name",ExampleValue:"sonar-administrators"
}

// Delete Delete a group. The default groups cannot be deleted.<br/>'name' must be provided.<br />Requires the following permission: 'Administer System'.
func (s *UserGroupsService) Delete(opt *UserGroupsDeleteOption) (resp *http.Response, err error) {
	err = s.ValidateDeleteOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "user_groups/delete", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type UserGroupsRemoveUserOption struct {
	Login string `url:"login,omitempty"` // Description:"User login",ExampleValue:"g.hopper"
	Name  string `url:"name,omitempty"`  // Description:"Group name",ExampleValue:"sonar-administrators"
}

// RemoveUser Remove a user from a group.<br />'name' must be provided.<br>Requires the following permission: 'Administer System'.
func (s *UserGroupsService) RemoveUser(opt *UserGroupsRemoveUserOption) (resp *http.Response, err error) {
	err = s.ValidateRemoveUserOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "user_groups/remove_user", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type UserGroupsSearchOption struct {
	F       string `url:"f,omitempty"`       // Description:"Comma-separated list of the fields to be returned in response. All the fields are returned by default.",ExampleValue:""
	Managed string `url:"managed,omitempty"` // Description:"Return managed or non-managed groups. Only available for managed instances, throws for non-managed instances.",ExampleValue:""
	P       string `url:"p,omitempty"`       // Description:"1-based page number",ExampleValue:"42"
	Ps      string `url:"ps,omitempty"`      // Description:"Page size. Must be greater than 0 and less or equal than 500",ExampleValue:"20"
	Q       string `url:"q,omitempty"`       // Description:"Limit search to names that contain the supplied string.",ExampleValue:"sonar-users"
}

// Search Search for user groups.<br>Requires the following permission: 'Administer System'.
func (s *UserGroupsService) Search(opt *UserGroupsSearchOption) (v *UserGroupsSearchObject, resp *http.Response, err error) {
	err = s.ValidateSearchOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "user_groups/search", opt)
	if err != nil {
		return
	}
	v = new(UserGroupsSearchObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type UserGroupsUpdateOption struct {
	CurrentName string `url:"currentName,omitempty"` // Description:"Name of the group to be updated.",ExampleValue:"AU-Tpxb--iU5OvuD2FLy"
	Description string `url:"description,omitempty"` // Description:"New optional description for the group. A group description cannot be larger than 200 characters. If value is not defined, then description is not changed.",ExampleValue:"Default group for new users"
	Name        string `url:"name,omitempty"`        // Description:"New optional name for the group. A group name cannot be larger than 255 characters and must be unique. Value 'anyone' (whatever the case) is reserved and cannot be used. If value is empty or not defined, then name is not changed.",ExampleValue:"my-group"
}

// Update Update a group.<br>Requires the following permission: 'Administer System'.
func (s *UserGroupsService) Update(opt *UserGroupsUpdateOption) (resp *http.Response, err error) {
	err = s.ValidateUpdateOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "user_groups/update", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type UserGroupsUsersOption struct {
	Name     string `url:"name,omitempty"`     // Description:"Group name",ExampleValue:"sonar-administrators"
	P        string `url:"p,omitempty"`        // Description:"1-based page number",ExampleValue:"42"
	Ps       string `url:"ps,omitempty"`       // Description:"Page size. Must be greater than 0.",ExampleValue:"20"
	Q        string `url:"q,omitempty"`        // Description:"Limit search to names or logins that contain the supplied string.",ExampleValue:"freddy"
	Selected string `url:"selected,omitempty"` // Description:"Depending on the value, show only selected items (selected=selected), deselected items (selected=deselected), or all items with their selection status (selected=all).",ExampleValue:""
}

// Users Search for users with membership information with respect to a group.<br>Requires the following permission: 'Administer System'.
func (s *UserGroupsService) Users(opt *UserGroupsUsersOption) (v *UserGroupsUsersObject, resp *http.Response, err error) {
	err = s.ValidateUsersOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "user_groups/users", opt)
	if err != nil {
		return
	}
	v = new(UserGroupsUsersObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}
