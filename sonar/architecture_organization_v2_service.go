package sonar

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// ArchitectureValueOther is the "OTHER" enumeration value.
	ArchitectureValueOther = "OTHER"
)

//nolint:gochecknoglobals,goconst // constant sets of allowed values
var (
	// allowedArchitectureListExternalInterfacesDirection is the set of allowed values for the corresponding option.
	allowedArchitectureListExternalInterfacesDirection = map[string]struct{}{
		"EXIT_POINT":  {},
		"ENTRY_POINT": {},
	}
	// allowedArchitectureUpdateBoundaryDescriptorDirection is the set of allowed values for the corresponding option.
	allowedArchitectureUpdateBoundaryDescriptorDirection = map[string]struct{}{
		"EXIT_POINT":  {},
		"ENTRY_POINT": {},
	}
	// allowedArchitectureUpdateBoundaryDescriptorEcosystem is the set of allowed values for the corresponding option.
	allowedArchitectureUpdateBoundaryDescriptorEcosystem = map[string]struct{}{
		LanguageJava: {},
		"js":         {},
		"ts":         {},
		"py":         {},
		"cs":         {},
	}
	// allowedArchitectureUpdateOrganizationPlaceholderType is the set of allowed values for the corresponding option.
	allowedArchitectureUpdateOrganizationPlaceholderType = map[string]struct{}{
		ProjectQualifierAPP:    {},
		"DB":                   {},
		"QUEUE":                {},
		ArchitectureValueOther: {},
	}
	// allowedArchitectureCreateBoundaryDescriptorEcosystem is the set of allowed values for the corresponding option.
	allowedArchitectureCreateBoundaryDescriptorEcosystem = map[string]struct{}{
		LanguageJava: {},
		"js":         {},
		"ts":         {},
		"py":         {},
		"cs":         {},
	}
	// allowedArchitectureCreateBoundaryDescriptorDirection is the set of allowed values for the corresponding option.
	allowedArchitectureCreateBoundaryDescriptorDirection = map[string]struct{}{
		"EXIT_POINT":  {},
		"ENTRY_POINT": {},
	}
	// allowedArchitectureCreateOrganizationPlaceholderType is the set of allowed values for the corresponding option.
	allowedArchitectureCreateOrganizationPlaceholderType = map[string]struct{}{
		ProjectQualifierAPP:    {},
		"DB":                   {},
		"QUEUE":                {},
		ArchitectureValueOther: {},
	}
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// ArchitectureListBoundaryDescriptorsOptions contains parameters for the ListBoundaryDescriptors method.
type ArchitectureListBoundaryDescriptorsOptions struct {
	// ProjectId project Id. This field is required.
	ProjectId string `json:"projectId"`
	// PageIndex page number for pagination (1-based). Specifies which page of results to return.
	// Minimum value: 1. Default: 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize number of items per page. Set to 0 to return only the total count without any items
	// (useful for counting).
	// Minimum value: 0. Default: 50.
	PageSize int32 `json:"pageSize,omitempty"`
}

// ArchitectureBoundaryDescriptor represents the BoundaryDescriptor object of the SonarQube V2 API.
type ArchitectureBoundaryDescriptor struct {
	// Direction of the boundary. Allowed values: EXIT_POINT, ENTRY_POINT.
	Direction string `json:"direction,omitempty"`
	// Name of the boundary descriptor.
	Name string `json:"name,omitempty"`
	// KeyTemplate template for generating boundary keys.
	KeyTemplate string `json:"keyTemplate,omitempty"`
	// Query to match boundaries.
	Query string `json:"query,omitempty"`
	// Id is the id.
	Id string `json:"id,omitempty"`
	// Ecosystem language ecosystem of the boundary. Allowed values: java, js, ts, py, cs.
	Ecosystem string `json:"ecosystem,omitempty"`
	// ProjectId is the project id.
	ProjectId string `json:"projectId,omitempty"`
}

// ArchitectureBoundaryDescriptorPagedResponse represents the BoundaryDescriptorPagedResponse object
// of the SonarQube V2 API.
type ArchitectureBoundaryDescriptorPagedResponse struct {
	// BoundaryDescriptors is the boundary descriptors.
	BoundaryDescriptors []ArchitectureBoundaryDescriptor `json:"boundaryDescriptors,omitempty"`
	// Page is the page.
	Page PageResponseV2 `json:"page,omitzero"`
}

// ArchitectureListExternalInterfacesOptions contains parameters for the ListExternalInterfaces method.
type ArchitectureListExternalInterfacesOptions struct {
	// OrganizationId UUID v4 of the organization. This field is required.
	OrganizationId string `json:"organizationId"`
	// ComponentId the component whose external interfaces are returned. When omitted, the external
	// interfaces of every component in the organization are returned.
	ComponentId string `json:"componentId,omitempty"`
	// Direction when set, only external interfaces of this direction (entry or exit) are returned.
	// Allowed values: EXIT_POINT, ENTRY_POINT.
	Direction string `json:"direction,omitempty"`
}

// ArchitectureExternalInterface represents a SonarQube V2 API object.
//
// An external interface of a component. Today this is an analyzer-found boundary; the abstraction
// will later also cover SDK-declared and default entry points.
type ArchitectureExternalInterface struct {
	// Direction of the external interface (entry or exit point). Allowed values: EXIT_POINT,
	// ENTRY_POINT.
	Direction string `json:"direction,omitempty"`
	// ComponentId id of the organization component this external interface belongs to.
	ComponentId string `json:"componentId,omitempty"`
	// Key identifying the external interface.
	Key string `json:"key,omitempty"`
}

// ArchitectureExternalInterfaceListResponse represents the ExternalInterfaceListResponse object of
// the SonarQube V2 API.
type ArchitectureExternalInterfaceListResponse struct {
	// ExternalInterfaces is the external interfaces.
	ExternalInterfaces []ArchitectureExternalInterface `json:"externalInterfaces,omitempty"`
}

// ArchitectureGetOrganizationArchitectureOptions contains parameters for the GetOrganizationArchitecture method.
type ArchitectureGetOrganizationArchitectureOptions struct {
	// OrganizationId UUID v4 of the organization. This field is required.
	OrganizationId string `json:"organizationId"`
}

// ArchitectureOrganizationArchitectureComponentMetadata represents the
// OrganizationArchitectureComponentMetadata object of the SonarQube V2 API.
type ArchitectureOrganizationArchitectureComponentMetadata struct {
	// UnmatchedExitPoints keys of the component's exit points that are not used by any project
	// relationship.
	UnmatchedExitPoints []string `json:"unmatchedExitPoints,omitempty"`
	// MatchedExitPoints keys of the component's exit points that are used by a project
	// relationship.
	MatchedExitPoints []string `json:"matchedExitPoints,omitempty"`
}

// ArchitectureOrganizationArchitectureComponent represents the OrganizationArchitectureComponent
// object of the SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type ArchitectureOrganizationArchitectureComponent struct {
	// Metadata is the metadata.
	Metadata ArchitectureOrganizationArchitectureComponentMetadata `json:"metadata,omitzero"`
	// Id is the id.
	Id string `json:"id,omitempty"`
}

