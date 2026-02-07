package sonar

import "net/http"

// ProjectDumpService handles communication with the project dump related methods
// of the SonarQube API.
// This service provides project export/import functionality.
type ProjectDumpService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// ProjectDumpExport represents the response from triggering a project export.
type ProjectDumpExport struct {
	// ProjectID is the project identifier.
	ProjectID string `json:"projectId,omitempty"`
	// ProjectKey is the project key.
	ProjectKey string `json:"projectKey,omitempty"`
	// ProjectName is the project name.
	ProjectName string `json:"projectName,omitempty"`
	// TaskID is the ID of the export task.
	TaskID string `json:"taskId,omitempty"`
}

// ProjectDumpStatus represents the response from getting the project dump status.
//
//nolint:govet // Field alignment is less important than logical grouping.
type ProjectDumpStatus struct {
	// CanBeExported indicates whether the project can be exported.
	CanBeExported bool `json:"canBeExported,omitempty"`
	// CanBeImported indicates whether a dump can be imported for this project.
	CanBeImported bool `json:"canBeImported,omitempty"`
	// DumpToImport is the path to the dump file that can be imported.
	DumpToImport string `json:"dumpToImport,omitempty"`
	// ExportedDump is the path to the exported dump file.
	ExportedDump string `json:"exportedDump,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// ProjectDumpExportOption contains parameters for the Export method.
type ProjectDumpExportOption struct {
	// Key is the project key.
	// This field is required.
	Key string `url:"key"`
}

// ProjectDumpStatusOption contains parameters for the Status method.
type ProjectDumpStatusOption struct {
	// ID is the project id.
	ID string `url:"id,omitempty"`
	// Key is the project key.
	Key string `url:"key,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateExportOpt validates the options for the Export method.
func (s *ProjectDumpService) ValidateExportOpt(opt *ProjectDumpExportOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	return nil
}

// ValidateStatusOpt validates the options for the Status method.
// Either ID or Key must be provided.
func (s *ProjectDumpService) ValidateStatusOpt(opt *ProjectDumpStatusOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	if opt.ID == "" && opt.Key == "" {
		return NewValidationError("ID or Key", "at least one of ID or Key is required", ErrMissingRequired)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Export triggers project dump so that the project can be imported to another SonarQube server.
// Requires the 'Administer' permission.
//
// API endpoint: POST /api/project_dump/export.
// Since: 1.0.
func (s *ProjectDumpService) Export(opt *ProjectDumpExportOption) (*ProjectDumpExport, *http.Response, error) {
	err := s.ValidateExportOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "project_dump/export", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectDumpExport)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Status provides the import and export status of a project.
// Permission 'Administer' is required. The project id or project key must be provided.
//
// API endpoint: GET /api/project_dump/status.
// Since: 1.0.
func (s *ProjectDumpService) Status(opt *ProjectDumpStatusOption) (*ProjectDumpStatus, *http.Response, error) {
	err := s.ValidateStatusOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "project_dump/status", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectDumpStatus)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
