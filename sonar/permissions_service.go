package sonar

import "net/http"

const (
	// MinPermissionQueryLength is the minimum required length for permission query strings.
	MinPermissionQueryLength = 3

	// PermissionGlobalAdmin is the global admin permission.
	PermissionGlobalAdmin = "admin"
	// PermissionGlobalGateAdmin is the global gate admin permission.
	PermissionGlobalGateAdmin = "gateadmin"
	// PermissionGlobalProfileAdmin is the global profile admin permission.
	PermissionGlobalProfileAdmin = "profileadmin"
	// PermissionGlobalProvisioning is the global provisioning permission.
	PermissionGlobalProvisioning = "provisioning"
	// PermissionGlobalScan is the global scan permission.
	PermissionGlobalScan = "scan"
	// PermissionGlobalApplicationCreator is the global application creator permission.
	PermissionGlobalApplicationCreator = "applicationcreator"
	// PermissionGlobalPortfolioCreator is the global portfolio creator permission.
	PermissionGlobalPortfolioCreator = "portfoliocreator"

	// PermissionProjectAdmin is the project admin permission.
	PermissionProjectAdmin = "admin"
	// PermissionProjectCodeViewer is the project code viewer permission.
	PermissionProjectCodeViewer = "codeviewer"
	// PermissionProjectIssueAdmin is the project issue admin permission.
	PermissionProjectIssueAdmin = "issueadmin"
	// PermissionProjectSecurityHotspotAdmin is the project security hotspot admin permission.
	PermissionProjectSecurityHotspotAdmin = "securityhotspotadmin"
	// PermissionProjectScan is the project scan permission.
	PermissionProjectScan = "scan"
	// PermissionProjectUser is the project user permission.
	PermissionProjectUser = "user"
)

// PermissionsService handles communication with the permissions related methods
// of the SonarQube API.
// This service manages permission templates and the granting/revoking of permissions
// at the global and project levels.
type PermissionsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	// allowedGlobalPermissions is the set of supported global permissions.
	allowedGlobalPermissions = map[string]struct{}{
		PermissionGlobalAdmin:              {},
		PermissionGlobalGateAdmin:          {},
		PermissionGlobalProfileAdmin:       {},
		PermissionGlobalProvisioning:       {},
		PermissionGlobalScan:               {},
		PermissionGlobalApplicationCreator: {},
		PermissionGlobalPortfolioCreator:   {},
	}

	// allowedProjectPermissions is the set of supported project permissions.
	allowedProjectPermissions = map[string]struct{}{
		PermissionProjectAdmin:                {},
		PermissionProjectCodeViewer:           {},
		PermissionProjectIssueAdmin:           {},
		PermissionProjectSecurityHotspotAdmin: {},
		PermissionProjectScan:                 {},
		PermissionProjectUser:                 {},
	}

	// allowedQualifiers is the set of supported qualifiers for permissions.
	allowedQualifiers = map[string]struct{}{
		"TRK": {},
	}
)

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// PermissionGroup represents a group with its permissions.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type PermissionGroup struct {
	// Description is the group description.
	Description string `json:"description,omitempty"`
	// ID is the deprecated unique identifier of the group.
	//
	// Deprecated: Since SonarQube 8.4 - use Name instead.
	ID string `json:"id,omitempty"`
	// Managed indicates if the group is externally managed.
	Managed bool `json:"managed,omitempty"`
	// Name is the group name.
	Name string `json:"name,omitempty"`
	// Permissions is the list of permissions granted to the group.
	Permissions []string `json:"permissions,omitempty"`
}

// PermissionUser represents a user with their permissions.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type PermissionUser struct {
	// Avatar is the user's avatar URL.
	Avatar string `json:"avatar,omitempty"`
	// Email is the user's email address.
	Email string `json:"email,omitempty"`
	// Login is the user's login.
	Login string `json:"login,omitempty"`
	// Managed indicates if the user is externally managed.
	Managed bool `json:"managed,omitempty"`
	// Name is the user's display name.
	Name string `json:"name,omitempty"`
	// Permissions is the list of permissions granted to the user.
	Permissions []string `json:"permissions,omitempty"`
}

// PermissionsPaging represents pagination information for permission queries.
type PermissionsPaging struct {
	// PageIndex is the current page index (1-based).
	PageIndex int64 `json:"pageIndex,omitempty"`
	// PageSize is the number of items per page.
	PageSize int64 `json:"pageSize,omitempty"`
	// Total is the total number of items.
	Total int64 `json:"total,omitempty"`
}

// PermissionTemplate represents a permission template.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type PermissionTemplate struct {
	// CreatedAt is the template creation date.
	CreatedAt string `json:"createdAt,omitempty"`
	// Description is the template description.
	Description string `json:"description,omitempty"`
	// ID is the unique identifier of the template.
	ID string `json:"id,omitempty"`
	// Name is the template name.
	Name string `json:"name,omitempty"`
	// Permissions is the list of permissions in the template.
	Permissions []TemplatePermission `json:"permissions,omitempty"`
	// ProjectKeyPattern is the regex pattern for matching project keys.
	ProjectKeyPattern string `json:"projectKeyPattern,omitempty"`
	// UpdatedAt is the template last update date.
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// TemplatePermission represents a permission entry in a template.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type TemplatePermission struct {
	// GroupsCount is the number of groups with this permission.
	GroupsCount int64 `json:"groupsCount,omitempty"`
	// Key is the permission key.
	Key string `json:"key,omitempty"`
	// UsersCount is the number of users with this permission.
	UsersCount int64 `json:"usersCount,omitempty"`
	// WithProjectCreator indicates if the project creator has this permission.
	WithProjectCreator bool `json:"withProjectCreator,omitempty"`
}