// ArchitectureOrganizationArchitectureRelationshipMetadata represents the
// OrganizationArchitectureRelationshipMetadata object of the SonarQube V2 API.
type ArchitectureOrganizationArchitectureRelationshipMetadata struct {
	// DependencyCount number of actual dependencies that exist between the 'from' and the 'to'
	// component.
	DependencyCount int32 `json:"dependencyCount,omitempty"`
}

// ArchitectureOrganizationArchitectureRelationship represents the
// OrganizationArchitectureRelationship object of the SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type ArchitectureOrganizationArchitectureRelationship struct {
	// From is the from.
	From string `json:"from,omitempty"`
	// Metadata is the metadata.
	Metadata ArchitectureOrganizationArchitectureRelationshipMetadata `json:"metadata,omitzero"`
	// To is the to.
	To string `json:"to,omitempty"`
}

// ArchitectureOrganizationArchitecture represents the OrganizationArchitecture object of the
// SonarQube V2 API.
type ArchitectureOrganizationArchitecture struct {
	// Components is the components.
	Components []ArchitectureOrganizationArchitectureComponent `json:"components,omitempty"`
	// Relationships is the relationships.
	Relationships []ArchitectureOrganizationArchitectureRelationship `json:"relationships,omitempty"`
}

// ArchitectureOrganizationArchitectureResponse represents the OrganizationArchitectureResponse
// object of the SonarQube V2 API.
type ArchitectureOrganizationArchitectureResponse struct {
	// OrganizationArchitectures is the organization architectures.
	OrganizationArchitectures []ArchitectureOrganizationArchitecture `json:"organizationArchitectures,omitempty"`
}

// ArchitectureListOrganizationComponentsOptions contains parameters for the ListOrganizationComponents method.
type ArchitectureListOrganizationComponentsOptions struct {
	// OrganizationId UUID v4 of the organization. This field is required.
	OrganizationId string `json:"organizationId"`
}

// ArchitectureOrganizationComponent represents the OrganizationComponent object of the SonarQube V2
// API.
type ArchitectureOrganizationComponent struct {
	// DisplayName is the display name.
	DisplayName string `json:"displayName,omitempty"`
	// Id component id. For PROJECT components, this is the project id.
	Id string `json:"id,omitempty"`
	// Type is the type. Allowed values: PROJECT, PLACEHOLDER.
	Type string `json:"type,omitempty"`
	// ProjectId set when type=PROJECT. Same value as id. Absent for PLACEHOLDER.
	ProjectId string `json:"projectId,omitempty"`
	// OrganizationId UUID v4 of the organization.
	OrganizationId string `json:"organizationId,omitempty"`
}

// ArchitectureOrganizationComponentListResponse represents the OrganizationComponentListResponse
// object of the SonarQube V2 API.
type ArchitectureOrganizationComponentListResponse struct {
	// Components is the components.
	Components []ArchitectureOrganizationComponent `json:"components,omitempty"`
}

// ArchitectureListOrganizationPlaceholdersOptions contains parameters for the ListOrganizationPlaceholders method.
type ArchitectureListOrganizationPlaceholdersOptions struct {
	// OrganizationId UUID v4 of the organization. This field is required.
	OrganizationId string `json:"organizationId"`
}

// ArchitectureOrganizationPlaceholder represents the OrganizationPlaceholder object of the
// SonarQube V2 API.
type ArchitectureOrganizationPlaceholder struct {
	// OrganizationId UUID v4 of the organization.
	OrganizationId string `json:"organizationId,omitempty"`
	// Id is the id.
	Id string `json:"id,omitempty"`
	// Type is the type. Allowed values: APP, DB, QUEUE, OTHER.
	Type string `json:"type,omitempty"`
	// DisplayName is the display name.
	DisplayName string `json:"displayName,omitempty"`
}

