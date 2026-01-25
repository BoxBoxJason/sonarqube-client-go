package sonargo

import "net/http"

// DuplicationsService handles communication with the duplications related methods
// of the SonarQube API.
// This service provides duplication information for projects.
type DuplicationsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// DuplicationsShow represents the response from showing duplications.
//
//nolint:govet // Field alignment is less important than logical grouping.
type DuplicationsShow struct {
	// Duplications is the list of duplication groups.
	Duplications []DuplicationGroup `json:"duplications,omitempty"`
	// Files is a map of file references to file information.
	// Keys are numeric strings like "1", "2", "3".
	Files map[string]DuplicatedFile `json:"files,omitempty"`
}

// DuplicationGroup represents a group of duplicated code blocks.
type DuplicationGroup struct {
	// Blocks is the list of code blocks in this duplication group.
	Blocks []DuplicationBlock `json:"blocks,omitempty"`
}

// DuplicationBlock represents a block of duplicated code.
type DuplicationBlock struct {
	// Ref is the reference to the file in the Files map.
	//nolint:tagliatelle // API uses _ref as the field name.
	Ref string `json:"_ref,omitempty"`
	// From is the starting line number.
	From int64 `json:"from,omitempty"`
	// Size is the number of lines in the duplicated block.
	Size int64 `json:"size,omitempty"`
}

// DuplicatedFile represents a file that contains duplicated code.
type DuplicatedFile struct {
	// Key is the file key.
	Key string `json:"key,omitempty"`
	// Name is the file name.
	Name string `json:"name,omitempty"`
	// ProjectName is the name of the project containing this file.
	ProjectName string `json:"projectName,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// DuplicationsShowOption contains parameters for the Show method.
type DuplicationsShowOption struct {
	// Branch key.
	// WARNING: This parameters is internal and may change without notice.
	Branch string `url:"branch,omitempty"`
	// Key is the file key.
	// This field is required.
	Key string `url:"key"`
	// PullRequest is the pull request id.
	// WARNING: This parameters is internal and may change without notice.
	PullRequest string `url:"pullRequest,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateShowOpt validates the options for the Show method.
func (s *DuplicationsService) ValidateShowOpt(opt *DuplicationsShowOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Show returns duplications for a file.
// Requires Browse permission on file's project.
//
// API endpoint: GET /api/duplications/show.
// Since: 4.4.
func (s *DuplicationsService) Show(opt *DuplicationsShowOption) (*DuplicationsShow, *http.Response, error) {
	err := s.ValidateShowOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "duplications/show", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(DuplicationsShow)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