// DefaultTemplate represents a default template mapping.
type DefaultTemplate struct {
	// Qualifier is the component qualifier (e.g., TRK for projects).
	Qualifier string `json:"qualifier,omitempty"`
	// TemplateID is the ID of the template set as default.
	TemplateID string `json:"templateId,omitempty"`
}

// TemplateGroup represents a group in a permission template.
type TemplateGroup struct {
	// Description is the group description.
	Description string `json:"description,omitempty"`
	// Name is the group name.
	Name string `json:"name,omitempty"`
	// Permissions is the list of permissions granted to the group.
	Permissions []string `json:"permissions,omitempty"`
}

// TemplateUser represents a user in a permission template.
type TemplateUser struct {
	// Avatar is the user's avatar URL.
	Avatar string `json:"avatar,omitempty"`
	// Email is the user's email address.
	Email string `json:"email,omitempty"`
	// Login is the user's login.
	Login string `json:"login,omitempty"`
	// Name is the user's display name.
	Name string `json:"name,omitempty"`
	// Permissions is the list of permissions granted to the user.
	Permissions []string `json:"permissions,omitempty"`
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// PermissionsCreateTemplate represents the response from creating a permission template.
type PermissionsCreateTemplate struct {
	// PermissionTemplate is the created permission template.
	PermissionTemplate PermissionTemplateBasic `json:"permissionTemplate,omitzero"`
}

// PermissionTemplateBasic represents basic permission template info returned on create.
type PermissionTemplateBasic struct {
	// Description is the template description.
	Description string `json:"description,omitempty"`
	// Name is the template name.
	Name string `json:"name,omitempty"`
	// ProjectKeyPattern is the regex pattern for matching project keys.
	ProjectKeyPattern string `json:"projectKeyPattern,omitempty"`
}

// PermissionsGroups represents the response from listing groups with permissions.
type PermissionsGroups struct {
	// Groups is the list of groups with their permissions.
	Groups []PermissionGroup `json:"groups,omitempty"`
	// Paging contains pagination information.
	Paging PermissionsPaging `json:"paging,omitzero"`
}

// PermissionsSearchTemplates represents the response from searching permission templates.
type PermissionsSearchTemplates struct {
	// DefaultTemplates is the list of default template mappings.
	DefaultTemplates []DefaultTemplate `json:"defaultTemplates,omitempty"`
	// PermissionTemplates is the list of permission templates.
	PermissionTemplates []PermissionTemplate `json:"permissionTemplates,omitempty"`
}

// PermissionsTemplateGroups represents the response from listing template groups.
type PermissionsTemplateGroups struct {
	// Groups is the list of groups in the template.
	Groups []TemplateGroup `json:"groups,omitempty"`
	// Paging contains pagination information.
	Paging PermissionsPaging `json:"paging,omitzero"`
}

// PermissionsTemplateUsers represents the response from listing template users.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type PermissionsTemplateUsers struct {
	// Paging contains pagination information.
	Paging PermissionsPaging `json:"paging,omitzero"`
	// Users is the list of users in the template.
	Users []TemplateUser `json:"users,omitempty"`
}

// PermissionsUpdateTemplate represents the response from updating a permission template.
type PermissionsUpdateTemplate struct {
	// PermissionTemplate is the updated permission template.
	PermissionTemplate PermissionTemplateUpdated `json:"permissionTemplate,omitzero"`
}