// ArchitectureOrganizationPlaceholderListResponse represents the
// OrganizationPlaceholderListResponse object of the SonarQube V2 API.
type ArchitectureOrganizationPlaceholderListResponse struct {
	// Placeholders is the placeholders.
	Placeholders []ArchitectureOrganizationPlaceholder `json:"placeholders,omitempty"`
}

// ArchitectureListPatternsOptions contains parameters for the ListPatterns method.
type ArchitectureListPatternsOptions struct {
	// OrganizationId UUID v4 of the organization. This field is required.
	OrganizationId string `json:"organizationId"`
}

// ArchitectureModelConstraint represents the ModelConstraint object of the SonarQube V2 API.
type ArchitectureModelConstraint struct {
	// To pattern identifying the target group of the dependency.
	To string `json:"to,omitempty"`
	// From pattern identifying the source group of the dependency.
	From string `json:"from,omitempty"`
}

// ArchitecturePatternData represents a SonarQube V2 API object.
//
// Structured pattern definition. Language-agnostic, no perspectives array — represents a single
// flat pattern. Reuses ModelGroup and ModelConstraint.
type ArchitecturePatternData struct {
	// Groups top-level components of the pattern.
	Groups []ArchitectureModelGroup `json:"groups,omitempty"`
	// Description of the pattern definition.
	Description string `json:"description,omitempty"`
	// Label display name of the pattern definition.
	Label string `json:"label,omitempty"`
	// Constraints dependency constraints between groups.
	Constraints []ArchitectureModelConstraint `json:"constraints,omitempty"`
}

// ArchitecturePattern represents a SonarQube V2 API object.
//
// A pattern resource.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type ArchitecturePattern struct {
	// OrganizationId UUID v4 of the organization.
	OrganizationId string `json:"organizationId,omitempty"`
	// Id is the id.
	Id string `json:"id,omitempty"`
	// Pattern is the pattern.
	Pattern ArchitecturePatternData `json:"pattern,omitzero"`
	// Name display name of the pattern.
	Name string `json:"name,omitempty"`
}

// ArchitectureListProjectRelationshipsOptions contains parameters for the ListProjectRelationships method.
type ArchitectureListProjectRelationshipsOptions struct {
	// ProjectId project Id. This field is required.
	ProjectId string `json:"projectId"`
	// PageIndex page number for pagination (1-based). Specifies which page of results to return.
	// Minimum value: 1. Default: 1.
	PageIndex int32 `json:"pageIndex,omitempty"`
	// PageSize number of items per page. Set to 0 to return only the total count without any items
	// (useful for counting).
	// Minimum value: 0. Default: 50.
	PageSize int32 `json:"pageSize,omitempty"`
}

// ArchitectureProjectRelationship represents the ProjectRelationship object of the SonarQube V2
// API.
type ArchitectureProjectRelationship struct {
	// BoundaryKey key identifying the exit boundary (consumer side) this relationship applies to.
	BoundaryKey string `json:"boundaryKey,omitempty"`
	// ProjectId id of the project this relationship belongs to (the consumer).
	ProjectId string `json:"projectId,omitempty"`
	// TargetEntryPointKey key of the target component's entry point this exit connects to.
	TargetEntryPointKey string `json:"targetEntryPointKey,omitempty"`
	// Origin whether this relationship was created manually via the API or auto-computed from
	// analysis. Allowed values: MANUAL, AUTOMAPPED.
	Origin string `json:"origin,omitempty"`
	// TargetComponentId id of the target organization component (the provider this exit connects
	// to).
	TargetComponentId string `json:"targetComponentId,omitempty"`
	// OrganizationId UUID v4 of the organization.
	OrganizationId string `json:"organizationId,omitempty"`
	// Id is the id.
	Id string `json:"id,omitempty"`
}

// ArchitectureProjectRelationshipPagedResponse represents the ProjectRelationshipPagedResponse
// object of the SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type ArchitectureProjectRelationshipPagedResponse struct {
	// Page is the page.
	Page PageResponseV2 `json:"page,omitzero"`
	// ProjectRelationships is the project relationships.
	ProjectRelationships []ArchitectureProjectRelationship `json:"projectRelationships,omitempty"`
}

// ArchitectureListSdksOptions contains parameters for the ListSdks method.
type ArchitectureListSdksOptions struct {
	// OrganizationId UUID v4 of the organization. This field is required.
	OrganizationId string `json:"organizationId"`
}

// ArchitectureSdkBoundaryDescriptor represents a SonarQube V2 API object.
//
// One ecosystem-specific exit boundary descriptor of an SDK. A project's analysis emits a boundary
// wherever this descriptor matches.
type ArchitectureSdkBoundaryDescriptor struct {
	// Ecosystem the ecosystem this boundary descriptor applies to. Allowed values: java, js, ts,
	// py, cs.
	Ecosystem string `json:"ecosystem,omitempty"`
	// Query matching the call into the SDK, in the boundary-descriptor query grammar.
	Query string `json:"query,omitempty"`
}

