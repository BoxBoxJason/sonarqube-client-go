package cli

// serviceDescriptions maps service names to their short descriptions.
//
//nolint:gochecknoglobals // static description registry
var serviceDescriptions = map[string]string{ //nolint:gosec // G101 false positive: these are API descriptions, not credentials
	"AlmIntegrations":    "Manages DevOps platform integrations and project imports",
	"AlmSettings":        "Manages DevOps platform settings and bindings",
	"AnalysisCache":      "Manages scanner analysis cache data",
	"AnalysisReports":    "Manages Compute Engine analysis report queue",
	"Authentication":     "Manages user authentication sessions",
	"Batch":              "Provides scanner batch operation data",
	"Ce":                 "Manages Compute Engine tasks and workers",
	"Components":         "Manages project components and navigation",
	"Developers":         "Provides developer-specific event data",
	"DismissMessage":     "Manages dismissible UI messages",
	"Duplications":       "Manages code duplication data",
	"Emails":             "Manages email sending functionality",
	"Favorites":          "Manages user favorite components",
	"Features":           "Lists available SonarQube features",
	"GithubProvisioning": "Manages GitHub provisioning status",
	"Hotspots":           "Manages security hotspots",
	"Issues":             "Manages code issues and their lifecycle",
	"L10N":               "Manages localization and internationalization",
	"Languages":          "Lists programming languages",
	"Measures":           "Manages project measures and metrics data",
	"Metrics":            "Manages metric definitions",
	"Monitoring":         "Provides system monitoring data",
	"Navigation":         "Provides navigation-related data",
	"NewCodePeriods":     "Manages new code period definitions",
	"Notifications":      "Manages user notifications",
	"Permissions":        "Manages project and global permissions",
	"Plugins":            "Manages SonarQube plugins",
	"ProjectAnalyses":    "Manages project analysis events",
	"ProjectBadges":      "Manages project badge generation",
	"ProjectBranches":    "Manages project branches",
	"ProjectDump":        "Manages project data export and import",
	"ProjectLinks":       "Manages project links",
	"ProjectTags":        "Manages project tags",
	"Projects":           "Manages SonarQube projects",
	"Push":               "Manages server-sent events for SonarLint",
	"Qualitygates":       "Manages quality gates and conditions",
	"Qualityprofiles":    "Manages quality profiles and rule activation",
	"Rules":              "Manages coding rules",
	"Server":             "Provides SonarQube server information",
	"Settings":           "Manages SonarQube settings",
	"Sources":            "Manages source code viewing",
	"System":             "Manages SonarQube system administration",
	"UserGroups":         "Manages user groups",
	"UserTokens":         "Manages user authentication tokens",
	"Users":              "Manages SonarQube users",
	"Webhooks":           "Manages webhooks",
	"Webservices":        "Provides API metadata",
}