// PermissionTemplateUpdated represents updated permission template info.
type PermissionTemplateUpdated struct {
	// CreatedAt is the template creation date.
	CreatedAt string `json:"createdAt,omitempty"`
	// Description is the template description.
	Description string `json:"description,omitempty"`
	// ID is the unique identifier of the template.
	ID string `json:"id,omitempty"`
	// Name is the template name.
	Name string `json:"name,omitempty"`
	// ProjectKeyPattern is the regex pattern for matching project keys.
	ProjectKeyPattern string `json:"projectKeyPattern,omitempty"`
	// UpdatedAt is the template last update date.
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// PermissionsUsers represents the response from listing users with permissions.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type PermissionsUsers struct {
	// Paging contains pagination information.
	Paging PermissionsPaging `json:"paging,omitzero"`
	// Users is the list of users with their permissions.
	Users []PermissionUser `json:"users,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// PermissionsAddGroupOptions contains parameters for the AddGroup method.
type PermissionsAddGroupOptions struct {
	// GroupName is the group name or 'anyone' (case insensitive).
	// This field is required.
	GroupName string `url:"groupName"`
	// Permission is the permission to grant.
	// This field is required.
	// Global permissions: admin, gateadmin, profileadmin, provisioning, scan, applicationcreator, portfoliocreator.
	// Project permissions: admin, codeviewer, issueadmin, securityhotspotadmin, scan, user.
	Permission string `url:"permission"`
	// ProjectID is the project id. Use either ProjectID or ProjectKey for project permissions.
	ProjectID string `url:"projectId,omitempty"`
	// ProjectKey is the project key. Use either ProjectID or ProjectKey for project permissions.
	ProjectKey string `url:"projectKey,omitempty"`
}

// PermissionsAddGroupToTemplateOptions contains parameters for the AddGroupToTemplate method.
type PermissionsAddGroupToTemplateOptions struct {
	// GroupName is the group name or 'anyone' (case insensitive).
	// This field is required.
	GroupName string `url:"groupName"`
	// Permission is the permission to grant.
	// This field is required.
	// Allowed values: admin, codeviewer, issueadmin, securityhotspotadmin, scan, user.
	Permission string `url:"permission"`
	// TemplateID is the template id. Use either TemplateID or TemplateName.
	TemplateID string `url:"templateId,omitempty"`
	// TemplateName is the template name. Use either TemplateID or TemplateName.
	TemplateName string `url:"templateName,omitempty"`
}

// PermissionsAddProjectCreatorToTemplateOptions contains parameters for the AddProjectCreatorToTemplate method.
type PermissionsAddProjectCreatorToTemplateOptions struct {
	// Permission is the permission to grant to the project creator.
	// This field is required.
	// Allowed values: admin, codeviewer, issueadmin, securityhotspotadmin, scan, user.
	Permission string `url:"permission"`
	// TemplateID is the template id. Use either TemplateID or TemplateName.
	TemplateID string `url:"templateId,omitempty"`
	// TemplateName is the template name. Use either TemplateID or TemplateName.
	TemplateName string `url:"templateName,omitempty"`
}

// PermissionsAddUserOptions contains parameters for the AddUser method.
type PermissionsAddUserOptions struct {
	// Login is the user login.
	// This field is required.
	Login string `url:"login"`
	// Permission is the permission to grant.
	// This field is required.
	// Global permissions: admin, gateadmin, profileadmin, provisioning, scan, applicationcreator, portfoliocreator.
	// Project permissions: admin, codeviewer, issueadmin, securityhotspotadmin, scan, user.
	Permission string `url:"permission"`
	// ProjectID is the project id. Use either ProjectID or ProjectKey for project permissions.
	ProjectID string `url:"projectId,omitempty"`
	// ProjectKey is the project key. Use either ProjectID or ProjectKey for project permissions.
	ProjectKey string `url:"projectKey,omitempty"`
}

// PermissionsAddUserToTemplateOptions contains parameters for the AddUserToTemplate method.
type PermissionsAddUserToTemplateOptions struct {
	// Login is the user login.
	// This field is required.
	Login string `url:"login"`
	// Permission is the permission to grant.
	// This field is required.
	// Allowed values: admin, codeviewer, issueadmin, securityhotspotadmin, scan, user.
	Permission string `url:"permission"`
	// TemplateID is the template id. Use either TemplateID or TemplateName.
	TemplateID string `url:"templateId,omitempty"`
	// TemplateName is the template name. Use either TemplateID or TemplateName.
	TemplateName string `url:"templateName,omitempty"`
}

// PermissionsApplyTemplateOptions contains parameters for the ApplyTemplate method.
type PermissionsApplyTemplateOptions struct {
	// ProjectID is the project id. Use either ProjectID or ProjectKey.
	ProjectID string `url:"projectId,omitempty"`
	// ProjectKey is the project key. Use either ProjectID or ProjectKey.
	ProjectKey string `url:"projectKey,omitempty"`
	// TemplateID is the template id. Use either TemplateID or TemplateName.
	TemplateID string `url:"templateId,omitempty"`
	// TemplateName is the template name. Use either TemplateID or TemplateName.
	TemplateName string `url:"templateName,omitempty"`
}

// PermissionsBulkApplyTemplateOptions contains parameters for the BulkApplyTemplate method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type PermissionsBulkApplyTemplateOptions struct {
	// AnalyzedBefore filters projects for which last analysis is older than the given date.
	// Either a date (server timezone) or datetime can be provided.
	AnalyzedBefore string `url:"analyzedBefore,omitempty"`
	// OnProvisionedOnly filters to only provisioned projects.
	OnProvisionedOnly bool `url:"onProvisionedOnly,omitempty"`
	// Projects is a comma-separated list of project keys.
	// Maximum 1000 values allowed.
	Projects []string `url:"projects,omitempty"`
	// Query limits search to project names containing the string or project keys matching exactly.
	Query string `url:"q,omitempty"`
	// Qualifiers filters by component qualifiers. Default is TRK (projects).
	Qualifiers string `url:"qualifiers,omitempty"`
	// TemplateID is the template id. Use either TemplateID or TemplateName.
	TemplateID string `url:"templateId,omitempty"`
	// TemplateName is the template name. Use either TemplateID or TemplateName.
	TemplateName string `url:"templateName,omitempty"`
}

// PermissionsCreateTemplateOptions contains parameters for the CreateTemplate method.
type PermissionsCreateTemplateOptions struct {
	// Description is the template description.
	Description string `url:"description,omitempty"`
	// Name is the template name.
	// This field is required.
	Name string `url:"name"`
	// ProjectKeyPattern is a project key pattern. Must be a valid Java regular expression.
	ProjectKeyPattern string `url:"projectKeyPattern,omitempty"`
}

// PermissionsDeleteTemplateOptions contains parameters for the DeleteTemplate method.
type PermissionsDeleteTemplateOptions struct {
	// TemplateID is the template id. Use either TemplateID or TemplateName.
	TemplateID string `url:"templateId,omitempty"`
	// TemplateName is the template name. Use either TemplateID or TemplateName.
	TemplateName string `url:"templateName,omitempty"`
}

// PermissionsGroupsOptions contains parameters for the Groups method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type PermissionsGroupsOptions struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs

	// Permission filters by specific permission.
	// Global permissions: admin, gateadmin, profileadmin, provisioning, scan, applicationcreator, portfoliocreator.
	// Project permissions: admin, codeviewer, issueadmin, securityhotspotadmin, scan, user.
	Permission string `url:"permission,omitempty"`
	// ProjectID is the project id for project permissions.
	ProjectID string `url:"projectId,omitempty"`
	// ProjectKey is the project key for project permissions.
	ProjectKey string `url:"projectKey,omitempty"`
	// Query limits search to group names containing the supplied string.
	Query string `url:"q,omitempty"`
}

// PermissionsRemoveGroupOptions contains parameters for the RemoveGroup method.
type PermissionsRemoveGroupOptions struct {
	// GroupName is the group name or 'anyone' (case insensitive).
	// This field is required.
	GroupName string `url:"groupName"`
	// Permission is the permission to revoke.
	// This field is required.
	// Global permissions: admin, gateadmin, profileadmin, provisioning, scan, applicationcreator, portfoliocreator.
	// Project permissions: admin, codeviewer, issueadmin, securityhotspotadmin, scan, user.
	Permission string `url:"permission"`
	// ProjectID is the project id. Use either ProjectID or ProjectKey for project permissions.
	ProjectID string `url:"projectId,omitempty"`
	// ProjectKey is the project key. Use either ProjectID or ProjectKey for project permissions.
	ProjectKey string `url:"projectKey,omitempty"`
}

// PermissionsRemoveGroupFromTemplateOptions contains parameters for the RemoveGroupFromTemplate method.
type PermissionsRemoveGroupFromTemplateOptions struct {
	// GroupName is the group name or 'anyone' (case insensitive).
	// This field is required.
	GroupName string `url:"groupName"`
	// Permission is the permission to revoke.
	// This field is required.
	// Allowed values: admin, codeviewer, issueadmin, securityhotspotadmin, scan, user.
	Permission string `url:"permission"`
	// TemplateID is the template id. Use either TemplateID or TemplateName.
	TemplateID string `url:"templateId,omitempty"`
	// TemplateName is the template name. Use either TemplateID or TemplateName.
	TemplateName string `url:"templateName,omitempty"`
}

// PermissionsRemoveProjectCreatorFromTemplateOptions contains parameters for the RemoveProjectCreatorFromTemplate method.
type PermissionsRemoveProjectCreatorFromTemplateOptions struct {
	// Permission is the permission to revoke from the project creator.
	// This field is required.
	// Allowed values: admin, codeviewer, issueadmin, securityhotspotadmin, scan, user.
	Permission string `url:"permission"`
	// TemplateID is the template id. Use either TemplateID or TemplateName.
	TemplateID string `url:"templateId,omitempty"`
	// TemplateName is the template name. Use either TemplateID or TemplateName.
	TemplateName string `url:"templateName,omitempty"`
}

// PermissionsRemoveUserOptions contains parameters for the RemoveUser method.
type PermissionsRemoveUserOptions struct {
	// Login is the user login.
	// This field is required.
	Login string `url:"login"`
	// Permission is the permission to revoke.
	// This field is required.
	// Global permissions: admin, gateadmin, profileadmin, provisioning, scan, applicationcreator, portfoliocreator.
	// Project permissions: admin, codeviewer, issueadmin, securityhotspotadmin, scan, user.
	Permission string `url:"permission"`
	// ProjectID is the project id. Use either ProjectID or ProjectKey for project permissions.
	ProjectID string `url:"projectId,omitempty"`
	// ProjectKey is the project key. Use either ProjectID or ProjectKey for project permissions.
	ProjectKey string `url:"projectKey,omitempty"`
}

// PermissionsRemoveUserFromTemplateOptions contains parameters for the RemoveUserFromTemplate method.
type PermissionsRemoveUserFromTemplateOptions struct {
	// Login is the user login.
	// This field is required.
	Login string `url:"login"`
	// Permission is the permission to revoke.
	// This field is required.
	// Allowed values: admin, codeviewer, issueadmin, securityhotspotadmin, scan, user.
	Permission string `url:"permission"`
	// TemplateID is the template id. Use either TemplateID or TemplateName.
	TemplateID string `url:"templateId,omitempty"`
	// TemplateName is the template name. Use either TemplateID or TemplateName.
	TemplateName string `url:"templateName,omitempty"`
}

// PermissionsSearchTemplatesOptions contains parameters for the SearchTemplates method.
type PermissionsSearchTemplatesOptions struct {
	// Query limits search to permission template names containing the supplied string.
	Query string `url:"q,omitempty"`
}

// PermissionsSetDefaultTemplateOptions contains parameters for the SetDefaultTemplate method.
type PermissionsSetDefaultTemplateOptions struct {
	// Qualifier is the project qualifier. Default is TRK (projects).
	Qualifier string `url:"qualifier,omitempty"`
	// TemplateID is the template id. Use either TemplateID or TemplateName.
	TemplateID string `url:"templateId,omitempty"`
	// TemplateName is the template name. Use either TemplateID or TemplateName.
	TemplateName string `url:"templateName,omitempty"`
}

// PermissionsTemplateGroupsOptions contains parameters for the TemplateGroups method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type PermissionsTemplateGroupsOptions struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs

	// Permission filters by specific permission.
	// Allowed values: admin, codeviewer, issueadmin, securityhotspotadmin, scan, user.
	Permission string `url:"permission,omitempty"`
	// Query limits search to group names containing the supplied string.
	Query string `url:"q,omitempty"`
	// TemplateID is the template id. Use either TemplateID or TemplateName.
	TemplateID string `url:"templateId,omitempty"`
	// TemplateName is the template name. Use either TemplateID or TemplateName.
	TemplateName string `url:"templateName,omitempty"`
}

// PermissionsTemplateUsersOptions contains parameters for the TemplateUsers method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type PermissionsTemplateUsersOptions struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs

	// Permission filters by specific permission.
	// Allowed values: admin, codeviewer, issueadmin, securityhotspotadmin, scan, user.
	Permission string `url:"permission,omitempty"`
	// Query limits search to user names containing the supplied string.
	Query string `url:"q,omitempty"`
	// TemplateID is the template id. Use either TemplateID or TemplateName.
	TemplateID string `url:"templateId,omitempty"`
	// TemplateName is the template name. Use either TemplateID or TemplateName.
	TemplateName string `url:"templateName,omitempty"`
}

// PermissionsUpdateTemplateOptions contains parameters for the UpdateTemplate method.
type PermissionsUpdateTemplateOptions struct {
	// Description is the template description.
	Description string `url:"description,omitempty"`
	// ID is the template id.
	// This field is required.
	ID string `url:"id"`
	// Name is the template name.
	Name string `url:"name,omitempty"`
	// ProjectKeyPattern is a project key pattern. Must be a valid Java regular expression.
	ProjectKeyPattern string `url:"projectKeyPattern,omitempty"`
}

// PermissionsUsersOptions contains parameters for the Users method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type PermissionsUsersOptions struct {
	// PaginationArgs contains pagination parameters.
	PaginationArgs

	// Permission filters by specific permission.
	// Global permissions: admin, gateadmin, profileadmin, provisioning, scan, applicationcreator, portfoliocreator.
	// Project permissions: admin, codeviewer, issueadmin, securityhotspotadmin, scan, user.
	Permission string `url:"permission,omitempty"`
	// ProjectID is the project id for project permissions.
	ProjectID string `url:"projectId,omitempty"`
	// ProjectKey is the project key for project permissions.
	ProjectKey string `url:"projectKey,omitempty"`
	// Query limits search to user names containing the supplied string.
	Query string `url:"q,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// isValidPermission checks if a permission is valid for either global or project scope.
func isValidPermission(permission string) bool {
	_, isGlobal := allowedGlobalPermissions[permission]
	_, isProject := allowedProjectPermissions[permission]

	return isGlobal || isProject
}

// ValidateAddGroupOpt validates the options for the AddGroup method.
func (s *PermissionsService) ValidateAddGroupOpt(opt *PermissionsAddGroupOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GroupName, "GroupName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Permission, "Permission")
	if err != nil {
		return err
	}

	if !isValidPermission(opt.Permission) {
		return NewValidationError("Permission", "must be a valid global or project permission", ErrInvalidValue)
	}

	return nil
}

// ValidateAddGroupToTemplateOpt validates the options for the AddGroupToTemplate method.
func (s *PermissionsService) ValidateAddGroupToTemplateOpt(opt *PermissionsAddGroupToTemplateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GroupName, "GroupName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Permission, "Permission")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Permission, allowedProjectPermissions, "Permission")
	if err != nil {
		return err
	}

	// Either TemplateID or TemplateName must be provided
	if opt.TemplateID == "" && opt.TemplateName == "" {
		return NewValidationError("TemplateID/TemplateName", "either TemplateID or TemplateName must be provided", ErrMissingRequired)
	}

	return nil
}

// ValidateAddProjectCreatorToTemplateOpt validates the options for the AddProjectCreatorToTemplate method.
func (s *PermissionsService) ValidateAddProjectCreatorToTemplateOpt(opt *PermissionsAddProjectCreatorToTemplateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Permission, "Permission")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Permission, allowedProjectPermissions, "Permission")
	if err != nil {
		return err
	}

	// Either TemplateID or TemplateName must be provided
	if opt.TemplateID == "" && opt.TemplateName == "" {
		return NewValidationError("TemplateID/TemplateName", "either TemplateID or TemplateName must be provided", ErrMissingRequired)
	}

	return nil
}