// ArchitectureSdk represents a SonarQube V2 API object.
//
// An organization-scoped SDK: one identity, target and key, grouping N ecosystem-specific exit
// boundary descriptors.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type ArchitectureSdk struct {
	// SdkKey unique key identifying this SDK within its organization. What edge matching compares
	// a found exit boundary's key against.
	SdkKey string `json:"sdkKey,omitempty"`
	// OrganizationId UUID v4 of the organization.
	OrganizationId string `json:"organizationId,omitempty"`
	// Id is the id.
	Id string `json:"id,omitempty"`
	// TargetComponentId id of the organization component the SDK leads into (the provider).
	TargetComponentId string `json:"targetComponentId,omitempty"`
	// SdkBoundaryDescriptors is the sdk boundary descriptors.
	SdkBoundaryDescriptors []ArchitectureSdkBoundaryDescriptor `json:"sdkBoundaryDescriptors,omitempty"`
	// Name display name of the SDK.
	Name string `json:"name,omitempty"`
}

// ArchitectureSdkListResponse represents the SdkListResponse object of the SonarQube V2 API.
type ArchitectureSdkListResponse struct {
	// Sdks is the sdks.
	Sdks []ArchitectureSdk `json:"sdks,omitempty"`
}

// ArchitectureUpdateBoundaryDescriptorOptions contains the request body for the
// UpdateBoundaryDescriptor method.
type ArchitectureUpdateBoundaryDescriptorOptions struct {
	// Direction of the boundary. Allowed values: EXIT_POINT, ENTRY_POINT.
	Direction *string `json:"direction,omitempty"`
	// Ecosystem language ecosystem of the boundary. Allowed values: java, js, ts, py, cs.
	Ecosystem *string `json:"ecosystem,omitempty"`
	// KeyTemplate template for generating boundary keys.
	KeyTemplate *string `json:"keyTemplate,omitempty"`
	// Query to match boundaries.
	Query *string `json:"query,omitempty"`
	// Name of the boundary descriptor.
	Name *string `json:"name,omitempty"`
}

// ArchitectureUpdateOrganizationPlaceholderOptions contains the request body for the
// UpdateOrganizationPlaceholder method.
type ArchitectureUpdateOrganizationPlaceholderOptions struct {
	// DisplayName is the display name.
	DisplayName *string `json:"displayName,omitempty"`
	// Type is the type. Allowed values: APP, DB, QUEUE, OTHER.
	Type *string `json:"type,omitempty"`
}

// ArchitectureUpdatePatternOptions contains the request body for the UpdatePattern method.
type ArchitectureUpdatePatternOptions struct {
	// Pattern is the pattern.
	Pattern *ArchitecturePatternData `json:"pattern,omitempty"`
	// Name updated display name.
	Name *string `json:"name,omitempty"`
}

// ArchitectureUpdateProjectRelationshipOptions contains the request body for the
// UpdateProjectRelationship method.
type ArchitectureUpdateProjectRelationshipOptions struct {
	// BoundaryKey key identifying the exit boundary (consumer side) this relationship applies to.
	BoundaryKey *string `json:"boundaryKey,omitempty"`
	// TargetComponentId id of the target organization component (the provider this exit connects
	// to).
	TargetComponentId *string `json:"targetComponentId,omitempty"`
	// TargetEntryPointKey key of the target component's entry point this exit connects to.
	TargetEntryPointKey *string `json:"targetEntryPointKey,omitempty"`
}

// ArchitectureUpdateSdkOptions contains the request body for the UpdateSdk method.
type ArchitectureUpdateSdkOptions struct {
	// TargetComponentId id of the organization component the SDK leads into (the provider).
	TargetComponentId *string `json:"targetComponentId,omitempty"`
	// Name is the name.
	Name *string `json:"name,omitempty"`
	// SdkBoundaryDescriptors is the sdk boundary descriptors.
	SdkBoundaryDescriptors []ArchitectureSdkBoundaryDescriptor `json:"sdkBoundaryDescriptors,omitempty"`
}

// ArchitectureCreateBoundaryDescriptorOptions contains the request body for the
// CreateBoundaryDescriptor method.
type ArchitectureCreateBoundaryDescriptorOptions struct {
	// Ecosystem language ecosystem of the boundary. This field is required. Allowed values: java,
	// js, ts, py, cs.
	Ecosystem string `json:"ecosystem"`
	// KeyTemplate template for generating boundary keys. This field is required.
	KeyTemplate string `json:"keyTemplate"`
	// Query to match boundaries. This field is required.
	Query string `json:"query"`
	// Direction of the boundary. This field is required. Allowed values: EXIT_POINT, ENTRY_POINT.
	Direction string `json:"direction"`
	// ProjectId is the project id. This field is required.
	ProjectId string `json:"projectId"`
	// Name of the boundary descriptor. This field is required.
	Name string `json:"name"`
}

// ArchitectureCreateOrganizationPlaceholderOptions contains the request body for the
// CreateOrganizationPlaceholder method.
type ArchitectureCreateOrganizationPlaceholderOptions struct {
	// OrganizationId UUID v4 of the organization.
	OrganizationId string `json:"organizationId,omitempty"`
	// DisplayName is the display name. This field is required.
	DisplayName string `json:"displayName"`
	// Type is the type. Allowed values: APP, DB, QUEUE, OTHER.
	Type string `json:"type,omitempty"`
}

// ArchitectureCreatePatternOptions contains the request body for the CreatePattern method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type ArchitectureCreatePatternOptions struct {
	// OrganizationId UUID v4 of the organization. This field is required.
	OrganizationId string `json:"organizationId"`
	// Pattern is the pattern. This field is required.
	Pattern ArchitecturePatternData `json:"pattern"`
	// Name display name of the pattern. This field is required.
	Name string `json:"name"`
}

