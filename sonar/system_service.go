package sonargo

import (
	"net/http"
)

// -----------------------------------------------------------------------------
// Constants
// -----------------------------------------------------------------------------

// Allowed log levels for the ChangeLogLevel method.
//
//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	// allowedLogLevels is the set of allowed log levels for system logging.
	allowedLogLevels = map[string]struct{}{
		"TRACE": {},
		"DEBUG": {},
		"INFO":  {},
	}

	// allowedLogNames is the set of allowed log names for the Logs method.
	allowedLogNames = map[string]struct{}{
		"access":      {},
		"app":         {},
		"ce":          {},
		"deprecation": {},
		"es":          {},
		"web":         {},
	}
)

// -----------------------------------------------------------------------------
// Service
// -----------------------------------------------------------------------------

// SystemService handles communication with the System related methods of the SonarQube API.
// This service provides endpoints for system management, health checks, and status monitoring.
type SystemService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// SystemDbMigrationStatus represents the response from getting the database migration status.
// State values indicate the current status of the database migration process.
type SystemDbMigrationStatus struct {
	// Message contains additional information about the migration status.
	Message string `json:"message,omitempty"`
	// StartedAt is the timestamp when the migration started.
	StartedAt string `json:"startedAt,omitempty"`
	// State is the current migration state.
	// Possible values: NO_MIGRATION, NOT_SUPPORTED, MIGRATION_RUNNING, MIGRATION_SUCCEEDED, MIGRATION_FAILED, MIGRATION_REQUIRED.
	State string `json:"state,omitempty"`
}

// SystemHealth represents the response from checking the health status of SonarQube.
type SystemHealth struct {
	// Causes contains the reasons for any health issues.
	Causes []HealthCause `json:"causes,omitempty"`
	// Health is the overall health status: GREEN, YELLOW, or RED.
	Health string `json:"health,omitempty"`
	// Nodes contains health information for each node in a cluster.
	Nodes []HealthNode `json:"nodes,omitempty"`
}

// HealthCause represents a reason for a health issue.
type HealthCause struct {
	// Message describes the cause of the health issue.
	Message string `json:"message,omitempty"`
}

// HealthNode represents health information for a single node in a cluster.
//
//nolint:govet // Field alignment is less important than logical grouping
type HealthNode struct {
	// Causes contains the reasons for any health issues on this node.
	Causes []HealthCause `json:"causes,omitempty"`
	// Health is the health status of this node: GREEN, YELLOW, or RED.
	Health string `json:"health,omitempty"`
	// Host is the hostname or IP address of this node.
	Host string `json:"host,omitempty"`
	// Name is the display name of this node.
	Name string `json:"name,omitempty"`
	// Port is the port number of this node.
	Port int64 `json:"port,omitempty"`
	// StartedAt is the timestamp when this node was started.
	StartedAt string `json:"startedAt,omitempty"`
	// Type is the type of this node (e.g., APPLICATION, SEARCH).
	Type string `json:"type,omitempty"`
}