// ValidateAddUserOpt validates the options for the AddUser method.
func (s *PermissionsService) ValidateAddUserOpt(opt *PermissionsAddUserOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Permission, "Permission")
	if err != nil {
		return err
	}

	if !isValidPermission(opt.Permission) {
		return NewValidationError("Permission", "must be a valid global or project permission", ErrInvalidValue)
	}

	return nil
}

// ValidateAddUserToTemplateOpt validates the options for the AddUserToTemplate method.
func (s *PermissionsService) ValidateAddUserToTemplateOpt(opt *PermissionsAddUserToTemplateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Permission, "Permission")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Permission, allowedProjectPermissions, "Permission")
	if err != nil {
		return err
	}

	// Either TemplateID or TemplateName must be provided
	if opt.TemplateID == "" && opt.TemplateName == "" {
		return NewValidationError("TemplateID/TemplateName", "either TemplateID or TemplateName must be provided", ErrMissingRequired)
	}

	return nil
}

// ValidateApplyTemplateOpt validates the options for the ApplyTemplate method.
func (s *PermissionsService) ValidateApplyTemplateOpt(opt *PermissionsApplyTemplateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	// Either ProjectID or ProjectKey must be provided
	if opt.ProjectID == "" && opt.ProjectKey == "" {
		return NewValidationError("ProjectID/ProjectKey", "either ProjectID or ProjectKey must be provided", ErrMissingRequired)
	}

	// Either TemplateID or TemplateName must be provided
	if opt.TemplateID == "" && opt.TemplateName == "" {
		return NewValidationError("TemplateID/TemplateName", "either TemplateID or TemplateName must be provided", ErrMissingRequired)
	}

	return nil
}