// ArchitectureCreateProjectRelationshipOptions contains the request body for the
// CreateProjectRelationship method.
type ArchitectureCreateProjectRelationshipOptions struct {
	// BoundaryKey key identifying the exit boundary (consumer side) this relationship applies to.
	// This field is required.
	BoundaryKey string `json:"boundaryKey"`
	// ProjectId id of the project this relationship belongs to (the consumer). This field is
	// required.
	ProjectId string `json:"projectId"`
	// TargetComponentId id of the target organization component (the provider this exit connects
	// to). This field is required.
	TargetComponentId string `json:"targetComponentId"`
	// TargetEntryPointKey key of the target component's entry point this exit connects to. This
	// field is required.
	TargetEntryPointKey string `json:"targetEntryPointKey"`
}

// ArchitectureCreateSdkOptions contains the request body for the CreateSdk method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type ArchitectureCreateSdkOptions struct {
	// TargetComponentId id of the organization component the SDK leads into (the provider). This
	// field is required.
	TargetComponentId string `json:"targetComponentId"`
	// SdkBoundaryDescriptors is the sdk boundary descriptors. This field is required.
	SdkBoundaryDescriptors []ArchitectureSdkBoundaryDescriptor `json:"sdkBoundaryDescriptors"`
	// Name is the name. This field is required.
	Name string `json:"name"`
	// SdkKey unique key identifying this SDK within its organization. This field is required.
	SdkKey string `json:"sdkKey"`
	// OrganizationId UUID v4 of the organization.
	OrganizationId string `json:"organizationId,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateListBoundaryDescriptorsOpt validates the options for the ListBoundaryDescriptors method.
func (s *ArchitectureService) ValidateListBoundaryDescriptorsOpt(opt *ArchitectureListBoundaryDescriptorsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectId, "ProjectId")
	if err != nil {
		return err
	}

	if opt.PageIndex < 0 {
		return NewValidationError("PageIndex", "must be greater than 0", ErrOutOfRange)
	}

	return nil
}

// ValidateListExternalInterfacesOpt validates the options for the ListExternalInterfaces method.
func (s *ArchitectureService) ValidateListExternalInterfacesOpt(opt *ArchitectureListExternalInterfacesOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Direction, allowedArchitectureListExternalInterfacesDirection, "Direction")
	if err != nil {
		return err
	}

	return nil
}

// ValidateGetOrganizationArchitectureOpt validates the options for the GetOrganizationArchitecture method.
func (s *ArchitectureService) ValidateGetOrganizationArchitectureOpt(opt *ArchitectureGetOrganizationArchitectureOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	return nil
}

// ValidateListOrganizationComponentsOpt validates the options for the ListOrganizationComponents method.
func (s *ArchitectureService) ValidateListOrganizationComponentsOpt(opt *ArchitectureListOrganizationComponentsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	return nil
}

// ValidateListOrganizationPlaceholdersOpt validates the options for the ListOrganizationPlaceholders method.
func (s *ArchitectureService) ValidateListOrganizationPlaceholdersOpt(opt *ArchitectureListOrganizationPlaceholdersOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	return nil
}

// ValidateListPatternsOpt validates the options for the ListPatterns method.
func (s *ArchitectureService) ValidateListPatternsOpt(opt *ArchitectureListPatternsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	return nil
}

// ValidateListProjectRelationshipsOpt validates the options for the ListProjectRelationships method.
func (s *ArchitectureService) ValidateListProjectRelationshipsOpt(opt *ArchitectureListProjectRelationshipsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectId, "ProjectId")
	if err != nil {
		return err
	}

	if opt.PageIndex < 0 {
		return NewValidationError("PageIndex", "must be greater than 0", ErrOutOfRange)
	}

	return nil
}

// ValidateListSdksOpt validates the options for the ListSdks method.
func (s *ArchitectureService) ValidateListSdksOpt(opt *ArchitectureListSdksOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	return nil
}

// ValidateUpdateBoundaryDescriptorOpt validates the options for the UpdateBoundaryDescriptor method.
func (s *ArchitectureService) ValidateUpdateBoundaryDescriptorOpt(opt *ArchitectureUpdateBoundaryDescriptorOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	var err error

	if opt.Direction != nil {
		err = IsValueAuthorized(*opt.Direction, allowedArchitectureUpdateBoundaryDescriptorDirection, "Direction")
		if err != nil {
			return err
		}
	}

	if opt.Ecosystem != nil {
		err = IsValueAuthorized(*opt.Ecosystem, allowedArchitectureUpdateBoundaryDescriptorEcosystem, "Ecosystem")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateUpdateOrganizationPlaceholderOpt validates the options for the UpdateOrganizationPlaceholder method.
func (s *ArchitectureService) ValidateUpdateOrganizationPlaceholderOpt(opt *ArchitectureUpdateOrganizationPlaceholderOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	var err error

	if opt.Type != nil {
		err = IsValueAuthorized(*opt.Type, allowedArchitectureUpdateOrganizationPlaceholderType, "Type")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateCreateBoundaryDescriptorOpt validates the options for the CreateBoundaryDescriptor method.
func (s *ArchitectureService) ValidateCreateBoundaryDescriptorOpt(opt *ArchitectureCreateBoundaryDescriptorOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Ecosystem, "Ecosystem")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Ecosystem, allowedArchitectureCreateBoundaryDescriptorEcosystem, "Ecosystem")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.KeyTemplate, "KeyTemplate")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Query, "Query")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Direction, "Direction")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Direction, allowedArchitectureCreateBoundaryDescriptorDirection, "Direction")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectId, "ProjectId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreateOrganizationPlaceholderOpt validates the options for the CreateOrganizationPlaceholder method.
func (s *ArchitectureService) ValidateCreateOrganizationPlaceholderOpt(opt *ArchitectureCreateOrganizationPlaceholderOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.DisplayName, "DisplayName")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Type, allowedArchitectureCreateOrganizationPlaceholderType, "Type")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreatePatternOpt validates the options for the CreatePattern method.
func (s *ArchitectureService) ValidateCreatePatternOpt(opt *ArchitectureCreatePatternOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreateProjectRelationshipOpt validates the options for the CreateProjectRelationship method.
func (s *ArchitectureService) ValidateCreateProjectRelationshipOpt(opt *ArchitectureCreateProjectRelationshipOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.BoundaryKey, "BoundaryKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectId, "ProjectId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.TargetComponentId, "TargetComponentId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.TargetEntryPointKey, "TargetEntryPointKey")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreateSdkOpt validates the options for the CreateSdk method.
func (s *ArchitectureService) ValidateCreateSdkOpt(opt *ArchitectureCreateSdkOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.TargetComponentId, "TargetComponentId")
	if err != nil {
		return err
	}

	if len(opt.SdkBoundaryDescriptors) == 0 {
		return NewValidationError("SdkBoundaryDescriptors", "is required", ErrMissingRequired)
	}

	err = ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.SdkKey, "SdkKey")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// DeleteBoundaryDescriptor deletes a boundary descriptor. Delete a boundary descriptor by id.
//
// API endpoint: DELETE /api/v2/architecture/boundary-descriptors/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) DeleteBoundaryDescriptor(ctx context.Context, boundaryDescriptorID string) (*http.Response, error) {
	err := ValidateRequired(boundaryDescriptorID, "boundaryDescriptorID")
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "architecture/boundary-descriptors/"+url.PathEscape(boundaryDescriptorID), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DeleteOrganizationPlaceholder deletes a placeholder. Delete the placeholder with the given id. A
// placeholder that is referenced by a project relationship cannot be deleted.
//
// API endpoint: DELETE /api/v2/architecture/organization-placeholders/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) DeleteOrganizationPlaceholder(ctx context.Context, organizationPlaceholderID string) (*http.Response, error) {
	err := ValidateRequired(organizationPlaceholderID, "organizationPlaceholderID")
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "architecture/organization-placeholders/"+url.PathEscape(organizationPlaceholderID), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DeletePattern deletes a pattern. Delete a pattern by id.
//
// API endpoint: DELETE /api/v2/architecture/patterns/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) DeletePattern(ctx context.Context, patternID string) (*http.Response, error) {
	err := ValidateRequired(patternID, "patternID")
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "architecture/patterns/"+url.PathEscape(patternID), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DeleteProjectRelationship deletes a project relationship. Delete a project relationship by id.
//
// API endpoint: DELETE /api/v2/architecture/project-relationships/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) DeleteProjectRelationship(ctx context.Context, projectRelationshipID string) (*http.Response, error) {
	err := ValidateRequired(projectRelationshipID, "projectRelationshipID")
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "architecture/project-relationships/"+url.PathEscape(projectRelationshipID), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DeleteSdk deletes an SDK. Delete the SDK with the given id, together with its boundary
// descriptors.
//
// API endpoint: DELETE /api/v2/architecture/sdks/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) DeleteSdk(ctx context.Context, sdkID string) (*http.Response, error) {
	err := ValidateRequired(sdkID, "sdkID")
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "architecture/sdks/"+url.PathEscape(sdkID), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// ListBoundaryDescriptors lists boundary descriptors for a project. Return all boundary descriptors
// for a project sorted by their uuid. May be empty.
//
// API endpoint: GET /api/v2/architecture/boundary-descriptors.
// Enterprise Edition only.
func (s *ArchitectureService) ListBoundaryDescriptors(ctx context.Context, opt *ArchitectureListBoundaryDescriptorsOptions) (*ArchitectureBoundaryDescriptorPagedResponse, *http.Response, error) {
	err := s.ValidateListBoundaryDescriptorsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/boundary-descriptors", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureBoundaryDescriptorPagedResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetBoundaryDescriptor gets a boundary descriptor by id. Return a single boundary descriptor by
// id.
//
// API endpoint: GET /api/v2/architecture/boundary-descriptors/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) GetBoundaryDescriptor(ctx context.Context, boundaryDescriptorID string) (*ArchitectureBoundaryDescriptor, *http.Response, error) {
	err := ValidateRequired(boundaryDescriptorID, "boundaryDescriptorID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/boundary-descriptors/"+url.PathEscape(boundaryDescriptorID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureBoundaryDescriptor)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListExternalInterfaces lists external interfaces. Allow searching for external interfaces within
// an organization, filtering by direction or organization component ID if provided.
//
// API endpoint: GET /api/v2/architecture/external-interfaces.
// Enterprise Edition only.
func (s *ArchitectureService) ListExternalInterfaces(ctx context.Context, opt *ArchitectureListExternalInterfacesOptions) (*ArchitectureExternalInterfaceListResponse, *http.Response, error) {
	err := s.ValidateListExternalInterfacesOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/external-interfaces", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureExternalInterfaceListResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetOrganizationArchitecture gets current architecture for an organization. fetch the current
// architecture for the organization.
//
// API endpoint: GET /api/v2/architecture/organization-architectures.
// Enterprise Edition only.
func (s *ArchitectureService) GetOrganizationArchitecture(ctx context.Context, opt *ArchitectureGetOrganizationArchitectureOptions) (*ArchitectureOrganizationArchitectureResponse, *http.Response, error) {
	err := s.ValidateGetOrganizationArchitectureOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/organization-architectures", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureOrganizationArchitectureResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListOrganizationComponents lists organization components. Return the organization components
// (projects and placeholders) for the given organization. Placeholders are included only when the
// placeholder feature is enabled for the organization.
//
// API endpoint: GET /api/v2/architecture/organization-components.
// Enterprise Edition only.
func (s *ArchitectureService) ListOrganizationComponents(ctx context.Context, opt *ArchitectureListOrganizationComponentsOptions) (*ArchitectureOrganizationComponentListResponse, *http.Response, error) {
	err := s.ValidateListOrganizationComponentsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/organization-components", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureOrganizationComponentListResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetOrganizationComponent gets an organization component by id. Return the organization component
// (project or placeholder) with the given id.
//
// API endpoint: GET /api/v2/architecture/organization-components/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) GetOrganizationComponent(ctx context.Context, organizationComponentID string) (*ArchitectureOrganizationComponent, *http.Response, error) {
	err := ValidateRequired(organizationComponentID, "organizationComponentID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/organization-components/"+url.PathEscape(organizationComponentID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureOrganizationComponent)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListOrganizationPlaceholders lists organization placeholders. Return the placeholders for the
// given organization.
//
// API endpoint: GET /api/v2/architecture/organization-placeholders.
// Enterprise Edition only.
func (s *ArchitectureService) ListOrganizationPlaceholders(ctx context.Context, opt *ArchitectureListOrganizationPlaceholdersOptions) (*ArchitectureOrganizationPlaceholderListResponse, *http.Response, error) {
	err := s.ValidateListOrganizationPlaceholdersOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/organization-placeholders", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureOrganizationPlaceholderListResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetOrganizationPlaceholder gets a placeholder by id. Return the placeholder with the given id.
//
// API endpoint: GET /api/v2/architecture/organization-placeholders/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) GetOrganizationPlaceholder(ctx context.Context, organizationPlaceholderID string) (*ArchitectureOrganizationPlaceholder, *http.Response, error) {
	err := ValidateRequired(organizationPlaceholderID, "organizationPlaceholderID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/organization-placeholders/"+url.PathEscape(organizationPlaceholderID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureOrganizationPlaceholder)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListPatterns lists patterns for an organization. Return all patterns for an organization. May be
// empty.
//
// API endpoint: GET /api/v2/architecture/patterns.
// Enterprise Edition only.
func (s *ArchitectureService) ListPatterns(ctx context.Context, opt *ArchitectureListPatternsOptions) ([]ArchitecturePattern, *http.Response, error) {
	err := s.ValidateListPatternsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/patterns", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []ArchitecturePattern

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetPattern gets a pattern by id. Return a single pattern by id.
//
// API endpoint: GET /api/v2/architecture/patterns/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) GetPattern(ctx context.Context, patternID string) (*ArchitecturePattern, *http.Response, error) {
	err := ValidateRequired(patternID, "patternID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/patterns/"+url.PathEscape(patternID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitecturePattern)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListProjectRelationships lists project relationships for a project. Return all project
// relationships for a project. May be empty.
//
// API endpoint: GET /api/v2/architecture/project-relationships.
// Enterprise Edition only.
func (s *ArchitectureService) ListProjectRelationships(ctx context.Context, opt *ArchitectureListProjectRelationshipsOptions) (*ArchitectureProjectRelationshipPagedResponse, *http.Response, error) {
	err := s.ValidateListProjectRelationshipsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/project-relationships", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureProjectRelationshipPagedResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetProjectRelationship gets a project relationship by id. Return a single project relationship by
// id.
//
// API endpoint: GET /api/v2/architecture/project-relationships/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) GetProjectRelationship(ctx context.Context, projectRelationshipID string) (*ArchitectureProjectRelationship, *http.Response, error) {
	err := ValidateRequired(projectRelationshipID, "projectRelationshipID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/project-relationships/"+url.PathEscape(projectRelationshipID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureProjectRelationship)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListSdks lists SDKs. Return the SDKs defined for the given organization.
//
// API endpoint: GET /api/v2/architecture/sdks.
// Enterprise Edition only.
func (s *ArchitectureService) ListSdks(ctx context.Context, opt *ArchitectureListSdksOptions) (*ArchitectureSdkListResponse, *http.Response, error) {
	err := s.ValidateListSdksOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/sdks", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureSdkListResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSdk gets an SDK by id. Return the SDK with the given id.
//
// API endpoint: GET /api/v2/architecture/sdks/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) GetSdk(ctx context.Context, sdkID string) (*ArchitectureSdk, *http.Response, error) {
	err := ValidateRequired(sdkID, "sdkID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/sdks/"+url.PathEscape(sdkID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureSdk)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateBoundaryDescriptor updates a boundary descriptor. Update a boundary descriptor. All mutable
// fields must be present in the body.
//
// API endpoint: PATCH /api/v2/architecture/boundary-descriptors/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) UpdateBoundaryDescriptor(ctx context.Context, boundaryDescriptorID string, opt *ArchitectureUpdateBoundaryDescriptorOptions) (*ArchitectureBoundaryDescriptor, *http.Response, error) {
	err := ValidateRequired(boundaryDescriptorID, "boundaryDescriptorID")
	if err != nil {
		return nil, nil, err
	}

	err = s.ValidateUpdateBoundaryDescriptorOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "architecture/boundary-descriptors/"+url.PathEscape(boundaryDescriptorID), nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureBoundaryDescriptor)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateOrganizationPlaceholder updates a placeholder. Update the mutable fields of a placeholder.
//
// API endpoint: PATCH /api/v2/architecture/organization-placeholders/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) UpdateOrganizationPlaceholder(ctx context.Context, organizationPlaceholderID string, opt *ArchitectureUpdateOrganizationPlaceholderOptions) (*ArchitectureOrganizationPlaceholder, *http.Response, error) {
	err := ValidateRequired(organizationPlaceholderID, "organizationPlaceholderID")
	if err != nil {
		return nil, nil, err
	}

	err = s.ValidateUpdateOrganizationPlaceholderOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "architecture/organization-placeholders/"+url.PathEscape(organizationPlaceholderID), nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureOrganizationPlaceholder)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdatePattern updates a pattern. Update pattern data and/or name.
//
// API endpoint: PATCH /api/v2/architecture/patterns/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) UpdatePattern(ctx context.Context, patternID string, opt *ArchitectureUpdatePatternOptions) (*ArchitecturePattern, *http.Response, error) {
	err := ValidateRequired(patternID, "patternID")
	if err != nil {
		return nil, nil, err
	}

	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "architecture/patterns/"+url.PathEscape(patternID), nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitecturePattern)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateProjectRelationship updates a project relationship. Update project relationship fields.
//
// API endpoint: PATCH /api/v2/architecture/project-relationships/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) UpdateProjectRelationship(ctx context.Context, projectRelationshipID string, opt *ArchitectureUpdateProjectRelationshipOptions) (*ArchitectureProjectRelationship, *http.Response, error) {
	err := ValidateRequired(projectRelationshipID, "projectRelationshipID")
	if err != nil {
		return nil, nil, err
	}

	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "architecture/project-relationships/"+url.PathEscape(projectRelationshipID), nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureProjectRelationship)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateSdk updates an SDK. Update the mutable fields of an SDK.
//
// API endpoint: PATCH /api/v2/architecture/sdks/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) UpdateSdk(ctx context.Context, sdkID string, opt *ArchitectureUpdateSdkOptions) (*ArchitectureSdk, *http.Response, error) {
	err := ValidateRequired(sdkID, "sdkID")
	if err != nil {
		return nil, nil, err
	}

	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "architecture/sdks/"+url.PathEscape(sdkID), nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureSdk)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateBoundaryDescriptor stores a boundary descriptor.
//
// API endpoint: POST /api/v2/architecture/boundary-descriptors.
// Enterprise Edition only.
func (s *ArchitectureService) CreateBoundaryDescriptor(ctx context.Context, opt *ArchitectureCreateBoundaryDescriptorOptions) (*ArchitectureBoundaryDescriptor, *http.Response, error) {
	err := s.ValidateCreateBoundaryDescriptorOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "architecture/boundary-descriptors", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureBoundaryDescriptor)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateOrganizationPlaceholder creates a placeholder. Create a new placeholder for the given
// organization.
//
// API endpoint: POST /api/v2/architecture/organization-placeholders.
// Enterprise Edition only.
func (s *ArchitectureService) CreateOrganizationPlaceholder(ctx context.Context, opt *ArchitectureCreateOrganizationPlaceholderOptions) (*ArchitectureOrganizationPlaceholder, *http.Response, error) {
	err := s.ValidateCreateOrganizationPlaceholderOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "architecture/organization-placeholders", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureOrganizationPlaceholder)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreatePattern creates a pattern. Create a pattern and return it including the generated id.
//
// API endpoint: POST /api/v2/architecture/patterns.
// Enterprise Edition only.
func (s *ArchitectureService) CreatePattern(ctx context.Context, opt *ArchitectureCreatePatternOptions) (*ArchitecturePattern, *http.Response, error) {
	err := s.ValidateCreatePatternOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "architecture/patterns", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitecturePattern)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateProjectRelationship stores a project relationship.
//
// API endpoint: POST /api/v2/architecture/project-relationships.
// Enterprise Edition only.
func (s *ArchitectureService) CreateProjectRelationship(ctx context.Context, opt *ArchitectureCreateProjectRelationshipOptions) (*ArchitectureProjectRelationship, *http.Response, error) {
	err := s.ValidateCreateProjectRelationshipOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "architecture/project-relationships", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureProjectRelationship)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateSdk creates an SDK. Create a new SDK for the given organization.
//
// API endpoint: POST /api/v2/architecture/sdks.
// Enterprise Edition only.
func (s *ArchitectureService) CreateSdk(ctx context.Context, opt *ArchitectureCreateSdkOptions) (*ArchitectureSdk, *http.Response, error) {
	err := s.ValidateCreateSdkOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "architecture/sdks", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureSdk)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