// SystemInfo represents the response from getting detailed system information.
// This contains extensive information about the SonarQube configuration and state.
//
//nolint:govet,tagliatelle // Field alignment is less important than logical grouping; JSON tags match SonarQube API
type SystemInfo struct {
	// ALMs contains ALM integration configuration.
	ALMs SystemInfoALMs `json:"ALMs,omitzero"`
	// Bundled contains versions of bundled plugins.
	Bundled SystemInfoBundled `json:"Bundled,omitzero"`
	// ComputeEngineDatabaseConnection contains CE database connection pool info.
	ComputeEngineDatabaseConnection DatabaseConnectionPool `json:"Compute Engine Database Connection,omitzero"`
	// ComputeEngineJVMProperties contains CE JVM properties.
	ComputeEngineJVMProperties map[string]any `json:"Compute Engine JVM Properties,omitempty"`
	// ComputeEngineJVMState contains CE JVM state information.
	ComputeEngineJVMState JVMState `json:"Compute Engine JVM State,omitzero"`
	// ComputeEngineLogging contains CE logging configuration.
	ComputeEngineLogging LoggingConfig `json:"Compute Engine Logging,omitzero"`
	// ComputeEngineTasks contains CE task information.
	ComputeEngineTasks ComputeEngineTasks `json:"Compute Engine Tasks,omitzero"`
	// Database contains database information.
	Database DatabaseInfo `json:"Database,omitzero"`
	// Health is the overall health status.
	Health string `json:"Health,omitempty"`
	// HealthCauses contains reasons for any health issues.
	HealthCauses []any `json:"Health Causes,omitempty"`
	// Plugins contains installed plugin information.
	Plugins map[string]string `json:"Plugins,omitempty"`
	// SearchIndexes contains search index statistics.
	SearchIndexes SearchIndexes `json:"Search Indexes,omitzero"`
	// SearchState contains search node state information.
	SearchState SearchState `json:"Search State,omitzero"`
	// ServerPushConnections contains server push connection information.
	ServerPushConnections ServerPushConnections `json:"Server Push Connections,omitzero"`
	// Settings contains SonarQube settings.
	Settings map[string]any `json:"Settings,omitempty"`
	// System contains core system information.
	System SystemInfoSystem `json:"System,omitzero"`
	// WebDatabaseConnection contains Web database connection pool info.
	WebDatabaseConnection DatabaseConnectionPool `json:"Web Database Connection,omitzero"`
	// WebJVMProperties contains Web JVM properties.
	WebJVMProperties map[string]any `json:"Web JVM Properties,omitempty"`
	// WebJVMState contains Web JVM state information.
	WebJVMState JVMState `json:"Web JVM State,omitzero"`
	// WebLogging contains Web logging configuration.
	WebLogging LoggingConfig `json:"Web Logging,omitzero"`
}

// SystemInfoALMs contains ALM integration configuration.
//
//nolint:tagliatelle // JSON tags match SonarQube API field names
type SystemInfoALMs struct {
	// GithubConfig contains GitHub configuration.
	GithubConfig string `json:"Github Config,omitempty"`
	// GitlabConfig contains GitLab configuration.
	GitlabConfig string `json:"Gitlab Config,omitempty"`
}

// SystemInfoBundled contains versions of bundled plugins.
type SystemInfoBundled struct {
	// Config is the config plugin version.
	Config string `json:"config,omitempty"`
	// Csharp is the C# plugin version.
	Csharp string `json:"csharp,omitempty"`
	// Flex is the Flex plugin version.
	Flex string `json:"flex,omitempty"`
	// Go is the Go plugin version.
	Go string `json:"go,omitempty"`
	// Iac is the IaC plugin version.
	Iac string `json:"iac,omitempty"`
	// Jacoco is the JaCoCo plugin version.
	Jacoco string `json:"jacoco,omitempty"`
	// Java is the Java plugin version.
	Java string `json:"java,omitempty"`
	// Javascript is the JavaScript plugin version.
	Javascript string `json:"javascript,omitempty"`
	// Kotlin is the Kotlin plugin version.
	Kotlin string `json:"kotlin,omitempty"`
	// Php is the PHP plugin version.
	Php string `json:"php,omitempty"`
	// Python is the Python plugin version.
	Python string `json:"python,omitempty"`
	// Ruby is the Ruby plugin version.
	Ruby string `json:"ruby,omitempty"`
	// Sonarscala is the Scala plugin version.
	Sonarscala string `json:"sonarscala,omitempty"`
	// Text is the Text plugin version.
	Text string `json:"text,omitempty"`
	// Vbnet is the VB.NET plugin version.
	Vbnet string `json:"vbnet,omitempty"`
	// Web is the Web plugin version.
	Web string `json:"web,omitempty"`
	// XML is the XML plugin version.
	XML string `json:"xml,omitempty"`
}

// DatabaseConnectionPool contains database connection pool information.
//
//nolint:tagliatelle // JSON tags match SonarQube API field names
type DatabaseConnectionPool struct {
	// PoolActiveConnections is the number of active connections.
	PoolActiveConnections int64 `json:"Pool Active Connections,omitempty"`
	// PoolIdleConnections is the number of idle connections.
	PoolIdleConnections int64 `json:"Pool Idle Connections,omitempty"`
	// PoolMaxConnections is the maximum number of connections.
	PoolMaxConnections int64 `json:"Pool Max Connections,omitempty"`
	// PoolMaxLifetimeMs is the maximum lifetime of a connection in milliseconds.
	PoolMaxLifetimeMs int64 `json:"Pool Max Lifetime (ms),omitempty"`
	// PoolMaxWaitMs is the maximum wait time for a connection in milliseconds.
	PoolMaxWaitMs int64 `json:"Pool Max Wait (ms),omitempty"`
	// PoolMinIdleConnections is the minimum number of idle connections.
	PoolMinIdleConnections int64 `json:"Pool Min Idle Connections,omitempty"`
	// PoolTotalConnections is the total number of connections.
	PoolTotalConnections int64 `json:"Pool Total Connections,omitempty"`
}