// ValidateBulkApplyTemplateOpt validates the options for the BulkApplyTemplate method.
func (s *PermissionsService) ValidateBulkApplyTemplateOpt(opt *PermissionsBulkApplyTemplateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	// Either TemplateID or TemplateName must be provided
	if opt.TemplateID == "" && opt.TemplateName == "" {
		return NewValidationError("TemplateID/TemplateName", "either TemplateID or TemplateName must be provided", ErrMissingRequired)
	}

	// Validate qualifiers if provided
	if opt.Qualifiers != "" {
		err := IsValueAuthorized(opt.Qualifiers, allowedQualifiers, "Qualifiers")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateCreateTemplateOpt validates the options for the CreateTemplate method.
func (s *PermissionsService) ValidateCreateTemplateOpt(opt *PermissionsCreateTemplateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	return nil
}

// ValidateDeleteTemplateOpt validates the options for the DeleteTemplate method.
func (s *PermissionsService) ValidateDeleteTemplateOpt(opt *PermissionsDeleteTemplateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	// Either TemplateID or TemplateName must be provided
	if opt.TemplateID == "" && opt.TemplateName == "" {
		return NewValidationError("TemplateID/TemplateName", "either TemplateID or TemplateName must be provided", ErrMissingRequired)
	}

	return nil
}

// ValidateGroupsOpt validates the options for the Groups method.
func (s *PermissionsService) ValidateGroupsOpt(opt *PermissionsGroupsOptions) error {
	// Options are optional
	if opt == nil {
		return nil
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	// Validate permission if provided
	if opt.Permission != "" && !isValidPermission(opt.Permission) {
		return NewValidationError("Permission", "must be a valid global or project permission", ErrInvalidValue)
	}

	// Validate query minimum length
	if opt.Query != "" {
		err := ValidateMinLength(opt.Query, MinPermissionQueryLength, "Query")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateRemoveGroupOpt validates the options for the RemoveGroup method.
func (s *PermissionsService) ValidateRemoveGroupOpt(opt *PermissionsRemoveGroupOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GroupName, "GroupName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Permission, "Permission")
	if err != nil {
		return err
	}

	if !isValidPermission(opt.Permission) {
		return NewValidationError("Permission", "must be a valid global or project permission", ErrInvalidValue)
	}

	return nil
}

// ValidateRemoveGroupFromTemplateOpt validates the options for the RemoveGroupFromTemplate method.
func (s *PermissionsService) ValidateRemoveGroupFromTemplateOpt(opt *PermissionsRemoveGroupFromTemplateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GroupName, "GroupName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Permission, "Permission")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Permission, allowedProjectPermissions, "Permission")
	if err != nil {
		return err
	}

	// Either TemplateID or TemplateName must be provided
	if opt.TemplateID == "" && opt.TemplateName == "" {
		return NewValidationError("TemplateID/TemplateName", "either TemplateID or TemplateName must be provided", ErrMissingRequired)
	}

	return nil
}

// ValidateRemoveProjectCreatorFromTemplateOpt validates the options for the RemoveProjectCreatorFromTemplate method.
func (s *PermissionsService) ValidateRemoveProjectCreatorFromTemplateOpt(opt *PermissionsRemoveProjectCreatorFromTemplateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Permission, "Permission")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Permission, allowedProjectPermissions, "Permission")
	if err != nil {
		return err
	}

	// Either TemplateID or TemplateName must be provided
	if opt.TemplateID == "" && opt.TemplateName == "" {
		return NewValidationError("TemplateID/TemplateName", "either TemplateID or TemplateName must be provided", ErrMissingRequired)
	}

	return nil
}

// ValidateRemoveUserOpt validates the options for the RemoveUser method.
func (s *PermissionsService) ValidateRemoveUserOpt(opt *PermissionsRemoveUserOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Permission, "Permission")
	if err != nil {
		return err
	}

	if !isValidPermission(opt.Permission) {
		return NewValidationError("Permission", "must be a valid global or project permission", ErrInvalidValue)
	}

	return nil
}

// ValidateRemoveUserFromTemplateOpt validates the options for the RemoveUserFromTemplate method.
func (s *PermissionsService) ValidateRemoveUserFromTemplateOpt(opt *PermissionsRemoveUserFromTemplateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Permission, "Permission")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Permission, allowedProjectPermissions, "Permission")
	if err != nil {
		return err
	}

	// Either TemplateID or TemplateName must be provided
	if opt.TemplateID == "" && opt.TemplateName == "" {
		return NewValidationError("TemplateID/TemplateName", "either TemplateID or TemplateName must be provided", ErrMissingRequired)
	}

	return nil
}