// methodDescriptions maps "ServiceName.MethodName" to short descriptions.
//
//nolint:gochecknoglobals // static description registry
var methodDescriptions = map[string]string{ //nolint:gosec // G101 false positive: these are API descriptions, not credentials
	// AlmIntegrations
	"AlmIntegrations.CheckPat":                     "Checks validity of a Personal Access Token",
	"AlmIntegrations.GetGithubClientId":            "Gets the client ID of a GitHub integration",
	"AlmIntegrations.ImportAzureProject":           "Creates a project from an Azure DevOps project",
	"AlmIntegrations.ImportBitbucketCloudRepo":     "Creates a project from a Bitbucket Cloud repo",
	"AlmIntegrations.ImportBitbucketServerProject": "Creates a project from a Bitbucket Server project",
	"AlmIntegrations.ImportGithubProject":          "Creates a project from a GitHub repository",
	"AlmIntegrations.ImportGitlabProject":          "Imports a GitLab project to SonarQube",
	"AlmIntegrations.ListAzureProjects":            "Lists Azure DevOps projects",
	"AlmIntegrations.ListBitbucketServerProjects":  "Lists Bitbucket Server projects",
	"AlmIntegrations.ListGithubOrganizations":      "Lists GitHub organizations",
	"AlmIntegrations.ListGithubRepositories":       "Lists GitHub repositories for an organization",
	"AlmIntegrations.SearchAzureRepos":             "Searches Azure DevOps repositories",
	"AlmIntegrations.SearchBitbucketCloudRepos":    "Searches Bitbucket Cloud repositories",
	"AlmIntegrations.SearchBitbucketServerRepos":   "Searches Bitbucket Server repositories",
	"AlmIntegrations.SearchGitlabRepos":            "Searches GitLab repositories",
	"AlmIntegrations.SetPat":                       "Sets a Personal Access Token",

	// AlmSettings
	"AlmSettings.CountBinding":         "Counts projects bound to a DevOps platform setting",
	"AlmSettings.CreateAzure":          "Creates an Azure DevOps instance setting",
	"AlmSettings.CreateBitbucket":      "Creates a Bitbucket Server instance setting",
	"AlmSettings.CreateBitbucketCloud": "Creates a Bitbucket Cloud instance setting",
	"AlmSettings.CreateGithub":         "Creates a GitHub instance setting",
	"AlmSettings.CreateGitlab":         "Creates a GitLab instance setting",
	"AlmSettings.Delete":               "Deletes a DevOps platform setting",
	"AlmSettings.GetBinding":           "Gets the DevOps binding of a project",
	"AlmSettings.List":                 "Lists DevOps platform settings",
	"AlmSettings.ListDefinitions":      "Lists DevOps platform setting definitions",
	"AlmSettings.UpdateAzure":          "Updates an Azure DevOps instance setting",
	"AlmSettings.UpdateBitbucket":      "Updates a Bitbucket Server instance setting",
	"AlmSettings.UpdateBitbucketCloud": "Updates a Bitbucket Cloud instance setting",
	"AlmSettings.UpdateGithub":         "Updates a GitHub instance setting",
	"AlmSettings.UpdateGitlab":         "Updates a GitLab instance setting",
	"AlmSettings.Validate":             "Validates a DevOps platform setting",

	// AnalysisCache
	"AnalysisCache.Clear": "Clears scanner cached data",
	"AnalysisCache.Get":   "Gets scanner cached data for a project",

	// AnalysisReports
	"AnalysisReports.IsQueueEmpty": "Checks if the CE queue is empty",

	// Authentication
	"Authentication.Login":    "Authenticates a user with credentials",
	"Authentication.Logout":   "Logs out the current user",
	"Authentication.Validate": "Checks if current credentials are valid",

	// Batch
	"Batch.File":    "Downloads a JAR file from the batch index",
	"Batch.Index":   "Lists JAR files for scanners",
	"Batch.Project": "Returns project repository info",

	// Ce
	"Ce.Activity":               "Searches CE activity history",
	"Ce.ActivityStatus":         "Returns CE activity metrics",
	"Ce.AnalysisStatus":         "Gets analysis status of a component",
	"Ce.Cancel":                 "Cancels a pending CE task",
	"Ce.CancelAll":              "Cancels all pending CE tasks",
	"Ce.Component":              "Gets CE activity for a component",
	"Ce.DismissAnalysisWarning": "Dismisses a specific analysis warning",
	"Ce.IndexationStatus":       "Returns indexation status",
	"Ce.Info":                   "Gets Compute Engine information",
	"Ce.Pause":                  "Pauses Compute Engine workers",
	"Ce.Resume":                 "Resumes paused CE workers",
	"Ce.Submit":                 "Submits an analysis report",
	"Ce.Task":                   "Gets details of a CE task",
	"Ce.TaskTypes":              "Lists available CE task types",
	"Ce.WorkerCount":            "Returns the number of CE workers",

	// Components
	"Components.App":            "Gets data for a component page",
	"Components.Search":         "Searches for components by name",
	"Components.SearchProjects": "Searches projects with facets",
	"Components.Show":           "Returns component details",
	"Components.Suggestions":    "Provides search autocomplete suggestions",
	"Components.Tree":           "Navigates component tree",

	// Developers
	"Developers.SearchEvents": "Searches developer events",

	// DismissMessage
	"DismissMessage.Check":   "Checks if a message has been dismissed",
	"DismissMessage.Dismiss": "Dismisses a UI message",

	// Duplications
	"Duplications.Show": "Shows duplications for a file",

	// Emails
	"Emails.Send": "Sends an email via SonarQube",

	// Favorites
	"Favorites.Add":    "Adds a component to favorites",
	"Favorites.Remove": "Removes a component from favorites",
	"Favorites.Search": "Searches favorite components",

	// Features
	"Features.List": "Lists all available features",

	// GithubProvisioning
	"GithubProvisioning.Check": "Checks GitHub provisioning status",

	// Hotspots
	"Hotspots.AddComment":    "Adds a comment to a hotspot",
	"Hotspots.Assign":        "Assigns a hotspot to a user",
	"Hotspots.ChangeStatus":  "Changes a hotspot's status",
	"Hotspots.DeleteComment": "Deletes a hotspot comment",
	"Hotspots.EditComment":   "Edits a hotspot comment",
	"Hotspots.List":          "Lists hotspots for a project",
	"Hotspots.Pull":          "Pulls hotspot data for sync",
	"Hotspots.Search":        "Searches for security hotspots",
	"Hotspots.Show":          "Shows hotspot details",

	// Issues
	"Issues.AddComment":             "Adds a comment to an issue",
	"Issues.AnticipatedTransitions": "Applies anticipated transitions",
	"Issues.Assign":                 "Assigns or unassigns an issue",
	"Issues.Authors":                "Searches SCM accounts",
	"Issues.BulkChange":             "Bulk changes up to 500 issues",
	"Issues.Changelog":              "Displays issue changelog",
	"Issues.ComponentTags":          "Lists tags for component issues",
	"Issues.DeleteComment":          "Deletes an issue comment",
	"Issues.DoTransition":           "Performs an issue transition",
	"Issues.EditComment":            "Edits an issue comment",
	"Issues.List":                   "Lists issues for a project",
	"Issues.Pull":                   "Pulls issue data for sync",
	"Issues.PullTaint":              "Pulls taint vulnerability data",
	"Issues.Reindex":                "Triggers issue reindexing",
	"Issues.Search":                 "Searches issues with filters",
	"Issues.SetSeverity":            "Changes an issue's severity",
	"Issues.SetTags":                "Sets tags on an issue",
	"Issues.SetType":                "Changes an issue's type",
	"Issues.Tags":                   "Lists rule tags",

	// L10N
	"L10N.Index": "Returns localized messages",

	// Languages
	"Languages.List": "Lists supported languages",

	// Measures
	"Measures.Component":     "Returns component measures",
	"Measures.ComponentTree": "Navigates component tree with measures",
	"Measures.Search":        "Searches project measures",
	"Measures.SearchHistory": "Searches measure history",

	// Metrics
	"Metrics.Search": "Searches for metrics",
	"Metrics.Types":  "Lists available metric types",

	// Monitoring
	"Monitoring.Metrics": "Returns monitoring metrics",

	// Navigation
	"Navigation.Component":   "Returns component navigation data",
	"Navigation.Global":      "Returns global navigation data",
	"Navigation.Marketplace": "Returns marketplace navigation data",
	"Navigation.Settings":    "Returns settings navigation data",

	// NewCodePeriods
	"NewCodePeriods.List":  "Lists new code period definitions",
	"NewCodePeriods.Set":   "Sets a new code period definition",
	"NewCodePeriods.Show":  "Shows the new code period definition",
	"NewCodePeriods.Unset": "Unsets the new code period definition",

	// Notifications
	"Notifications.Add":    "Adds a user notification",
	"Notifications.List":   "Lists user notifications",
	"Notifications.Remove": "Removes a user notification",

	// Permissions
	"Permissions.AddGroup":                         "Adds permission to a group",
	"Permissions.AddGroupToTemplate":               "Adds a group to a permission template",
	"Permissions.AddProjectCreatorToTemplate":      "Adds project creator to a template",
	"Permissions.AddUser":                          "Adds permission to a user",
	"Permissions.AddUserToTemplate":                "Adds a user to a permission template",
	"Permissions.ApplyTemplate":                    "Applies a permission template",
	"Permissions.BulkApplyTemplate":                "Bulk applies a permission template",
	"Permissions.CreateTemplate":                   "Creates a permission template",
	"Permissions.DeleteTemplate":                   "Deletes a permission template",
	"Permissions.Groups":                           "Lists groups with permissions",
	"Permissions.RemoveGroup":                      "Removes permission from a group",
	"Permissions.RemoveGroupFromTemplate":          "Removes a group from a template",
	"Permissions.RemoveProjectCreatorFromTemplate": "Removes project creator from a template",
	"Permissions.RemoveUser":                       "Removes permission from a user",
	"Permissions.RemoveUserFromTemplate":           "Removes a user from a template",
	"Permissions.SearchTemplates":                  "Searches permission templates",
	"Permissions.SetDefaultTemplate":               "Sets a default permission template",
	"Permissions.TemplateGroups":                   "Lists groups in a template",
	"Permissions.TemplateUsers":                    "Lists users in a template",
	"Permissions.UpdateTemplate":                   "Updates a permission template",
	"Permissions.Users":                            "Lists users with permissions",

	// Plugins
	"Plugins.Available": "Lists plugins available for installation",
	"Plugins.CancelAll": "Cancels pending plugin operations",
	"Plugins.Download":  "Downloads a plugin JAR file",
	"Plugins.Install":   "Installs a plugin from marketplace",
	"Plugins.Installed": "Lists installed plugins",
	"Plugins.Pending":   "Lists plugins with pending operations",
	"Plugins.Uninstall": "Uninstalls a plugin",
	"Plugins.Update":    "Updates an installed plugin",
	"Plugins.Updates":   "Lists plugins with available updates",

	// ProjectAnalyses
	"ProjectAnalyses.CreateEvent": "Creates a project analysis event",
	"ProjectAnalyses.Delete":      "Deletes a project analysis",
	"ProjectAnalyses.DeleteEvent": "Deletes an analysis event",
	"ProjectAnalyses.Search":      "Searches project analyses",
	"ProjectAnalyses.SearchAll":   "Searches all project analyses",
	"ProjectAnalyses.UpdateEvent": "Updates an analysis event",

	// ProjectBadges
	"ProjectBadges.Measure":     "Generates a measure badge",
	"ProjectBadges.QualityGate": "Generates a quality gate badge",
	"ProjectBadges.RenewToken":  "Renews a project badge token",
	"ProjectBadges.Token":       "Returns a project badge token",

	// ProjectBranches
	"ProjectBranches.Delete":                         "Deletes a project branch",
	"ProjectBranches.List":                           "Lists project branches",
	"ProjectBranches.Rename":                         "Renames a project branch",
	"ProjectBranches.SetAutomaticDeletionProtection": "Sets deletion protection on a branch",
	"ProjectBranches.SetMain":                        "Sets the main branch",

	// ProjectDump
	"ProjectDump.Export": "Exports a project dump",
	"ProjectDump.Status": "Gets project dump status",

	// ProjectLinks
	"ProjectLinks.Create": "Creates a project link",
	"ProjectLinks.Delete": "Deletes a project link",
	"ProjectLinks.Search": "Searches project links",

	// ProjectTags
	"ProjectTags.Search": "Searches project tags",
	"ProjectTags.Set":    "Sets tags on a project",

	// Projects
	"Projects.BulkDelete":                "Deletes multiple projects in bulk",
	"Projects.Create":                    "Creates a new project",
	"Projects.Delete":                    "Deletes a project",
	"Projects.Search":                    "Searches for projects",
	"Projects.SearchMyProjects":          "Searches projects the user administers",
	"Projects.SearchMyScannableProjects": "Searches scannable projects",
	"Projects.UpdateDefaultVisibility":   "Updates default project visibility",
	"Projects.UpdateKey":                 "Updates a project key",
	"Projects.UpdateVisibility":          "Updates a project's visibility",

	// Push
	"Push.SonarlintEvents": "Streams server-sent events for SonarLint",

	// Qualitygates
	"Qualitygates.AddGroup":        "Adds a group to a quality gate",
	"Qualitygates.AddUser":         "Adds a user to a quality gate",
	"Qualitygates.Copy":            "Copies a quality gate",
	"Qualitygates.Create":          "Creates a new quality gate",
	"Qualitygates.CreateCondition": "Adds a condition to a quality gate",
	"Qualitygates.DeleteCondition": "Deletes a quality gate condition",
	"Qualitygates.Deselect":        "Removes quality gate from a project",
	"Qualitygates.Destroy":         "Deletes a quality gate",
	"Qualitygates.GetByProject":    "Gets a project's quality gate",
	"Qualitygates.List":            "Lists all quality gates",
	"Qualitygates.ProjectStatus":   "Gets quality gate status of a project",
	"Qualitygates.RemoveGroup":     "Removes a group from a quality gate",
	"Qualitygates.RemoveUser":      "Removes a user from a quality gate",
	"Qualitygates.Rename":          "Renames a quality gate",
	"Qualitygates.Search":          "Searches quality gate projects",
	"Qualitygates.SearchGroups":    "Searches quality gate groups",
	"Qualitygates.SearchUsers":     "Searches quality gate users",
	"Qualitygates.Select":          "Associates project with a quality gate",
	"Qualitygates.SetAsDefault":    "Sets a default quality gate",
	"Qualitygates.Show":            "Shows quality gate details",
	"Qualitygates.UpdateCondition": "Updates a quality gate condition",

	// Qualityprofiles
	"Qualityprofiles.ActivateRule":    "Activates a rule in a profile",
	"Qualityprofiles.ActivateRules":   "Bulk activates rules in a profile",
	"Qualityprofiles.AddGroup":        "Adds a group to a profile",
	"Qualityprofiles.AddProject":      "Associates project with a profile",
	"Qualityprofiles.AddUser":         "Adds a user to a profile",
	"Qualityprofiles.Backup":          "Backs up a profile in XML format",
	"Qualityprofiles.ChangeParent":    "Changes a profile's parent",
	"Qualityprofiles.Changelog":       "Gets profile change history",
	"Qualityprofiles.Compare":         "Compares two profiles",
	"Qualityprofiles.Copy":            "Copies a quality profile",
	"Qualityprofiles.Create":          "Creates a new quality profile",
	"Qualityprofiles.DeactivateRule":  "Deactivates a rule in a profile",
	"Qualityprofiles.DeactivateRules": "Bulk deactivates rules in a profile",
	"Qualityprofiles.Delete":          "Deletes a quality profile",
	"Qualityprofiles.Export":          "Exports a profile definition",
	"Qualityprofiles.Exporters":       "Lists profile exporters",
	"Qualityprofiles.Importers":       "Lists profile importers",
	"Qualityprofiles.Inheritance":     "Shows profile ancestors and children",
	"Qualityprofiles.Projects":        "Lists projects on a profile",
	"Qualityprofiles.RemoveGroup":     "Removes a group from a profile",
	"Qualityprofiles.RemoveProject":   "Removes project from a profile",
	"Qualityprofiles.RemoveUser":      "Removes a user from a profile",
	"Qualityprofiles.Rename":          "Renames a quality profile",
	"Qualityprofiles.Restore":         "Restores a profile from backup",
	"Qualityprofiles.Search":          "Searches quality profiles",
	"Qualityprofiles.SearchGroups":    "Searches profile groups",
	"Qualityprofiles.SearchUsers":     "Searches profile users",
	"Qualityprofiles.SetDefault":      "Sets default profile for a language",
	"Qualityprofiles.Show":            "Shows profile details",

	// Rules
	"Rules.App":          "Gets Coding Rules page data",
	"Rules.Create":       "Creates a custom coding rule",
	"Rules.Delete":       "Deletes a custom coding rule",
	"Rules.List":         "Lists rules excluding external ones",
	"Rules.Repositories": "Lists rule repositories",
	"Rules.Search":       "Searches for coding rules",
	"Rules.Show":         "Shows details of a rule",
	"Rules.Tags":         "Lists available rule tags",
	"Rules.Update":       "Updates a coding rule",

	// Server
	"Server.Version": "Returns the SonarQube server version",

	// Settings
	"Settings.CheckSecretKey":    "Checks for a secret key",
	"Settings.Encrypt":           "Encrypts a value",
	"Settings.GenerateSecretKey": "Generates a new secret key",
	"Settings.ListDefinitions":   "Lists setting definitions",
	"Settings.LoginMessage":      "Gets the login page message",
	"Settings.Reset":             "Resets a setting to default",
	"Settings.Set":               "Sets a setting value",
	"Settings.Values":            "Returns setting values",

	// Sources
	"Sources.Index":         "Gets source code index",
	"Sources.IssueSnippets": "Gets source snippets for issues",
	"Sources.Lines":         "Returns source code lines",
	"Sources.Raw":           "Returns raw source code",
	"Sources.Scm":           "Returns SCM information",
	"Sources.Show":          "Shows source code of a file",

	// System
	"System.ChangeLogLevel":    "Changes the system log level",
	"System.DbMigrationStatus": "Gets DB migration status",
	"System.Health":            "Returns system health status",
	"System.Info":              "Returns system information",
	"System.Liveness":          "Checks system liveness",
	"System.Logs":              "Returns system log content",
	"System.MigrateDb":         "Triggers database migration",
	"System.Ping":              "Pings the server",
	"System.Restart":           "Restarts the SonarQube server",
	"System.Status":            "Returns server status",
	"System.Upgrades":          "Lists available upgrades",

	// UserGroups
	"UserGroups.AddUser":    "Adds a user to a group",
	"UserGroups.Create":     "Creates a user group",
	"UserGroups.Delete":     "Deletes a user group",
	"UserGroups.RemoveUser": "Removes a user from a group",
	"UserGroups.Search":     "Searches for user groups",
	"UserGroups.Update":     "Updates a user group",
	"UserGroups.Users":      "Lists users in a group",

	// UserTokens
	"UserTokens.Generate": "Generates a new user token",
	"UserTokens.Revoke":   "Revokes a user token",
	"UserTokens.Search":   "Searches for user tokens",

	// Users
	"Users.Anonymize":              "Anonymizes a deactivated user",
	"Users.ChangePassword":         "Changes a user password",
	"Users.Create":                 "Creates a new user",
	"Users.Current":                "Returns the current user",
	"Users.Deactivate":             "Deactivates a user account",
	"Users.DismissNotice":          "Dismisses a notice for the user",
	"Users.Groups":                 "Lists groups of a user",
	"Users.IdentityProviders":      "Lists identity providers",
	"Users.Search":                 "Searches for users",
	"Users.SetHomepage":            "Sets the user's homepage",
	"Users.Update":                 "Updates a user",
	"Users.UpdateIdentityProvider": "Updates a user's identity provider",
	"Users.UpdateLogin":            "Updates a user's login",

	// Webhooks
	"Webhooks.Create":     "Creates a webhook",
	"Webhooks.Delete":     "Deletes a webhook",
	"Webhooks.Deliveries": "Lists webhook deliveries",
	"Webhooks.Delivery":   "Returns webhook delivery details",
	"Webhooks.List":       "Lists webhooks",
	"Webhooks.Update":     "Updates a webhook",

	// Webservices
	"Webservices.List":            "Lists web services and actions",
	"Webservices.ResponseExample": "Returns a response example",
}

// GetServiceDescription returns the description for a service, or a default.
func GetServiceDescription(serviceName string) string {
	if desc, ok := serviceDescriptions[serviceName]; ok {
		return desc
	}

	return "Commands for the " + serviceName + " service"
}

// GetMethodDescription returns the description for a method, or a default.
func GetMethodDescription(serviceName, methodName string) string {
	key := serviceName + "." + methodName
	if desc, ok := methodDescriptions[key]; ok {
		return desc
	}

	return "Call " + serviceName + "." + methodName
}