// JVMState contains JVM state information.
//
//nolint:govet,tagliatelle // Field alignment is less important than logical grouping; JSON tags match SonarQube API
type JVMState struct {
	// FreeMemoryMB is the free memory in megabytes.
	FreeMemoryMB int64 `json:"Free Memory (MB),omitempty"`
	// HeapCommittedMB is the committed heap memory in megabytes.
	HeapCommittedMB int64 `json:"Heap Committed (MB),omitempty"`
	// HeapInitMB is the initial heap memory in megabytes.
	HeapInitMB int64 `json:"Heap Init (MB),omitempty"`
	// HeapMaxMB is the maximum heap memory in megabytes.
	HeapMaxMB int64 `json:"Heap Max (MB),omitempty"`
	// HeapUsedMB is the used heap memory in megabytes.
	HeapUsedMB int64 `json:"Heap Used (MB),omitempty"`
	// MaxMemoryMB is the maximum memory in megabytes.
	MaxMemoryMB int64 `json:"Max Memory (MB),omitempty"`
	// NonHeapCommittedMB is the committed non-heap memory in megabytes.
	NonHeapCommittedMB int64 `json:"Non Heap Committed (MB),omitempty"`
	// NonHeapInitMB is the initial non-heap memory in megabytes.
	NonHeapInitMB int64 `json:"Non Heap Init (MB),omitempty"`
	// NonHeapUsedMB is the used non-heap memory in megabytes.
	NonHeapUsedMB int64 `json:"Non Heap Used (MB),omitempty"`
	// SystemLoadAverage is the system load average.
	SystemLoadAverage string `json:"System Load Average,omitempty"`
	// Threads is the number of threads.
	Threads int64 `json:"Threads,omitempty"`
}

// LoggingConfig contains logging configuration.
//
//nolint:tagliatelle // JSON tags match SonarQube API field names
type LoggingConfig struct {
	// LogsDir is the directory where logs are stored.
	LogsDir string `json:"Logs Dir,omitempty"`
	// LogsLevel is the current log level.
	LogsLevel string `json:"Logs Level,omitempty"`
}

// ComputeEngineTasks contains Compute Engine task information.
//
//nolint:tagliatelle // JSON tags match SonarQube API field names
type ComputeEngineTasks struct {
	// InProgress is the number of tasks currently in progress.
	InProgress int64 `json:"In Progress,omitempty"`
	// LongestTimePendingMs is the longest time a task has been pending in milliseconds.
	LongestTimePendingMs int64 `json:"Longest Time Pending (ms),omitempty"`
	// MaxWorkerCount is the maximum number of worker threads.
	MaxWorkerCount int64 `json:"Max Worker Count,omitempty"`
	// Pending is the number of pending tasks.
	Pending int64 `json:"Pending,omitempty"`
	// ProcessedWithError is the number of tasks that completed with errors.
	ProcessedWithError int64 `json:"Processed With Error,omitempty"`
	// ProcessedWithSuccess is the number of tasks that completed successfully.
	ProcessedWithSuccess int64 `json:"Processed With Success,omitempty"`
	// ProcessingTimeMs is the total processing time in milliseconds.
	ProcessingTimeMs int64 `json:"Processing Time (ms),omitempty"`
	// WorkerCount is the current number of worker threads.
	WorkerCount int64 `json:"Worker Count,omitempty"`
	// WorkersPaused indicates if workers are currently paused.
	WorkersPaused bool `json:"Workers Paused,omitempty"`
}