// ValidateSearchTemplatesOpt validates the options for the SearchTemplates method.
func (s *PermissionsService) ValidateSearchTemplatesOpt(opt *PermissionsSearchTemplatesOptions) error {
	// Options are optional; nothing to validate.
	return nil
}

// ValidateSetDefaultTemplateOpt validates the options for the SetDefaultTemplate method.
func (s *PermissionsService) ValidateSetDefaultTemplateOpt(opt *PermissionsSetDefaultTemplateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	// Either TemplateID or TemplateName must be provided
	if opt.TemplateID == "" && opt.TemplateName == "" {
		return NewValidationError("TemplateID/TemplateName", "either TemplateID or TemplateName must be provided", ErrMissingRequired)
	}

	// Validate qualifier if provided
	if opt.Qualifier != "" {
		err := IsValueAuthorized(opt.Qualifier, allowedQualifiers, "Qualifier")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateTemplateGroupsOpt validates the options for the TemplateGroups method.
func (s *PermissionsService) ValidateTemplateGroupsOpt(opt *PermissionsTemplateGroupsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	// Either TemplateID or TemplateName must be provided
	if opt.TemplateID == "" && opt.TemplateName == "" {
		return NewValidationError("TemplateID/TemplateName", "either TemplateID or TemplateName must be provided", ErrMissingRequired)
	}

	// Validate permission if provided
	if opt.Permission != "" {
		err := IsValueAuthorized(opt.Permission, allowedProjectPermissions, "Permission")
		if err != nil {
			return err
		}
	}

	// Validate query minimum length
	if opt.Query != "" {
		err := ValidateMinLength(opt.Query, MinPermissionQueryLength, "Query")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateTemplateUsersOpt validates the options for the TemplateUsers method.
func (s *PermissionsService) ValidateTemplateUsersOpt(opt *PermissionsTemplateUsersOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	// Either TemplateID or TemplateName must be provided
	if opt.TemplateID == "" && opt.TemplateName == "" {
		return NewValidationError("TemplateID/TemplateName", "either TemplateID or TemplateName must be provided", ErrMissingRequired)
	}

	// Validate permission if provided
	if opt.Permission != "" {
		err := IsValueAuthorized(opt.Permission, allowedProjectPermissions, "Permission")
		if err != nil {
			return err
		}
	}

	// Validate query minimum length
	if opt.Query != "" {
		err := ValidateMinLength(opt.Query, MinPermissionQueryLength, "Query")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateUpdateTemplateOpt validates the options for the UpdateTemplate method.
func (s *PermissionsService) ValidateUpdateTemplateOpt(opt *PermissionsUpdateTemplateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ID, "ID")
	if err != nil {
		return err
	}

	return nil
}

// ValidateUsersOpt validates the options for the Users method.
func (s *PermissionsService) ValidateUsersOpt(opt *PermissionsUsersOptions) error {
	// Options are optional
	if opt == nil {
		return nil
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	// Validate permission if provided
	if opt.Permission != "" && !isValidPermission(opt.Permission) {
		return NewValidationError("Permission", "must be a valid global or project permission", ErrInvalidValue)
	}

	// Validate query minimum length
	if opt.Query != "" {
		err := ValidateMinLength(opt.Query, MinPermissionQueryLength, "Query")
		if err != nil {
			return err
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// AddGroup adds a permission to a group.
// This service defaults to global permissions, but can be limited to project permissions
// by providing project id or project key.
// Requires one of the following permissions:
//   - 'Administer System'
//   - 'Administer' rights on the specified project
//
// API endpoint: POST /api/permissions/add_group.
// Since: 5.2.
func (s *PermissionsService) AddGroup(opt *PermissionsAddGroupOptions) (*http.Response, error) {
	err := s.ValidateAddGroupOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/add_group", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// AddGroupToTemplate adds a group to a permission template.
// The group name must be provided.
// Requires the following permission: 'Administer System'.
//
// API endpoint: POST /api/permissions/add_group_to_template.
// Since: 5.2.
func (s *PermissionsService) AddGroupToTemplate(opt *PermissionsAddGroupToTemplateOptions) (*http.Response, error) {
	err := s.ValidateAddGroupToTemplateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/add_group_to_template", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// AddProjectCreatorToTemplate adds a project creator to a permission template.
// Requires the following permission: 'Administer System'.
//
// API endpoint: POST /api/permissions/add_project_creator_to_template.
// Since: 6.0.
func (s *PermissionsService) AddProjectCreatorToTemplate(opt *PermissionsAddProjectCreatorToTemplateOptions) (*http.Response, error) {
	err := s.ValidateAddProjectCreatorToTemplateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/add_project_creator_to_template", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// AddUser adds permission to a user.
// This service defaults to global permissions, but can be limited to project permissions
// by providing project id or project key.
// Requires one of the following permissions:
//   - 'Administer System'
//   - 'Administer' rights on the specified project
//
// API endpoint: POST /api/permissions/add_user.
// Since: 5.2.
func (s *PermissionsService) AddUser(opt *PermissionsAddUserOptions) (*http.Response, error) {
	err := s.ValidateAddUserOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/add_user", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// AddUserToTemplate adds a user to a permission template.
// Requires the following permission: 'Administer System'.
//
// API endpoint: POST /api/permissions/add_user_to_template.
// Since: 5.2.
func (s *PermissionsService) AddUserToTemplate(opt *PermissionsAddUserToTemplateOptions) (*http.Response, error) {
	err := s.ValidateAddUserToTemplateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/add_user_to_template", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// ApplyTemplate applies a permission template to one project.
// The project id or project key must be provided.
// The template id or name must be provided.
// Requires the following permission: 'Administer System'.
//
// API endpoint: POST /api/permissions/apply_template.
// Since: 5.2.
func (s *PermissionsService) ApplyTemplate(opt *PermissionsApplyTemplateOptions) (*http.Response, error) {
	err := s.ValidateApplyTemplateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/apply_template", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// BulkApplyTemplate applies a permission template to several components.
// Managed projects will be ignored.
// The template id or name must be provided.
// Requires the following permission: 'Administer System'.
//
// API endpoint: POST /api/permissions/bulk_apply_template.
// Since: 5.5.
func (s *PermissionsService) BulkApplyTemplate(opt *PermissionsBulkApplyTemplateOptions) (*http.Response, error) {
	err := s.ValidateBulkApplyTemplateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/bulk_apply_template", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// CreateTemplate creates a permission template.
// Requires the following permission: 'Administer System'.
//
// API endpoint: POST /api/permissions/create_template.
// Since: 5.2.
func (s *PermissionsService) CreateTemplate(opt *PermissionsCreateTemplateOptions) (*PermissionsCreateTemplate, *http.Response, error) {
	err := s.ValidateCreateTemplateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/create_template", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(PermissionsCreateTemplate)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteTemplate deletes a permission template.
// Requires the following permission: 'Administer System'.
//
// API endpoint: POST /api/permissions/delete_template.
// Since: 5.2.
func (s *PermissionsService) DeleteTemplate(opt *PermissionsDeleteTemplateOptions) (*http.Response, error) {
	err := s.ValidateDeleteTemplateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/delete_template", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Groups lists the groups with their permissions.
// This service defaults to global permissions, but can be limited to project permissions
// by providing project id or project key.
// This service defaults to all groups, but can be limited to groups with a specific permission
// by providing the desired permission.
// Requires one of the following permissions:
//   - 'Administer System'
//   - 'Administer' rights on the specified project
//
// Note: This is an internal API.
//
// API endpoint: GET /api/permissions/groups.
// Since: 5.2.
func (s *PermissionsService) Groups(opt *PermissionsGroupsOptions) (*PermissionsGroups, *http.Response, error) {
	err := s.ValidateGroupsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodGet, "permissions/groups", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(PermissionsGroups)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// RemoveGroup removes a permission from a group.
// This service defaults to global permissions, but can be limited to project permissions
// by providing project id or project key.
// The group name must be provided.
// Requires one of the following permissions:
//   - 'Administer System'
//   - 'Administer' rights on the specified project
//
// API endpoint: POST /api/permissions/remove_group.
// Since: 5.2.
func (s *PermissionsService) RemoveGroup(opt *PermissionsRemoveGroupOptions) (*http.Response, error) {
	err := s.ValidateRemoveGroupOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/remove_group", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveGroupFromTemplate removes a group from a permission template.
// The group name must be provided.
// Requires the following permission: 'Administer System'.
//
// API endpoint: POST /api/permissions/remove_group_from_template.
// Since: 5.2.
func (s *PermissionsService) RemoveGroupFromTemplate(opt *PermissionsRemoveGroupFromTemplateOptions) (*http.Response, error) {
	err := s.ValidateRemoveGroupFromTemplateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/remove_group_from_template", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveProjectCreatorFromTemplate removes a project creator from a permission template.
// Requires the following permission: 'Administer System'.
//
// API endpoint: POST /api/permissions/remove_project_creator_from_template.
// Since: 6.0.
func (s *PermissionsService) RemoveProjectCreatorFromTemplate(opt *PermissionsRemoveProjectCreatorFromTemplateOptions) (*http.Response, error) {
	err := s.ValidateRemoveProjectCreatorFromTemplateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/remove_project_creator_from_template", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveUser removes permission from a user.
// This service defaults to global permissions, but can be limited to project permissions
// by providing project id or project key.
// Requires one of the following permissions:
//   - 'Administer System'
//   - 'Administer' rights on the specified project
//
// API endpoint: POST /api/permissions/remove_user.
// Since: 5.2.
func (s *PermissionsService) RemoveUser(opt *PermissionsRemoveUserOptions) (*http.Response, error) {
	err := s.ValidateRemoveUserOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/remove_user", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveUserFromTemplate removes a user from a permission template.
// Requires the following permission: 'Administer System'.
//
// API endpoint: POST /api/permissions/remove_user_from_template.
// Since: 5.2.
func (s *PermissionsService) RemoveUserFromTemplate(opt *PermissionsRemoveUserFromTemplateOptions) (*http.Response, error) {
	err := s.ValidateRemoveUserFromTemplateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/remove_user_from_template", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SearchTemplates lists permission templates.
// Requires the following permission: 'Administer System'.
//
// API endpoint: GET /api/permissions/search_templates.
// Since: 5.2.
func (s *PermissionsService) SearchTemplates(opt *PermissionsSearchTemplatesOptions) (*PermissionsSearchTemplates, *http.Response, error) {
	err := s.ValidateSearchTemplatesOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodGet, "permissions/search_templates", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(PermissionsSearchTemplates)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SetDefaultTemplate sets a permission template as default.
// Requires the following permission: 'Administer System'.
//
// API endpoint: POST /api/permissions/set_default_template.
// Since: 5.2.
func (s *PermissionsService) SetDefaultTemplate(opt *PermissionsSetDefaultTemplateOptions) (*http.Response, error) {
	err := s.ValidateSetDefaultTemplateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/set_default_template", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// TemplateGroups lists the groups with their permission as individual groups
// rather than through user affiliation on the chosen template.
// This service defaults to all groups, but can be limited to groups with a specific permission
// by providing the desired permission.
// Requires the following permission: 'Administer System'.
//
// Note: This is an internal API.
//
// API endpoint: GET /api/permissions/template_groups.
// Since: 5.2.
func (s *PermissionsService) TemplateGroups(opt *PermissionsTemplateGroupsOptions) (*PermissionsTemplateGroups, *http.Response, error) {
	err := s.ValidateTemplateGroupsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodGet, "permissions/template_groups", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(PermissionsTemplateGroups)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// TemplateUsers lists the users with their permission as individual users
// rather than through group affiliation on the chosen template.
// This service defaults to all users, but can be limited to users with a specific permission
// by providing the desired permission.
// Requires the following permission: 'Administer System'.
//
// Note: This is an internal API.
//
// API endpoint: GET /api/permissions/template_users.
// Since: 5.2.
func (s *PermissionsService) TemplateUsers(opt *PermissionsTemplateUsersOptions) (*PermissionsTemplateUsers, *http.Response, error) {
	err := s.ValidateTemplateUsersOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodGet, "permissions/template_users", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(PermissionsTemplateUsers)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateTemplate updates a permission template.
// Requires the following permission: 'Administer System'.
//
// API endpoint: POST /api/permissions/update_template.
// Since: 5.2.
func (s *PermissionsService) UpdateTemplate(opt *PermissionsUpdateTemplateOptions) (*PermissionsUpdateTemplate, *http.Response, error) {
	err := s.ValidateUpdateTemplateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodPost, "permissions/update_template", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(PermissionsUpdateTemplate)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Users lists the users with their permissions as individual users
// rather than through group affiliation.
// This service defaults to global permissions, but can be limited to project permissions
// by providing project id or project key.
// This service defaults to all users, but can be limited to users with a specific permission
// by providing the desired permission.
// Requires one of the following permissions:
//   - 'Administer System'
//   - 'Administer' rights on the specified project
//
// Note: This is an internal API.
//
// API endpoint: GET /api/permissions/users.
// Since: 5.2.
func (s *PermissionsService) Users(opt *PermissionsUsersOptions) (*PermissionsUsers, *http.Response, error) {
	err := s.ValidateUsersOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(http.MethodGet, "permissions/users", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(PermissionsUsers)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