// DatabaseInfo contains database information.
//
//nolint:tagliatelle // JSON tags match SonarQube API field names
type DatabaseInfo struct {
	// Database is the database name.
	Database string `json:"Database,omitempty"`
	// DatabaseVersion is the database version.
	DatabaseVersion string `json:"Database Version,omitempty"`
	// DefaultTransactionIsolation is the default transaction isolation level.
	DefaultTransactionIsolation string `json:"Default transaction isolation,omitempty"`
	// Driver is the database driver name.
	Driver string `json:"Driver,omitempty"`
	// DriverVersion is the database driver version.
	DriverVersion string `json:"Driver Version,omitempty"`
	// URL is the database connection URL.
	URL string `json:"URL,omitempty"`
	// Username is the database username.
	Username string `json:"Username,omitempty"`
}

// SearchIndexes contains search index statistics.
//
//nolint:govet,tagliatelle // Field alignment is less important than logical grouping; JSON tags match SonarQube API
type SearchIndexes struct {
	// IndexComponentsDocs is the number of documents in the components index.
	IndexComponentsDocs int64 `json:"Index components - Docs,omitempty"`
	// IndexComponentsShards is the number of shards in the components index.
	IndexComponentsShards int64 `json:"Index components - Shards,omitempty"`
	// IndexComponentsStoreSize is the store size of the components index.
	IndexComponentsStoreSize string `json:"Index components - Store Size,omitempty"`
	// IndexIssuesDocs is the number of documents in the issues index.
	IndexIssuesDocs int64 `json:"Index issues - Docs,omitempty"`
	// IndexIssuesShards is the number of shards in the issues index.
	IndexIssuesShards int64 `json:"Index issues - Shards,omitempty"`
	// IndexIssuesStoreSize is the store size of the issues index.
	IndexIssuesStoreSize string `json:"Index issues - Store Size,omitempty"`
	// IndexMetadatasDocs is the number of documents in the metadatas index.
	IndexMetadatasDocs int64 `json:"Index metadatas - Docs,omitempty"`
	// IndexMetadatasShards is the number of shards in the metadatas index.
	IndexMetadatasShards int64 `json:"Index metadatas - Shards,omitempty"`
	// IndexMetadatasStoreSize is the store size of the metadatas index.
	IndexMetadatasStoreSize string `json:"Index metadatas - Store Size,omitempty"`
	// IndexProjectmeasuresDocs is the number of documents in the projectmeasures index.
	IndexProjectmeasuresDocs int64 `json:"Index projectmeasures - Docs,omitempty"`
	// IndexProjectmeasuresShards is the number of shards in the projectmeasures index.
	IndexProjectmeasuresShards int64 `json:"Index projectmeasures - Shards,omitempty"`
	// IndexProjectmeasuresStoreSize is the store size of the projectmeasures index.
	IndexProjectmeasuresStoreSize string `json:"Index projectmeasures - Store Size,omitempty"`
	// IndexRulesDocs is the number of documents in the rules index.
	IndexRulesDocs int64 `json:"Index rules - Docs,omitempty"`
	// IndexRulesShards is the number of shards in the rules index.
	IndexRulesShards int64 `json:"Index rules - Shards,omitempty"`
	// IndexRulesStoreSize is the store size of the rules index.
	IndexRulesStoreSize string `json:"Index rules - Store Size,omitempty"`
	// IndexUsersDocs is the number of documents in the users index.
	IndexUsersDocs int64 `json:"Index users - Docs,omitempty"`
	// IndexUsersShards is the number of shards in the users index.
	IndexUsersShards int64 `json:"Index users - Shards,omitempty"`
	// IndexUsersStoreSize is the store size of the users index.
	IndexUsersStoreSize string `json:"Index users - Store Size,omitempty"`
	// IndexViewsDocs is the number of documents in the views index.
	IndexViewsDocs int64 `json:"Index views - Docs,omitempty"`
	// IndexViewsShards is the number of shards in the views index.
	IndexViewsShards int64 `json:"Index views - Shards,omitempty"`
	// IndexViewsStoreSize is the store size of the views index.
	IndexViewsStoreSize string `json:"Index views - Store Size,omitempty"`
}

// SearchState contains search node state information.
//
//nolint:govet,tagliatelle // Field alignment is less important than logical grouping; JSON tags match SonarQube API
type SearchState struct {
	// CPUUsage is the CPU usage percentage.
	CPUUsage int64 `json:"CPU Usage (%),omitempty"`
	// DiskAvailable is the available disk space.
	DiskAvailable string `json:"Disk Available,omitempty"`
	// FieldDataCircuitBreakerEstimation is the field data circuit breaker estimation.
	FieldDataCircuitBreakerEstimation string `json:"Field Data Circuit Breaker Estimation,omitempty"`
	// FieldDataCircuitBreakerLimit is the field data circuit breaker limit.
	FieldDataCircuitBreakerLimit string `json:"Field Data Circuit Breaker Limit,omitempty"`
	// FieldDataMemory is the field data memory usage.
	FieldDataMemory string `json:"Field Data Memory,omitempty"`
	// JVMHeapMax is the maximum JVM heap size.
	JVMHeapMax string `json:"JVM Heap Max,omitempty"`
	// JVMHeapUsage is the JVM heap usage percentage.
	JVMHeapUsage string `json:"JVM Heap Usage,omitempty"`
	// JVMHeapUsed is the used JVM heap size.
	JVMHeapUsed string `json:"JVM Heap Used,omitempty"`
	// JVMNonHeapUsed is the used JVM non-heap size.
	JVMNonHeapUsed string `json:"JVM Non Heap Used,omitempty"`
	// JVMThreads is the number of JVM threads.
	JVMThreads int64 `json:"JVM Threads,omitempty"`
	// MaxFileDescriptors is the maximum number of file descriptors.
	MaxFileDescriptors int64 `json:"Max File Descriptors,omitempty"`
	// OpenFileDescriptors is the number of open file descriptors.
	OpenFileDescriptors int64 `json:"Open File Descriptors,omitempty"`
	// QueryCacheMemory is the query cache memory usage.
	QueryCacheMemory string `json:"Query Cache Memory,omitempty"`
	// RequestCacheMemory is the request cache memory usage.
	RequestCacheMemory string `json:"Request Cache Memory,omitempty"`
	// RequestCircuitBreakerEstimation is the request circuit breaker estimation.
	RequestCircuitBreakerEstimation string `json:"Request Circuit Breaker Estimation,omitempty"`
	// RequestCircuitBreakerLimit is the request circuit breaker limit.
	RequestCircuitBreakerLimit string `json:"Request Circuit Breaker Limit,omitempty"`
	// State is the search node state (GREEN, YELLOW, RED).
	State string `json:"State,omitempty"`
	// StoreSize is the total store size.
	StoreSize string `json:"Store Size,omitempty"`
	// TranslogSize is the translog size.
	TranslogSize string `json:"Translog Size,omitempty"`
}

// ServerPushConnections contains server push connection information.
//
//nolint:tagliatelle // JSON tags match SonarQube API field names
type ServerPushConnections struct {
	// SonarLintConnectedClients is the number of connected SonarLint clients.
	SonarLintConnectedClients int64 `json:"SonarLint Connected Clients,omitempty"`
}

// SystemInfoSystem contains core system information.
//
//nolint:govet,tagliatelle // Field alignment is less important than logical grouping; JSON tags match SonarQube API
type SystemInfoSystem struct {
	// AcceptedExternalIdentityProviders lists accepted external identity providers.
	AcceptedExternalIdentityProviders string `json:"Accepted external identity providers,omitempty"`
	// DataDir is the data directory path.
	DataDir string `json:"Data Dir,omitempty"`
	// Docker indicates if running in Docker.
	Docker bool `json:"Docker,omitempty"`
	// Edition is the SonarQube edition.
	Edition string `json:"Edition,omitempty"`
	// ExternalIdentityProvidersWhoseUsersAreAllowedToSignThemselvesUp describes which providers allow self-signup.
	ExternalIdentityProvidersWhoseUsersAreAllowedToSignThemselvesUp string `json:"External identity providers whose users are allowed to sign themselves up,omitempty"`
	// ForceAuthentication indicates if authentication is forced.
	ForceAuthentication bool `json:"Force authentication,omitempty"`
	// HighAvailability indicates if high availability mode is enabled.
	HighAvailability bool `json:"High Availability,omitempty"`
	// HomeDir is the home directory path.
	HomeDir string `json:"Home Dir,omitempty"`
	// OfficialDistribution indicates if this is an official distribution.
	OfficialDistribution bool `json:"Official Distribution,omitempty"`
	// Processors is the number of available processors.
	Processors int64 `json:"Processors,omitempty"`
	// ServerID is the server identifier.
	ServerID string `json:"Server ID,omitempty"`
	// TempDir is the temp directory path.
	TempDir string `json:"Temp Dir,omitempty"`
	// Version is the SonarQube version.
	Version string `json:"Version,omitempty"`
}

// SystemLiveness represents the response from the liveness check.
// Note: The API does not return a response body, liveness is indicated by HTTP status.
type SystemLiveness struct{}

// SystemMigrateDb represents the response from initiating a database migration.
// State values indicate the current status of the database migration process.
type SystemMigrateDb struct {
	// Message contains additional information about the migration.
	Message string `json:"message,omitempty"`
	// StartedAt is the timestamp when the migration started.
	StartedAt string `json:"startedAt,omitempty"`
	// State is the current migration state.
	State string `json:"state,omitempty"`
}

// SystemStatus represents the response from getting the system status.
type SystemStatus struct {
	// ID is the server identifier.
	ID string `json:"id,omitempty"`
	// Status is the running status.
	// Possible values: STARTING, UP, DOWN, RESTARTING, DB_MIGRATION_NEEDED, DB_MIGRATION_RUNNING.
	Status string `json:"status,omitempty"`
	// Version is the SonarQube version.
	Version string `json:"version,omitempty"`
}

// SystemUpgrades represents the response from checking for available upgrades.
//
//nolint:govet,tagliatelle // Field alignment is less important than logical grouping; JSON tags match SonarQube API
type SystemUpgrades struct {
	// InstalledVersionActive indicates if the installed version is still active.
	InstalledVersionActive bool `json:"installedVersionActive,omitempty"`
	// LatestLTA is the latest Long-Term Active version.
	LatestLTA string `json:"latestLTA,omitempty"`
	// LatestLTS is the latest Long-Term Support version (deprecated, use LatestLTA).
	LatestLTS string `json:"latestLTS,omitempty"`
	// UpdateCenterRefresh is the timestamp when Update Center was last refreshed.
	UpdateCenterRefresh string `json:"updateCenterRefresh,omitempty"`
	// Upgrades is the list of available upgrades.
	Upgrades []Upgrade `json:"upgrades,omitempty"`
}

// Upgrade represents an available upgrade.
//
//nolint:govet // Field alignment is less important than logical grouping
type Upgrade struct {
	// ChangeLogURL is the URL to the release changelog.
	ChangeLogURL string `json:"changeLogUrl,omitempty"`
	// Description describes the upgrade.
	Description string `json:"description,omitempty"`
	// DownloadURL is the URL to download the upgrade.
	DownloadURL string `json:"downloadUrl,omitempty"`
	// Plugins contains plugin compatibility information.
	Plugins UpgradePlugins `json:"plugins,omitzero"`
	// ReleaseDate is the release date of the upgrade.
	ReleaseDate string `json:"releaseDate,omitempty"`
	// Version is the version of the upgrade.
	Version string `json:"version,omitempty"`
}

// UpgradePlugins contains plugin compatibility information for an upgrade.
type UpgradePlugins struct {
	// Incompatible lists plugins that are incompatible with the upgrade.
	Incompatible []IncompatiblePlugin `json:"incompatible,omitempty"`
	// RequireUpdate lists plugins that require updating.
	RequireUpdate []PluginUpdate `json:"requireUpdate,omitempty"`
}

// IncompatiblePlugin represents a plugin that is incompatible with an upgrade.
//
//nolint:govet // Field alignment is less important than logical grouping
type IncompatiblePlugin struct {
	// Category is the plugin category.
	Category string `json:"category,omitempty"`
	// Description is the plugin description.
	Description string `json:"description,omitempty"`
	// EditionBundled indicates if the plugin is bundled with an edition.
	EditionBundled bool `json:"editionBundled,omitempty"`
	// Key is the plugin key.
	Key string `json:"key,omitempty"`
	// License is the plugin license.
	License string `json:"license,omitempty"`
	// Name is the plugin name.
	Name string `json:"name,omitempty"`
	// OrganizationName is the organization name.
	OrganizationName string `json:"organizationName,omitempty"`
	// OrganizationURL is the organization URL.
	OrganizationURL string `json:"organizationUrl,omitempty"`
}

// PluginUpdate represents a plugin that requires updating for an upgrade.
//
//nolint:govet // Field alignment is less important than logical grouping
type PluginUpdate struct {
	// Category is the plugin category.
	Category string `json:"category,omitempty"`
	// Description is the plugin description.
	Description string `json:"description,omitempty"`
	// EditionBundled indicates if the plugin is bundled with an edition.
	EditionBundled bool `json:"editionBundled,omitempty"`
	// Key is the plugin key.
	Key string `json:"key,omitempty"`
	// License is the plugin license.
	License string `json:"license,omitempty"`
	// Name is the plugin name.
	Name string `json:"name,omitempty"`
	// OrganizationName is the organization name.
	OrganizationName string `json:"organizationName,omitempty"`
	// OrganizationURL is the organization URL.
	OrganizationURL string `json:"organizationUrl,omitempty"`
	// TermsAndConditionsURL is the URL to terms and conditions.
	TermsAndConditionsURL string `json:"termsAndConditionsUrl,omitempty"`
	// Version is the required version.
	Version string `json:"version,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// SystemChangeLogLevelOption contains options for the ChangeLogLevel method.
type SystemChangeLogLevelOption struct {
	// Level is the new log level.
	// Possible values: TRACE, DEBUG, INFO.
	// Be cautious: DEBUG, and even more TRACE, may have performance impacts.
	Level string `url:"level,omitempty"`
}

// SystemLogsOption contains options for the Logs method.
type SystemLogsOption struct {
	// Name is the name of the logs to retrieve.
	// Possible values: access, app, ce, deprecation, es, web.
	// Default: app.
	Name string `url:"name,omitempty"`
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// ChangeLogLevel temporarily changes the level of logs.
// The new level is not persistent and is lost when restarting the server.
// Requires system administration permission.
//
// API Docs: https://next.sonarqube.com/sonarqube/web_api/api/system/change_log_level
func (s *SystemService) ChangeLogLevel(opt *SystemChangeLogLevelOption) (*http.Response, error) {
	err := s.ValidateChangeLogLevelOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "system/change_log_level", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DbMigrationStatus displays the database migration status of SonarQube.
// State values are:
//   - NO_MIGRATION: DB is up to date with current version of SonarQube.
//   - NOT_SUPPORTED: Migration is not supported on embedded databases.
//   - MIGRATION_RUNNING: DB migration is underway.
//   - MIGRATION_SUCCEEDED: DB migration has run and has been successful.
//   - MIGRATION_FAILED: DB migration has run and failed.
//   - MIGRATION_REQUIRED: DB migration is required.
//
// Deprecated: since 10.6. Use the API v2 version /api/v2/system/migrations-status instead.
//
// API Docs: https://next.sonarqube.com/sonarqube/web_api/api/system/db_migration_status
func (s *SystemService) DbMigrationStatus() (*SystemDbMigrationStatus, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "system/db_migration_status", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SystemDbMigrationStatus)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Health provides the health status of SonarQube.
// Although global health is calculated based on both application and search nodes,
// detailed information is returned only for application nodes.
//
// Health status values:
//   - GREEN: SonarQube is fully operational.
//   - YELLOW: SonarQube is usable, but it needs attention.
//   - RED: SonarQube is not operational.
//
// Requires the 'Administer System' permission or system passcode.
// When SonarQube is in safe mode, only authentication with a system passcode is supported.
//
// API Docs: https://next.sonarqube.com/sonarqube/web_api/api/system/health
func (s *SystemService) Health() (*SystemHealth, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "system/health", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SystemHealth)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Info gets detailed information about system configuration.
// Requires 'Administer' permissions.
//
// API Docs: https://next.sonarqube.com/sonarqube/web_api/api/system/info
func (s *SystemService) Info() (*SystemInfo, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "system/info", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SystemInfo)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Liveness provides the liveness status of SonarQube, meant to be used for a liveness probe on Kubernetes.
// Requires 'Administer System' permission or authentication with passcode.
//
// When SonarQube is fully started, liveness checks for database connectivity, Compute Engine status,
// and (except for DataCenter Edition) if ElasticSearch is Green or Yellow.
//
// When SonarQube is on Safe Mode, liveness checks only for database connectivity.
//
// Response codes:
//   - HTTP 204: this SonarQube node is alive.
//   - Any other HTTP code: this SonarQube node is not alive and should be rescheduled.
//
// API Docs: https://next.sonarqube.com/sonarqube/web_api/api/system/liveness
func (s *SystemService) Liveness() (*SystemLiveness, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "system/liveness", nil)
	if err != nil {
		return nil, nil, err
	}

	// The API indicates liveness via the HTTP status (204 No Content) and
	// does not return a response body. Avoid attempting to decode an empty
	// body which would result in EOF. We still return an empty result object
	// for API compatibility.
	result := new(SystemLiveness)

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Logs gets system logs in plain-text format.
// Requires system administration permission.
//
// API Docs: https://next.sonarqube.com/sonarqube/web_api/api/system/logs
func (s *SystemService) Logs(opt *SystemLogsOption) (*string, *http.Response, error) {
	err := s.ValidateLogsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "system/logs", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(string)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// MigrateDb migrates the database to match the current version of SonarQube.
// It is strongly advised to make a database backup before invoking this method.
//
// State values are:
//   - NO_MIGRATION: DB is up to date with current version of SonarQube.
//   - NOT_SUPPORTED: Migration is not supported on embedded databases.
//   - MIGRATION_RUNNING: DB migration is underway.
//   - MIGRATION_SUCCEEDED: DB migration has run and has been successful.
//   - MIGRATION_FAILED: DB migration has run and failed.
//   - MIGRATION_REQUIRED: DB migration is required.
//
// API Docs: https://next.sonarqube.com/sonarqube/web_api/api/system/migrate_db
func (s *SystemService) MigrateDb() (*SystemMigrateDb, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "system/migrate_db", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SystemMigrateDb)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Ping answers "pong" as plain-text.
// This can be used for health checks.
//
// API Docs: https://next.sonarqube.com/sonarqube/web_api/api/system/ping
func (s *SystemService) Ping() (*string, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "system/ping", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(string)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Restart restarts the server.
// Requires 'Administer System' permission.
// Performs a full restart of the Web, Search and Compute Engine Servers processes.
// Does not reload sonar.properties.
//
// API Docs: https://next.sonarqube.com/sonarqube/web_api/api/system/restart
func (s *SystemService) Restart() (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "system/restart", nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Status gets state information about SonarQube.
//
// Status values:
//   - STARTING: SonarQube Web Server is up and serving some Web Services but initialization is still ongoing.
//   - UP: SonarQube instance is up and running.
//   - DOWN: SonarQube instance is up but not running because migration has failed or some other reason.
//   - RESTARTING: SonarQube instance is still up but a restart has been requested.
//   - DB_MIGRATION_NEEDED: Database migration is required.
//   - DB_MIGRATION_RUNNING: DB migration is running.
//
// API Docs: https://next.sonarqube.com/sonarqube/web_api/api/system/status
func (s *SystemService) Status() (*SystemStatus, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "system/status", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SystemStatus)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Upgrades lists available upgrades for the SonarQube instance (if any).
// For each upgrade, it lists incompatible plugins and plugins requiring upgrade.
// Plugin information is retrieved from Update Center.
//
// API Docs: https://next.sonarqube.com/sonarqube/web_api/api/system/upgrades
func (s *SystemService) Upgrades() (*SystemUpgrades, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "system/upgrades", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SystemUpgrades)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateChangeLogLevelOpt validates the options for ChangeLogLevel.
func (s *SystemService) ValidateChangeLogLevelOpt(opt *SystemChangeLogLevelOption) error {
	if opt == nil {
		return NewValidationError("opt", "options cannot be nil", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Level, "Level")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Level, allowedLogLevels, "Level")
	if err != nil {
		return err
	}

	return nil
}

// ValidateLogsOpt validates the options for Logs.
func (s *SystemService) ValidateLogsOpt(opt *SystemLogsOption) error {
	// opt can be nil (uses defaults)
	if opt == nil {
		return nil
	}

	// If Name is provided, validate it
	if opt.Name != "" {
		err := IsValueAuthorized(opt.Name, allowedLogNames, "Name")
		if err != nil {
			return err
		}
	}

	return nil
}
